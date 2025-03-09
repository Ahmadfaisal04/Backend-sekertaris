package model

type SuratMasuk struct {
	Id      int    `json:"id"`
	Nomor   string `json:"nomor"`
	Tanggal string `json:"tanggal"`
	Perihal string `json:"perihal"`
	Asal    string `json:"asal"`
	Title   string `json:"title"`
	File    string `json:"file"`
}
