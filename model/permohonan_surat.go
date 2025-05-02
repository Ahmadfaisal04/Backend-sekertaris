package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

// JenisKelamin defines the possible values for gender
type JenisKelamin string

const (
	LakiLaki  JenisKelamin = "Laki-laki"
	Perempuan JenisKelamin = "Perempuan"
)

// Status defines the possible values for permohonan status
type Status string

const (
	Diproses Status = "Diproses"
	Selesai  Status = "Selesai"
	Ditolak  Status = "Ditolak"
)

// PermohonanSurat represents the structure of a permohonan surat record
type PermohonanSurat struct {
	ID               int64          `json:"id"`
	NIK              string         `json:"nik"`
	NamaLengkap      string         `json:"nama_lengkap"`
	TempatLahir      string         `json:"tempat_lahir"`
	TanggalLahir     time.Time      `json:"tanggal_lahir"`
	JenisKelamin     JenisKelamin   `json:"jenis_kelamin"`
	Pendidikan       string         `json:"pendidikan"`
	Pekerjaan        string         `json:"pekerjaan"`
	Agama            string         `json:"agama"`
	StatusPernikahan string         `json:"status_pernikahan"`
	Kewarganegaraan  string         `json:"kewarganegaraan"`
	AlamatLengkap    string         `json:"alamat_lengkap"`
	JenisSurat       string         `json:"jenis_surat"`
	Keterangan       string         `json:"keterangan"`
	NomorHP          string         `json:"nomor_hp"`
	DokumenURL       sql.NullString `json:"dokumen_url"`
	NamaUsaha        sql.NullString `json:"nama_usaha"`
	JenisUsaha       sql.NullString `json:"jenis_usaha"`
	AlamatUsaha      sql.NullString `json:"alamat_usaha"`
	AlamatTujuan     sql.NullString `json:"alamat_tujuan"`
	AlasanPindah     sql.NullString `json:"alasan_pindah"`
	NamaAyah         sql.NullString `json:"nama_ayah"`
	NamaIbu          sql.NullString `json:"nama_ibu"`
	TglKematian      sql.NullTime   `json:"tgl_kematian"`
	PenyebabKematian sql.NullString `json:"penyebab_kematian"`
	Ditujukan        sql.NullString `json:"ditujukan"`
	Status           Status         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// MarshalJSON customizes JSON serialization for PermohonanSurat.
func (p PermohonanSurat) MarshalJSON() ([]byte, error) {
	type Alias PermohonanSurat
	return json.Marshal(&struct {
		*Alias
		DokumenURL       *string    `json:"dokumen_url"`
		NamaUsaha        *string    `json:"nama_usaha"`
		JenisUsaha       *string    `json:"jenis_usaha"`
		AlamatUsaha      *string    `json:"alamat_usaha"`
		AlamatTujuan     *string    `json:"alamat_tujuan"`
		AlasanPindah     *string    `json:"alasan_pindah"`
		NamaAyah         *string    `json:"nama_ayah"`
		NamaIbu          *string    `json:"nama_ibu"`
		TglKematian      *time.Time `json:"tgl_kematian"`
		PenyebabKematian *string    `json:"penyebab_kematian"`
	}{
		Alias:            (*Alias)(&p),
		DokumenURL:       nullStringToPtr(p.DokumenURL),
		NamaUsaha:        nullStringToPtr(p.NamaUsaha),
		JenisUsaha:       nullStringToPtr(p.JenisUsaha),
		AlamatUsaha:      nullStringToPtr(p.AlamatUsaha),
		AlamatTujuan:     nullStringToPtr(p.AlamatTujuan),
		AlasanPindah:     nullStringToPtr(p.AlasanPindah),
		NamaAyah:         nullStringToPtr(p.NamaAyah),
		NamaIbu:          nullStringToPtr(p.NamaIbu),
		TglKematian:      nullTimeToPtr(p.TglKematian),
		PenyebabKematian: nullStringToPtr(p.PenyebabKematian),
	})
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
