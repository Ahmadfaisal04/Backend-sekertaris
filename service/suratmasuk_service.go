package service

import (
	"Sekertaris/model"
	"Sekertaris/repository"
	"database/sql"
	"errors"
	"fmt"
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

func (s *SuratMasukService) GetSuratById(id int) ([]model.SuratMasuk, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID harus lebih besar dari 0")
	}

	// Panggil repository untuk mengambil data surat masuk
	surat, err := s.repo.GetSuratById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("surat dengan ID %d tidak ditemukan", id)
		}
		log.Printf("Error retrieving surat masuk by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil surat masuk: %v", err)
	}

	// Bungkus data dalam slice (array)
	return []model.SuratMasuk{*surat}, nil
}


func (s *SuratMasukService) UpdateSuratMasukByID(id int, surat model.SuratMasuk) error {
	err := s.repo.UpdateSuratMasukByID(id, surat)
	if err != nil {
		log.Println("Error updating surat masuk:", err)
		return err
	}
	return nil
}

func (s *SuratMasukService) DeleteSuratMasuk(id int) error {
	err := s.repo.DeleteSuratMasuk(id)
	if err != nil {
		log.Println("Error deleting surat masuk:", err)
		return err
	}
	return nil
}

