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

type SuratMasukService struct {
	repo *repository.SuratMasukRepository
}

func NewSuratMasukService(repo *repository.SuratMasukRepository) *SuratMasukService {
	return &SuratMasukService{repo: repo}
}

func AddSuratMasuk(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {

		var surat model.SuratMasuk
		surat.Nomor = r.FormValue("nomor")
		surat.Tanggal = r.FormValue("tanggal")
		surat.Perihal = r.FormValue("perihal")
		surat.Asal = r.FormValue("asal")

		// Validasi bahwa tanggal tidak boleh kosong
		if surat.Tanggal == "" {
			http.Error(w, `{"Error Message": "Tanggal is required"}`, http.StatusBadRequest)
			return
		}

		parsedDate, err := time.Parse("2006-01-02", surat.Tanggal)
		if err != nil {
			panic(err)
			// http.Error(w, `{"Error Message": "Invalid date format, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			// return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Menentukan path penyimpanan file di direktori static
		staticPath := "./static/suratmasuk/"
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

		// Menambahkan judul file dan path file ke dalam struktur SuratMasuk
		surat.Title = header.Filename
		surat.File = filePath

		// Panggil repository untuk menyimpan data surat masuk
		repository.AddSuratMasuk(db, w, r, surat, parsedDate)
	}
}

func GetSuratMasuk(w http.ResponseWriter, db *sql.DB) {
	suratMasukList, err := repository.GetSuratMasuk(db)
	if err != nil {
		http.Error(w, `{"Error Message": "Error retrieving data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, surat := range suratMasukList {
		idStr := strconv.Itoa(surat.Id) // Konversi Id ke string
		w.Write([]byte(`{"id":` + idStr + `,"nomor":"` + surat.Nomor + `","tanggal":"` + surat.Tanggal + `","perihal":"` + surat.Perihal + `","asal":"` + surat.Asal + `","title":"` + surat.Title + `","file":"` + surat.File + `"}`))
	}
}

// GetSuratById mengambil data surat masuk berdasarkan ID
func (s *SuratMasukService) GetSuratById(id int) (*model.SuratMasuk, error) {
	surat, err := s.repo.GetSuratById(id)
	if err != nil {
		log.Println("Error retrieving surat masuk by ID from repository:", err)
		return nil, err
	}
	return surat, nil
}

func (s *SuratMasukService) GetCountSuratMasuk() (int, error) {
	count, err := s.repo.GetCountSuratMasuk()
	if err != nil {
		log.Println("Error retrieving count surat masuk from repository:", err)
		return 0, err
	}
	return count, nil
}

// UpdateSuratMasukByID memperbarui data surat masuk berdasarkan ID
func (s *SuratMasukService) UpdateSuratMasukByID(id int, surat model.SuratMasuk) error {
	err := s.repo.UpdateSuratMasukByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk from repository:", err)
		return err
	}
	return nil
}

// DeleteSuratMasuk menghapus data surat masuk berdasarkan nomor dan perihal
func (s *SuratMasukService) DeleteSuratMasuk(nomor, perihal string) error {
	err := s.repo.DeleteSuratMasuk(nomor, perihal)
	if err != nil {
		log.Println("Error deleting surat masuk from repository:", err)
		return err
	}
	return nil
}
