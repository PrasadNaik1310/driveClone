package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	newFileName := c.PostForm("newfilename")

	var fileHeaderName string
	var fileErr error
	var file = new(struct {
		Filename string
	})

	for _, key := range []string{"file", "files", "upload", "document"} {
		fh, err := c.FormFile(key)
		if err == nil {
			fileHeaderName = key
			file.Filename = fh.Filename
			fileErr = nil
			break
		}
		fileErr = err
	}

	if fileHeaderName == "" {
		form, _ := c.MultipartForm()
		receivedKeys := make([]string, 0)
		if form != nil {
			for k := range form.Value {
				receivedKeys = append(receivedKeys, k)
			}
			for k := range form.File {
				receivedKeys = append(receivedKeys, k)
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":                "File not found in request. Send multipart/form-data with a file field named 'file'",
			"received_content_type": c.ContentType(),
			"received_form_keys":    receivedKeys,
			"tried_file_fields":     []string{"file", "files", "upload", "document"},
			"details":              fileErr.Error(),
		})
		return
	}

	fh, err := c.FormFile(fileHeaderName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read uploaded file"})
		return
	}
	folderId := c.PostForm("folder_id")
	// creating a storage for each file because i cant store raw files as there as chances of file duplicate names .
	ext := filepath.Ext(fh.Filename)
	fileID := uuid.New().String()
	storageName := fileID + ext

	os.MkdirAll("uploads", os.ModePerm)
	filePath := filepath.Join("uploads", storageName)
	err = c.SaveUploadedFile(fh, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to save file"})
		return
	}
	log.Printf("Save file successfull at %s", filePath)

	if strings.TrimSpace(newFileName) == "" {
		newFileName = fh.Filename
	}

	log.Printf("Got file :-> %s (field: %s) from client %s ", storageName, fileHeaderName, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"file":        fh.Filename,
		"newfilename": newFileName,
		"storageName": storageName,
		"folder_id":   folderId,
		"file_field":  fileHeaderName,
	})

}
