package service

import (
	"errors"
	"log"
	"Sekertaris/dto"
	"Sekertaris/model"
	"Sekertaris/repository"
	"Sekertaris/util"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) Signup(req dto.SignupRequest) error {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("error processing request")
	}

	user := model.User{
		Email:     req.Email,
		Nama:      req.Nama,
		Handphone: req.Handphone,
		Angkatan:  req.Angkatan,
		Password:  hashedPassword,
		Role:      "admin",
	}

	return s.UserRepo.CreateUser(user)
}

func (s *UserService) Signin(req dto.LoginRequest) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	return util.GenerateJWT(user.Email, user.Role)
}