package timeit

import (
	"fmt"
	"time"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Timer struct {
	msg         string
	ops         []typelog.LogType
	timeStarted time.Time

	ResultErr error
}

type TimeOption func(m *Timer)

func NewTimer(msg string, opts ...TimeOption) *Timer {
	return NewTimerMain(append([]TimeOption{WithMsg(msg)}, opts...)...)
}

func NewTimerL(msg string, opts ...typelog.LogType) *Timer {
	return NewTimer(msg, WithLogs(opts...))
}

func NewTimerMain(opts ...TimeOption) *Timer {
	m := &Timer{
		timeStarted: time.Now(),
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithMsg(msg string) TimeOption {
	return func(m *Timer) { m.msg = msg }
}

func WithLogs(log_types ...typelog.LogType) TimeOption {
	return func(m *Timer) { m.ops = log_types }
}

func (m *Timer) Close() {
	utils_logus.Log.Debug(fmt.Sprintf("time_measure %v | %s", time.Since(m.timeStarted), m.msg), m.ops...)
}

func NewTimerF(callback func(m *Timer), opts ...TimeOption) *Timer {
	m := NewTimerMain(opts...)
	defer m.Close()
	callback(m)
	return m
}

func NewTimerMF(msg string, callback func(m *Timer), opts ...TimeOption) *Timer {
	return NewTimerF(callback, append([]TimeOption{WithMsg(msg)}, opts...)...)
}

func NewTimerMFL(msg string, callback func(m *Timer), opts ...typelog.LogType) *Timer {
	return NewTimerMF(msg, callback, WithLogs(opts...))
}
