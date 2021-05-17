package controllers

import (
	"api/lib"
	"api/models"
	"database/sql"
	"net/http"
	"strconv"

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
