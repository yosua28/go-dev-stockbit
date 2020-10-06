package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetAdminCmsPostList(c echo.Context) error {
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

	paramType := make(map[string]string)
	field := c.QueryParam("field")
	if field != "" {
		keyStr := c.QueryParam("key")
		if keyStr == "" {
			log.Error("Wrong value for parameter: key")
			return lib.CustomError(http.StatusBadRequest, "parameter key not allowed empty if parameter field not empty", "parameter key not allowed empty if parameter field not empty")
		}
		key, _ := strconv.ParseUint(keyStr, 10, 64)
		if key == 0 {
			return lib.CustomError(http.StatusNotFound)
		}

		if field == "type" {
			paramType["post_type_key"] = keyStr
		} else if field == "subtype" {
			paramType["post_subtype_key"] = keyStr
		} else {
			log.Error("Wrong value for parameter: field")
			return lib.CustomError(http.StatusBadRequest, "Wrong value for parameter: field", "Wrong value for parameter: field")
		}
	}

	params := make(map[string]string)

	items := []string{"post_title", "post_sub_title", "post_content_author", "post_publish_thru", "post_publish_start"}

	// Get parameter order_by
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
	}

	postSubtypeData := make(map[uint64]models.CmsPostSubtype)
	var subtypeIdsParamCount []string
	var postSubtypeIDs []string
	var postTypeIDs []string
	postTypeIDs = append(postTypeIDs, strconv.FormatUint(0, 10))

	var posts []models.CmsPost

	if len(paramType) > 0 {
		var postSubtypeDB []models.CmsPostSubtype
		status, err = models.GetAllCmsPostSubtype(&postSubtypeDB, limit, offset, paramType, true)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if len(postSubtypeDB) < 1 {
			log.Error("post not found")
			return lib.CustomError(http.StatusNotFound, "Post subtype not found", "Post subtype not found")
		}

		for _, postSubtyp := range postSubtypeDB {
			postSubtypeData[postSubtyp.PostSubtypeKey] = postSubtyp
			postSubtypeIDs = append(postSubtypeIDs, strconv.FormatUint(postSubtyp.PostSubtypeKey, 10))
			subtypeIdsParamCount = append(subtypeIdsParamCount, strconv.FormatUint(postSubtyp.PostSubtypeKey, 10))
			postTypeIDs = append(postTypeIDs, strconv.FormatUint(postSubtyp.PostTypeKey, 10))
		}

		status, err = models.GetAdminCmsPostListIn(&posts, limit, offset, noLimit, params, postSubtypeIDs)
		if err != nil {
			return lib.CustomError(status)
		}
		if len(posts) < 1 {
			return lib.CustomError(http.StatusNotFound)
		}
	} else {
		status, err = models.GetAdminCmsPostListIn(&posts, limit, offset, noLimit, params, postSubtypeIDs)
		if err != nil {
			return lib.CustomError(status)
		}
		if len(posts) < 1 {
			return lib.CustomError(http.StatusNotFound)
		}
		for _, postData := range posts {
			postSubtypeIDs = append(postSubtypeIDs, strconv.FormatUint(postData.PostSubtypeKey, 10))
		}

		var postSubtypeDB []models.CmsPostSubtype
		status, err = models.GetPostSubtypeIn(&postSubtypeDB, postSubtypeIDs, "post_subtype_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		} else {
			for _, postSubtyp := range postSubtypeDB {
				postSubtypeData[postSubtyp.PostSubtypeKey] = postSubtyp
				postTypeIDs = append(postTypeIDs, strconv.FormatUint(postSubtyp.PostTypeKey, 10))
			}
		}

	}

	ptData := make(map[uint64]models.CmsPostType)
	var postTypeDB []models.CmsPostType
	status, err = models.GetCmsPostTypeIn(&postTypeDB, postTypeIDs, "post_type_key")
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		for _, pt := range postTypeDB {
			ptData[pt.PostTypeKey] = pt
		}
	}

	var responseData []models.AdminCmsPostList
	for _, post := range posts {
		var data models.AdminCmsPostList

		data.PostKey = post.PostKey

		var dir string

		if n, ok := postSubtypeData[post.PostSubtypeKey]; ok {
			data.PostSubtypeKey = n.PostSubtypeKey
			data.PostSubtypeCode = n.PostSubtypeCode
			data.PostSubtypeName = n.PostSubtypeName

			if p, ok := ptData[n.PostTypeKey]; ok {
				data.PostTypeKey = p.PostTypeKey
				data.PostTypeCode = p.PostTypeCode
				data.PostTypeName = p.PostTypeName
				data.PostTypeDesc = p.PostTypeDesc
				dir = strings.ToLower(p.PostTypeCode)
			}
		}

		data.PostTitle = post.PostTitle
		data.PostSubTitle = post.PostSubTitle
		data.PostContentAuthor = post.PostContentAuthor
		data.PostContentSources = post.PostContentSources
		if post.PostPinned > 0 {
			data.PostPinned = true
		}
		if post.RecImage1 != nil && *post.RecImage1 != "" {
			data.RecImage1 = config.BaseUrl + "/images/post/" + dir + "/" + *post.RecImage1
		} else {
			data.RecImage1 = config.BaseUrl + "/images/post/default.png"
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, post.PostPublishStart)
		data.PostPublishStart = date.Format(newLayout)
		date, _ = time.Parse(layout, post.PostPublishThru)
		data.PostPublishThru = date.Format(newLayout)
		responseData = append(responseData, data)
	}

	var countData models.CmsPostCount
	var pagination int
	if limit > 0 {
		status, err = models.GetCountCmsPost(&countData, params, subtypeIdsParamCount)
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
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
