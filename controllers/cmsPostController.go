package controllers

import (
	"api/models"
	"api/config"
	"api/lib"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetCmsPostList(c echo.Context) error {
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
			return lib.CustomError(http.StatusBadRequest)
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
			return lib.CustomError(http.StatusBadRequest)
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
			return lib.CustomError(http.StatusBadRequest)
		}
	} else {
		noLimit = false
	}

	params := make(map[string]string)
	field := c.Param("field")
	if field == "" {
		log.Error("Missing required parameters")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameters","Missing required parameters")
	}
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}
	var postTypeDB models.CmsPostType
	paramType := make(map[string]string)
	if field == "type" {	
		status, err = models.GetCmsPostType(&postTypeDB, "post_type_key", keyStr)
		paramType["post_type_key"] = keyStr
		if err != nil {
			return lib.CustomError(status)
		}
	}else if field == "subtype" {
		paramType["post_subtype_key"] = keyStr
	}else{
		return lib.CustomError(http.StatusBadRequest)
	}
	
	var postSubtypeDB []models.CmsPostSubtype
	status, err = models.GetAllCmsPostSubtype(&postSubtypeDB, limit, offset, paramType, noLimit)
	if err != nil {
		return lib.CustomError(status, "Error get data", "Failed get data")
	}
	if len(postSubtypeDB) < 1 {
		return lib.CustomError(http.StatusNotFound, "Post not found", "Post not found")
	}
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
	status, err = models.GetCmsPostIn(&posts, postSubtypeIDs, "post_subtype_key")
	if err != nil {
		return lib.CustomError(status)
	}
	if len(posts) < 1 {
		return lib.CustomError(http.StatusNotFound)
	}
	
	var responseData models.CmsPostTypeData
	if field == "type"{
		responseData.PostTypeKey = postTypeDB.PostTypeKey
		responseData.PostTypeCode = postTypeDB.PostTypeCode
		responseData.PostTypeName = postTypeDB.PostTypeName
		responseData.PostTypeDesc = postTypeDB.PostTypeDesc
		responseData.PostTypeGroup = postTypeDB.PostTypeGroup
	}
	var postData []models.CmsPostList
	for _, post := range posts {
		var data models.CmsPostList
	
		data.PostKey = post.PostKey
		data.PostSubtype.PostSubtypeKey = postSubtypeData[post.PostSubtypeKey].PostSubtypeKey
		data.PostSubtype.PostSubtypeCode = postSubtypeData[post.PostSubtypeKey].PostSubtypeCode
		if postSubtypeData[post.PostSubtypeKey].PostSubtypeName != nil {
			data.PostSubtype.PostSubtypeName = *postSubtypeData[post.PostSubtypeKey].PostSubtypeName
		}
		data.PostTitle = post.PostTitle
		if post.PostSubTitle != nil {
			data.PostSubTitle = *post.PostSubTitle
		}
		if post.PostContentAuthor != nil {
			data.PostContentAuthor = *post.PostContentAuthor
		}
		
		if post.PostContentSources != nil {
			data.PostContentSources = *post.PostContentSources
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, post.PostPublishStart)
		data.PostPublishStart = date.Format(newLayout)
		date, _ = time.Parse(layout, post.PostPublishThru)
		data.PostPublishThru = date.Format(newLayout)

		if post.PostPinned > 0 {
			data.PostPinned = true
		}

		if post.RecImage1 != nil && *post.RecImage1 != "" {
			data.RecImage1 = config.BaseUrl + "/images/post/" + *post.RecImage1
		}else{
			data.RecImage1 = config.BaseUrl + "/images/post/default.png"
		}
		postData = append(postData, data)
	}

	responseData.PostList = postData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)
}

func GetCmsPostData(c echo.Context) error {
	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var post models.CmsPost
	status, err = models.GetCmsPost(&post, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var postSubtype models.CmsPostSubtype
	status, err = models.GetCmsPostSubtype(&postSubtype, strconv.FormatUint(post.PostSubtypeKey, 10))

	var responseData models.CmsPostData

	responseData.PostKey = post.PostKey
	responseData.PostSubtype.PostSubtypeKey = postSubtype.PostSubtypeKey
	responseData.PostSubtype.PostSubtypeCode = postSubtype.PostSubtypeCode
	if postSubtype.PostSubtypeName != nil {
		responseData.PostSubtype.PostSubtypeName = *postSubtype.PostSubtypeName
	}
	responseData.PostTitle = post.PostTitle
	if post.PostSubTitle != nil {
		responseData.PostSubTitle = *post.PostSubTitle
	}
	if post.PostContentAuthor != nil {
		responseData.PostContentAuthor = *post.PostContentAuthor
	}
	if post.PostContentSources != nil {
		responseData.PostContentSources = *post.PostContentSources
	}
	if post.PostContent != nil {
		responseData.PostContent = *post.PostContent
	}
	

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, post.PostPublishStart)
	responseData.PostPublishStart = date.Format(newLayout)
	date, _ = time.Parse(layout, post.PostPublishThru)
	responseData.PostPublishThru = date.Format(newLayout)

	if post.PostPageAllowed > 0 {
		responseData.PostPageAllowed = true
	}
	if post.PostCommentAllowed > 0 {
		responseData.PostCommentAllowed = true
	}
	if post.PostCommentDisplayed > 0 {
		responseData.PostCommentDisplayed = true
	}
	if post.PostFilesAllowed > 0 {
		responseData.PostFilesAllowed = true
	}
	if post.PostVideoAllowed > 0 {
		responseData.PostVideoAllowed = true
	}
	if post.PostVideoUrl != nil {
		responseData.PostVideoUrl = *post.PostVideoUrl
	}
	if post.PostPinned > 0 {
		responseData.PostPinned = true
	}
	if post.RecImage1 != nil && *post.RecImage1 != "" {
		responseData.RecImage1 = config.BaseUrl + "/images/post/" + *post.RecImage1
	}else{
		responseData.RecImage1 = config.BaseUrl + "/images/post/default.png"
	}	
	if post.RecImage2 != nil && *post.RecImage2 != "" {
		responseData.RecImage1 = config.BaseUrl + "/images/post/" + *post.RecImage1
	}else{
		responseData.RecImage2 = config.BaseUrl + "/images/post/default.png"
	}	

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
