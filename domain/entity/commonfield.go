package entity

import (
	"github.com/zedisdog/sweetbean/tools/snowflake"
	"gorm.io/gorm"
)

type CommonField struct {
	ID        uint64 `json:"id,string" gorm:"primary"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a CommonField) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == 0 {
		a.ID, err = snowflake.NextID()
	}
	return
}
