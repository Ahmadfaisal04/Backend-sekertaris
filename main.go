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
	"github.com/rs/cors"
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

	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	

	// Surat Keluar Routes
	router.POST("/api/suratkeluar", controller.AddSuratKeluar(db))
	router.GET("/api/suratkeluar", suratKeluarController.GetAllSuratKeluar)
	router.GET("/api/suratkeluar/count", suratKeluarController.GetCountSuratKeluar)
	router.GET("/api/suratkeluar/get/:id", suratKeluarController.GetSuratKeluarById)
	router.PUT("/api/suratkeluar/:id", suratKeluarController.UpdateSuratKeluarByID)
	router.DELETE("/api/suratkeluar/delete/:id", suratKeluarController.DeleteSuratKeluar)

	// Surat Masuk Routes
	router.POST("/api/suratmasuk", suratMasukController.AddSuratMasuk)
	router.GET("/api/suratmasuk", suratMasukController.GetSuratMasuk)
	router.GET("/api/suratmasuk/get/:id", suratMasukController.GetSuratById)
	router.GET("/api/suratmasuk/count", suratMasukController.GetCountSuratMasuk)
	router.PUT("/api/suratmasuk/update/:id", suratMasukController.UpdateSuratMasukByID)
	router.DELETE("/api/suratmasuk/delete/:id", suratMasukController.DeleteSuratMasuk)

	// Enable CORS for all routes
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5800"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Content-Disposition", // Penting untuk file upload
	},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(errServer)
	}
}
