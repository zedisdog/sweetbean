{{- /* Go Template file */ -}}
package seed

import (
	"{{.Project}}/internal/module/{{.Module}}/domain/entity"
	"gorm.io/gorm"
)

func {{.SeederName}}Seed(db *gorm.DB) error {
	ent := entity.Position{
		Name:      consts.OrdinaryEmployees,
		DeptId:    0,
		IsAppoint: false,
	}

	return db.FirstOrCreate(&ent).Error
}