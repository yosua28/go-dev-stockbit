package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListMenuLogin(c echo.Context) error {

	var err error

	//mapping role management
	var listMenuRoleUser []models.ListMenuRoleUser
	_, err = models.AdminGetMenuListRoleLogin(&listMenuRoleUser, strconv.FormatUint(lib.Profile.RoleKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var menuParentIds []string
	for _, mn := range listMenuRoleUser {
		if mn.MenuParent != nil {
			if _, ok := lib.Find(menuParentIds, strconv.FormatUint(*mn.MenuParent, 10)); !ok {
				menuParentIds = append(menuParentIds, strconv.FormatUint(*mn.MenuParent, 10))
			}
		}
	}

	var listParentMenuRoleUser []models.ListParentMenuRoleUser
	if len(menuParentIds) > 0 {
		_, err = models.AdminGetParentMenuListRoleLogin(&listParentMenuRoleUser, menuParentIds)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
			}
		}
	}

	var responseData []models.MenuUserRole

	//manual dashboard
	var dashboard models.MenuUserRole
	var cName string
	cName = "CSidebarNavItem"
	dashboard.ClassName = &cName
	var name string
	name = "Dashboard"
	dashboard.Name = &name
	var to string
	to = "/dashboard"
	dashboard.To = &to
	var icon string
	icon = "cil-speedometer"
	dashboard.Icon = &icon
	responseData = append(responseData, dashboard)
	//end manual dashboard

	for _, parent := range listParentMenuRoleUser {
		var data models.MenuUserRole
		data.ClassName = parent.ClassName
		data.Name = parent.MenuPage
		data.Icon = parent.Icon

		var child []models.MenuChild
		for _, c := range listMenuRoleUser {
			if parent.MenuKey == *c.MenuParent {
				var cc models.MenuChild
				cc.Name = c.MenuPage
				cc.To = c.MenuURL
				cc.Icon = c.Icon
				child = append(child, cc)
			}
		}
		data.Items = &child

		if len(child) > 0 {
			responseData = append(responseData, data)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMenu(c echo.Context) error {

	var err error
	var status int
	//Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			log.Error("Limit should be number")
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			log.Error("Page should be number")
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	noLimitStr := c.QueryParam("nolimit")
	var noLimit bool
	if noLimitStr != "" {
		noLimit, err = strconv.ParseBool(noLimitStr)
		if err != nil {
			log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"menu_key", "parent_name", "menu_code", "menu_name", "menu_page", "menu_desc", "menu_url", "app_module_name", "menu_type_name"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			params["orderBy"] = orderBy
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "menu_code"
		params["orderType"] = "ASC"
	}

	searchLike := c.QueryParam("search_like")

	appModulKey := c.QueryParam("app_module_key")
	if appModulKey != "" {
		params["m.app_module_key"] = appModulKey
	}
	menuTypeKey := c.QueryParam("menu_type_key")
	if menuTypeKey != "" {
		params["m.menu_type_key"] = menuTypeKey
	}

	var menu []models.ListMenuAdmin

	status, err = models.AdminGetListMenu(&menu, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(menu) < 1 {
		log.Error("Menu not found")
		return lib.CustomError(http.StatusNotFound, "Menu not found", "Menu not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListMenu(&countData, params, searchLike)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			calc := math.Ceil(float64(countData.CountData) / float64(limit))
			pagination = int(calc)
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = menu

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMenu(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("menu_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: menu_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: menu_key", "Missing required parameter: menu_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["menu_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScMenu(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
