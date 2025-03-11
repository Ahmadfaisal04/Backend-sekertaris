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

	//Surat Masuk
	suratMasukRepo := repository.NewSuratMasukRepository(db)
	suratMasukService := service.NewSuratMasukService(suratMasukRepo)
	suratMasukController := controller.NewSuratMasukController(suratMasukService)

	//Surat Keluar
	suratKeluarRepo := repository.NewSuratKeluarRepository(db)
	suratKeluarService := service.NewSuratKeluarService(suratKeluarRepo)
	suratKeluarController := controller.NewSuratKeluarController(suratKeluarService)

	router := httprouter.New()

	// Surat Keluar Routes
	router.POST("/api/suratkeluar", controller.AddSuratKeluar(db))
	router.GET("/api/suratkeluar", suratKeluarController.GetAllSuratKeluar) 
	router.GET("/api/suratkeluar/count", suratKeluarController.GetCountSuratKeluar)
	router.GET("/api/suratkeluar/get/:id", suratKeluarController.GetSuratKeluarById)
	router.PUT("/api/suratkeluar/:id", suratKeluarController.UpdateSuratKeluarByID)

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
