package models

import (
	"database/sql/driver"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID             int64   `json:"id" db:"id" validate:"omitempty"`
	Email          string  `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password       string  `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	Name           string  `json:"name" db:"name" validate:"required,gte=2"`
	Surname        string  `json:"surname" db:"surname" validate:"required,gte=2,lte=60"`
	Patronymic     *string `json:"patronymic" db:"patronymic" validate:"omitempty"`
	Address        string  `json:"address" db:"address" validate:"required,lte=100"`
	Admin          bool    `json:"admin,omitempty" db:"admin" validate:"omitempty"`
	PassportNumber int     `json:"passport_number,omitempty" db:"passport_number" validate:"required"`
	PassportSeries int     `json:"passport_series,omitempty" db:"passport_series" validate:"required"`
}

type UserWithToken struct {
	*User
	Token string `json:"token"`
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// PrepareCreate prepare user for registration
func (user *User) PrepareCreate() error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Password = strings.TrimSpace(user.Password)

	if err := user.HashPassword(); err != nil {
		return err
	}
	return nil
}

// PrepareUpdate prepare user for update
func (user *User) PrepareUpdate() error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Password = strings.TrimSpace(user.Password)

	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}
	return nil
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) Sanitize() {
	user.PassportNumber = 0
	user.PassportSeries = 0
}

func (user *User) Columns() []string {
	return []string{"id", "email", "password", "name", "surname", "patronymic",
		"address", "admin", "passport_number", "passport_series"}
}

func (user *User) Rows() []driver.Value {
	return []driver.Value{user.ID, user.Email, user.Password, user.Name, user.Surname,
		user.Patronymic, user.Address, user.Admin, user.PassportNumber, user.PassportSeries}
}
