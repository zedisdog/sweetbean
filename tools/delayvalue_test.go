package tools

import (
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	value := NewDelayValue(false, 2)
	value.ThisTime(true)
	if value.ChangeRecently() {
		t.Fatal("expect false, actual true")
	}
	time.Sleep(3 * time.Second)
	if value.ChangeRecently() {
		t.Fatal("expect true, actual false")
	}
}

func TestLockValue(t *testing.T) {
	value := NewDelayValue(false, 2)
	value.Lock(true, 3)
	if !value.Current().(bool) {
		t.Fatal("expect true, actual false")
	}

	//等2.5秒设置成false 正常情况不会变
	time.Sleep(2500 * time.Millisecond)
	value.ThisTime(false)
	if value.ChangeRecently() || !value.Current().(bool) {
		t.Fatal("expect value not change, actual change")
	}

	//等一秒 上面的锁定3秒已经过去 这个时候观察时间变为开始的2秒
	// 重点并且没有刷新时间 就是说如果开始设置的时间比锁定设置的时间少 那么锁定以后可以马上修改值
	time.Sleep(1 * time.Second)
	value.ThisTime(false)
	if !value.ChangeRecently() || value.Current().(bool) {
		t.Fatal("expect value change, actual not change")
	}
}
