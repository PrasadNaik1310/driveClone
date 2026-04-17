package handlers

import (
	"log"
	"net/http"
	"time"

	//	"os"

	"github.com/PrasadNaik1310/driveClone/db"
	"github.com/PrasadNaik1310/driveClone/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(c *gin.Context) {
	var folder struct {
		Foldername string `json:"foldername" binding:"required"`
		//including parent as json for testing on postman
		Parentid *string `json:"parentid"`
	}
	//var user models.User
	var parent models.Folder
	var newFolder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error , request not formatted well , check the following error <:::": err.Error()})
		log.Println("Bad request recieved !! ")
		return
	}
	UserIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthenticated user "})
		return
	}
	UserId, ok := UserIdRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	if len(folder.Foldername) < 2 {
		c.JSON(http.StatusNotAcceptable, gin.H{"Error": "Folder name too short"})
		log.Println("Folder name too short")
		c.Abort()
		return
	}
	//root allocation logif

	if folder.Parentid != nil {
		if err := db.DB.Where("id=? AND ownerId=?", *folder.Parentid, UserId).First(&parent).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Parent user not found or user not permitted"})
			return
		}

	}

	newFolder.FolderId = uuid.New().String()
	newFolder.CreatedAt = time.Now()
	newFolder.FolderName = folder.Foldername
	newFolder.OwnerId = UserId
	if folder.Parentid != nil {
		newFolder.ParentId = *folder.Parentid
	}
	if err := db.DB.Create(&newFolder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create new folder"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"newfFolderid": newFolder.FolderId,
		"name":         newFolder.FolderName,
		"CreatedAt":    newFolder.CreatedAt,
	})
}

func FolderContents(c *gin.Context) {
	/*var folder struct {
		Foldername string `json:"foldername" binding:"required"`
		//including parent as json for testing on postman
		Parentid *string `json:"parentid"`
	}*/
	//var user models.User
	//	var parent models.Folder
	var FetchFolder []models.Folder
	var CurrentFolder models.Folder
	var files []models.File
	childid := c.Param("id")
	if childid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Id not found"})
		return
	}

	UserIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthenticated user "})
		return
	}
	UserId, ok := UserIdRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	if err := db.DB.Where("id = ? AND ownerId = ? ", childid, UserId).First(&CurrentFolder).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"Error": err})
		log.Print("Status Forbidden or folder not fonud")
		log.Print(err)
		return
	}
	if err := db.DB.Where("parent_id =? AND ownerId=? ", childid, UserId).Find(&FetchFolder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Folder not found"})
		log.Println(err)
		return
	}
	if err := db.DB.Where("folder_id=? AND ownerId=?", childid, UserId).Find(&files).Error; err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "Fiesl not fonud"})
		return
	}

	//retrun the final fuckinggg reposnse

	c.JSON(http.StatusOK, gin.H{
		"folders": FetchFolder,
		"files":   files,
	})

}

func DeleteFolder(c *gin.Context) {
window.alert("Under construction , will be shipped soooooon!!")
return
}
