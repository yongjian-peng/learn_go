package mediator

// AirportMediator 机场调度中介者
type AirportMediator interface {
	CanLandAirport(aircraft Aircraft) bool // 确认是否可以降落
	NotifyWaitingAircraft()                // 通知等待降落的其他飞机
}

// ApproachTower 机场塔台
type ApproachTower struct {
	hasFreeAirstrip bool
	waitingQueue    []Aircraft // 等待降落的飞机队列
}

func (a *ApproachTower) CanLandAirport(aircraft Aircraft) bool {
	if a.hasFreeAirstrip {
		a.hasFreeAirstrip = false
		return true
	}
	// 没有空余的跑道，加入等待队列
	a.waitingQueue = append(a.waitingQueue, aircraft)
	return false
}

func (a *ApproachTower) NotifyWaitingAircraft() {
	if !a.hasFreeAirstrip {
		a.hasFreeAirstrip = true
	}
	if len(a.waitingQueue) > 0 {
		// 如果存在等待降落的飞机，通知第一个降落
		first := a.waitingQueue[0]
		a.waitingQueue = a.waitingQueue[1:]
		first.ApproachAirport()
	}
}
