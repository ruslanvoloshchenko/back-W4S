package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"w4s/models"
)

//Find all users on the database
func FindUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User

	if err := db.Where("deleted = ? AND actived = ?", "0", true).Preload("Profile").Preload("Tables").Find(&users).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": users,
	})
	return
}

//Find a user by his(her) nickname/Encontrando um usuario pelo seu nick(url)
func FindUserByNick(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.Profile
	if err := db.Where("nickname = ? AND deleted = ?", c.Query("e"), "0", "1").Preload("User").Preload("Tables").First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Registro não encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": user,
	})
}
