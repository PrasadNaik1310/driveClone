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

func CreateFolder(c *gin.Context){
	var folder struct{
		Foldername string `json:"foldername" binding:"required"`
		//including parent as json for testing on postman
		Parentid *string `json:"parentid"`
		
	}
	//var user models.User
	var parent models.Folder
	var newFolder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error , request not formatted well , check the following error <:::":err.Error()})
		log.Println("Bad request recieved !! ")
		return
	}
	userIDRaw,exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"Error":"Unauthenticated user "})
		return
	}
	userId,ok := userIDRaw.(uint)
	if !ok{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
    return
	}


	if  len(folder.Foldername) < 2 {
		c.JSON(http.StatusNotAcceptable,gin.H{"Error":"Folder name too short"})
		log.Println("Folder name too short")
		c.Abort()
		return
	}
	//root allocation logif
	
	if(folder.Parentid !=nil){
		if err := db.DB.Where("id=? AND ownerId=?",*folder.Parentid,userId).First(&parent).Error; err!= nil{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Parent user not found or user not permitted"})
			return
		}

	}

	newFolder.FolderId = uuid.New().String() 
	newFolder.CreatedAt = time.Now()
	newFolder.FolderName = folder.Foldername
	newFolder.OwnerId = userId
	if folder.Parentid != nil{
	newFolder.ParentId = *folder.Parentid
	}
	if err := db.DB.Create(&newFolder); err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"cant create new folder"})
		log.Print(err)
		return
	}




	c.JSON(http.StatusOK,gin.H{
	"newfFolderid":newFolder.FolderId,
		"name":newFolder.FolderName,
		"CreatedAt":newFolder.CreatedAt,


})
}
func FolderContents(c *gin.Context){
	id:= c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest,gin.H{"Error":"Id not found"})
		return
	}
	
}
func DeleteFolder(c *gin.Context){

}