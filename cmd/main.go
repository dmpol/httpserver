package main

import (
	"github.com/joho/godotenv"
	"log"
	"myhttpserver/db/connection"
	router_main "myhttpserver/router"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	_, err = connection.NewDbConnectPostgres(connection.Config{
		"postgres",
		"5432",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		"disable",
	})
	if err != nil {
		log.Fatalln(err)
	}

	router := router_main.Setup()
	router.Run("0.0.0.0:8080")
}
