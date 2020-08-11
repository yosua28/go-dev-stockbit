package controllers

import (
	"api/models"
	_ "api/config"
	"api/lib"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetCmsPostTypeData(c echo.Context) error {
	var err error
	var status int
	// Get parameter limit
	// limitStr := c.QueryParam("limit")
	// var limit uint64
	// if limitStr != "" {
	// 	limit, err = strconv.ParseUint(limitStr, 10, 64)
	// 	if err == nil {
	// 		if (limit == 0) || (limit > config.LimitQuery) {
	// 			limit = config.LimitQuery
	// 		}
	// 	} else {
	// 		return lib.CustomError(http.StatusBadRequest)
	// 	}
	// } else {
	// 	limit = config.LimitQuery
	// }
	// // Get parameter page
	// pageStr := c.QueryParam("page")
	// var page uint64
	// if pageStr != "" {
	// 	page, err = strconv.ParseUint(pageStr, 10, 64)
	// 	if err == nil {
	// 		if page == 0 {
	// 			page = 1
	// 		}
	// 	} else {
	// 		return lib.CustomError(http.StatusBadRequest)
	// 	}
	// } else {
	// 	page = 1
	// }
	// var offset uint64
	// if page > 1 {
	// 	offset = limit * (page - 1)
	// }

	// noLimitStr := c.QueryParam("nolimit")
	// var noLimit bool
	// if noLimitStr != "" {
	// 	noLimit, err = strconv.ParseBool(noLimitStr)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest)
	// 	}
	// } else {
	// 	noLimit = false
	// }

	params := make(map[string]string)
	// Get parameter post type
	postType := c.QueryParam("post_type")
	if postType == "" {
		return lib.CustomError(http.StatusBadRequest)
	}
	var postTypeDB models.CmsPostType
	status, err = models.GetCmsPostType(&postTypeDB, "post_type_code", postType)
	if err != nil {
		return lib.CustomError(status)
	}

	var postSubtypeDB []models.CmsPostSubtype
	status, err = models.GetAllCmsPostSubtype(&postSubtypeDB, 0, 0, nil, true)

	postSubtypeData := make(map[uint64]models.CmsPostSubtype)
	var postSubtypeIDs []string
	for _, postSubtyp := range postSubtypeDB {
		postSubtypeData[postSubtyp.PostSubtypeKey] = postSubtyp
		postSubtypeIDs = append(postSubtypeIDs, strconv.FormatUint(postSubtyp.PostSubtypeKey, 10))
	}

	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "post_title") || (orderBy == "post_publish_thru") || (orderBy == "post_publish_start") {
			params["orderBy"] = orderBy
		}else{
			return lib.CustomError(http.StatusBadRequest)
		}
	}
	// Get parameter order_type
	orderType := c.QueryParam("order_type")
	if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
		params["orderType"] = orderType
	}
	var posts []models.CmsPost
	status, err = models.GetCmsPostIn(&posts,  postSubtypeIDs, "post_subtype_key")
	if err != nil {
		return lib.CustomError(status)
	}
	if len(posts) < 1 {
		return lib.CustomError(http.StatusNotFound)
	}
	var data models.CmsPostTypeData

	data.PostTypeKey = postTypeDB.PostTypeKey
	data.PostTypeCode = postTypeDB.PostTypeCode
	data.PostTypeName = postTypeDB.PostTypeName
	data.PostTypeDesc = postTypeDB.PostTypeDesc
	data.PostTypeGroup = postTypeDB.PostTypeGroup

	var postData []models.CmsPostData
	for _, post := range posts {
		var data models.CmsPostData
	
		data.PostKey = post.PostKey
		data.PostSubtype.PostSubtypeKey = postSubtypeData[post.PostSubtypeKey].PostSubtypeKey
		data.PostSubtype.PostSubtypeCode = postSubtypeData[post.PostSubtypeKey].PostSubtypeCode
		data.PostSubtype.PostSubtypeName = postSubtypeData[post.PostSubtypeKey].PostSubtypeName
		data.PostTitle = post.PostTitle
		data.PostSubTitle = post.PostSubTitle
		data.PostContent = post.PostContent
		data.PostContentAuthor = post.PostContentAuthor
		data.PostContentSources = post.PostContentSources

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, post.PostPublishStart)
		data.PostPublishStart = date.Format(newLayout)
		date, _ = time.Parse(layout, post.PostPublishThru)
		data.PostPublishStart = date.Format(newLayout)

		if post.PostPageAllowed > 0 {
			data.PostPageAllowed = true
		}
		if post.PostCommentAllowed > 0 {
			data.PostCommentAllowed = true
		}
		if post.PostCommentDisplayed > 0 {
			data.PostCommentDisplayed = true
		}
		if post.PostFilesAllowed > 0 {
			data.PostFilesAllowed = true
		}
		if post.PostVideoAllowed > 0 {
			data.PostVideoAllowed = true
		}

		data.PostVideoUrl = post.PostVideoUrl

		if post.PostPinned > 0 {
			data.PostPinned = true
		}

		postData = append(postData, data)
	}

	data.PostList = postData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)
	
}