package controller

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type SuratKeluarController interface {
	CreateSuratKeluar(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateSuratKeluar(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
