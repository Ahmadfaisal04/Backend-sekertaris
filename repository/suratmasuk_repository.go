package repository

import (
	"database/sql"
	"Sekertaris/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user model.User) error {
	_, err := r.DB.Exec("INSERT INTO user (email, nama, handphone, angkatan, password, role) VALUES (?, ?, ?, ?, ?, ?)",
		user.Email, user.Nama, user.Handphone, user.Angkatan, user.Password, user.Role)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT email, password, role FROM user WHERE email = ?", email).Scan(&user.Email, &user.Password, &user.Role)
	return &user, err
}