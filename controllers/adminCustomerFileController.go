package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func CustomerUpdateFile(c echo.Context) error {
	maxCountFile := 15
	// maxFileSize = 1000

	customer_key := c.FormValue("customer_key")
	if customer_key == "" {
		log.Println("customer_key required")
		return lib.CustomError(http.StatusBadRequest, "customer_key required", "customer_key required")
	}

	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByCustomerKey(&userData, customer_key)
	if err != nil {
		log.Println("customer not found")
		return lib.CustomError(http.StatusBadRequest, "customer not found", "customer not found")
	}

	var deleteIndex []string
	var updateIndex []string
	var createIndex []string

	for i := 1; i <= maxCountFile; i++ {
		fileKey := c.FormValue("file_key_" + strconv.FormatUint(uint64(i), 10))
		fileNotes := c.FormValue("file_notes_" + strconv.FormatUint(uint64(i), 10))
		fileDelete := c.FormValue("file_delete_" + strconv.FormatUint(uint64(i), 10))
		file, err := c.FormFile("file_upload_" + strconv.FormatUint(uint64(i), 10))

		if fileKey != "" || fileNotes != "" || fileDelete != "" || file != nil {
			//delete
			if fileDelete == "1" {
				if fileKey != "" {
					deleteIndex = append(deleteIndex, strconv.FormatUint(uint64(i), 10))
				}
			} else {
				if fileNotes == "" {
					log.Println("file_notes_" + strconv.FormatUint(uint64(i), 10) + " required")
					return lib.CustomError(http.StatusBadRequest, "file_notes_"+strconv.FormatUint(uint64(i), 10)+" required", "Missing required parameter: file_notes_"+strconv.FormatUint(uint64(i), 10)+" required")
				}
				if fileKey != "" { //update
					if file != nil {
						if err != nil {
							log.Println("file_upload_" + strconv.FormatUint(uint64(i), 10) + " : " + err.Error())
							return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_"+strconv.FormatUint(uint64(i), 10))
						}
						updateIndex = append(updateIndex, strconv.FormatUint(uint64(i), 10))
					} else {
						updateIndex = append(updateIndex, strconv.FormatUint(uint64(i), 10))
					}
				} else { //create
					if file != nil {
						if err != nil {
							log.Println("file_upload_" + strconv.FormatUint(uint64(i), 10) + " required")
							return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_"+strconv.FormatUint(uint64(i), 10))
						}
						createIndex = append(createIndex, strconv.FormatUint(uint64(i), 10))
					}
				}
			}
		}

	}

	dateLayout := "2006-01-02 15:04:05"
	err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(userData.UserLoginKey, 10)+"/ms_file", 0755)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	} else {
		for i := 1; i <= maxCountFile; i++ {
			fileKey := c.FormValue("file_key_" + strconv.FormatUint(uint64(i), 10))
			fileNotes := c.FormValue("file_notes_" + strconv.FormatUint(uint64(i), 10))
			fileDelete := c.FormValue("file_delete_" + strconv.FormatUint(uint64(i), 10))
			file, err := c.FormFile("file_upload_" + strconv.FormatUint(uint64(i), 10))

			if fileKey != "" || fileNotes != "" || fileDelete != "" || file != nil {
				_, foundDelete := lib.Find(deleteIndex, strconv.FormatUint(uint64(i), 10))
				_, foundUpdate := lib.Find(updateIndex, strconv.FormatUint(uint64(i), 10))
				_, foundCreate := lib.Find(createIndex, strconv.FormatUint(uint64(i), 10))

				if foundDelete {
					log.Println("del")
					//update ms_file -> delete
					deleteFile := make(map[string]string)
					deleteFile["file_key"] = fileKey
					deleteFile["file_notes"] = fileNotes
					deleteFile["rec_status"] = "0"
					deleteFile["rec_modified_date"] = time.Now().Format(dateLayout)
					deleteFile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

					_, err = models.UpdateMsFile(deleteFile)
					if err != nil {
						log.Error("Error update personal data delete")
						return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
					}
				} else if foundUpdate {
					log.Println("upd")
					//update ms_file -> update
					updateFile := make(map[string]string)
					updateFile["file_key"] = fileKey
					updateFile["file_notes"] = fileNotes
					updateFile["rec_status"] = "1"
					updateFile["rec_modified_date"] = time.Now().Format(dateLayout)
					updateFile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					if file != nil {
						if err != nil {
							return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_file_upload_"+strconv.FormatUint(uint64(i), 10))
						}
						// Get file extension
						extension := filepath.Ext(file.Filename)
						// Generate filename
						filename := lib.RandStringBytesMaskImprSrc(50)
						updateFile["file_name"] = filename + extension
						updateFile["file_ext"] = extension
						updateFile["file_path"] = config.BasePath + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/ms_file/" + filename + extension

						err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(userData.UserLoginKey, 10)+"/ms_file/"+filename+extension)
						if err != nil {
							log.Println(err)
							return lib.CustomError(http.StatusInternalServerError)
						}
					}

					_, err = models.UpdateMsFile(updateFile)
					if err != nil {
						log.Error("Error update ms_file")
						return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data ms_file")
					}
				} else if foundCreate {
					log.Println("cre")
					//create
					//update ms_file -> update
					createFile := make(map[string]string)
					createFile["ref_fk_key"] = customer_key
					createFile["ref_fk_domain"] = "ms_customer"
					createFile["file_notes"] = fileNotes
					createFile["rec_status"] = "1"
					createFile["rec_created_date"] = time.Now().Format(dateLayout)
					createFile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					if file != nil {
						if err != nil {
							return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_file_upload_"+strconv.FormatUint(uint64(i), 10))
						}
						// Get file extension
						extension := filepath.Ext(file.Filename)
						// Generate filename
						filename := lib.RandStringBytesMaskImprSrc(50)
						createFile["file_name"] = filename + extension
						createFile["file_ext"] = extension
						createFile["blob_mode"] = "0"
						createFile["file_path"] = config.BasePath + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/ms_file/" + filename + extension

						err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(userData.UserLoginKey, 10)+"/ms_file/"+filename+extension)
						if err != nil {
							log.Println(err)
							return lib.CustomError(http.StatusInternalServerError)
						}
					}

					_, err = models.CreateMsFile(createFile)
					if err != nil {
						log.Error("Error create ms_file")
						return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed create data ms_file")
					}
				}
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func AdminGetDetailCustomerDocument(c echo.Context) error {
	customerKeyStr := c.Param("customer_key")
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err != nil || customerKey == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByCustomerKey(&userData, customerKeyStr)
	if err != nil {
		log.Println("customer not found")
		return lib.CustomError(http.StatusBadRequest, "customer not found", "customer not found")
	}

	//customer
	params := make(map[string]string)
	params["c.customer_key"] = customerKeyStr
	params["c.investor_type"] = "263"

	var customer models.CustomerIndividuStatusSuspend
	_, err = models.AdminGetDetailCustomerStatusSuspend(&customer, params)
	if err != nil {
		log.Error("Error get data ms_customer")
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	//document
	paramsDoc := make(map[string]string)
	paramsDoc["ref_fk_key"] = customerKeyStr
	paramsDoc["ref_fk_domain"] = "ms_customer"
	paramsDoc["rec_status"] = "1"
	paramsDoc["orderBy"] = "file_key"
	paramsDoc["orderType"] = "ASC"

	var document []models.MsFileDetail
	var responseDocument []models.MsFileDetail
	_, err = models.GetALlDetailMsFile(&document, paramsDoc)
	if err == nil {
		for _, doc := range document {
			var data models.MsFileDetail
			data.FileKey = doc.FileKey
			data.RefFkKey = doc.RefFkKey
			data.FileName = doc.FileName
			data.Path = config.BaseUrl + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/ms_file/" + doc.FileName
			data.FileExt = doc.FileExt
			data.FileNotes = doc.FileNotes
			data.RecCreatedDate = doc.RecCreatedDate
			responseDocument = append(responseDocument, data)
		}
	}

	var responseData models.CustomerDocumentDetail

	responseData.Customer = customer
	responseData.Document = responseDocument

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)

}
