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

func (c *PermohonanSuratController) AddPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
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
	// Kolom baru
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
			http.Error(w, `{"error": "Invalid date format for tgl_kematian, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		permohonan.TglKematian = sql.NullTime{Time: parsedTglKematian, Valid: true}
	}
	if penyebabKematian := r.FormValue("penyebab_kematian"); penyebabKematian != "" {
		permohonan.PenyebabKematian = sql.NullString{String: penyebabKematian, Valid: true}
	}

	// Validasi kolom wajib
	if permohonan.NIK == "" || permohonan.NamaLengkap == "" || tanggalLahirStr == "" {
		http.Error(w, `{"error": "NIK, Nama Lengkap, dan Tanggal Lahir wajib diisi"}`, http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", tanggalLahirStr)
	if err != nil {
		http.Error(w, `{"error": "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}
	permohonan.TanggalLahir = parsedDate

	// Validasi JenisKelamin
	if permohonan.JenisKelamin != model.LakiLaki && permohonan.JenisKelamin != model.Perempuan {
		http.Error(w, `{"error": "Jenis Kelamin harus 'Laki-laki' atau 'Perempuan'"}`, http.StatusBadRequest)
		return
	}

	// Validasi Status
	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if permohonan.Status != "" && !validStatuses[permohonan.Status] {
		http.Error(w, `{"error": "Status harus 'Diproses', 'Selesai', atau 'Ditolak'"}`, http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()

		staticPath := "./static/suratkeluar/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v", err)
			http.Error(w, `{"error": "Gagal membuat direktori"}`, http.StatusInternalServerError)
			return
		}

		filePath := staticPath + header.Filename
		outFile, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			http.Error(w, `{"error": "Gagal membuat file"}`, http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Printf("Error saving file: %v", err)
			http.Error(w, `{"error": "Gagal menyimpan file"}`, http.StatusInternalServerError)
			return
		}

		permohonan.DokumenURL = sql.NullString{String: filePath, Valid: true}
	}

	newPermohonan, err := c.service.AddPermohonanSurat(permohonan)
	if err != nil {
		log.Printf("Error adding permohonan surat: %v", err)
		http.Error(w, `{"error": "Gagal menambahkan permohonan surat: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPermohonan)
}

func (c *PermohonanSuratController) GetPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	permohonanSuratList, err := c.service.GetPermohonanSurat()
	if err != nil {
		log.Printf("Error retrieving permohonan surat: %v", err)
		http.Error(w, `{"error": "Gagal mengambil data permohonan surat"}`, http.StatusInternalServerError)
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
		log.Printf("Invalid ID: %v", err)
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	permohonan, err := c.service.GetPermohonanSuratByID(id)
	if err != nil {
		log.Printf("Error retrieving permohonan surat by ID %d: %v", id, err)
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Permohonan surat tidak ditemukan"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Gagal mengambil permohonan surat"}`, http.StatusInternalServerError)
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
		log.Printf("Invalid ID: %v", err)
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		http.Error(w, `{"error": "Gagal memproses data form"}`, http.StatusBadRequest)
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
	// Kolom baru
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
			http.Error(w, `{"error": "Invalid date format for tgl_kematian, expected YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		permohonan.TglKematian = sql.NullTime{Time: parsedTglKematian, Valid: true}
	}
	if penyebabKematian := r.FormValue("penyebab_kematian"); penyebabKematian != "" {
		permohonan.PenyebabKematian = sql.NullString{String: penyebabKematian, Valid: true}
	}

	if tanggalLahirStr != "" {
		parsedDate, err := time.Parse("2006-01-02", tanggalLahirStr)
		if err != nil {
			http.Error(w, `{"error": "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		permohonan.TanggalLahir = parsedDate
	}

	// Validasi JenisKelamin
	if permohonan.JenisKelamin != "" && permohonan.JenisKelamin != model.LakiLaki && permohonan.JenisKelamin != model.Perempuan {
		http.Error(w, `{"error": "Jenis Kelamin harus 'Laki-laki' atau 'Perempuan'"}`, http.StatusBadRequest)
		return
	}

	// Validasi Status
	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if permohonan.Status != "" && !validStatuses[permohonan.Status] {
		http.Error(w, `{"error": "Status harus 'Diproses', 'Selesai', atau 'Ditolak'"}`, http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, handler, err := r.FormFile("dokumen")
	if err == nil {
		defer file.Close()

		staticPath := "./static/permohonansurat/"
		err = os.MkdirAll(staticPath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v", err)
			http.Error(w, `{"error": "Gagal membuat direktori"}`, http.StatusInternalServerError)
			return
		}

		filePath := staticPath + handler.Filename
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			http.Error(w, `{"error": "Gagal membuat file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Printf("Error copying file: %v", err)
			http.Error(w, `{"error": "Gagal menyimpan file"}`, http.StatusInternalServerError)
			return
		}

		permohonan.DokumenURL = sql.NullString{String: filePath, Valid: true}
	} else {
		existingDokumenURL := r.FormValue("existing_dokumen_url")
		if existingDokumenURL != "" {
			permohonan.DokumenURL = sql.NullString{String: existingDokumenURL, Valid: true}
		}
	}

	err = c.service.UpdatePermohonanSuratByID(id, permohonan)
	if err != nil {
		log.Printf("Error updating permohonan surat: %v", err)
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Permohonan surat tidak ditemukan"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Gagal memperbarui permohonan surat: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Permohonan surat berhasil diperbarui"})
}

func (c *PermohonanSuratController) UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `{"error": "Metode tidak diizinkan"}`, http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Status model.Status `json:"status"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("Error parsing JSON body: %v", err)
		http.Error(w, `{"error": "Body JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if !validStatuses[requestBody.Status] {
		http.Error(w, `{"error": "Status harus 'Diproses', 'Selesai', atau 'Ditolak'"}`, http.StatusBadRequest)
		return
	}

	err = c.service.UpdateStatusByID(id, requestBody.Status)
	if err != nil {
		log.Printf("Error updating status for ID %d: %v", id, err)
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Permohonan surat tidak ditemukan"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Gagal memperbarui status"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status berhasil diperbarui",
	})
}

func (c *PermohonanSuratController) ProcessNextPermohonan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Metode tidak diizinkan"}`, http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Status model.Status `json:"status"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("Error parsing JSON body: %v", err)
		http.Error(w, `{"error": "Body JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	validStatuses := map[model.Status]bool{
		model.Selesai: true,
		model.Ditolak: true,
	}
	if !validStatuses[requestBody.Status] {
		http.Error(w, `{"error": "Status harus 'Selesai' atau 'Ditolak'"}`, http.StatusBadRequest)
		return
	}

	err := c.service.ProcessNextPermohonan(requestBody.Status)
	if err != nil {
		log.Printf("Error processing next permohonan: %v", err)
		if strings.Contains(err.Error(), "tidak ada permohonan") {
			http.Error(w, `{"error": "Tidak ada permohonan untuk diproses"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Gagal memproses permohonan: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Permohonan berhasil diproses",
	})
}

func (c *PermohonanSuratController) DeletePermohonanSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	err = c.service.DeletePermohonanSurat(id)
	if err != nil {
		log.Printf("Error deleting permohonan surat with ID %d: %v", id, err)
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, `{"error": "Permohonan surat tidak ditemukan"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Gagal menghapus permohonan surat"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Permohonan surat berhasil dihapus",
	})
}
