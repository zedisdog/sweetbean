package gormplugs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zedisdog/sweetbean/tools/snowflake"
	"gorm.io/gorm"
)

//Deprecated:
//AutoSetIDUsingSnowFlake
func AutoSetIDUsingSnowFlake(db *gorm.DB) {
	if db.Statement.Schema != nil && db.Statement.Schema.LookUpField("id") != nil {
		fmt.Printf("%+v\n", db.Statement.ReflectValue)
		if db.Statement.ReflectValue.Kind() == reflect.Slice {
			for _, v := range db.Statement.ReflectValue.Interface().([]reflect.Value) {
				_, isZero := db.Statement.Schema.LookUpField("id").ValueOf(context.Background(), v)
				if isZero {
					id, err := snowflake.NextID()
					if err != nil {
						panic(err)
					}
					db.Statement.Schema.LookUpField("id").Set(context.Background(), db.Statement.ReflectValue, id)
				}
			}
		} else {
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
}
