package jtt

import (
	"math"
	"sync"
)

type counter struct {
	mu  sync.Mutex // 互斥锁，用于并发安全
	val uint16     // 当前计数值

	min uint16 // 最小计数值
	max uint16 // 最大计数值
}

// 循环自增计数器，取值范围 [min, max}
func newCounter(min, max uint16) *counter {
	return &counter{val: min, min: min, max: max}
}

func (c *counter) Get() uint16 {
	c.mu.Lock()         // 获取锁，确保并发安全
	defer c.mu.Unlock() // 函数返回前释放锁

	val := c.val
	c.val++
	if c.val > c.max {
		c.val = c.min // 超过最大值后，循环到 min
	}
	return val
}

var (
	serialNumber = newCounter(1, math.MaxUint16)
)

// GenerateSerialNumber 生成消息序列号，取值范围 [1, 65535}
func GenerateSerialNumber() uint16 {
	return serialNumber.Get()
}
