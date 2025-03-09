package repository

import (
	"Sekertaris/model"
	"database/sql"
	"net/http"
	"time"
)

func AddSuratKeluar(w http.ResponseWriter, r *http.Request, surat model.SuratKeluar, parsedDate time.Time, db *sql.DB) {
	query := `INSERT INTO suratkeluar (nomor, tanggal, perihal, ditujukan, file, title_file) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, surat.Nomor, parsedDate, surat.Perihal, surat.Ditujukan, surat.File, surat.Title)
	if err != nil {
		http.Error(w, `{"Error Message": "Error inserting data"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Surat Keluar created successfully"}`))
}
