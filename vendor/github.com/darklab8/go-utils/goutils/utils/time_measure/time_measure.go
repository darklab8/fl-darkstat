package time_measure

import (
	"fmt"
	"time"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
)

type TimeMeasurer struct {
	msg         string
	ops         []typelog.LogType
	timeStarted time.Time

	ResultErr error
}

type TimeOption func(m *TimeMeasurer)

func NewTimeMeasure(opts ...TimeOption) *TimeMeasurer {
	m := &TimeMeasurer{
		timeStarted: time.Now(),
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithMsg(msg string) TimeOption {
	return func(m *TimeMeasurer) { m.msg = msg }
}

func WithLogMsgs(log_types ...typelog.LogType) TimeOption {
	return func(m *TimeMeasurer) { m.ops = log_types }
}

func (m *TimeMeasurer) Close() {
	utils_logus.Log.Debug(fmt.Sprintf("time_measure %v | %s", time.Since(m.timeStarted), m.msg), m.ops...)
}

func TimeMeasure(callback func(m *TimeMeasurer), opts ...TimeOption) *TimeMeasurer {
	m := NewTimeMeasure(opts...)
	defer m.Close()
	callback(m)
	return m
}
