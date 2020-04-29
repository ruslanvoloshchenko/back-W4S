package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/smtp"
	"os"
	"w4s/authc"
	"w4s/models"
)

func ConfirmationEmail(userEmail string, c *gin.Context) error {
	userSingUpToken, err := authc.GenerateJWT(userEmail, 14400)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return err
	}
	if err := SendEmail(userEmail, userSingUpToken); err != nil {
		//Declaring a new user to be populated/
		//Declarando um novo usuario que será populado
		var user models.User
		//Declaring and inicializing a userAccountCreatedToken
		//Declarando e inicializando uma variavel userAccount
		userAccountCreatedToken := models.AccountCreatedToken{
			Token: userSingUpToken,
		}
		db := c.MustGet("db").(*gorm.DB) //Establish conection with database/Estabelecendo conexão com o banco de dados
		if err := db.Create(userAccountCreatedToken).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return err
		}
		if err := db.Where("email = ?", userEmail).Find(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return err
		}
		db.Delete(&user)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return err
	}
	return nil
}
func SendEmail(userEmail string, userSingupToken string) error {
	userURL := "http://localhost:8080/user/confirm?e=" + userEmail + "&t=" + userSingupToken
	//Thanks https://blog.mailtrap.io/golang-send-email/#Sending_emails_with_smtpSendMail for the code
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "findatablew4s@gmail.com", os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com")
	// Here we do it all: connect to our server, set up a message and send it
	to := []string{userEmail}
	msg := []byte("To: " + userEmail + "\r\n" +
		"Subject: NÂO RESPONDA ESTE EMAIL - Confirmação de Conta \r\n" +
		"\r\n" +
		"Bem vindo ! obrigado por criar uma conta no Find a Table - RPG !\r\n" +
		"\r\n" +
		"Clique aqui para confirmar sua conta = " + userURL +
		"\r\n" +
		"Caso não consiga, é só copiar o link e colar no navegador ! " +
		"\r\n" +
		"Caso não tenha criado, por favor entrar em contato com findatablew4s@gmail.com, com o Assunto : Conta Criada Indevidamente ")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "findatablew4s@gmail.com", to, msg)
	if err != nil {
		return err
	}
	return nil
}
