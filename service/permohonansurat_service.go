package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"fmt"
	"log"
	"time"
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
	if data.JenisKelamin != model.LakiLaki && data.JenisKelamin != model.Perempuan {
		return fmt.Errorf("jenis kelamin harus Laki-laki atau Perempuan")
	}
	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if !validStatuses[data.Status] {
		return fmt.Errorf("status harus salah satu dari: Diproses, Selesai, Ditolak")
	}

	// Validasi field unik berdasarkan jenis surat
	switch data.JenisSurat {
	case "Surat Keterangan Usaha":
		if !data.NamaUsaha.Valid || !data.JenisUsaha.Valid || !data.AlamatUsaha.Valid {
			return fmt.Errorf("nama_usaha, jenis_usaha, dan alamat_usaha wajib diisi untuk Surat Keterangan Usaha")
		}
	case "Surat Keterangan Pindah":
		if !data.AlamatTujuan.Valid || !data.AlasanPindah.Valid {
			return fmt.Errorf("alamat_tujuan dan alasan_pindah wajib diisi untuk Surat Keterangan Pindah")
		}
	case "Surat Keterangan Kelahiran":
		if !data.NamaAyah.Valid || !data.NamaIbu.Valid {
			return fmt.Errorf("nama_ayah dan nama_ibu wajib diisi untuk Surat Keterangan Kelahiran")
		}
	case "Surat Keterangan Kematian":
		if !data.TglKematian.Valid || !data.PenyebabKematian.Valid {
			return fmt.Errorf("tgl_kematian dan penyebab_kematian wajib diisi untuk Surat Keterangan Kematian")
		}
	}
	return nil
}

func (s *PermohonanSuratService) AddPermohonanSurat(permohonan model.PermohonanSurat) (*model.PermohonanSurat, error) {
	// Validasi input
	if err := s.validatePermohonanSurat(permohonan); err != nil {
		log.Printf("Validation error adding permohonan surat: %v", err)
		return nil, err
	}

	// Set status default jika kosong
	if permohonan.Status == "" {
		permohonan.Status = model.Diproses
	}

	// Tidak perlu mengisi CreatedAt dan UpdatedAt, karena diatur otomatis oleh database
	newPermohonan, err := s.repo.AddPermohonanSurat(permohonan)
	if err != nil {
		log.Printf("Error adding permohonan surat: %v", err)
		return nil, fmt.Errorf("gagal menambahkan permohonan surat: %v", err)
	}
	return newPermohonan, nil
}

func (s *PermohonanSuratService) GetPermohonanSurat() ([]model.PermohonanSurat, error) {
	permohonanSuratList, err := s.repo.GetPermohonanSurat()
	if err != nil {
		log.Printf("Error retrieving permohonan surat: %v", err)
		return nil, fmt.Errorf("gagal mengambil daftar permohonan surat: %v", err)
	}
	return permohonanSuratList, nil
}

func (s *PermohonanSuratService) GetPermohonanSuratByID(id int64) (*model.PermohonanSurat, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID harus lebih besar dari 0")
	}

	permohonan, err := s.repo.GetPermohonanSuratByID(id)
	if err != nil {
		log.Printf("Error retrieving permohonan surat by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil permohonan surat dengan ID %d: %v", id, err)
	}

	return permohonan, nil
}

func (s *PermohonanSuratService) UpdatePermohonanSuratByID(id int64, permohonan model.PermohonanSurat) error {
	if id <= 0 {
		return fmt.Errorf("ID harus lebih besar dari 0")
	}

	// Validasi input
	if err := s.validatePermohonanSurat(permohonan); err != nil {
		log.Printf("Validation error updating permohonan surat: %v", err)
		return err
	}

	err := s.repo.UpdatePermohonanSuratByID(id, permohonan)
	if err != nil {
		log.Printf("Error updating permohonan surat with ID %d: %v", id, err)
		return fmt.Errorf("gagal memperbarui permohonan surat dengan ID %d: %v", id, err)
	}
	return nil
}

func (s *PermohonanSuratService) UpdateStatusByID(id int64, status model.Status) error {
	if id <= 0 {
		return fmt.Errorf("ID harus lebih besar dari 0")
	}

	validStatuses := map[model.Status]bool{
		model.Diproses: true,
		model.Selesai:  true,
		model.Ditolak:  true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("status harus salah satu dari: Diproses, Selesai, Ditolak")
	}

	updatedAt := time.Now()
	err := s.repo.UpdateStatusByID(id, status, updatedAt)
	if err != nil {
		log.Printf("Error updating status permohonan surat with ID %d: %v", id, err)
		return fmt.Errorf("gagal memperbarui status permohonan surat dengan ID %d: %v", id, err)
	}
	return nil
}

func (s *PermohonanSuratService) ProcessNextPermohonan(status model.Status) error {
	permohonan, err := s.repo.GetOldestPendingPermohonan()
	if err != nil {
		log.Printf("Error retrieving oldest pending permohonan: %v", err)
		return fmt.Errorf("gagal mengambil permohonan tertua: %v", err)
	}
	if permohonan == nil {
		return fmt.Errorf("tidak ada permohonan dengan status Diproses untuk diproses")
	}

	// Validasi status
	validStatuses := map[model.Status]bool{
		model.Selesai: true,
		model.Ditolak: true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("status harus Selesai atau Ditolak untuk memproses permohonan")
	}

	// Update status permohonan tertua
	err = s.UpdateStatusByID(permohonan.ID, status)
	if err != nil {
		log.Printf("Error processing permohonan ID %d: %v", permohonan.ID, err)
		return fmt.Errorf("gagal memproses permohonan ID %d: %v", permohonan.ID, err)
	}

	log.Printf("Successfully processed permohonan ID %d with status %s", permohonan.ID, status)
	return nil
}

func (s *PermohonanSuratService) DeletePermohonanSurat(id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID harus lebih besar dari 0")
	}

	err := s.repo.DeletePermohonanSurat(id)
	if err != nil {
		log.Printf("Error deleting permohonan surat with ID %d: %v", id, err)
		return fmt.Errorf("gagal menghapus permohonan surat dengan ID %d: %v", id, err)
	}
	return nil
}
