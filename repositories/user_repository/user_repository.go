package user_repository

import (
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	GetByID(id uint) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Create(user *user.User) error
	Update(user *user.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetByID(id uint) (*user.User, error) {
	var user user.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetByEmail(email string) (*user.User, error) {
	var user user.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Create(user *user.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) Update(user *user.User) error {
	return u.db.Save(user).Error
}
