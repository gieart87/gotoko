package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	RoleID        string `gorm:"size:36;index"`
	Role          Role
	Addresses     []Address
	FirstName     string `gorm:"size:100;not null"`
	LastName      string `gorm:"size:100;not null"`
	Email         string `gorm:"size:100;not null;uniqueIndex"`
	Password      string `gorm:"size:255;not null"`
	RememberToken string `gorm:"size:255;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	var user User

	err = db.Debug().Model(User{}).Where("LOWER(email) = ?", strings.ToLower(email)).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) FindByID(db *gorm.DB, userID string) (*User, error) {
	var err error
	var user User

	err = db.Debug().Preload("Role").Model(User{}).Where("id = ?", userID).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) CreateUser(db *gorm.DB, param *User) (*User, error) {
	user := &User{
		ID:        param.ID,
		FirstName: param.FirstName,
		LastName:  param.LastName,
		Email:     param.Email,
		Password:  param.Password,
	}

	err := db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
