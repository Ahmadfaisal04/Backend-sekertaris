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

func (s *PermohonanSuratService) AddPermohonanSurat(permohonan model.PermohonanSurat) (*model.PermohonanSurat, error) {
	newPermohonan, err := s.repo.AddPermohonanSurat(permohonan)
	if err != nil {
		log.Println("Error adding permohonan surat:", err)
		return nil, err
	}
	return newPermohonan, nil
}

func (s *PermohonanSuratService) GetPermohonanSurat() ([]model.PermohonanSurat, error) {
	permohonanSuratList, err := s.repo.GetPermohonanSurat()
	if err != nil {
		log.Println("Error retrieving permohonan surat:", err)
		return nil, err
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
		return nil, err
	}

	return permohonan, nil
}

func (s *PermohonanSuratService) UpdatePermohonanSuratByID(id int64, permohonan model.PermohonanSurat) error {
	err := s.repo.UpdatePermohonanSuratByID(id, permohonan)
	if err != nil {
		log.Println("Error updating permohonan surat:", err)
		return err
	}
	return nil
}

func (s *PermohonanSuratService) UpdateStatusByID(id int64, status model.Status) error {
	
	updatedAt := time.Now()
	err := s.repo.UpdateStatusByID(id, status, updatedAt)
	if err != nil {
		log.Println("Error updating status permohonan surat:", err)
		return err
	}
	return nil
}

func (s *PermohonanSuratService) DeletePermohonanSurat(id int64) error {
	err := s.repo.DeletePermohonanSurat(id)
	if err != nil {
		log.Println("Error deleting permohonan surat:", err)
		return err
	}
	return nil
}
