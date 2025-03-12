package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"log"
	"time"
)

type SuratMasukService struct {
	repo *repository.SuratMasukRepository
}

func NewSuratMasukService(repo *repository.SuratMasukRepository) *SuratMasukService {
	return &SuratMasukService{repo: repo}
}

func (s *SuratMasukService) AddSuratMasuk(surat model.SuratMasuk, parsedDate time.Time) (*model.SuratMasuk, error) {
	newSurat, err := s.repo.AddSuratMasuk(surat, parsedDate)
	if err != nil {
		log.Println("Error adding surat masuk:", err)
		return nil, err
	}
	return newSurat, nil
}

func (s *SuratMasukService) GetSuratMasuk() ([]model.SuratMasuk, error) {
	suratMasukList, err := s.repo.GetSuratMasuk()
	if err != nil {
		log.Println("Error retrieving surat masuk:", err)
		return nil, err
	}
	return suratMasukList, nil
}

func (s *SuratMasukService) GetSuratById(id int) (*model.SuratMasuk, error) {
	surat, err := s.repo.GetSuratById(id)
	if err != nil {
		log.Println("Error retrieving surat masuk by ID:", err)
		return nil, err
	}
	return surat, nil
}

func (s *SuratMasukService) GetCountSuratMasuk() (int, error) {
	count, err := s.repo.GetCountSuratMasuk()
	if err != nil {
		log.Println("Error retrieving count surat masuk:", err)
		return 0, err
	}
	return count, nil
}

func (s *SuratMasukService) UpdateSuratMasukByID(id int, surat model.SuratMasuk) error {
	err := s.repo.UpdateSuratMasukByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		return err
	}
	return nil
}

func (s *SuratMasukService) DeleteSuratMasuk(nomor, perihal string) error {
	err := s.repo.DeleteSuratMasuk(nomor, perihal)
	if err != nil {
		log.Println("Error deleting surat masuk:", err)
		return err
	}
	return nil
}
