package controller

import (
	"Sekertaris/service"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddSuratMasuk(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		service.AddSuratMasuk(w, r, db)
		
	}
}
