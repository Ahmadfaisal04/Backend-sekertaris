package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"encoding/json"
	"fmt"
	"strings"
)

type PermohonanSuratService struct {
	repo *repository.PermohonanSuratRepository
}

func NewPermohonanSuratService(repo *repository.PermohonanSuratRepository) *PermohonanSuratService {
	return &PermohonanSuratService{repo: repo}
}

func (s *PermohonanSuratService) validatePermohonanSurat(data model.PermohonanSurat) error {
	if data.NIK == "" || len(data.NIK) != 16 {
		return fmt.Errorf("NIK harus 16 digit")
	}
	if data.NamaLengkap == "" {
		return fmt.Errorf("nama lengkap wajib diisi")
	}
	if data.TanggalLahir.IsZero() {
		return fmt.Errorf("tanggal lahir wajib diisi")
	}

	if data.JenisSurat == "" {
		return fmt.Errorf("jenis surat wajib diisi")
	}
	// Validasi khusus berdasarkan jenis surat
	switch strings.ToLower(data.JenisSurat) {
	case "surat keterangan domisili":
		if data.AlamatLengkap == "" {
			return fmt.Errorf("alamat_lengkap wajib diisi untuk Surat Keterangan Domisili")
		}
	case "surat keterangan usaha":
		if !data.NamaUsaha.Valid || !data.JenisUsaha.Valid || !data.AlamatUsaha.Valid {
			return fmt.Errorf("nama_usaha, jenis_usaha, dan alamat_usaha wajib diisi untuk Surat Keterangan Usaha")
		}
	case "surat keterangan pindah":
		if !data.AlamatTujuan.Valid || !data.AlasanPindah.Valid {
			return fmt.Errorf("alamat_tujuan dan alasan_pindah wajib diisi untuk Surat Keterangan Pindah")
		}
	case "surat keterangan kematian":
		if !data.TglKematian.Valid || !data.PenyebabKematian.Valid {
			return fmt.Errorf("tgl_kematian dan penyebab_kematian wajib diisi untuk Surat Keterangan Kematian")
		}
	}
	// Validasi status
	// validStatuses := map[model.Status]bool{
	// 	model.Diproses: true,
	// 	model.Selesai:  true,
	// 	model.Ditolak:  true,
	// }
	// if data.Status != "" && !validStatuses[data.Status] {
	// 	return fmt.Errorf("status harus 'Diproses', 'Selesai', atau 'Ditolak'")
	// }
	return nil
}

func (s *PermohonanSuratService) AddPermohonanSuratJSON(data model.PermohonanSurat) ([]byte, error) {
	if err := s.validatePermohonanSurat(data); err != nil {
		return nil, err
	}
	if data.Status == "" {
		data.Status = model.Diproses
	}
	id, err := s.repo.AddPermohonanSurat(data)
	if err != nil {
		return nil, err
	}
	data.ID = id
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	return jsonData, nil
}

func (s *PermohonanSuratService) GetPermohonanSurat() ([]model.PermohonanSurat, error) {
	return s.repo.GetPermohonanSurat()
}

func (s *PermohonanSuratService) GetPermohonanSuratByID(id int64) (model.PermohonanSurat, error) {
	return s.repo.GetPermohonanSuratByID(id)
}

func (s *PermohonanSuratService) UpdatePermohonanSurat(data model.PermohonanSurat) error {
	if err := s.validatePermohonanSurat(data); err != nil {
		return err
	}
	return s.repo.UpdatePermohonanSurat(data)
}

func (s *PermohonanSuratService) DeletePermohonanSurat(id int64) error {
	return s.repo.DeletePermohonanSurat(id)
}

func (s *PermohonanSuratService) UpdateStatus(id int64, status model.Status) error {
	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("status harus 'Diproses', 'Selesai', atau 'Ditolak'")
	}
	return s.repo.UpdateStatus(id, status)
}
