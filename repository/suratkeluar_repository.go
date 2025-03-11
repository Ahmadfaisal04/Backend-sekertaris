package repository

import (
	"Sekertaris/model"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func AddSuratKeluar(w http.ResponseWriter, r *http.Request, surat model.SuratKeluar, parsedDate time.Time, db *sql.DB) {
	query := "INSERT INTO suratkeluar (nomor, tanggal, perihal, ditujukan, file, title_file) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, surat.Nomor, parsedDate, surat.Perihal, surat.Ditujukan, surat.File, surat.Title)
	if err != nil {
		http.Error(w, `{"Error Message": "Error inserting data"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Surat Keluar created successfully"}`))
}
func GetCountSuratKeluar(w http.ResponseWriter, R *http.Request, db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT count(*) FROM suratkeluar").Scan(&count)
if err != nil {
    if err == sql.ErrNoRows {
        http.Error(w, `{"error": "Data not found"}`, http.StatusNotFound)
    } else {
        log.Println("Error retrieving count surat keluar:", err)
        http.Error(w, `{"error": "Error retrieving count surat keluar"}`, http.StatusInternalServerError)
    }
    return
}

	response := map[string]int{"count": count}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, `{"error": "Error creating JSON response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetSuratKeluarByid(db *sql.DB, w http.ResponseWriter, nomor int) model.SuratKeluar {
	var existingSurat model.SuratKeluar
	err := db.QueryRow("SELECT nomor, tanggal, perihal, ditujukan, title_file FROM suratkeluar WHERE nomor = ?", nomor).Scan(
		&existingSurat.Nomor, &existingSurat.Tanggal, &existingSurat.Perihal, &existingSurat.Ditujukan, &existingSurat.Title)
	if err != nil {
		log.Println("Error fetching existing surat keluar:", err)
		http.Error(w, `{"Error Message": "Surat Keluar not found"}`, http.StatusNotFound)
		return existingSurat
	}
	return existingSurat
}

func GetSuratKeluar(w http.ResponseWriter, r *http.Request, db *sql.DB, nomor string) []model.SuratKeluar {
	var suratMasukList []model.SuratKeluar
	rows, err := db.Query("SELECT id, nomor, tanggal, perihal, ditujukan, title_file FROM suratkeluar")
	if err != nil {
		log.Println("Error retrieving surat masuk:", err)
		http.Error(w, `{"Error Message": "Error retrieving surat masuk"}`, http.StatusInternalServerError)
		return suratMasukList
	}
	defer rows.Close()

	for rows.Next() {
		var surat model.SuratKeluar
		if err != nil {
			log.Println("Error retrieving surat keluar:", err)
			http.Error(w, `{"Error Message": "Surat Keluar not found"}`, http.StatusNotFound)
			return []model.SuratKeluar{} // Pastikan return nilai kosong
		}

		suratMasukList = append(suratMasukList, surat)
	}
	return suratMasukList
}

func UpdateSuratKeluar(w http.ResponseWriter, db *sql.DB, existingSurat model.SuratKeluar, nomor string) model.SuratKeluar {
	result, err := db.Exec("UPDATE suratkeluar SET nomor = ?, tanggal = ?, perihal = ?, ditujukan = ?, title_file =? WHERE nomor = ?",
		existingSurat.Nomor, existingSurat.Tanggal, existingSurat.Perihal, existingSurat.Ditujukan, existingSurat.Title, nomor)
	if err != nil {
		log.Println("Error updating surat keluar:", err)
		http.Error(w, `{"Error Message": "Error updating surat keluar"}`, http.StatusInternalServerError)
		return existingSurat
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error fetching rows affected:", err)
		http.Error(w, `{"Error Message": "Error processing request"}`, http.StatusInternalServerError)
		return existingSurat
	}

	if rowsAffected == 0 {
		http.Error(w, `{"Error Message": "Surat Keluar not found"}`, http.StatusNotFound)
		return existingSurat
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Surat Keluar updated successfully"}`))
	return existingSurat
}
