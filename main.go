package main

import (
	"Sekertaris/config"
	"Sekertaris/controller"
	"Sekertaris/repository"
	"Sekertaris/service"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Koneksi ke database
	db := config.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	fmt.Println("Connected to database")

	// Inisialisasi repository, service, dan controller
	suratKeluarRepository := repository.NewSuratKeluarRepository(db)
	suratKeluarService := service.AddSuratKeluar(suratKeluarRepository)
	suratKeluarController := controller.NewSuratKeluarController(suratKeluarService)

	// Router menggunakan httprouter
	router := httprouter.New()

	// Rute untuk Surat Keluar
	router.POST("/surat-keluar", suratKeluarController.CreateSuratKeluar)
	router.PUT("/surat-keluar/:id", suratKeluarController.UpdateSuratKeluar)

	// Jalankan server
	port := ":8080"
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
