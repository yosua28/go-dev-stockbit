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
		log.Error("Error delete sc_menu")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMenu(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	menuParent := c.FormValue("menu_parent")
	if menuParent != "" {
		n, err := strconv.ParseUint(menuParent, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_parent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_parent", "Wrong input for parameter: menu_parent")
		}
		params["menu_parent"] = menuParent
		params["rec_attribute_id2"] = "CSidebarNavDropdown"
	}

	appModuleKey := c.FormValue("app_module_key")
	if appModuleKey == "" {
		log.Error("Missing required parameter: app_module_key")
		return lib.CustomError(http.StatusBadRequest, "app_module_key can not be blank", "app_module_key can not be blank")
	} else {
		n, err := strconv.ParseUint(appModuleKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: app_module_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: app_module_key", "Wrong input for parameter: app_module_key")
		}
		params["app_module_key"] = appModuleKey
	}

	menuCode := c.FormValue("menu_code")
	if menuCode == "" {
		log.Error("Missing required parameter: menu_code")
		return lib.CustomError(http.StatusBadRequest, "menu_code can not be blank", "menu_code can not be blank")
	} else {
		//validate unique menu_code
		var countData models.CountData
		status, err = models.CountScMenuValidateUnique(&countData, "menu_code", menuCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: menu_code")
			return lib.CustomError(http.StatusBadRequest, "menu_code already used", "menu_code already used")
		}
		params["menu_code"] = menuCode
	}

	menuName := c.FormValue("menu_name")
	if menuName == "" {
		log.Error("Missing required parameter: menu_name")
		return lib.CustomError(http.StatusBadRequest, "menu_name can not be blank", "menu_name can not be blank")
	} else {
		//validate unique menu_name
		var countData models.CountData
		status, err = models.CountScMenuValidateUnique(&countData, "menu_name", menuName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: menu_name")
			return lib.CustomError(http.StatusBadRequest, "menu_name already used", "menu_name already used")
		}
		params["menu_name"] = menuName
	}

	menuPage := c.FormValue("menu_page")
	if menuPage != "" {
		params["menu_page"] = menuPage
	}

	menuUrl := c.FormValue("menu_url")
	if menuUrl != "" {
		params["menu_url"] = menuUrl
	}

	menuTypeKey := c.FormValue("menu_type_key")
	if menuTypeKey == "" {
		log.Error("Missing required parameter: menu_type_key")
		return lib.CustomError(http.StatusBadRequest, "menu_type_key can not be blank", "menu_type_key can not be blank")
	} else {
		n, err := strconv.ParseUint(menuTypeKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_type_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_type_key", "Wrong input for parameter: menu_type_key")
		}
		params["menu_type_key"] = menuTypeKey
	}
	params["has_endpoint"] = "0"

	menuDesc := c.FormValue("menu_desc")
	if menuDesc != "" {
		params["menu_desc"] = menuDesc
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	icon := c.FormValue("icon")
	if icon != "" {
		params["rec_attribute_id1"] = icon
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateScMenu(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateMenu(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	menuKey := c.FormValue("menu_key")
	if menuKey == "" {
		log.Error("Missing required parameter: menu_key")
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	} else {
		n, err := strconv.ParseUint(menuKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_key", "Wrong input for parameter: menu_key")
		}
		var menu models.ScMenu
		status, err = models.GetScMenu(&menu, menuKey)
		if err != nil {
			log.Error("Menu not found")
			return lib.CustomError(http.StatusBadRequest, "Menu not found", "Menu not found")
		}
		params["menu_key"] = menuKey
	}

	menuParent := c.FormValue("menu_parent")
	if menuParent != "" {
		n, err := strconv.ParseUint(menuParent, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_parent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_parent", "Wrong input for parameter: menu_parent")
		}
		params["menu_parent"] = menuParent
		params["rec_attribute_id2"] = "CSidebarNavDropdown"
	} else {
		params["rec_attribute_id2"] = ""
	}

	appModuleKey := c.FormValue("app_module_key")
	if appModuleKey == "" {
		log.Error("Missing required parameter: app_module_key")
		return lib.CustomError(http.StatusBadRequest, "app_module_key can not be blank", "app_module_key can not be blank")
	} else {
		n, err := strconv.ParseUint(appModuleKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: app_module_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: app_module_key", "Wrong input for parameter: app_module_key")
		}
		params["app_module_key"] = appModuleKey
	}

	menuCode := c.FormValue("menu_code")
	if menuCode == "" {
		log.Error("Missing required parameter: menu_code")
		return lib.CustomError(http.StatusBadRequest, "menu_code can not be blank", "menu_code can not be blank")
	} else {
		//validate unique menu_code
		var countData models.CountData
		status, err = models.CountScMenuValidateUnique(&countData, "menu_code", menuCode, menuKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: menu_code")
			return lib.CustomError(http.StatusBadRequest, "menu_code already used", "menu_code already used")
		}
		params["menu_code"] = menuCode
	}

	menuName := c.FormValue("menu_name")
	if menuName == "" {
		log.Error("Missing required parameter: menu_name")
		return lib.CustomError(http.StatusBadRequest, "menu_name can not be blank", "menu_name can not be blank")
	} else {
		//validate unique menu_name
		var countData models.CountData
		status, err = models.CountScMenuValidateUnique(&countData, "menu_name", menuName, menuKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: menu_name")
			return lib.CustomError(http.StatusBadRequest, "menu_name already used", "menu_name already used")
		}
		params["menu_name"] = menuName
	}

	menuPage := c.FormValue("menu_page")
	if menuPage != "" {
		params["menu_page"] = menuPage
	}

	menuUrl := c.FormValue("menu_url")
	if menuUrl != "" {
		params["menu_url"] = menuUrl
	}

	menuTypeKey := c.FormValue("menu_type_key")
	if menuTypeKey == "" {
		log.Error("Missing required parameter: menu_type_key")
		return lib.CustomError(http.StatusBadRequest, "menu_type_key can not be blank", "menu_type_key can not be blank")
	} else {
		n, err := strconv.ParseUint(menuTypeKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_type_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_type_key", "Wrong input for parameter: menu_type_key")
		}
		params["menu_type_key"] = menuTypeKey
	}
	params["has_endpoint"] = "0"

	menuDesc := c.FormValue("menu_desc")
	if menuDesc != "" {
		params["menu_desc"] = menuDesc
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	icon := c.FormValue("icon")
	if icon != "" {
		params["rec_attribute_id1"] = icon
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateScMenu(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminDetailMenu(c echo.Context) error {
	var err error

	menuKey := c.Param("menu_key")
	if menuKey == "" {
		log.Error("Missing required parameter: menu_key")
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	} else {
		n, err := strconv.ParseUint(menuKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_key", "Wrong input for parameter: menu_key")
		}
	}
	var menu models.ScMenu
	_, err = models.GetScMenu(&menu, menuKey)
	if err != nil {
		log.Error("Menu not found")
		return lib.CustomError(http.StatusBadRequest, "Menu not found", "Menu not found")
	}

	responseData := make(map[string]interface{})
	responseData["menu_key"] = menu.MenuKey
	if menu.MenuParent != nil {
		responseData["menu_parent"] = *menu.MenuParent
	} else {
		responseData["menu_parent"] = ""
	}
	responseData["app_module_key"] = menu.AppModuleKey
	responseData["menu_code"] = menu.MenuCode
	responseData["menu_name"] = menu.MenuName
	if menu.MenuPage != nil {
		responseData["menu_page"] = *menu.MenuPage
	} else {
		responseData["menu_page"] = ""
	}
	if menu.MenuURL != nil {
		responseData["menu_url"] = *menu.MenuURL
	} else {
		responseData["menu_url"] = ""
	}
	responseData["menu_type_key"] = menu.MenuTypeKey
	if menu.MenuDesc != nil {
		responseData["menu_desc"] = *menu.MenuDesc
	} else {
		responseData["menu_desc"] = ""
	}
	if menu.RecOrder != nil {
		responseData["rec_order"] = *menu.RecOrder
	} else {
		responseData["rec_order"] = ""
	}
	if menu.RecAttributeID1 != nil {
		responseData["icon"] = *menu.RecAttributeID1
	} else {
		responseData["icon"] = ""
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
