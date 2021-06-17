package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func CustomerList(c echo.Context) error {
	maxCountFile := 15
	// maxFileSize = 1000

	var deleteIndex []string
	var updateIndex []string
	var createIndex []string

	for i := 1; i <= maxCountFile; i++ {
		fileKey := c.FormValue("file_key_" + strconv.FormatUint(uint64(i), 10))
		// fileNotes := c.FormValue("file_notes_" + strconv.FormatUint(uint64(i), 10))
		fileDelete := c.FormValue("file_delete_" + strconv.FormatUint(uint64(i), 10))
		file, err := c.FormFile("file_upload_" + strconv.FormatUint(uint64(i), 10))

		//delete
		if fileDelete == "1" {
			if fileKey != "" {
				deleteIndex = append(deleteIndex, fileKey)
			}
		} else {
			if fileKey != "" { //update
				if file != nil {
					if err != nil {
						log.Println("file_upload_" + strconv.FormatUint(uint64(i), 10) + " : " + err.Error())
						return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_"+strconv.FormatUint(uint64(i), 10))
					}
					updateIndex = append(updateIndex, fileKey)
				} else {
					updateIndex = append(updateIndex, fileKey)
				}
			} else { //create
				if file != nil {
					if err != nil {
						log.Println("file_upload_" + strconv.FormatUint(uint64(i), 10) + " : " + err.Error())
						return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: file_upload_"+strconv.FormatUint(uint64(i), 10))
					}
					createIndex = append(createIndex, fileKey)
				}
			}
		}
	}

	dateLayout := "2006-01-02 15:04:05"
	err := os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/ms_file", 0755)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	} else {
		for i := 1; i <= maxCountFile; i++ {

			fileKey := c.FormValue("file_key_" + strconv.FormatUint(uint64(i), 10))
			fileNotes := c.FormValue("file_notes_" + strconv.FormatUint(uint64(i), 10))
			// fileDelete := c.FormValue("file_delete_" + strconv.FormatUint(uint64(i), 10))
			// file, err := c.FormFile("file_upload_" + strconv.FormatUint(uint64(i), 10))

			_, foundDelete := lib.Find(deleteIndex, strconv.FormatUint(uint64(i), 10))
			if foundDelete {
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
			}
			_, foundUpdate := lib.Find(updateIndex, strconv.FormatUint(uint64(i), 10))
			if foundUpdate {
				//update data
			}
			_, foundCreate := lib.Find(createIndex, strconv.FormatUint(uint64(i), 10))
			if foundCreate {
				//create
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
