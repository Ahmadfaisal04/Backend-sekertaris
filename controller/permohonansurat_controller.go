package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
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

// response adalah struktur untuk respons JSON yang konsisten
type response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (c *PermohonanSuratController) AddPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	var permohonan model.PermohonanSurat
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		// Tangani input JSON
		var requestBody struct {
			NIK              string  `json:"nik"`
			NamaLengkap      string  `json:"nama_lengkap"`
			TempatLahir      string  `json:"tempat_lahir"`
			TanggalLahir     string  `json:"tanggal_lahir"` // Format: YYYY-MM-DD
			JenisKelamin     string  `json:"jenis_kelamin"`
			Pendidikan       string  `json:"pendidikan"`
			Pekerjaan        string  `json:"pekerjaan"`
			Agama            string  `json:"agama"`
			StatusPernikahan string  `json:"status_pernikahan"`
			Kewarganegaraan  string  `json:"kewarganegaraan"`
			AlamatLengkap    string  `json:"alamat_lengkap"`
			JenisSurat       string  `json:"jenis_surat"`
			Keterangan       string  `json:"keterangan"`
			NomorHP          string  `json:"nomor_hp"`
			Status           string  `json:"status"`
			Ditujukan        string  `json:"ditujukan"`
			NamaUsaha        *string `json:"nama_usaha"`
			JenisUsaha       *string `json:"jenis_usaha"`
			AlamatUsaha      *string `json:"alamat_usaha"`
			AlamatTujuan     *string `json:"alamat_tujuan"`
			AlasanPindah     *string `json:"alasan_pindah"`
			NamaAyah         *string `json:"nama_ayah"`
			NamaIbu          *string `json:"nama_ibu"`
			TglKematian      *string `json:"tgl_kematian"` // Format: YYYY-MM-DD
			PenyebabKematian *string `json:"penyebab_kematian"`
			DokumenURL       *string `json:"dokumen_url"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			log.Printf("Error parsing JSON body: %v", err)
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Body JSON tidak valid"})
			return
		}

		// Validasi kolom wajib untuk JSON
		if requestBody.NIK == "" || requestBody.NamaLengkap == "" || requestBody.TanggalLahir == "" {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "NIK, Nama Lengkap, dan Tanggal Lahir wajib diisi"})
			return
		}

		// Mapping JSON ke struct PermohonanSurat
		permohonan.NIK = requestBody.NIK
		permohonan.NamaLengkap = requestBody.NamaLengkap
		permohonan.TempatLahir = requestBody.TempatLahir
		parsedDate, err := time.Parse("2006-01-02", requestBody.TanggalLahir)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		permohonan.TanggalLahir = parsedDate
		permohonan.JenisKelamin = model.JenisKelamin(requestBody.JenisKelamin)
		permohonan.Pendidikan = requestBody.Pendidikan
		permohonan.Pekerjaan = requestBody.Pekerjaan
		permohonan.Agama = requestBody.Agama
		permohonan.StatusPernikahan = requestBody.StatusPernikahan
		permohonan.Kewarganegaraan = requestBody.Kewarganegaraan
		permohonan.AlamatLengkap = requestBody.AlamatLengkap
		permohonan.JenisSurat = requestBody.JenisSurat
		permohonan.Keterangan = requestBody.Keterangan
		permohonan.NomorHP = requestBody.NomorHP
		permohonan.Status = model.Status(requestBody.Status)
		permohonan.Ditujukan = requestBody.Ditujukan
		if requestBody.NamaUsaha != nil {
			permohonan.NamaUsaha = sql.NullString{String: *requestBody.NamaUsaha, Valid: true}
		}
		if requestBody.JenisUsaha != nil {
			permohonan.JenisUsaha = sql.NullString{String: *requestBody.JenisUsaha, Valid: true}
		}
		if requestBody.AlamatUsaha != nil {
			permohonan.AlamatUsaha = sql.NullString{String: *requestBody.AlamatUsaha, Valid: true}
		}
		if requestBody.AlamatTujuan != nil {
			permohonan.AlamatTujuan = sql.NullString{String: *requestBody.AlamatTujuan, Valid: true}
		}
		if requestBody.AlasanPindah != nil {
			permohonan.AlasanPindah = sql.NullString{String: *requestBody.AlasanPindah, Valid: true}
		}
		if requestBody.NamaAyah != nil {
			permohonan.NamaAyah = sql.NullString{String: *requestBody.NamaAyah, Valid: true}
		}
		if requestBody.NamaIbu != nil {
			permohonan.NamaIbu = sql.NullString{String: *requestBody.NamaIbu, Valid: true}
		}
		if requestBody.TglKematian != nil {
			parsedTglKematian, err := time.Parse("2006-01-02", *requestBody.TglKematian)
			if err != nil {
				writeJSONResponse(w, http.StatusBadRequest, response{Error: "Invalid date format for tgl_kematian, expected YYYY-MM-DD"})
				return
			}
			permohonan.TglKematian = sql.NullTime{Time: parsedTglKematian, Valid: true}
		}
		if requestBody.PenyebabKematian != nil {
			permohonan.PenyebabKematian = sql.NullString{String: *requestBody.PenyebabKematian, Valid: true}
		}
		if requestBody.DokumenURL != nil {
			permohonan.DokumenURL = sql.NullString{String: *requestBody.DokumenURL, Valid: true}
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Tangani input multipart/form-data
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Printf("Error parsing form data: %v", err)
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Error parsing form data"})
			return
		}

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
		permohonan.Ditujukan = r.FormValue("ditujukan")

		// Kolom nullable
		if namaUsaha := r.FormValue("nama_usaha"); namaUsaha != "" {
			permohonan.NamaUsaha = sql.NullString{String: namaUsaha, Valid: true}
		}
		if jenisUsaha := r.FormValue("jenis_usaha"); jenisUsaha != "" {
			permohonan.JenisUsaha = sql.NullString{String: jenisUsaha, Valid: true}
		}
		if alamatUsaha := r.FormValue("alamat_usaha"); alamatUsaha != "" {
			permohonan.AlamatUsaha = sql.NullString{String: alamatUsaha, Valid: true}
		}
		if alamatTujuan := r.FormValue("alamat_tujuan"); alamatTujuan != "" {
			permohonan.AlamatTujuan = sql.NullString{String: alamatTujuan, Valid: true}
		}
		if alasanPindah := r.FormValue("alasan_pindah"); alasanPindah != "" {
			permohonan.AlasanPindah = sql.NullString{String: alasanPindah, Valid: true}
		}
		if namaAyah := r.FormValue("nama_ayah"); namaAyah != "" {
			permohonan.NamaAyah = sql.NullString{String: namaAyah, Valid: true}
		}
		if namaIbu := r.FormValue("nama_ibu"); namaIbu != "" {
			permohonan.NamaIbu = sql.NullString{String: namaIbu, Valid: true}
		}
		if tglKematian := r.FormValue("tgl_kematian"); tglKematian != "" {
			parsedTglKematian, err := time.Parse("2006-01-02", tglKematian)
			if err != nil {
				writeJSONResponse(w, http.StatusBadRequest, response{Error: "Invalid date format for tgl_kematian, expected YYYY-MM-DD"})
				return
			}
			permohonan.TglKematian = sql.NullTime{Time: parsedTglKematian, Valid: true}
		}
		if penyebabKematian := r.FormValue("penyebab_kematian"); penyebabKematian != "" {
			permohonan.PenyebabKematian = sql.NullString{String: penyebabKematian, Valid: true}
		}

		// Validasi kolom wajib untuk multipart
		if permohonan.NIK == "" || permohonan.NamaLengkap == "" || tanggalLahirStr == "" {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "NIK, Nama Lengkap, dan Tanggal Lahir wajib diisi"})
			return
		}

		parsedDate, err := time.Parse("2006-01-02", tanggalLahirStr)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		permohonan.TanggalLahir = parsedDate

		// Handle file upload
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()

			staticPath := "./static/suratkeluar/"
			err = os.MkdirAll(staticPath, os.ModePerm)
			if err != nil {
				log.Printf("Error creating directory: %v", err)
				writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal membuat direktori"})
				return
			}

			filePath := staticPath + header.Filename
			outFile, err := os.Create(filePath)
			if err != nil {
				log.Printf("Error creating file: %v", err)
				writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal membuat file"})
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, file)
			if err != nil {
				log.Printf("Error saving file: %v", err)
				writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal menyimpan file"})
				return
			}

			permohonan.DokumenURL = sql.NullString{String: filePath, Valid: true}
		}
	} else {
		writeJSONResponse(w, http.StatusUnsupportedMediaType, response{Error: "Content-Type harus application/json atau multipart/form-data"})
		return
	}

	// Validasi umum
	if permohonan.JenisKelamin != model.LakiLaki && permohonan.JenisKelamin != model.Perempuan {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "Jenis Kelamin harus 'Laki-laki' atau 'Perempuan'"})
		return
	}

	// Gunakan fungsi AddPermohonanSuratJSON dari service
	jsonData, err := c.service.AddPermohonanSuratJSON(permohonan)
	if err != nil {
		log.Printf("Error adding permohonan surat: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal menambahkan permohonan surat: " + err.Error()})
		return
	}

	// Kirim respons JSON
	writeJSONResponse(w, http.StatusCreated, response{
		Message: "Permohonan surat berhasil ditambahkan",
		Data:    json.RawMessage(jsonData),
	})
}

func (c *PermohonanSuratController) GetPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	permohonans, err := c.service.GetPermohonanSurat()
	if err != nil {
		log.Printf("Error getting permohonan surat: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal mengambil data permohonan surat"})
		return
	}

	writeJSONResponse(w, http.StatusOK, response{
		Message: "Berhasil mengambil data permohonan surat",
		Data:    permohonans,
	})
}

func (c *PermohonanSuratController) GetPermohonanSuratByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "ID tidak valid"})
		return
	}

	permohonan, err := c.service.GetPermohonanSuratByID(id)
	if err != nil {
		log.Printf("Error getting permohonan surat by ID: %v", err)
		writeJSONResponse(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, response{
		Message: "Berhasil mengambil data permohonan surat",
		Data:    permohonan,
	})
}

func (c *PermohonanSuratController) UpdatePermohonanSuratByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPut {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "ID tidak valid"})
		return
	}

	var permohonan model.PermohonanSurat
	permohonan.ID = id

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var requestBody struct {
			NIK              string  `json:"nik"`
			NamaLengkap      string  `json:"nama_lengkap"`
			TempatLahir      string  `json:"tempat_lahir"`
			TanggalLahir     string  `json:"tanggal_lahir"`
			JenisKelamin     string  `json:"jenis_kelamin"`
			Pendidikan       string  `json:"pendidikan"`
			Pekerjaan        string  `json:"pekerjaan"`
			Agama            string  `json:"agama"`
			StatusPernikahan string  `json:"status_pernikahan"`
			Kewarganegaraan  string  `json:"kewarganegaraan"`
			AlamatLengkap    string  `json:"alamat_lengkap"`
			JenisSurat       string  `json:"jenis_surat"`
			Keterangan       string  `json:"keterangan"`
			NomorHP          string  `json:"nomor_hp"`
			Status           string  `json:"status"`
			Ditujukan        string  `json:"ditujukan"`
			NamaUsaha        *string `json:"nama_usaha"`
			JenisUsaha       *string `json:"jenis_usaha"`
			AlamatUsaha      *string `json:"alamat_usaha"`
			AlamatTujuan     *string `json:"alamat_tujuan"`
			AlasanPindah     *string `json:"alasan_pindah"`
			NamaAyah         *string `json:"nama_ayah"`
			NamaIbu          *string `json:"nama_ibu"`
			TglKematian      *string `json:"tgl_kematian"`
			PenyebabKematian *string `json:"penyebab_kematian"`
			DokumenURL       *string `json:"dokumen_url"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			log.Printf("Error parsing JSON body: %v", err)
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Body JSON tidak valid"})
			return
		}

		// Validasi kolom wajib
		if requestBody.NIK == "" || requestBody.NamaLengkap == "" || requestBody.TanggalLahir == "" {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "NIK, Nama Lengkap, dan Tanggal Lahir wajib diisi"})
			return
		}

		// Mapping
		permohonan.NIK = requestBody.NIK
		permohonan.NamaLengkap = requestBody.NamaLengkap
		permohonan.TempatLahir = requestBody.TempatLahir
		parsedDate, err := time.Parse("2006-01-02", requestBody.TanggalLahir)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, response{Error: "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		permohonan.TanggalLahir = parsedDate
		permohonan.JenisKelamin = model.JenisKelamin(requestBody.JenisKelamin)
		permohonan.Pendidikan = requestBody.Pendidikan
		permohonan.Pekerjaan = requestBody.Pekerjaan
		permohonan.Agama = requestBody.Agama
		permohonan.StatusPernikahan = requestBody.StatusPernikahan
		permohonan.Kewarganegaraan = requestBody.Kewarganegaraan
		permohonan.AlamatLengkap = requestBody.AlamatLengkap
		permohonan.JenisSurat = requestBody.JenisSurat
		permohonan.Keterangan = requestBody.Keterangan
		permohonan.NomorHP = requestBody.NomorHP
		permohonan.Status = model.Status(requestBody.Status)
		permohonan.Ditujukan = requestBody.Ditujukan
		if requestBody.NamaUsaha != nil {
			permohonan.NamaUsaha = sql.NullString{String: *requestBody.NamaUsaha, Valid: true}
		}
		if requestBody.JenisUsaha != nil {
			permohonan.JenisUsaha = sql.NullString{String: *requestBody.JenisUsaha, Valid: true}
		}
		if requestBody.AlamatUsaha != nil {
			permohonan.AlamatUsaha = sql.NullString{String: *requestBody.AlamatUsaha, Valid: true}
		}
		if requestBody.AlamatTujuan != nil {
			permohonan.AlamatTujuan = sql.NullString{String: *requestBody.AlamatTujuan, Valid: true}
		}
		if requestBody.AlasanPindah != nil {
			permohonan.AlasanPindah = sql.NullString{String: *requestBody.AlasanPindah, Valid: true}
		}
		if requestBody.NamaAyah != nil {
			permohonan.NamaAyah = sql.NullString{String: *requestBody.NamaAyah, Valid: true}
		}
		if requestBody.NamaIbu != nil {
			permohonan.NamaIbu = sql.NullString{String: *requestBody.NamaIbu, Valid: true}
		}
		if requestBody.TglKematian != nil {
			parsedTglKematian, err := time.Parse("2006-01-02", *requestBody.TglKematian)
			if err != nil {
				writeJSONResponse(w, http.StatusBadRequest, response{Error: "Invalid date format for tgl_kematian, expected YYYY-MM-DD"})
				return
			}
			permohonan.TglKematian = sql.NullTime{Time: parsedTglKematian, Valid: true}
		}
		if requestBody.PenyebabKematian != nil {
			permohonan.PenyebabKematian = sql.NullString{String: *requestBody.PenyebabKematian, Valid: true}
		}
		if requestBody.DokumenURL != nil {
			permohonan.DokumenURL = sql.NullString{String: *requestBody.DokumenURL, Valid: true}
		}
	} else {
		writeJSONResponse(w, http.StatusUnsupportedMediaType, response{Error: "Content-Type harus application/json"})
		return
	}

	err = c.service.UpdatePermohonanSurat(permohonan)
	if err != nil {
		log.Printf("Error updating permohonan surat: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal memperbarui permohonan surat: " + err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, response{
		Message: "Permohonan surat berhasil diperbarui",
	})
}

func (c *PermohonanSuratController) DeletePermohonanSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodDelete {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "ID tidak valid"})
		return
	}

	err = c.service.DeletePermohonanSurat(id)
	if err != nil {
		log.Printf("Error deleting permohonan surat: %v", err)
		writeJSONResponse(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, response{
		Message: "Permohonan surat berhasil dihapus",
	})
}

func (c *PermohonanSuratController) UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPatch {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "ID tidak valid"})
		return
	}

	var requestBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error parsing JSON body: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "Body JSON tidak valid"})
		return
	}

	status := model.Status(requestBody.Status)
	err = c.service.UpdateStatus(id, status)
	if err != nil {
		log.Printf("Error updating status: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal memperbarui status: " + err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, response{
		Message: "Status permohonan surat berhasil diperbarui",
	})
}

// writeJSONResponse adalah helper untuk mengirim respons JSON yang konsisten
func writeJSONResponse(w http.ResponseWriter, status int, resp response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
