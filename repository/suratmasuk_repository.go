package repository

import (
	"Sekertaris/model"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SuratMasukRepository struct {
	db *sql.DB
}

func NewSuratMasukRepository(db *sql.DB) *SuratMasukRepository {
	return &SuratMasukRepository{db: db}
}

func AddSuratMasuk(db *sql.DB, w http.ResponseWriter, r *http.Request, surat model.SuratMasuk, parsedDate time.Time) {
	// Query untuk memasukkan data
	query := `INSERT INTO suratmasuk (nomor, tanggal, perihal, asal, title, file) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, surat.Nomor, parsedDate, surat.Perihal, surat.Asal, surat.Title, surat.File)
	if err != nil {
		http.Error(w, `{"error": "Error inserting data"}`, http.StatusInternalServerError)
		return
	}

	// Ambil ID dari data yang baru saja dimasukkan
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "Error getting last insert ID"}`, http.StatusInternalServerError)
		return
	}

	// Query untuk mengambil data yang baru saja dimasukkan
	var newSurat model.SuratMasuk
	query = `SELECT id, nomor, tanggal, perihal, asal, title, file FROM suratmasuk WHERE id = ?`
	err = db.QueryRow(query, lastInsertID).Scan(&newSurat.Id, &newSurat.Nomor, &newSurat.Tanggal, &newSurat.Perihal, &newSurat.Asal, &newSurat.Title, &newSurat.File)
	if err != nil {
		http.Error(w, `{"error": "Error retrieving inserted data"}`, http.StatusInternalServerError)
		return
	}

	// Set header dan kembalikan data sebagai JSON
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newSurat)
}

func GetSuratMasuk(db *sql.DB) ([]model.SuratMasuk, error) {
	rows, err := db.Query("SELECT id, nomor, tanggal, perihal, asal, title, file FROM suratmasuk")
	if err != nil {
		log.Println("Error retrieving surat masuk:", err)
		return nil, err
	}
	defer rows.Close()

	var suratMasukList []model.SuratMasuk
	for rows.Next() {
		var surat model.SuratMasuk
		if err := rows.Scan(&surat.Id, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Asal, &surat.Title, &surat.File); err != nil {
			log.Println("Error scanning surat masuk row:", err)
			return nil, err
		}
		suratMasukList = append(suratMasukList, surat)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after retrieving surat masuk:", err)
		return nil, err
	}

	return suratMasukList, nil
}

// GetSuratById mengambil data surat masuk berdasarkan ID
func (r *SuratMasukRepository) GetSuratById(id int) (*model.SuratMasuk, error) {
	var surat model.SuratMasuk
	err := r.db.QueryRow("SELECT id, nomor, tanggal, perihal, asal, title, file FROM suratmasuk WHERE id = ?", id).
		Scan(&surat.Id, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Asal, &surat.Title, &surat.File)
	if err != nil {
		if err == sql.ErrNoRows {
			// Data tidak ditemukan
			return nil, nil
		}
		log.Println("Error retrieving surat masuk by ID:", err)
		return nil, err
	}
	return &surat, nil
}

func (r *SuratMasukRepository) GetCountSuratMasuk() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM suratmasuk").Scan(&count)
	if err != nil {
		log.Println("Error retrieving count surat masuk:", err)
		return 0, err
	}
	return count, nil
}

// UpdateSuratMasukByID memperbarui data surat masuk berdasarkan ID
func (r *SuratMasukRepository) UpdateSuratMasukByID(id int, surat model.SuratMasuk) error {
	query := `
		UPDATE suratmasuk 
		SET nomor = ?, tanggal = ?, perihal = ?, asal = ?, title = ?, file = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, surat.Nomor, surat.Tanggal, surat.Perihal, surat.Asal, surat.Title, surat.File, id)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		return err
	}
	return nil
}

// DeleteSuratMasuk menghapus data surat masuk berdasarkan nomor dan perihal
func (r *SuratMasukRepository) DeleteSuratMasuk(nomor, perihal string) error {
	query := "DELETE FROM suratmasuk WHERE nomor = ? AND perihal = ?"
	result, err := r.db.Exec(query, nomor, perihal)
	if err != nil {
		log.Println("Error deleting surat masuk:", err)
		return err
	}

	// Cek apakah ada baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		// Tidak ada data yang dihapus
		return sql.ErrNoRows
	}

	return nil

}
