package drivers

import (
	"gorm.io/gorm"
)

func NewGorm(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(dialector, config)
}
