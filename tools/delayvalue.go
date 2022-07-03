package tools

import (
	"sync"
	"time"

	"github.com/zedisdog/sweetbean/errx"
)

func NewDelayValue(initValue interface{}, duration time.Duration) *DelayValue {
	return &DelayValue{
		currentValue: initValue,
		time:         time.Now(),
		duration:     duration,
		lockValue:    initValue,
	}
}

//DelayValue 延迟duration改变值
type DelayValue struct {
	currentValue interface{}   //当前值
	time         time.Time     //值变更时间
	duration     time.Duration //持续观察时间
	lock         sync.RWMutex
	recently     bool //是否刚改变

	//下面是锁定状态的设置
	lockValue   interface{} //锁定的值
	valueLocked bool        //值是否被手动锁定
}

//Lock 锁定状态
func (t *DelayValue) Lock(state interface{}, duration time.Duration) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.valueLocked {
		return errx.New("value locked")
	} else {
		t.valueLocked = true
	}
	backDuration := t.duration

	t.lockValue = state
	t.currentValue = state
	t.duration = duration
	t.time = time.Now()

	go func(t *DelayValue) {
		t.lock.Lock()
		defer t.lock.Unlock()
		t.duration = backDuration
		t.valueLocked = false
	}(t)

	return nil
}

// ThisTime 设置当前状态
func (t *DelayValue) ThisTime(state interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if state != t.currentValue {
		t.time = time.Now()
		t.currentValue = state
	} else {
		if t.time.Add(t.duration).Before(time.Now()) {
			t.time = time.Now()
			if t.currentValue != t.lockValue {
				t.lockValue = state
				t.recently = true
			}
		}
	}
}

//Current 获取当前状态
func (t *DelayValue) Current() interface{} {
	t.lock.RLock()
	defer t.lock.RUnlock()
	t.recently = false

	return t.lockValue
}

//CurrentBool 获取当前状态
func (t *DelayValue) CurrentBool() bool {
	return t.Current().(bool)
}

//CurrentInt 获取当前状态
func (t *DelayValue) CurrentInt() int {
	return t.Current().(int)
}

//ChangeRecently 是否最近一次的改变,当至少一次通过Current方法获取值后，recently change就变为false
func (t *DelayValue) ChangeRecently() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.recently
}
