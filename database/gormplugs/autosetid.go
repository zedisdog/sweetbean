package gormplugs

import (
	"context"

	"github.com/zedisdog/sweetbean/tools/snowflake"
	"gorm.io/gorm"
)

func AutoSetIDUsingSnowFlake(db *gorm.DB) {
	if db.Statement.Schema != nil && db.Statement.Schema.LookUpField("id") != nil {
		_, isZero := db.Statement.Schema.LookUpField("id").ValueOf(context.Background(), db.Statement.ReflectValue)
		if isZero {
			id, err := snowflake.NextID()
			if err != nil {
				panic(err)
			}
			db.Statement.Schema.LookUpField("id").Set(context.Background(), db.Statement.ReflectValue, id)
		}
	}
}
