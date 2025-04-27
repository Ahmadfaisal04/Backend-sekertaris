package repository

import (
	"Sekertaris/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type SuratKeluarRepository struct {
	db *sql.DB
}

func NewSuratKeluarRepository(db *sql.DB) *SuratKeluarRepository {
	return &SuratKeluarRepository{db: db}
}

func (r *SuratKeluarRepository) AddSuratKeluar(surat *model.SuratKeluar) error {
	query := `INSERT INTO suratkeluar (nomor, tanggal, perihal, ditujukan, title, file, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		surat.Nomor,
		surat.Tanggal,
		surat.Perihal,
		surat.Ditujukan,
		surat.Title,
		surat.File,
		surat.CreatedAt,
		surat.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("nomor surat sudah digunakan")
		}
		return fmt.Errorf("gagal menyimpan surat: %v", err)
	}
	return nil
}

func (r *SuratKeluarRepository) GetAllSuratKeluar() ([]model.SuratKeluar, error) {
	rows, err := r.db.Query("SELECT id, nomor, tanggal, perihal, ditujukan, title, file, created_at, updated_at FROM suratkeluar ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error retrieving all surat keluar:", err)
		return nil, err
	}
	defer rows.Close()

	var suratKeluarList []model.SuratKeluar
	for rows.Next() {
		var surat model.SuratKeluar
		if err := rows.Scan(&surat.ID, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Ditujukan, &surat.Title, &surat.File, &surat.CreatedAt, &surat.UpdatedAt); err != nil {
			log.Println("Error scanning surat keluar row:", err)
			return nil, err
		}
		suratKeluarList = append(suratKeluarList, surat)
	}

	return suratKeluarList, nil
}

func (r *SuratKeluarRepository) GetSuratKeluarById(id int) (*model.SuratKeluar, error) {
	var surat model.SuratKeluar

	query := "SELECT id, nomor, tanggal, perihal, ditujukan, title, file, created_at, updated_at FROM suratkeluar WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&surat.ID, &surat.Nomor, &surat.Tanggal, &surat.Perihal, &surat.Ditujukan, &surat.Title, &surat.File, &surat.CreatedAt, &surat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("surat keluar dengan ID %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal mengambil surat keluar: %v", err)
	}
	return &surat, nil
}

func (r *SuratKeluarRepository) UpdateSuratKeluarByID(id int, surat *model.SuratKeluar) error {
	query := `UPDATE suratkeluar 
	          SET nomor = ?, tanggal = ?, perihal = ?, ditujukan = ?, title = ?, file = ?, updated_at = ?
	          WHERE id = ?`

	_, err := r.db.Exec(query,
		surat.Nomor,
		surat.Tanggal,
		surat.Perihal,
		surat.Ditujukan,
		surat.Title,
		surat.File,
		surat.UpdatedAt,
		id,
	)

	if err != nil {
		log.Println("Error updating surat keluar:", err)
		return err
	}
	return nil
}

func (r *SuratKeluarRepository) DeleteSuratKeluar(id int) error {
	query := "DELETE FROM suratkeluar WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada surat dengan id %d yang ditemukan", id)
	}
	return nil
}
