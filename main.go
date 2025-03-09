package main

import (
	"Sekertaris/config"
	"Sekertaris/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	port := os.Getenv("APP_PORT")

	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	fmt.Print("running on port: ", port)

	router := httprouter.New()
	router.POST("/api/suratmasuk", controller.AddSuratMasuk(db))

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(errServer)
	}
}
