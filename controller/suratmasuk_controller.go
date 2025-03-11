package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SuratMasukController struct {
	service *service.SuratMasukService
}

func NewSuratMasukController(service *service.SuratMasukService) *SuratMasukController {
	return &SuratMasukController{service: service}
}

func AddSuratMasuk(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.AddSuratMasuk(w, r, db)
	}
}

func GetSuratMasuk(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.GetSuratMasuk(w, db)
	}
}

// GetSuratById menangani request untuk mendapatkan surat masuk berdasarkan ID
func (c *SuratMasukController) GetSuratById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ambil ID dari parameter URL
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Panggil service untuk mendapatkan data surat masuk
	surat, err := c.service.GetSuratById(id)
	if err != nil {
		log.Println("Error retrieving surat masuk by ID:", err)
		http.Error(w, `{"error": "Error retrieving surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Jika data tidak ditemukan
	if surat == nil {
		http.Error(w, `{"error": "Surat masuk not found"}`, http.StatusNotFound)
		return
	}

	// Buat response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(surat); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, `{"error": "Error processing request"}`, http.StatusInternalServerError)
	}
}

func (c *SuratMasukController) GetCountSuratMasuk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.service.GetCountSuratMasuk()
	if err != nil {
		log.Println("Error retrieving count surat masuk:", err)
		http.Error(w, `{"error": "Error retrieving count surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Buat response JSON
	response := map[string]int{
		"jumlah surat": count,
	}

	// Set header dan tulis response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, `{"error": "Error processing request"}`, http.StatusInternalServerError)
	}
}

// UpdateSuratMasukByID menangani request untuk memperbarui surat masuk berdasarkan ID
func (c *SuratMasukController) UpdateSuratMasukByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ambil ID dari parameter URL
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Parse form data (termasuk file)
	err = r.ParseMultipartForm(10 << 20) // Batas ukuran file: 10 MB
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
		return
	}

	// Ambil nilai dari form
	nomor := r.FormValue("nomor")
	tanggal := r.FormValue("tanggal")
	perihal := r.FormValue("perihal")
	asal := r.FormValue("asal")
	title := r.FormValue("title")

	// Ambil file dari form
	file, handler, err := r.FormFile("file")
	var filePath string
	if err == nil {
		defer file.Close()

		// Simpan file ke folder static/suratmasuk
		filePath = fmt.Sprintf("static/suratmasuk/%s", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error saving file:", err)
			http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Salin file ke lokasi tujuan
		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Jika tidak ada file yang diunggah, gunakan file yang sudah ada
		filePath = r.FormValue("existing_file")
	}

	// Buat struct SuratMasuk dengan data dari form
	surat := model.SuratMasuk{
		Nomor:   nomor,
		Tanggal: tanggal,
		Perihal: perihal,
		Asal:    asal,
		Title:   title,
		File:    filePath, // Simpan path file
	}

	// Panggil service untuk memperbarui data surat masuk
	err = c.service.UpdateSuratMasukByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		http.Error(w, `{"error": "Error updating surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Surat masuk updated successfully"}`))
}


// DeleteSuratMasuk menangani request untuk menghapus surat masuk berdasarkan nomor dan perihal
func (c *SuratMasukController) DeleteSuratMasuk(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ambil nomor dan perihal dari parameter URL
	nomor := ps.ByName("nomor")
	perihal := ps.ByName("perihal")

	// Panggil service untuk menghapus data surat masuk
	err := c.service.DeleteSuratMasuk(nomor, perihal)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Surat masuk not found:", err)
			http.Error(w, `{"error": "Surat masuk not found"}`, http.StatusNotFound)
			return
		}
		log.Println("Error deleting surat masuk:", err)
		http.Error(w, `{"error": "Error deleting surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Surat masuk deleted successfully"}`))
}
