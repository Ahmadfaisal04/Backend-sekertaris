package controller

import (
	"Sekertaris/model"
	"Sekertaris/service"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SuratKeluarController struct {
	service *service.SuratKeluarService
}

func NewSuratKeluarController(service *service.SuratKeluarService) *SuratKeluarController {
	return &SuratKeluarController{service: service}
}

func (c *SuratKeluarController) AddSuratKeluar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(10 << 20)

	surat := &model.SuratKeluar{
		Nomor:     r.FormValue("nomor"),
		Tanggal:   r.FormValue("tanggal"),
		Perihal:   r.FormValue("perihal"),
		Ditujukan: r.FormValue("ditujukan"),
		Title:     r.FormValue("title"),
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File wajib diupload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = c.service.AddSuratKeluar(surat, file, fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Surat keluar berhasil ditambahkan"})
}

func (c *SuratKeluarController) GetAllSuratKeluar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	surat, err := c.service.GetAllSuratKeluar()
	if err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surat)
}

func (c *SuratKeluarController) GetSuratKeluarById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	surat, err := c.service.GetSuratKeluarById(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surat)
}

func (c *SuratKeluarController) UpdateSuratKeluarByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	r.ParseMultipartForm(10 << 20)

	surat := &model.SuratKeluar{
		Nomor:     r.FormValue("nomor"),
		Tanggal:   r.FormValue("tanggal"),
		Perihal:   r.FormValue("perihal"),
		Ditujukan: r.FormValue("ditujukan"),
		Title:     r.FormValue("title"),
		File:      r.FormValue("existing_file"),
	}

	var file io.Reader
	var fileName string

	fileUpload, fileHeader, err := r.FormFile("file")
	if err == nil {
		file = fileUpload
		fileName = fileHeader.Filename
		defer fileUpload.Close()
	}

	err = c.service.UpdateSuratKeluarByID(id, surat, file, fileName)
	if err != nil {
		http.Error(w, "Gagal update data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Surat keluar berhasil diupdate"})
}

func (c *SuratKeluarController) DeleteSuratKeluar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	err := c.service.DeleteSuratKeluar(id)
	if err != nil {
		http.Error(w, "Gagal menghapus data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Surat keluar berhasil dihapus"})
}
