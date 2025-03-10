package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"io"
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
