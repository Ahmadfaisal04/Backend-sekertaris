package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type SuratKeluarService struct {
	repo *repository.SuratKeluarRepository
}

func NewSuratKeluarService(repo *repository.SuratKeluarRepository) *SuratKeluarService {
	return &SuratKeluarService{repo: repo}
}

func AddSuratKeluar(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"Error Message": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var surat model.SuratKeluar
	surat.Nomor = r.FormValue("nomor")
	surat.Tanggal = r.FormValue("tanggal")
	surat.Perihal = r.FormValue("perihal")
	surat.Ditujukan = r.FormValue("ditujukan")

	if surat.Tanggal == "" {
		http.Error(w, `{"Error Message": "Tanggal is required"}`, http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
	if err != nil {
		http.Error(w, `{"Error Message": "Invalid date format (YYYY-MM-DD)"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"Error Message": "Unable to get file from form"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Tentukan direktori penyimpanan
	staticPath := "./static/suratkeluar/"
	if err := os.MkdirAll(staticPath, os.ModePerm); err != nil {
		http.Error(w, `{"Error Message": "Unable to create static directory"}`, http.StatusInternalServerError)
		return
	}

	// Simpan file ke disk
	filePath := filepath.Join(staticPath, header.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"Error Message": "Unable to create file"}`, http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, `{"Error Message": "Error saving file"}`, http.StatusInternalServerError)
		return
	}

	// Simpan detail ke struct
	surat.Title = header.Filename
	surat.File = filePath

	// Panggil repository (jika ingin diganti ke service, tinggal ubah)
	repository.AddSuratKeluar(db, w, r, surat, parsedDate)
}

func GetSuratKeluar(w http.ResponseWriter, db *sql.DB) {
	suratKeluarList, err := repository.GetSuratKeluar(db)
	if err != nil {
		http.Error(w, `{"Error Message": "Error retrieving data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, surat := range suratKeluarList {
		idStr := strconv.Itoa(surat.ID) // Konversi Id ke string
		w.Write([]byte(`{"id":` + idStr + `,"nomor":"` + surat.Nomor + `","tanggal":"` + surat.Tanggal + `","perihal":"` + surat.Perihal + `","ditujukan":"` + surat.Ditujukan + `","title":"` + surat.Title + `","file":"` + surat.File + `"}`))
	}
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

func (s *SuratKeluarService) GetCountSuratKeluar() (int, error) {
	count, err := s.repo.GetCountSuratKeluar()
	if err != nil {
		log.Println("Error retrieving count surat keluar from repository:", err)
		return 0, err
	}
	return count, nil
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
