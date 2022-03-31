package tools

import (
	"sync"
	"time"
)

func NewDelayValue(initValue interface{}, duration time.Duration) *DelayValue {
	return &DelayValue{
		currentValue: initValue,
		time:         time.Now(),
		duration:     duration,
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
	lockValue    interface{}   //锁定的值
	lockDuration time.Duration //锁定时间
	lockAt       time.Time     //锁定开始时间
}

//Lock 锁定状态
func (t *DelayValue) Lock(state interface{}, duration time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.lockValue = state
	t.lockDuration = duration
	t.lockAt = time.Now()
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
			t.currentValue = state
			t.time = time.Now()

			//没锁定才改变状态
			if t.lockValue == nil || t.lockAt.Add(t.lockDuration).Before(time.Now()) {
				t.lockValue = nil
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

	//判断锁定状态,如果有就返回锁定的值
	if t.lockValue != nil && t.lockAt.Add(t.lockDuration).After(time.Now()) {
		return t.lockValue
	} else {
		t.lockValue = nil
	}

	return t.currentValue
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
