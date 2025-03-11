package main

import (
	"Sekertaris/config"
	"Sekertaris/controller"
	"Sekertaris/repository"
	"Sekertaris/service"
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

	suratMasukRepo := repository.NewSuratMasukRepository(db)
	suratMasukService := service.NewSuratMasukService(suratMasukRepo)
	suratMasukController := controller.NewSuratMasukController(suratMasukService)

	router := httprouter.New()

	//Surat Keluar
	router.POST("/api/suratkeluar", controller.AddSuratKeluar(db))
	router.GET("/api/suratkeluar/:id", controller.GetSuratKeluarByid(db))
	router.PUT("/api/suratkeluar/:nomor", controller.UpdateSuratKeluar(db))

	//Surat Masuk
	router.POST("/api/suratmasuk", controller.AddSuratMasuk(db))
	router.GET("/api/suratmasuk/get", controller.GetSuratMasuk(db))
	router.GET("/api/suratmasuk/get/:id", suratMasukController.GetSuratById)
	router.GET("/api/suratmasuk/count", suratMasukController.GetCountSuratMasuk)
	router.PUT("/api/suratmasuk/update/:id", suratMasukController.UpdateSuratMasukByID)
	router.DELETE("/api/suratmasuk/delete/:nomor/:perihal", suratMasukController.DeleteSuratMasuk)

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(errServer)
	}
}
