package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func AddSuratKeluar(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {

		var surat model.SuratKeluar
		surat.Nomor = r.FormValue("nomor")
		surat.Tanggal = r.FormValue("tanggal")
		surat.Perihal = r.FormValue("perihal")
		surat.Ditujukan = r.FormValue("ditujukan")

		parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
		if err != nil {
			http.Error(w, `{"Error Message": "Invalid date format, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}

		// Mengambil file yang diunggah
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Menentukan path penyimpanan file di direktori static
		staticPath := "./static/suratkeluar/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create static directory", http.StatusInternalServerError)
			return
		}

		// Membuat path lengkap untuk menyimpan file
		filePath := filepath.Join(staticPath, header.Filename)

		// Membuat file di path yang telah ditentukan
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Menambahkan judul file dan path file ke dalam struktur SuratKeluar
		surat.Title = header.Filename
		surat.File = filePath
		repository.AddSuratKeluar(w, r, surat, parsedDate, db)
	}
}
func UpdateSuratKeluar(w http.ResponseWriter, r *http.Request, db *sql.DB, nomor string) {
	existingSurat := repository.GetSuratKeluar(w, r, db, nomor)

	// Parse form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, `{"Error Message": "Invalid form data"}`, http.StatusBadRequest)
		return
	}

	// Update field only if there is new input, otherwise keep the existing data
	if nomorInput := r.FormValue("nomor"); nomorInput != "" {
		existingSurat.Nomor = nomorInput
	}
	if tanggalInput := r.FormValue("tanggal"); tanggalInput != "" {
		existingSurat.Tanggal = tanggalInput
	}
	if perihalInput := r.FormValue("perihal"); perihalInput != "" {
		existingSurat.Perihal = perihalInput
	}
	if ditujukanInput := r.FormValue("ditujukan"); ditujukanInput != "" {
		existingSurat.Ditujukan = ditujukanInput
	}
	if titleFileInput := r.FormValue("title_file"); titleFileInput != "" {
		existingSurat.Title = titleFileInput
	}
	repository.UpdateSuratKeluar(w, db, existingSurat, nomor)
}
func GetSuratKeluarByid(w http.ResponseWriter, r *http.Request, id string, db *sql.DB) {

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM suratkeluar WHERE id = ?", id).Scan(&count)
	if err != nil {
		log.Println("Error retrieving count surat keluar:", err)
		http.Error(w, `{"Error Message": "Error retrieving count surat keluar"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"count": "%d"}`, count)))
}

func GetSuratKeluar(w http.ResponseWriter, r *http.Request) {
	
	if err := rows.Err(); err != nil {
		log.Println("Error after retrieving surat masuk:", err)
		http.Error(w, `{"Error Message": "Error processing request"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suratMasukList); err != nil {
		log.Println("Error encoding surat masuk list to JSON:", err)
		http.Error(w, `{"Error Message": "Error processing request"}`, http.StatusInternalServerError)
	}
}
