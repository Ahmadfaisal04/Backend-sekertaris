package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// RequestBody adalah struct untuk parsing JSON
type RequestBody struct {
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
	Ditujukan        *string `json:"ditujukan"`
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

// stringToNullString mengonversi *string ke sql.NullString
func stringToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// stringToNullTime mengonversi *string (tanggal) ke sql.NullTime
func stringToNullTime(s *string) sql.NullTime {
	if s == nil || *s == "" {
		return sql.NullTime{Time: time.Time{}, Valid: false}
	}
	parsedTime, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return sql.NullTime{Time: time.Time{}, Valid: false}
	}
	return sql.NullTime{Time: parsedTime, Valid: true}
}

// AddPermohonanSurat menangani POST /permohonan-surat
func (c *PermohonanSuratController) AddPermohonanSurat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	// Validasi Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONResponse(w, http.StatusUnsupportedMediaType, response{Error: "Content-Type harus application/json"})
		return
	}

	var requestBody RequestBody
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

	// Parsing tanggal lahir
	parsedTanggalLahir, err := time.Parse("2006-01-02", requestBody.TanggalLahir)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
		return
	}





	// Mapping ke PermohonanSurat
	permohonan := model.PermohonanSurat{
		NIK:              requestBody.NIK,
		NamaLengkap:      requestBody.NamaLengkap,
		TempatLahir:      requestBody.TempatLahir,
		TanggalLahir:     parsedTanggalLahir,
		JenisKelamin:     model.JenisKelamin(requestBody.JenisKelamin),
		Pendidikan:       requestBody.Pendidikan,
		Pekerjaan:        requestBody.Pekerjaan,
		Agama:            requestBody.Agama,
		StatusPernikahan: requestBody.StatusPernikahan,
		Kewarganegaraan:  requestBody.Kewarganegaraan,
		AlamatLengkap:    requestBody.AlamatLengkap,
		JenisSurat:       requestBody.JenisSurat,
		Keterangan:       requestBody.Keterangan,
		NomorHP:          requestBody.NomorHP,
		DokumenURL:       stringToNullString(requestBody.DokumenURL),
		NamaUsaha:        stringToNullString(requestBody.NamaUsaha),
		JenisUsaha:       stringToNullString(requestBody.JenisUsaha),
		AlamatUsaha:      stringToNullString(requestBody.AlamatUsaha),
		AlamatTujuan:     stringToNullString(requestBody.AlamatTujuan),
		AlasanPindah:     stringToNullString(requestBody.AlasanPindah),
		NamaAyah:         stringToNullString(requestBody.NamaAyah),
		NamaIbu:          stringToNullString(requestBody.NamaIbu),
		TglKematian:      stringToNullTime(requestBody.TglKematian),
		PenyebabKematian: stringToNullString(requestBody.PenyebabKematian),
		Ditujukan:        stringToNullString(requestBody.Ditujukan),
		Status:           model.Status(requestBody.Status),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Simpan ke database melalui service
	jsonData, err := c.service.AddPermohonanSuratJSON(permohonan)
	if err != nil {
		log.Printf("Error adding permohonan surat: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, response{Error: "Gagal menambahkan permohonan surat: " + err.Error()})
		return
	}

	// Kirim respons
	writeJSONResponse(w, http.StatusCreated, response{
		Message: "Permohonan surat berhasil ditambahkan",
		Data:    json.RawMessage(jsonData),
	})
}

// GetPermohonanSurat menangani GET /permohonan-surat
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

// GetPermohonanSuratByID menangani GET /permohonan-surat/:id
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

// UpdatePermohonanSuratByID menangani PUT /permohonan-surat/:id
func (c *PermohonanSuratController) UpdatePermohonanSuratByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPut {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	// Validasi Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONResponse(w, http.StatusUnsupportedMediaType, response{Error: "Content-Type harus application/json"})
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "ID tidak valid"})
		return
	}

	var requestBody RequestBody
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

	// Parsing tanggal lahir
	parsedTanggalLahir, err := time.Parse("2006-01-02", requestBody.TanggalLahir)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, response{Error: "Format tanggal_lahir tidak valid, gunakan YYYY-MM-DD"})
		return
	}



	// Mapping ke PermohonanSurat
	permohonan := model.PermohonanSurat{
		ID:               id,
		NIK:              requestBody.NIK,
		NamaLengkap:      requestBody.NamaLengkap,
		TempatLahir:      requestBody.TempatLahir,
		TanggalLahir:     parsedTanggalLahir,
		JenisKelamin:     model.JenisKelamin(requestBody.JenisKelamin),
		Pendidikan:       requestBody.Pendidikan,
		Pekerjaan:        requestBody.Pekerjaan,
		Agama:            requestBody.Agama,
		StatusPernikahan: requestBody.StatusPernikahan,
		Kewarganegaraan:  requestBody.Kewarganegaraan,
		AlamatLengkap:    requestBody.AlamatLengkap,
		JenisSurat:       requestBody.JenisSurat,
		Keterangan:       requestBody.Keterangan,
		NomorHP:          requestBody.NomorHP,
		DokumenURL:       stringToNullString(requestBody.DokumenURL),
		NamaUsaha:        stringToNullString(requestBody.NamaUsaha),
		JenisUsaha:       stringToNullString(requestBody.JenisUsaha),
		AlamatUsaha:      stringToNullString(requestBody.AlamatUsaha),
		AlamatTujuan:     stringToNullString(requestBody.AlamatTujuan),
		AlasanPindah:     stringToNullString(requestBody.AlasanPindah),
		NamaAyah:         stringToNullString(requestBody.NamaAyah),
		NamaIbu:          stringToNullString(requestBody.NamaIbu),
		TglKematian:      stringToNullTime(requestBody.TglKematian),
		PenyebabKematian: stringToNullString(requestBody.PenyebabKematian),
		Ditujukan:        stringToNullString(requestBody.Ditujukan),
		Status:           model.Status(requestBody.Status),
		UpdatedAt:        time.Now(),
	}

	// Update melalui service
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

// DeletePermohonanSurat menangani DELETE /permohonan-surat/:id
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

// UpdateStatus menangani PATCH /permohonan-surat/:id/status
func (c *PermohonanSuratController) UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPatch {
		writeJSONResponse(w, http.StatusMethodNotAllowed, response{Error: "Method not allowed"})
		return
	}

	// Validasi Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONResponse(w, http.StatusUnsupportedMediaType, response{Error: "Content-Type harus application/json"})
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
