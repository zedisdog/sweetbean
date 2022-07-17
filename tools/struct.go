package tools

import (
	"errors"
	"reflect"
	"strings"
	"unsafe"

	"github.com/zedisdog/sweetbean/errx"
)

//CopyFields copy fields from src to dest, note: dest mast be point
//	params:
//		src       source object
//		dest      point of dest object
//		notCopyZero  if copy zero field too
//Deprecated: use Convert instead
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

//Deprecated: use Convert instead
func CopyStructFields(src interface{}, dest interface{}, copyZero ...bool) (dirty bool, err error) {

	// dest必须为指针
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		err = errx.New("need dest ptr")
		return
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
			((len(copyZero) <= 0 || !copyZero[0]) && srcValueField.IsZero()) || //判断是否要复制零值 不要并且是零值就跳过
			srcValueField.Kind() != destValueField.Kind() { //判断同类型
			continue
		}
		if !reflect.DeepEqual(destValueField.Interface(), srcValueField.Interface()) {
			dirty = true
		}
		destValueField.Set(srcValueField)
	}

	return
}

func Convert(src interface{}, dest interface{}) (err error) {
	key := "from"
	tagInDest := true
	sType := TypeOf(src)
	sValue := ValueOf(src)
	dType := TypeOf(dest)
	dValue := ValueOf(dest)

	tags := GetTags(dType, key, true)
	if len(tags) < 1 {
		tags = GetTags(sType, key, true)
		tagInDest = false
	}

	for _, value := range tags {
		if value == "" {
			return errors.New("tag can not be emtpy")
		}
	}

	for key, value := range parseFromTag(tags) {
		dField := dValue
		sField := sValue

		if tagInDest {
			for _, name := range strings.Split(key, ".") {
				dField = dField.FieldByName(name)
			}
			for _, name := range value.Names() {
				sField = sField.FieldByName(name)
			}
		} else {
			for _, name := range value.Names() {
				dField = dField.FieldByName(name)
			}
			for _, name := range strings.Split(key, ".") {
				sField = sField.FieldByName(name)
			}
		}

		if !dField.CanSet() {
			ptr := reflect.NewAt(dField.Type(), unsafe.Pointer(dField.UnsafeAddr()))
			ptr.Elem().Set(sField)
		} else {
			dField.Set(sField)
		}
	}

	return nil
}

func ValueOf(target interface{}) reflect.Value {
	v := reflect.ValueOf(target)
	return ElemOfValue(v)
}

func TypeOf(target interface{}) reflect.Type {
	t := reflect.TypeOf(target)
	return ElemOfType(t)
}

func ElemOfType(target reflect.Type) reflect.Type {
	for target.Kind() == reflect.Pointer {
		target = target.Elem()
	}
	return target
}

func ElemOfValue(target reflect.Value) reflect.Value {
	for target.Type().Kind() == reflect.Pointer {
		target = target.Elem()
	}
	return target
}

func GetTags(t reflect.Type, key string, recursion bool) map[string]string {
	tags := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr, ok := field.Tag.Lookup(key)
		if ok {
			tags[field.Name] = tagStr
		} else if recursion && (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer) {
			sub := ElemOfType(field.Type)
			subMap := GetTags(sub, key, recursion)
			for key, value := range subMap {
				tags[field.Name+"."+key] = value
			}
		}
	}
	return tags
}

func parseFromTag(tags map[string]string) (m map[string]tagFrom) {
	m = make(map[string]tagFrom, len(tags))
	for key, value := range tags {
		s := strings.Split(value, ",")
		if len(s) < 1 || s[0] == "" {
			continue
		}
		tag := tagFrom{
			Name: s[0],
		}
		if len(s) > 1 {
			tag.Type = s[1]
		}
		m[key] = tag
	}
	return
}

type tagFrom struct {
	Name string
	Type string
}

//Names returns string slice splited from tagFrom.Name by ","
func (t tagFrom) Names() []string {
	return strings.Split(t.Name, ".")
}
