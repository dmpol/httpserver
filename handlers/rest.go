package handlers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"myhttpserver/db/connection"
	"net/http"
)

func GetPersons(c *gin.Context) {

	rows, err := connection.GetConnect().Query("SELECT user_id, user_name, password_hash, email FROM users")
	if err != nil {
		log.Printf("Ошибка запроса: %s", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Ошибка закрытия: %s", err)
		}
	}()

	var persons []PersonWithId

	for rows.Next() {
		var person PersonWithId
		err := rows.Scan(&person.Id, &person.Username, &person.Password, &person.Email)
		if err != nil {
			log.Printf("Ошибка сканирования строки DB: %s", err)
		}
		persons = append(persons, person)
	}
	c.JSON(http.StatusOK, persons)
}

func GetPerson(c *gin.Context) {
	id := c.Param("id")

	var person PersonWithId
	err := connection.GetConnect().QueryRow("SELECT user_id, user_name, password_hash, email FROM users WHERE user_id = $1", id).
		Scan(&person.Id, &person.Username, &person.Password, &person.Email)
	if err != nil {
		log.Printf("Ошибка запроса к DB: %s", err)
	}

	c.JSON(http.StatusOK, person)
}

func CreatePerson(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error Create": err.Error()})
		return
	}

	personHashPassword, err := hashPassword(person.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %s", err)
	}

	_, err = connection.GetConnect().Exec("INSERT INTO users (user_name, password_hash, email) VALUES ($1, $2, $3)", person.Username, personHashPassword, person.Email)
	if err != nil {
		log.Printf("Ошибка при вставке данных: %s", err)
	}

	c.JSON(http.StatusCreated, person)
}

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	personHashPassword, err := hashPassword(person.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %s", err)
	}

	_, err = connection.GetConnect().Exec("UPDATE users SET user_name = $1, password_hash = $2, email = $3 WHERE user_id = $4", person.Username, personHashPassword, person.Email, id)
	if err != nil {
		log.Printf("Ошибка обновления данных: %s", err)
	}

	c.JSON(http.StatusOK, person)
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	_, err := connection.GetConnect().Exec("DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		log.Printf("Ошибка удаления данных: %s", err)
	}

	c.String(http.StatusOK, "Person deleted")
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

//"$2a$10$8nNnxb/.jmdwzF3W0Jkskuzxw08P7lROJqUzm1VIogrlVQu2BDG.2"
