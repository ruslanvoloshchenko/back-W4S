package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"w4s/models"
)

func CreateTable(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.TableInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var table models.Table
	if db.Where("name = ?", input.Name).Find(&table).RecordNotFound() {
		var user models.User
		if err := db.Where("nickname = ? AND actived = ?", c.Query("nickname"), true).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "não encontrado o nickname",
			})
			return
		}
		table.Name = input.Name
		table.Description = input.Description
		table.NumberOfParticipants = 1
		table.Thumbnail = input.Thumbnail
		table.AdventureLink = input.AdventureLink
		table.MaxOfParticipants = input.MaxOfParticipants

		if err := db.Create(&table).Error; err != nil { //Return the error by JSON / Retornando o erro por JSON
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		db.Model(table).Association("User").Append([]*models.User{&user})
		/* if err:=db.Model(table).Association("User").Append([]*models.User{&user}).Error;err!=nil{
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}*/

		c.JSON(http.StatusOK, gin.H{"success": "table created"})
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table name already exists"})
	return
}
func UserJoinTable(c *gin.Context) {
	//Empty parametrs error checking
	if c.Query("nickname") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not inform"})
		return
	}
	if c.Query("table") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table not inform"})
		return
	}
	//=========================

	db := c.MustGet("db").(*gorm.DB)
	var userToADD models.User
	if err := db.Where("deleted = ? AND actived = ?", "0", true).Preload("Profile").Preload("Tables").Find(&userToADD).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	if userToADD.Profile.IDUser == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not have profile"})
		return
	}
	var table models.Table
	if err := db.Where("name = ?", c.Query("table")).Preload("User").Find(&table).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	for i := 0; i < len(table.User); i++ {
		if table.User[i].ID == userToADD.ID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already is in the table"})
			return
		}
	}
	if table.NumberOfParticipants != table.MaxOfParticipants {
		//.Where("name = ? ", c.Query("table"))
		db.Model(&table).Association("User").Append([]*models.User{&userToADD})
		db.Model(&table).Update("numberofparticipants", table.NumberOfParticipants+1)
		c.JSON(http.StatusOK, gin.H{"success": "join in the table"})
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table full"})
	return
}
func FindAllTables(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tables []models.Table

	if err := db.Preload("User").Preload("User.Profile").Find(&tables).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": tables,
	})
	return
}

/*func insertPictures(c *gin.Context, TableId uint) {
	db := c.MustGet("db").(*gorm.DB)
	var pictures models.Picture
	if err := c.BindJSON(pictures); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	pictures.TableID = TableId
	split := strings.Split(pictures.PictureFile, " ")
	for i := 0; i < len(split); i++ {
		if err := db.Create(&pictures).Error; err != nil { //Return the error by JSON / Retornando o erro por JSON
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	return
}
*/
