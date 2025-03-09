package repository

import (
	"Sekertaris/model"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func AddSuratMasuk(db *sql.DB, w http.ResponseWriter, r *http.Request, surat model.SuratMasuk, parsedDate time.Time) {
	query := `INSERT INTO suratmasuk (nomor, tanggal, perihal, asal, title, file) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, surat.Nomor, parsedDate, surat.Perihal, surat.Asal, surat.Title, surat.File)
	if err != nil {
		http.Error(w, `{"Error Message": "Error inserting data"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Surat Masuk created successfully"}`))
}

func GetSuratMasuk(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nomor, tanggal, perihal, asal, title FROM suratmasuk")
	if err != nil {
		log.Println("Error retrieving surat masuk:", err)
		http.Error(w, `{"Error Message": "Error retrieving surat masuk"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
}
