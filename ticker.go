package SryTask

import(
   "time"
)

type RealTimeManager struct{}

type MockTimeManager struct {
   C chan time.Time
}

func(rtm *RealTimeManager) NewTicker(duration time.Duration) (*time.Ticker) {
   return time.NewTicker(duration)
}

func NewMockTimeManager() (*MockTimeManager) {
   return &MockTimeManager{C: make(chan time.Time)}
}

func(mtm *MockTimeManager) NewTicker(_ time.Duration)(*time.Ticker) {
   return &time.Ticker{C: mtm.C}
}
