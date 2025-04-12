package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type SuratKeluarController struct {
	service *service.SuratKeluarService
}

func NewSuratKeluarController(service *service.SuratKeluarService) *SuratKeluarController {
	return &SuratKeluarController{service: service}
}

func AddSuratKeluar(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.AddSuratKeluar(w, r, db)
	}
}

func GetSuratKeluar(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.GetSuratKeluar(w, db)
	}
}

// GetAllSuratKeluar menangani request untuk mendapatkan semua surat keluar
func (c *SuratKeluarController) GetAllSuratKeluar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	suratKeluarList, err := c.service.GetAllSuratKeluar()
	if err != nil {
		log.Println("Error retrieving all surat keluar:", err)
		http.Error(w, `{"error": "Error retrieving all surat keluar"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(suratKeluarList); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, `{"error": "Error processing request"}`, http.StatusInternalServerError)
	}
}

// GetSuratKeluarById menangani request untuk mendapatkan surat keluar berdasarkan ID
func (c *SuratKeluarController) GetSuratKeluarById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ambil ID dari parameter URL
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Panggil service untuk mengambil data surat keluar
	surat, err := c.service.GetSuratKeluarById(id)
	if err != nil {
		// Jika data tidak ditemukan
		if strings.Contains(err.Error(), "tidak ditemukan") {
			log.Printf("Surat keluar with ID %d not found: %v", id, err)
			http.Error(w, `{"error": "Surat keluar not found"}`, http.StatusNotFound)
			return
		}

		// Jika terjadi error lain
		log.Printf("Error retrieving surat keluar by ID %d: %v", id, err)
		http.Error(w, `{"error": "Failed to retrieve surat keluar"}`, http.StatusInternalServerError)
		return
	}

	// Jika berhasil, kirim response JSON dalam bentuk array
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(surat); err != nil {
		log.Printf("Error encoding JSON response for surat keluar ID %d: %v", id, err)
		http.Error(w, `{"error": "Failed to process response"}`, http.StatusInternalServerError)
	}
}

func (c *SuratKeluarController) GetCountSuratKeluar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.service.GetCountSuratKeluar()
	if err != nil {
		log.Println("Error retrieving count surat keluar:", err)
		http.Error(w, `{"error": "Error retrieving count surat keluar"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"jumlah surat": count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, `{"error": "Error processing request"}`, http.StatusInternalServerError)
	}
}

// UpdateSuratKeluarByID menangani request untuk memperbarui surat keluar berdasarkan ID
// func (c *SuratKeluarController) UpdateSuratKeluarByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	idStr := ps.ByName("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		log.Println("Invalid ID:", err)
// 		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
// 		return
// 	}

// 	err = r.ParseMultipartForm(10 << 20) // Batas ukuran file: 10 MB
// 	if err != nil {
// 		log.Println("Error parsing form data:", err)
// 		http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
// 		return
// 	}

// 	nomor := r.FormValue("nomor")
// 	tanggal := r.FormValue("tanggal")
// 	perihal := r.FormValue("perihal")
// 	ditujukan := r.FormValue("ditujukan")
// 	title := r.FormValue("title")

// 	file, handler, err := r.FormFile("file")
// 	var filePath string
// 	if err == nil {
// 		defer file.Close()

// 		filePath = fmt.Sprintf("static/suratkeluar/%s", handler.Filename)
// 		dst, err := os.Create(filePath)
// 		if err != nil {
// 			log.Println("Error saving file:", err)
// 			http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
// 			return
// 		}
// 		defer dst.Close()

// 		_, err = io.Copy(dst, file)
// 		if err != nil {
// 			log.Println("Error copying file:", err)
// 			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		filePath = r.FormValue("existing_file")
// 	}

// 	surat := model.SuratKeluar{
// 		Nomor:     nomor,
// 		Tanggal:   tanggal,
// 		Perihal:   perihal,
// 		Ditujukan: ditujukan,
// 		Title:     title,
// 		File:      filePath,
// 	}

// 	err = c.service.UpdateSuratKeluarByID(id, surat)
// 	if err != nil {
// 		log.Println("Error updating surat keluar:", err)
// 		http.Error(w, `{"error": "Error updating surat keluar"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "Surat keluar updated successfully"}`))
// }

func (c *SuratKeluarController) UpdateSuratKeluarByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Parse form data
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
		return
	}

	// Ambil data dari form
	surat := model.SuratKeluar{
		Nomor:     r.FormValue("nomor"),
		Tanggal:   r.FormValue("tanggal"),
		Perihal:   r.FormValue("perihal"),
		Ditujukan: r.FormValue("ditujukan"),
		Title:     r.FormValue("title"),
	}

	// Handle file upload
	file, handler, err := r.FormFile("file")
	if err == nil {
		defer file.Close()

		// Buat direktori jika belum ada
		staticPath := "./static/suratkeluar/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			log.Println("Error creating directory:", err)
			http.Error(w, `{"error": "Error creating directory"}`, http.StatusInternalServerError)
			return
		}

		// Buat file baru
		filePath := staticPath + handler.Filename
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, `{"error": "Error creating file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Salin file
		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
			return
		}

		// Set path file baru
		surat.File = filePath
		surat.Title = handler.Filename
	} else {
		// Jika tidak ada file baru diupload, gunakan file yang sudah ada
		existingFile := r.FormValue("existing_file")
		existingTitle := r.FormValue("existing_title")

		if existingFile != "" {
			surat.File = existingFile
			surat.Title = existingTitle
		}
	}

	// Update data
	err = c.service.UpdateSuratKeluarByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		http.Error(w, `{"error": "Error updating surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Surat masuk updated successfully"})
}
