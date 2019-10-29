package timedtask

import (
	"fmt"
	"time"
)

type TimeTaskTicker struct {
	ticker *time.Ticker
}

func New(sec time.Duration) TimeTaskTicker {

	ticker := time.NewTicker(time.Second * sec)
	return TimeTaskTicker{ticker: ticker}
}

func (ttt *TimeTaskTicker) RunOnce(fc func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			//lib.SendErrorMesg("panic")
		}
	}()
	<-ttt.ticker.C

	fc()
	ttt.ticker.Stop()
}

func (ttt *TimeTaskTicker) Run(fc func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			//lib.SendErrorMesg("panic")
		}
	}()
	for range ttt.ticker.C {
		fc()
	}
}
