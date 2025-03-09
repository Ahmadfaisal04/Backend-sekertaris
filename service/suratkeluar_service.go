package service

import (
	"Sekertaris/model"
	"database/sql"
	"encoding/json"
	"net/http"
)

func AddSuratKeluar(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var surat model.SuratKeluar
	surat.Nomor = r.FormValue("nomor")
	surat.Tanggal = r.FormValue("tanggal")
	surat.Ditujukan = r.FormValue("tujuan")
	surat.Perihal = r.FormValue("perihal")
	surat.Title_file = r.FormValue("title_file")
	surat.File = r.FormValue("file")

	if surat.Nomor == "" || surat.Tanggal == "" || surat.Ditujukan == "" {
		http.Error(w, "Data tidak boleh kosong", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO surat_keluar (nomor, tanggal, tujuan, perihal, lampiran) VALUES (?, ?, ?, ?, ?)",
		surat.Nomor, surat.Tanggal, surat.Ditujukan, surat.Perihal, surat.Title_file, surat.File)

	if err != nil {
		http.Error(w, "Gagal menyimpan surat keluar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Surat keluar berhasil disimpan"})
}
