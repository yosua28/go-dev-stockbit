package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func GetAdminCmsPostData(c echo.Context) error {
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

	var postTypeDB models.CmsPostType
	status, err = models.GetCmsPostType(&postTypeDB, "post_type_key", strconv.FormatUint(postSubtype.PostTypeKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get post type data")
	}

	dir := strings.ToLower(postTypeDB.PostTypeCode)

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
		responseData.RecImage1 = config.BaseUrl + "/images/post/" + dir + "/" + *post.RecImage1
	} else {
		responseData.RecImage1 = config.BaseUrl + "/images/post/default.png"
	}
	if post.RecImage2 != nil && *post.RecImage2 != "" {
		responseData.RecImage2 = config.BaseUrl + "/images/post/" + dir + "/" + *post.RecImage2
	} else {
		responseData.RecImage2 = config.BaseUrl + "/images/post/default.png"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CreateAdminCmsPost(c echo.Context) error {
	var err error
	params := make(map[string]string)

	postsubtypekey := c.FormValue("post_subtype_key")
	if postsubtypekey == "" {
		log.Error("Missing required parameter: post_subtype_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_subtype_key", "Missing required parameter: post_subtype_key")
	}
	sub, err := strconv.ParseUint(postsubtypekey, 10, 64)
	if err == nil && sub > 0 {
		params["post_subtype_key"] = postsubtypekey
	} else {
		log.Error("Wrong input for parameter: post_subtype_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_subtype_key", "Missing required parameter: post_subtype_key")
	}

	posttitle := c.FormValue("post_title")
	if posttitle == "" {
		log.Error("Missing required parameter: post_title")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_title", "Missing required parameter: post_title")
	}
	params["post_title"] = posttitle

	postsubtitle := c.FormValue("post_sub_title")
	params["post_sub_title"] = postsubtitle

	postcontent := c.FormValue("post_content")
	params["post_content"] = postcontent

	postcontentauthor := c.FormValue("post_content_author")
	params["post_content_author"] = postcontentauthor

	postcontentsources := c.FormValue("post_content_sources")
	params["post_content_sources"] = postcontentsources

	//date
	postpublishstart := c.FormValue("post_publish_start")
	if postpublishstart == "" {
		log.Error("Missing required parameter: post_publish_start")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_publish_start", "Missing required parameter: post_publish_start")
	}
	params["post_publish_start"] = postpublishstart
	//date
	postpublishthru := c.FormValue("post_publish_thru")
	if postpublishthru == "" {
		log.Error("Missing required parameter: post_publish_thru")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_publish_thru", "Missing required parameter: post_publish_thru")
	}
	params["post_publish_thru"] = postpublishthru

	postpageallowed := c.FormValue("post_page_allowed")
	var postpageallowedBool bool
	if postpageallowed != "" {
		postpageallowedBool, err = strconv.ParseBool(postpageallowed)
		if err != nil {
			log.Error("post_page_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_page_allowed parameter should be true/false", "post_page_allowed parameter should be true/false")
		}
		if postpageallowedBool == true {
			params["post_page_allowed"] = "1"
		} else {
			params["post_page_allowed"] = "0"
		}
	} else {
		log.Error("post_page_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_page_allowed parameter should be true/false", "post_page_allowed parameter should be true/false")
	}

	postcommentallowed := c.FormValue("post_comment_allowed")
	var postcommentallowedBool bool
	if postcommentallowed != "" {
		postcommentallowedBool, err = strconv.ParseBool(postcommentallowed)
		if err != nil {
			log.Error("post_comment_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_comment_allowed parameter should be true/false", "post_comment_allowed parameter should be true/false")
		}
		if postcommentallowedBool == true {
			params["post_comment_allowed"] = "1"
		} else {
			params["post_comment_allowed"] = "0"
		}
	} else {
		log.Error("post_comment_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_comment_allowed parameter should be true/false", "post_comment_allowed parameter should be true/false")
	}

	postcommentdisplayed := c.FormValue("post_comment_displayed")
	var postcommentdisplayedBool bool
	if postcommentdisplayed != "" {
		postcommentdisplayedBool, err = strconv.ParseBool(postcommentdisplayed)
		if err != nil {
			log.Error("post_comment_displayed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_comment_displayed parameter should be true/false", "post_comment_displayed parameter should be true/false")
		}
		if postcommentdisplayedBool == true {
			params["post_comment_displayed"] = "1"
		} else {
			params["post_comment_displayed"] = "0"
		}
	} else {
		log.Error("post_comment_displayed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_comment_displayed parameter should be true/false", "post_comment_displayed parameter should be true/false")
	}

	postfilesallowed := c.FormValue("post_files_allowed")
	var postfilesallowedBool bool
	if postfilesallowed != "" {
		postfilesallowedBool, err = strconv.ParseBool(postfilesallowed)
		if err != nil {
			log.Error("post_files_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_files_allowed parameter should be true/false", "post_files_allowed parameter should be true/false")
		}
		if postfilesallowedBool == true {
			params["post_files_allowed"] = "1"
		} else {
			params["post_files_allowed"] = "0"
		}
	} else {
		log.Error("post_files_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_files_allowed parameter should be true/false", "post_files_allowed parameter should be true/false")
	}

	postvideoallowed := c.FormValue("post_video_allowed")
	var postvideoallowedBool bool
	if postvideoallowed != "" {
		postvideoallowedBool, err = strconv.ParseBool(postvideoallowed)
		if err != nil {
			log.Error("post_video_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_video_allowed parameter should be true/false", "post_video_allowed parameter should be true/false")
		}
		if postvideoallowedBool == true {
			params["post_video_allowed"] = "1"
		} else {
			params["post_video_allowed"] = "0"
		}
	} else {
		log.Error("post_video_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_video_allowed parameter should be true/false", "post_video_allowed parameter should be true/false")
	}

	postvideourl := c.FormValue("post_video_url")
	params["post_video_url"] = postvideourl

	postpinned := c.FormValue("post_pinned")
	var postpinnedBool bool
	if postpinned != "" {
		postpinnedBool, err = strconv.ParseBool(postpinned)
		if err != nil {
			log.Error("post_pinned parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_pinned parameter should be true/false", "post_pinned parameter should be true/false")
		}
		if postpinnedBool == true {
			params["post_pinned"] = "1"
		} else {
			params["post_pinned"] = "0"
		}
	} else {
		log.Error("post_pinned parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_pinned parameter should be true/false", "post_pinned parameter should be true/false")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	var postSubType models.CmsPostSubtype
	status, err := models.GetCmsPostSubtype(&postSubType, postsubtypekey)
	if err != nil {
		return lib.CustomError(status)
	}

	var postType models.CmsPostType
	strTypeKey := strconv.FormatUint(postSubType.PostTypeKey, 10)
	status, err = models.GetCmsPostType(&postType, "post_type_key", strTypeKey)
	if err != nil {
		return lib.CustomError(status)
	}

	pathType := strings.ToLower(postType.PostTypeCode)
	randInt := rand.Intn(9999)
	randName := strconv.FormatUint(uint64(randInt), 10)

	err = os.MkdirAll(config.BasePath+"/images/post/"+pathType, 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("rec_image1")

		items := []string{"1", "3", "4", "5"}
		strType := strconv.FormatUint(postSubType.PostTypeKey, 10)

		_, found := lib.Find(items, strType)
		if found {
			if file == nil {
				log.Error("Wrong input for parameter rec_image1")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter rec_image1", "Wrong input for parameter rec_image1")
			}
		}

		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/post/"+pathType+"/"+randName+"_"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = randName + "_" + filename + extension
		}

		file, err = c.FormFile("rec_image2")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/post/"+pathType+"/"+randName+"_"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image2"] = randName + "_" + filename + extension
		}
	}

	status, err = models.CreatePost(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminCmsPost(c echo.Context) error {
	var err error
	params := make(map[string]string)

	postkey := c.FormValue("post_key")
	if postkey == "" {
		log.Error("Missing required parameter: post_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_key", "Missing required parameter: post_key")
	}
	strPostKey, err := strconv.ParseUint(postkey, 10, 64)
	if err == nil && strPostKey > 0 {
		params["post_key"] = postkey
	} else {
		log.Error("Wrong input for parameter: post_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_key", "Missing required parameter: post_key")
	}

	var post models.CmsPost
	status, err := models.GetCmsPost(&post, postkey)
	if err != nil {
		return lib.CustomError(status)
	}

	postsubtypekey := c.FormValue("post_subtype_key")
	if postsubtypekey == "" {
		log.Error("Missing required parameter: post_subtype_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_subtype_key", "Missing required parameter: post_subtype_key")
	}
	sub, err := strconv.ParseUint(postsubtypekey, 10, 64)
	if err == nil && sub > 0 {
		params["post_subtype_key"] = postsubtypekey
	} else {
		log.Error("Wrong input for parameter: post_subtype_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_subtype_key", "Missing required parameter: post_subtype_key")
	}

	posttitle := c.FormValue("post_title")
	if posttitle == "" {
		log.Error("Missing required parameter: post_title")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_title", "Missing required parameter: post_title")
	}
	params["post_title"] = posttitle

	postsubtitle := c.FormValue("post_sub_title")
	params["post_sub_title"] = postsubtitle

	postcontent := c.FormValue("post_content")
	params["post_content"] = postcontent

	postcontentauthor := c.FormValue("post_content_author")
	params["post_content_author"] = postcontentauthor

	postcontentsources := c.FormValue("post_content_sources")
	params["post_content_sources"] = postcontentsources

	//date
	postpublishstart := c.FormValue("post_publish_start")
	if postpublishstart == "" {
		log.Error("Missing required parameter: post_publish_start")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_publish_start", "Missing required parameter: post_publish_start")
	}
	params["post_publish_start"] = postpublishstart
	//date
	postpublishthru := c.FormValue("post_publish_thru")
	if postpublishthru == "" {
		log.Error("Missing required parameter: post_publish_thru")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_publish_thru", "Missing required parameter: post_publish_thru")
	}
	params["post_publish_thru"] = postpublishthru

	postpageallowed := c.FormValue("post_page_allowed")
	var postpageallowedBool bool
	if postpageallowed != "" {
		postpageallowedBool, err = strconv.ParseBool(postpageallowed)
		if err != nil {
			log.Error("post_page_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_page_allowed parameter should be true/false", "post_page_allowed parameter should be true/false")
		}
		if postpageallowedBool == true {
			params["post_page_allowed"] = "1"
		} else {
			params["post_page_allowed"] = "0"
		}
	} else {
		log.Error("post_page_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_page_allowed parameter should be true/false", "post_page_allowed parameter should be true/false")
	}

	postcommentallowed := c.FormValue("post_comment_allowed")
	var postcommentallowedBool bool
	if postcommentallowed != "" {
		postcommentallowedBool, err = strconv.ParseBool(postcommentallowed)
		if err != nil {
			log.Error("post_comment_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_comment_allowed parameter should be true/false", "post_comment_allowed parameter should be true/false")
		}
		if postcommentallowedBool == true {
			params["post_comment_allowed"] = "1"
		} else {
			params["post_comment_allowed"] = "0"
		}
	} else {
		log.Error("post_comment_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_comment_allowed parameter should be true/false", "post_comment_allowed parameter should be true/false")
	}

	postcommentdisplayed := c.FormValue("post_comment_displayed")
	var postcommentdisplayedBool bool
	if postcommentdisplayed != "" {
		postcommentdisplayedBool, err = strconv.ParseBool(postcommentdisplayed)
		if err != nil {
			log.Error("post_comment_displayed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_comment_displayed parameter should be true/false", "post_comment_displayed parameter should be true/false")
		}
		if postcommentdisplayedBool == true {
			params["post_comment_displayed"] = "1"
		} else {
			params["post_comment_displayed"] = "0"
		}
	} else {
		log.Error("post_comment_displayed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_comment_displayed parameter should be true/false", "post_comment_displayed parameter should be true/false")
	}

	postfilesallowed := c.FormValue("post_files_allowed")
	var postfilesallowedBool bool
	if postfilesallowed != "" {
		postfilesallowedBool, err = strconv.ParseBool(postfilesallowed)
		if err != nil {
			log.Error("post_files_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_files_allowed parameter should be true/false", "post_files_allowed parameter should be true/false")
		}
		if postfilesallowedBool == true {
			params["post_files_allowed"] = "1"
		} else {
			params["post_files_allowed"] = "0"
		}
	} else {
		log.Error("post_files_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_files_allowed parameter should be true/false", "post_files_allowed parameter should be true/false")
	}

	postvideoallowed := c.FormValue("post_video_allowed")
	var postvideoallowedBool bool
	if postvideoallowed != "" {
		postvideoallowedBool, err = strconv.ParseBool(postvideoallowed)
		if err != nil {
			log.Error("post_video_allowed parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_video_allowed parameter should be true/false", "post_video_allowed parameter should be true/false")
		}
		if postvideoallowedBool == true {
			params["post_video_allowed"] = "1"
		} else {
			params["post_video_allowed"] = "0"
		}
	} else {
		log.Error("post_video_allowed parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_video_allowed parameter should be true/false", "post_video_allowed parameter should be true/false")
	}

	postvideourl := c.FormValue("post_video_url")
	params["post_video_url"] = postvideourl

	postpinned := c.FormValue("post_pinned")
	var postpinnedBool bool
	if postpinned != "" {
		postpinnedBool, err = strconv.ParseBool(postpinned)
		if err != nil {
			log.Error("post_pinned parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "post_pinned parameter should be true/false", "post_pinned parameter should be true/false")
		}
		if postpinnedBool == true {
			params["post_pinned"] = "1"
		} else {
			params["post_pinned"] = "0"
		}
	} else {
		log.Error("post_pinned parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "post_pinned parameter should be true/false", "post_pinned parameter should be true/false")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	var postSubType models.CmsPostSubtype
	status, err = models.GetCmsPostSubtype(&postSubType, postsubtypekey)
	if err != nil {
		return lib.CustomError(status)
	}

	var postType models.CmsPostType
	strTypeKey := strconv.FormatUint(postSubType.PostTypeKey, 10)
	status, err = models.GetCmsPostType(&postType, "post_type_key", strTypeKey)
	if err != nil {
		return lib.CustomError(status)
	}

	pathType := strings.ToLower(postType.PostTypeCode)
	randInt := rand.Intn(9999)
	randName := strconv.FormatUint(uint64(randInt), 10)

	err = os.MkdirAll(config.BasePath+"/images/post/"+pathType, 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("rec_image1")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/post/"+pathType+"/"+randName+"_"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = randName + "_" + filename + extension
		}

		file, err = c.FormFile("rec_image2")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/post/"+pathType+"/"+randName+"_"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image2"] = randName + "_" + filename + extension
		}
	}

	status, err = models.UpdateCmsPost(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func DeleteAdminCmsPost(c echo.Context) error {
	var err error
	params := make(map[string]string)

	postkey := c.FormValue("post_key")
	if postkey == "" {
		log.Error("Missing required parameter: post_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_key", "Missing required parameter: post_key")
	}
	strPostKey, err := strconv.ParseUint(postkey, 10, 64)
	if err == nil && strPostKey > 0 {
		params["post_key"] = postkey
	} else {
		log.Error("Wrong input for parameter: post_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_key", "Missing required parameter: post_key")
	}

	var post models.CmsPost
	status, err := models.GetCmsPost(&post, postkey)
	if err != nil {
		return lib.CustomError(status)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateCmsPost(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}
