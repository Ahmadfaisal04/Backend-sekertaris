package repository

import (
	"Sekertaris/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type SuratKeluarRepository struct {
	db *sql.DB
}

func NewSuratKeluarRepository(db *sql.DB) *SuratKeluarRepository {
	return &SuratKeluarRepository{db: db}
}

func (r *SuratKeluarRepository) AddSuratKeluar(surat *model.SuratKeluar, parsedDate time.Time) error {
	query := `INSERT INTO suratkeluar (id, nomor, tanggal, perihal, ditujukan, title, file) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, surat.ID, surat.Nomor, parsedDate, surat.Perihal, surat.Ditujukan, surat.Title, surat.File,
		time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("nomor surat sudah digunakan")
		}
		return fmt.Errorf("gagal menyimpan surat: %v", err)
	}
	return nil
}



// GetAllSuratKeluar mengambil semua data surat keluar dari database
func (r *SuratKeluarRepository) GetAllSuratKeluar() ([]model.SuratKeluar, error) {
	// Tambahkan ORDER BY created_at DESC

	rows, err := r.db.Query("SELECT id, nomor, tanggal, perihal, ditujukan, title, file FROM suratkeluar ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error retrieving all surat keluar:", err)
		return nil, err
	}
	defer rows.Close()

	var suratKeluarList []model.SuratKeluar
	for rows.Next() {
		var surat model.SuratKeluar
		if err := rows.Scan(&surat.ID, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Ditujukan, &surat.Title, &surat.File); err != nil {
			log.Println("Error scanning surat keluar row:", err)
			return nil, err
		}
		suratKeluarList = append(suratKeluarList, surat)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after retrieving all surat keluar:", err)
		return nil, err
	}

	return suratKeluarList, nil
}

// GetSuratKeluarById mengambil data surat keluar berdasarkan ID
func (r *SuratKeluarRepository) GetSuratKeluarById(id int) (*model.SuratKeluar, error) {
	var surat model.SuratKeluar

	query := "SELECT id, nomor, tanggal, perihal, ditujukan, title, file, created_at, updated_at FROM suratkeluar WHERE id = ? "
	err := r.db.QueryRow(query, id).Scan(&surat.ID, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Ditujukan, &surat.Title, &surat.File)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("surat keluar dengan ID %d tidak ditemukan", id)
		}
		log.Printf("Error retrieving surat keluar by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil surat keluar: %v", err)
	}
	return &surat, nil
}

// UpdateSuratKeluarByID memperbarui data surat keluar berdasarkan ID
func (r *SuratKeluarRepository) UpdateSuratKeluarByID(id int, surat model.SuratKeluar) error {

	query := `UPDATE suratkeluar 
          SET nomor = ?, tanggal = ?, perihal = ?, ditujukan = ?, title = ?, file = ?, updated_at = ?
          WHERE id = ?`
	_, err := r.db.Exec(query, surat.Nomor, surat.Tanggal, surat.Perihal, surat.Ditujukan, surat.Title, surat.File,
		time.Now().Format("2006-01-02 15:04:05"), id)

	if err != nil {
		log.Println("Error updating surat keluar:", err)
		return err
	}
	return nil
}

func (r *SuratKeluarRepository) DeleteSuratKeluar(id int) error {
	query := "DELETE FROM suratkeluar WHERE id = ? "
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting surat keluar:", err)
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
