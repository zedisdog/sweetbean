package tools

import (
	"github.com/zedisdog/sweetbean/errx"
	"reflect"
)

//CopyFields copy fields from src to dest, note: dest mast be point
//	params:
//		src       source object
//		dest      point of dest object
//		notCopyZero  if copy zero field too
func CopyFields(src interface{}, dest interface{}, notCopyZero ...bool) error {

	// dest必须为指针
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errx.New("need dest ptr")
	} else {
		destType = destType.Elem()
	}

	// 取src的value, 如果是指针就避开指针
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	destValue := reflect.ValueOf(dest).Elem()
	for i := 0; i < destType.NumField(); i++ {
		destTypeField := destType.Field(i)
		srcValueField := srcValue.FieldByName(destTypeField.Name)
		destValueField := destValue.Field(i)
		if !srcValueField.IsValid() || //判断同名
			(len(notCopyZero) > 0 && notCopyZero[0] && srcValueField.IsZero()) || //判断是否要复制零值 不要并且是零值就跳过
			srcValueField.Kind() != destValueField.Kind() { //判断同类型
			continue
		}
		destValueField.Set(srcValueField)
	}

	return nil
}
