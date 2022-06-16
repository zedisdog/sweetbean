package domain

import (
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

func NewEntity(db *gorm.DB, snowflake *sonyflake.Sonyflake) *Entity {
	return &Entity{
		DB:        db,
		snowflake: snowflake,
	}
}

type Entity struct {
	*gorm.DB
	snowflake *sonyflake.Sonyflake
	ID        uint64 `json:"id,string" gorm:"primaryKey"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == 0 {
		a.ID, err = a.snowflake.NextID()
	}
	return
}
