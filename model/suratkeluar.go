package model

type SuratKeluar struct {
	ID        int    `json:"id"`
	Nomor     string `json:"nomor"`
	Tanggal   string `json:"tanggal"`
	Perihal   string `json:"perihal"`
	Ditujukan string `json:"ditujukan"`
	Title     string `json:"title"`
	File      string `json:"file"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
