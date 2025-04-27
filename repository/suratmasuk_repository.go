package repository

import (
	"Sekertaris/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SuratMasukRepository struct {
	db *sql.DB
}

func NewSuratMasukRepository(db *sql.DB) *SuratMasukRepository {
	return &SuratMasukRepository{db: db}
}

func (r *SuratMasukRepository) AddSuratMasuk(surat model.SuratMasuk, parsedDate time.Time) (*model.SuratMasuk, error) {
	query := `INSERT INTO suratmasuk (nomor, tanggal, perihal, asal, title, file) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, surat.Nomor, parsedDate, surat.Perihal, surat.Asal, surat.Title, surat.File)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var newSurat model.SuratMasuk
	query = `SELECT nomor, tanggal, perihal, asal, title, file FROM suratmasuk WHERE id = ?`
	err = r.db.QueryRow(query, lastInsertID).Scan(&newSurat.Id, &newSurat.Nomor, &newSurat.Tanggal, &newSurat.Perihal, &newSurat.Asal, &newSurat.Title, &newSurat.File)
	if err != nil {
		return nil, err
	}

	return &newSurat, nil
}

func (r *SuratMasukRepository) GetSuratMasuk() ([]model.SuratMasuk, error) {
	rows, err := r.db.Query("SELECT nomor, tanggal, perihal, asal, title, file FROM suratmasuk ORDER BY created_at DESC" )
	if err != nil {
		log.Println("Error retrieving surat masuk:", err)
		return nil, err
	}
	defer rows.Close()

	var suratMasukList []model.SuratMasuk
	for rows.Next() {
		var surat model.SuratMasuk
		if err := rows.Scan(&surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Asal, &surat.Title, &surat.File); err != nil {
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

func (r *SuratMasukRepository) GetSuratById(id int) (*model.SuratMasuk, error) {
	var surat model.SuratMasuk
	query := "SELECT id, nomor, tanggal, perihal, asal, title, file FROM suratmasuk WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&surat.Id, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Asal, &surat.Title, &surat.File)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("surat dengan ID %d tidak ditemukan", id)
		}
		log.Printf("Error retrieving surat masuk by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil surat masuk: %v", err)
	}
	return &surat, nil
}



func (r *SuratMasukRepository) UpdateSuratMasukByID(id int, surat model.SuratMasuk) error {
	query := `
		UPDATE suratmasuk 
		SET nomor = ?, tanggal = ?, perihal = ?, asal = ?, title = ?, file = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, surat.Nomor, surat.Tanggal, surat.Perihal, surat.Asal, surat.Title, surat.File, id)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada surat dengan ID %d yang ditemukan", id)
	}

	return nil
}

func (r *SuratMasukRepository) DeleteSuratMasuk(id int) error {
	query := "DELETE FROM suratmasuk WHERE id = ? "
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting surat masuk:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada surat dengan id %s yang ditemukan", id)
	}

	return nil
}