package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type PermohonanSuratController struct {
	service *service.PermohonanSuratService
}

func NewPermohonanSuratController(service *service.PermohonanSuratService) *PermohonanSuratController {
	return &PermohonanSuratController{service: service}
}

func (c *PermohonanSuratController) AddPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
			return
		}

		var permohonan model.PermohonanSurat
		permohonan.NIK = r.FormValue("nik")
		permohonan.NamaLengkap = r.FormValue("nama_lengkap")
		permohonan.TempatLahir = r.FormValue("tempat_lahir")
		tanggalLahirStr := r.FormValue("tanggal_lahir")
		permohonan.JenisKelamin = model.JenisKelamin(r.FormValue("jenis_kelamin"))
		permohonan.Pendidikan = r.FormValue("pendidikan")
		permohonan.Pekerjaan = r.FormValue("pekerjaan")
		permohonan.Agama = r.FormValue("agama")
		permohonan.StatusPernikahan = r.FormValue("status_pernikahan")
		permohonan.Kewarganegaraan = r.FormValue("kewarganegaraan")
		permohonan.AlamatLengkap = r.FormValue("alamat_lengkap")
		permohonan.JenisSurat = r.FormValue("jenis_surat")
		permohonan.Keterangan = r.FormValue("keterangan")
		permohonan.NomorHP = r.FormValue("nomor_hp")
		permohonan.Status = model.Status(r.FormValue("status"))
		permohonan.CreatedAt = time.Now()
		permohonan.UpdatedAt = time.Now()

		// Validate required fields
		if permohonan.NIK == "" || permohonan.NamaLengkap == "" || tanggalLahirStr == "" {
			http.Error(w, `{"error": "NIK, Nama Lengkap, and Tanggal Lahir are required"}`, http.StatusBadRequest)
			return
		}

		parsedDate, err := time.Parse("2006-01-02", tanggalLahirStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format for tanggal_lahir, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		permohonan.TanggalLahir = parsedDate

		// Validate JenisKelamin
		if permohonan.JenisKelamin != model.LakiLaki && permohonan.JenisKelamin != model.Perempuan {
			http.Error(w, `{"error": "Jenis Kelamin must be 'Laki-laki' or 'Perempuan'"}`, http.StatusBadRequest)
			return
		}

		// Validate Status
		if permohonan.Status != model.Pending && permohonan.Status != model.Diproses && permohonan.Status != model.Selesai {
			http.Error(w, `{"error": "Status must be 'Pending', 'Diproses', or 'Selesai'"}`, http.StatusBadRequest)
			return
		}

		// Handle file upload
		file, header, err := r.FormFile("dokumen")
		if err == nil {
			defer file.Close()

			staticPath := "./static/permohonansurat/"
			err = os.MkdirAll(staticPath, os.ModePerm)
			if err != nil {
				http.Error(w, `{"error": "Unable to create static directory"}`, http.StatusInternalServerError)
				return
			}

			filePath := staticPath + header.Filename
			outFile, err := os.Create(filePath)
			if err != nil {
				http.Error(w, `{"error": "Unable to create file"}`, http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, file)
			if err != nil {
				http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
				return
			}

			permohonan.DokumenURL = sql.NullString{String: filePath, Valid: true}
		}

		newPermohonan, err := c.service.AddPermohonanSurat(permohonan)
		if err != nil {
			http.Error(w, `{"error": "Error adding permohonan surat"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPermohonan)
	}
}

func (c *PermohonanSuratController) GetPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	permohonanSuratList, err := c.service.GetPermohonanSurat()
	if err != nil {
		http.Error(w, `{"error": "Error retrieving data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permohonanSuratList)
}

func (c *PermohonanSuratController) GetPermohonanSuratByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	permohonan, err := c.service.GetPermohonanSuratByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Permohonan surat not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Failed to retrieve permohonan surat"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permohonan)
}

func (c *PermohonanSuratController) UpdatePermohonanSuratByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, `{"error": "Error parsing form data"}`, http.StatusBadRequest)
		return
	}

	var permohonan model.PermohonanSurat
	permohonan.NIK = r.FormValue("nik")
	permohonan.NamaLengkap = r.FormValue("nama_lengkap")
	permohonan.TempatLahir = r.FormValue("tempat_lahir")
	tanggalLahirStr := r.FormValue("tanggal_lahir")
	permohonan.JenisKelamin = model.JenisKelamin(r.FormValue("jenis_kelamin"))
	permohonan.Pendidikan = r.FormValue("pendidikan")
	permohonan.Pekerjaan = r.FormValue("pekerjaan")
	permohonan.Agama = r.FormValue("agama")
	permohonan.StatusPernikahan = r.FormValue("status_pernikahan")
	permohonan.Kewarganegaraan = r.FormValue("kewarganegaraan")
	permohonan.AlamatLengkap = r.FormValue("alamat_lengkap")
	permohonan.JenisSurat = r.FormValue("jenis_surat")
	permohonan.Keterangan = r.FormValue("keterangan")
	permohonan.NomorHP = r.FormValue("nomor_hp")
	permohonan.Status = model.Status(r.FormValue("status"))
	permohonan.UpdatedAt = time.Now()

	if tanggalLahirStr != "" {
		parsedDate, err := time.Parse("2006-01-02", tanggalLahirStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format for tanggal_lahir, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		permohonan.TanggalLahir = parsedDate
	}

	// Handle file upload
	file, handler, err := r.FormFile("dokumen")
	if err == nil {
		defer file.Close()

		staticPath := "./static/permohonansurat/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			log.Println("Error creating directory:", err)
			http.Error(w, `{"error": "Error creating directory"}`, http.StatusInternalServerError)
			return
		}

		filePath := staticPath + handler.Filename
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, `{"error": "Error creating file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
			return
		}

		permohonan.DokumenURL = sql.NullString{String: filePath, Valid: true}
	} else {
		existingDokumenURL := r.FormValue("existing_dokumen_url")
		if existingDokumenURL != "" {
			permohonan.DokumenURL = sql.NullString{String: existingDokumenURL, Valid: true}
		} else {
			permohonan.DokumenURL = sql.NullString{Valid: false}
		}
	}

	err = c.service.UpdatePermohonanSuratByID(id, permohonan)
	if err != nil {
		log.Println("Error updating permohonan surat:", err)
		http.Error(w, `{"error": "Error updating permohonan surat"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Permohonan surat updated successfully"})
}

func (c *PermohonanSuratController) DeletePermohonanSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	err = c.service.DeletePermohonanSurat(id)
	if err != nil {
		log.Println("Error deleting permohonan surat:", err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error": "Permohonan surat not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Failed to delete permohonan surat"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Permohonan surat deleted successfully"}`))
}