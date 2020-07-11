package password_reset

import (
	pwd "github.com/alonelegion/go_graphql_api/internal/models/reset_password"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	GetOneByToken(token string) (*pwd.ResetPassword, error)
	Create(pr *pwd.ResetPassword) error
	Delete(id uint) error
}

type pwdRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) Repository {
	return &pwdRepository{
		db: db,
	}
}

func (p *pwdRepository) GetOneByToken(token string) (*pwd.ResetPassword, error) {
	var pass pwd.ResetPassword
	if err := p.db.Where("token = ?", token).First(&pass).Error; err != nil {
		return nil, err
	}
	return &pass, nil
}

func (p *pwdRepository) Create(pr *pwd.ResetPassword) error {
	return p.db.Create(pr).Error
}

func (p *pwdRepository) Delete(id uint) error {
	pass := pwd.ResetPassword{
		Model: gorm.Model{ID: id},
	}
	return p.db.Delete(&pass).Error
}
