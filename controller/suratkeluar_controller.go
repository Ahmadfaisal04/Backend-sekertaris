package controller

import (
	"Sekertaris/service"
	"database/sql"
	"net/http"
	"github.com/julienschmidt/httprouter"
	
)

func AddSuratKeluar(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.AddSuratKeluar(w, r, db)

	}
}

// Handler untuk mendapatkan semua surat keluar (GET)
func UpdateSuratKeluar(db *sql.DB) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        service.UpdateSuratKeluar(w, r, db, ps.ByName("nomor")) // Hanya meneruskan (w, r, ps)
    }
}

func GetSuratKeluarByid(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
        service.GetSuratKeluarByid(w, r, id, db) // Hanya meneruskan (w, r, ps)
    }
}
