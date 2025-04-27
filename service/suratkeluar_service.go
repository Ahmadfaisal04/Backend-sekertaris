package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"fmt"
	"io"
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
	// Validasi input
	if surat.Nomor == "" || surat.Tanggal == "" || surat.Perihal == "" || surat.Ditujukan == "" || surat.Title == "" {
		return fmt.Errorf("semua field wajib diisi")
	}
	if filepath.Ext(fileName) != ".pdf" {
		return fmt.Errorf("file harus berupa PDF")
	}

	// Simpan file
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

	// Isi data
	now := time.Now().Format("2006-01-02 15:04:05")
	surat.File = filePath
	surat.CreatedAt = now
	surat.UpdatedAt = now

	return s.repo.AddSuratKeluar(surat)
}

func (s *SuratKeluarService) GetAllSuratKeluar() ([]model.SuratKeluar, error) {
	return s.repo.GetAllSuratKeluar()
}

func (s *SuratKeluarService) GetSuratKeluarById(id int) (*model.SuratKeluar, error) {
	return s.repo.GetSuratKeluarById(id)
}

func (s *SuratKeluarService) UpdateSuratKeluarByID(id int, surat *model.SuratKeluar, file io.Reader, fileName string) error {
	if file != nil {
		// Kalau ada file baru
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
	}

	surat.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return s.repo.UpdateSuratKeluarByID(id, surat)
}

func (s *SuratKeluarService) DeleteSuratKeluar(id int) error {
	return s.repo.DeleteSuratKeluar(id)
}
