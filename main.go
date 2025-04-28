package main

import (
	"Sekertaris/config"
	"Sekertaris/controller"
	"Sekertaris/repository"
	"Sekertaris/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("Running on port: %s\n", port)

	// Permohonan Surat
	permohonanSuratRepo := repository.NewPermohonanSuratRepository(db)
	permohonanSuratService := service.NewPermohonanSuratService(permohonanSuratRepo)
	permohonanSuratController := controller.NewPermohonanSuratController(permohonanSuratService)

	// Surat Masuk
	suratMasukRepo := repository.NewSuratMasukRepository(db)
	suratMasukService := service.NewSuratMasukService(suratMasukRepo)
	suratMasukController := controller.NewSuratMasukController(suratMasukService)

	// Surat Keluar
	suratKeluarRepo := repository.NewSuratKeluarRepository(db)
	suratKeluarService := service.NewSuratKeluarService(suratKeluarRepo)
	suratKeluarController := controller.NewSuratKeluarController(suratKeluarService)

	// Setup router
	router := httprouter.New()
	router.HandleOPTIONS = true 

	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	// Permohonan Surat Routes
	router.POST("/api/permohonansurat", permohonanSuratController.AddPermohonanSurat)
	router.GET("/api/permohonansurat", permohonanSuratController.GetPermohonanSurat)
	router.GET("/api/permohonansurat/get/:id", permohonanSuratController.GetPermohonanSuratByID)
	router.PUT("/api/permohonansurat/update/:id", permohonanSuratController.UpdatePermohonanSuratByID)
	router.DELETE("/api/permohonansurat/delete/:id", permohonanSuratController.DeletePermohonanSurat)
	router.PATCH("/api/permohonansurat/patch/:id", permohonanSuratController.UpdateStatus)

	// Surat Masuk Routes
	router.POST("/api/suratmasuk", suratMasukController.AddSuratMasuk)
	router.GET("/api/suratmasuk", suratMasukController.GetSuratMasuk)
	router.GET("/api/suratmasuk/get/:id", suratMasukController.GetSuratById)
	router.PUT("/api/suratmasuk/update/:id", suratMasukController.UpdateSuratMasukByID)
	router.DELETE("/api/suratmasuk/delete/:id", suratMasukController.DeleteSuratMasuk)

	// Surat Keluar Routes
	router.POST("/api/suratkeluar", suratKeluarController.AddSuratKeluar)
	router.GET("/api/suratkeluar", suratKeluarController.GetAllSuratKeluar)
	router.GET("/api/suratkeluar/get/:id", suratKeluarController.GetSuratKeluarById)
	router.PUT("/api/suratkeluar/update/:id", suratKeluarController.UpdateSuratKeluarByID)
	router.DELETE("/api/suratkeluar/delete/:id", suratKeluarController.DeleteSuratKeluar)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5800",
			"http://192.168.1.7:5800",
			"http://192.168.1.85:5800",
		},
		 // Tambah semua kemungkinan origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Content-Disposition",
			"ngrok-skip-browser-warning",
		},
		AllowCredentials: true,
		Debug:            true,
	})

	// Wrap router with CORS
	handler := c.Handler(router)

	// Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("Server running on http://192.168.1.85:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}
	fmt.Println("Server stopped gracefully")
}
