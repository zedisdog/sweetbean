package tools

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/atomic"
	"log"
	"time"
)

var (
	ErrClosed = errors.New("chan is closed")
)

type Storage interface {
	HasMore(ctx context.Context) bool
	SaveMany(messages ...any) (err error)
	PullByLimit(ctx context.Context, i int) (messages []any, err error)
}

func WithLoadDuration(d time.Duration) func(*MemQueue) {
	return func(queue *MemQueue) {
		queue.loadDuration = d
	}
}

//NewMemQueue 创建队列
//  ctx: 将会在调用 Storage.HasMore 和 Storage.PullByLimit 时原样返回，可用来传递数据
//  storage: 实现了 Storage 接口的存储对象
//  size: chan大小
func NewMemQueue(ctx context.Context, storage Storage, size int, setters ...func(*MemQueue)) *MemQueue {
	queue := &MemQueue{
		memQueue: make(chan any, size),
		storage:  storage,
		ctx:      ctx,
		running:  atomic.NewBool(true),
		size:     size,
	}
	for _, set := range setters {
		set(queue)
	}
	return queue
}

//MemQueue 简易内存队列，数据先通过缓冲chan存在内存中，当chan存满后通过 Storage 存入任何持久化存储中
type MemQueue struct {
	memQueue     chan interface{}
	storage      Storage
	ctx          context.Context
	running      *atomic.Bool
	size         int
	loadDuration time.Duration //读取存储的间隔时间
}

func (m MemQueue) Put(msg any) (err error) {
	if !m.running.Load() {
		err = ErrClosed
		return
	}

	if m.storage.HasMore(m.ctx) || len(m.memQueue) >= m.size {
		err = m.storage.SaveMany(msg)
		return
	}

	m.memQueue <- msg
	return
}

//Pull 从通道中获取数据
func (m MemQueue) Pull() (msg any, err error) {
	if !m.running.Load() {
		err = ErrClosed
		return
	}
	for {
		select {
		case msg, ok := <-m.memQueue:
			if !ok {
				return nil, errors.New("queue closed")
			} else {
				return msg, nil
			}
		default:
			_ = m.replenish()
			time.Sleep(m.loadDuration)
		}
	}
}

//Out 以channel的形式获取数据
func (m MemQueue) Out() chan any {
	var c = make(chan any)
	go func() {
		for {
			msg, err := m.Pull()
			if err != nil {
				log.Printf("message: pull failed, logger: sweetbean.tools.MemQueue, error: %s\n", err.Error())
				close(c)
				break
			}
			c <- msg
		}
	}()
	return c
}

func (m MemQueue) replenish() (err error) {
	need := m.size - len(m.memQueue)
	if m.storage.HasMore(m.ctx) && need > 0 {
		var msgs []any
		msgs, err = m.storage.PullByLimit(m.ctx, need)
		if err != nil {
			err = fmt.Errorf("read message from storage failed: %w", err)
			return
		}
		for _, msg := range msgs {
			m.memQueue <- msg
		}
	}
	return
}

func (m *MemQueue) Close() {
	m.running.Store(false)
	close(m.memQueue)
	var msgs []any
	for m := range m.memQueue {
		msgs = append(msgs, m)
	}
	if len(msgs) > 0 {
		err := m.storage.SaveMany(msgs...)
		if err != nil {
			log.Printf("save msg from mem to storage failed: %s", err.Error())
		}
	}
}
