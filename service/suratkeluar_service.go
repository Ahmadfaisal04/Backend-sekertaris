package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
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
	if r.Method == "POST" {
		var surat model.SuratKeluar
		surat.Nomor = r.FormValue("nomor")
		surat.Tanggal = r.FormValue("tanggal")
		surat.Perihal = r.FormValue("perihal")
		surat.Ditujukan = r.FormValue("ditujukan") // Sesuaikan dengan field di model

		// Validasi bahwa tanggal tidak boleh kosong
		if surat.Tanggal == "" {
			http.Error(w, `{"Error Message": "Tanggal is required"}`, http.StatusBadRequest)
			return
		}

		parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
		if err != nil {
			panic(err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Menentukan path penyimpanan file di direktori static
		staticPath := "./static/suratkeluar/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create static directory", http.StatusInternalServerError)
			return
		}

		// Membuat path lengkap untuk menyimpan file
		filePath := staticPath + header.Filename

		// Membuat file di path yang telah ditentukan
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Menambahkan judul file dan path file ke dalam struktur SuratKeluar
		surat.Title = header.Filename
		surat.File = filePath

		// Panggil repository untuk menyimpan data surat keluar
		repository.AddSuratKeluar(db, w, r, surat, parsedDate)
	}
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
func (s *SuratKeluarService) GetSuratKeluarById(id int) (*model.SuratKeluar, error) {
	surat, err := s.repo.GetSuratKeluarById(id)
	if err != nil {
		log.Println("Error retrieving surat keluar by ID from repository:", err)
		return nil, err
	}
	return surat, nil
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