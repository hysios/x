package flow

import "time"

type Event struct {
	FlowId  string // flow id 流程 id
	Flow    string // flow name 流程名称
	EventId string // event id 事件 id
	Event   string // event name 事件名称
	// 参与用户
	Participants []string
	// 发起用户
	Owner string
	// 观察者
	Watchs []string
	// 事件发生时间
	StartAt time.Time
	// 事件结束时间
	EndAt time.Time
	// 事件状态
	Status    string
	UpdatedAt time.Time
}

type Flow struct {
	Name string
}

type Step struct {
	Name string
}

// RegisterEvent 注册事件
