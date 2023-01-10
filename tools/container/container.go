package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var stuff = make([]interface{}, 0)
var instances = new(sync.Map)
var lock sync.RWMutex

func Get[T any]() T {
	//先获取泛型的类型
	var t *T
	tt := reflect.TypeOf(t).Elem()
	fmt.Printf("%+v\n", tt)
	//看map缓存里面有没有，有就返回
	if instance, ok := instances.Load(tt); ok {
		return instance.(T)
	}

	//FIXME: 设a b 两个接口， b接口集成a接口，c d 对象分别继承 a b 两个接口，则此时取a接口的对象 c 和 d 都有可能取到
	//没有就从slice里面读
	lock.RLock()
	defer lock.RUnlock()
	for _, item := range stuff {
		if t, ok := item.(T); ok {
			instances.Store(tt, item)
			return t
		}
	}

	panic(errors.New("instance not found"))
}

// Set
//
// Deprecated: use SetT instead.
func Set(instance any) {
	lock.Lock()
	defer lock.Unlock()
	stuff = append(stuff, instance)
}

func SetT[T any](instance T) {
	//先获取泛型的类型
	var t *T
	tt := reflect.TypeOf(t).Elem()
	fmt.Printf("%+v\n", tt)
	instances.Store(tt, instance)
}
