package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
	"log"
	"myhttpserver/db/models"
	"net/http"
	"strconv"
)

func GetPersons(c *gin.Context) {
	users, err := models.Users().All(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, users)
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetPerson(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Переданный ID не является числом: %s", err)
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	user, err := models.Users(models.UserWhere.UserID.EQ(userID)).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}
	if user == nil {
		log.Printf("User %s не найден", id)
		c.JSON(http.StatusNotFound, user)
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreatePerson(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error Create": err.Error()})
		return
	}

	user := models.User{
		UserName:     person.Username,
		PasswordHash: person.Password,
		Email:        person.Email,
	}

	if err := user.Insert(c.Request.Context(), boil.GetContextDB(), boil.Infer()); err != nil {
		log.Printf("Ошибка создания User: %s", err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Переданный ID не является числом: %s", err)
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	user, err := models.Users(models.UserWhere.UserID.EQ(userID)).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}
	if user == nil {
		log.Printf("User %s не найден", id)
		c.JSON(http.StatusNotFound, user)
		return
	}

	if person.Username != "" {
		user.UserName = person.Username
	}
	if person.Password != "" {
		personHashPassword, err := hashPassword(person.Password)
		if err != nil {
			log.Printf("Ошибка хеширования пароля: %s", err)
			c.JSON(http.StatusInternalServerError, user)
			return
		}

		user.PasswordHash = personHashPassword
	}
	if person.Email != "" {
		user.Email = person.Email
	}

	if _, err := user.Update(c.Request.Context(), boil.GetContextDB(), boil.Infer()); err != nil {
		log.Printf("Ошибка обновления в базе данных: %s", err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Переданный ID не является числом: %s", err)
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	user, err := models.Users(models.UserWhere.UserID.EQ(userID)).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}
	if user == nil {
		log.Printf("User %s не найден", id)
		c.JSON(http.StatusNotFound, user)
		return
	}

	if _, err := user.Delete(c.Request.Context(), boil.GetContextDB()); err != nil {
		log.Printf("Ошибка удаления из базы данных: %s", err)
		c.JSON(http.StatusInternalServerError, user)
		return
	}

	c.String(http.StatusOK, "User deleted")
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
