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
