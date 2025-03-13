package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type SuratMasukController struct {
	service *service.SuratMasukService
}

func NewSuratMasukController(service *service.SuratMasukService) *SuratMasukController {
	return &SuratMasukController{service: service}
}

func (c *SuratMasukController) AddSuratMasuk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		var surat model.SuratMasuk
		surat.Nomor = r.FormValue("nomor")
		surat.Tanggal = r.FormValue("tanggal")
		surat.Perihal = r.FormValue("perihal")
		surat.Asal = r.FormValue("asal")

		// Validasi bahwa tanggal tidak boleh kosong
		if surat.Tanggal == "" {
			http.Error(w, `{"error": "Tanggal is required"}`, http.StatusBadRequest)
			return
		}

		parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, `{"error": "Unable to get file from form"}`, http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Menentukan path penyimpanan file di direktori static
		staticPath := "./static/suratmasuk/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			http.Error(w, `{"error": "Unable to create static directory"}`, http.StatusInternalServerError)
			return
		}

		// Membuat path lengkap untuk menyimpan file
		filePath := staticPath + header.Filename

		// Membuat file di path yang telah ditentukan
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, `{"error": "Unable to create file"}`, http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
			return
		}

		// Menambahkan judul file dan path file ke dalam struktur SuratMasuk
		surat.Title = header.Filename
		surat.File = filePath

		// Panggil service untuk menyimpan data surat masuk
		newSurat, err := c.service.AddSuratMasuk(surat, parsedDate)
		if err != nil {
			http.Error(w, `{"error": "Error adding surat masuk"}`, http.StatusInternalServerError)
			return
		}

		// Kirim response JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newSurat)
	}
}

func (c *SuratMasukController) GetSuratMasuk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	suratMasukList, err := c.service.GetSuratMasuk()
	if err != nil {
		http.Error(w, `{"error": "Error retrieving data"}`, http.StatusInternalServerError)
		return
	}

	// Kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suratMasukList)
}

func (c *SuratMasukController) GetSuratById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Panggil service untuk mengambil data surat masuk
	surat, err := c.service.GetSuratById(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Surat masuk not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Failed to retrieve surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Jika berhasil, kirim response JSON dalam bentuk array
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(surat) // surat adalah slice, sehingga di-encode sebagai array JSON
}

func (c *SuratMasukController) GetCountSuratMasuk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.service.GetCountSuratMasuk()
	if err != nil {
		log.Println("Error retrieving count surat masuk:", err)
		http.Error(w, `{"error": "Error retrieving count surat masuk"}`, http.StatusInternalServerError)
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

func (c *SuratMasukController) UpdateSuratMasukByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Batas ukuran file: 10 MB
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
		return
	}

	nomor := r.FormValue("nomor")
	tanggal := r.FormValue("tanggal")
	perihal := r.FormValue("perihal")
	asal := r.FormValue("asal")
	title := r.FormValue("title")

	file, handler, err := r.FormFile("file")
	var filePath string
	if err == nil {
		defer file.Close()

		filePath = fmt.Sprintf("static/suratmasuk/%s", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error saving file:", err)
			http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
			return
		}
	} else {
		filePath = r.FormValue("existing_file")
	}

	surat := model.SuratMasuk{
		Nomor:   nomor,
		Tanggal: tanggal,
		Perihal: perihal,
		Asal:    asal,
		Title:   title,
		File:    filePath,
	}

	err = c.service.UpdateSuratMasukByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		http.Error(w, `{"error": "Error updating surat masuk"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Surat masuk updated successfully"}`))
}

func (c *SuratMasukController) DeleteSuratMasuk(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Ambil ID dari parameter URL
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Panggil service untuk menghapus surat masuk
	err = c.service.DeleteSuratMasuk(id)
	if err != nil {
		log.Println("Error deleting surat masuk:", err)

		// Jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error": "Surat masuk not found"}`, http.StatusNotFound)
			return
		}

		// Jika terjadi error lain
		http.Error(w, `{"error": "Failed to delete surat masuk"}`, http.StatusInternalServerError)
		return
	}

	// Jika berhasil, kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Surat masuk deleted successfully"}`))
}

