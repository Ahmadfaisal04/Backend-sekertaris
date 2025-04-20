package model

import (
	"database/sql"
	"time"
)

type JenisKelamin string
type Status string

const (
	LakiLaki  JenisKelamin = "Laki-laki"
	Perempuan JenisKelamin = "Perempuan"

	Pending  Status = "Pending"
	Diproses Status = "Diproses"
	Selesai  Status = "Selesai"
)

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
	Status           Status         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
