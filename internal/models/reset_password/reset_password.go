package reset_password

import "github.com/jinzhu/gorm"

type ResetPassword struct {
	gorm.Model
	UserID uint   `gorm:"NOT NULL"`
	Token  string `gorm:"NOT NULL;UNIQUE_INDEX"`
}
