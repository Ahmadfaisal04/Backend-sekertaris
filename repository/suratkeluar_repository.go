package repository

import (
	"Sekertaris/model"
	"context"
	"database/sql"
)

type SuratKeluarRepository interface {
    CreateSuratKeluar(ctx context.Context, tx *sql.Tx, surat model.SuratKeluar) (model.SuratKeluar, error)
    UpdateSuratKeluar(ctx context.Context, tx *sql.Tx, surat model.SuratKeluar) (model.SuratKeluar, error)
}
