package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type SuratKeluarService struct {
	repo *repository.SuratKeluarRepository
}

func NewSuratKeluarService(repo *repository.SuratKeluarRepository) *SuratKeluarService {
	return &SuratKeluarService{repo: repo}
}

func (s *SuratKeluarService) AddSuratKeluar(surat *model.SuratKeluar, file io.Reader, fileName string) error {
	// Validasi field
	if surat.ID == 0 || surat.Nomor == "" || surat.Tanggal == "" || surat.Perihal == "" || surat.Ditujukan == "" || surat.Title == "" {
		return fmt.Errorf("semua field wajib diisi")
	}
	if filepath.Ext(fileName) != ".pdf" {
		return fmt.Errorf("file harus berupa PDF")
	}

	// Parse tanggal
	parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
	if err != nil {
		return fmt.Errorf("format tanggal tidak valid: %v", err)
	}

	// Simpan file PDF
	staticPath := "./static/suratkeluar/"
	fileID := uuid.New().String()
	filePath := filepath.Join(staticPath, fileID+filepath.Ext(fileName))
	if err := os.MkdirAll(staticPath, os.ModePerm); err != nil {
		return fmt.Errorf("gagal membuat direktori: %v", err)
	}
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %v", err)
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, file); err != nil {
		return fmt.Errorf("gagal menyimpan file: %v", err)
	}
	surat.File = filePath

	// Simpan ke repository
	return s.repo.AddSuratKeluar(surat, parsedDate)
}


// GetAllSuratKeluar mengambil semua data surat keluar dari repository
func (s *SuratKeluarService) GetAllSuratKeluar() ([]model.SuratKeluar, error) {
	return s.repo.GetAllSuratKeluar()
}

// GetSuratKeluarById mengambil data surat keluar berdasarkan ID
func (s *SuratKeluarService) GetSuratKeluarById(id int) ([]model.SuratKeluar, error) {
	// Validasi ID
	if id <= 0 {
		return nil, fmt.Errorf("ID harus lebih besar dari 0")
	}

	// Panggil repository untuk mengambil data surat keluar
	surat, err := s.repo.GetSuratKeluarById(id)
	if err != nil {
		// Jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("surat keluar dengan ID %d tidak ditemukan", id)
		}

		// Jika terjadi error lain
		log.Printf("Error retrieving surat keluar by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil surat keluar: %v", err)
	}

	// Bungkus data dalam slice (array)
	return []model.SuratKeluar{*surat}, nil
}

// UpdateSuratKeluarByID memperbarui data surat keluar berdasarkan ID
func (s *SuratKeluarService) UpdateSuratKeluarByID(id int, surat model.SuratKeluar) error {
	err := s.repo.UpdateSuratKeluarByID(id, surat)
	if err != nil {
		log.Println("Error updating surat keluar from repository:", err)
		return err
	}
	return nil
}

func (s *SuratKeluarService) DeleteSuratKeluar(id int) error {
	err := s.repo.DeleteSuratKeluar(id)
	if err != nil {
		log.Println("Error deleting surat keluar:", err)
		return err
	}
	return nil
}
