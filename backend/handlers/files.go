package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {

	var newFile struct {
		NewFileName string `json:"newfilename"`
	}
	if err := c.ShouldBindJSON(&newFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in the request"})
		return
	}
	folderId := c.PostForm("folder_id")
	// creating a storage for each file because i cant store raw files as there as chances of file duplicate names .
	ext := filepath.Ext(file.Filename)
	fileID := uuid.New().String()
	storageName := fileID + ext

	os.MkdirAll("uploads", os.ModePerm)
	filePath := filepath.Join("uploads", storageName)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to save file"})
		return
	}
	log.Printf("Save file successfull at %s", filePath)

	log.Printf("Got file :-> %s  from client %s ", storageName, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"file":        file.Filename,
		"storageName": storageName,
		"folder_id":   folderId,
	})

}
