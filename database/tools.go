package database

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zedisdog/sweetbean/errx"
	"gorm.io/gorm"
)

type Condition = []interface{}
type Conditions = []Condition

func ParseConditionGorm(q *gorm.DB, conditions Conditions) (query *gorm.DB, err error) {
	query = q
	for _, condition := range conditions {
		switch len(condition) {
		case 0, 1:
			err = errx.New("condition is invalid", 1)
			return
		case 2:
			v := reflect.ValueOf(condition[1])
			if v.Kind() == reflect.Slice {
				query = query.Where(fmt.Sprintf("%s IN ?", condition[0]), condition[1])
			} else {
				if s, ok := condition[1].(string); ok && strings.HasPrefix(strings.ToLower(s), "is") {
					query = query.Where(fmt.Sprintf("%s %s", condition[0], condition[1]))
				} else {
					query = query.Where(fmt.Sprintf("%s = ?", condition[0]), condition[1])
				}
			}
		case 3:
			query = query.Where(fmt.Sprintf("%s %s ?", condition[0], condition[1]), condition[2])
		case 4:
			switch strings.Replace(strings.ToUpper(condition[1].(string)), " ", "", -1) {
			case "BETWEEN", "NOTBETWEEN":
				query = query.Where(fmt.Sprintf("%s %s ? AND ?", condition[0], condition[1]), condition[2], condition[3])
			default:
				err = errx.New("condition is invalid", 1)
				return
			}
		default:
			err = errx.New("condition is invalid", 1)
			return
		}
	}
	return
}
