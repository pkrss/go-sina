package kline

import (
	"container/list"
	"time"

	pkTime "github.com/pkrss/go-utils/time"
)

func handlePeriodKData(origList []K_MACDV, period string) (retList []K_MACDV) {

	l := list.New()
	var oldItem *K_MACDV
	var samePeriod bool
	for i, _ := range origList {
		origItem := origList[i]
		samePeriod = false

		if oldItem != nil {
			samePeriod = pkTime.CheckSamePeriod(period, time.Unix(oldItem.Time/1000, 0), time.Unix(origItem.Time/1000, 0))
		}

		if samePeriod {
			if oldItem.High < origItem.High {
				oldItem.High = origItem.High
			}
			if oldItem.Low > origItem.Low {
				oldItem.Low = origItem.Low
			}
			oldItem.Close = origItem.Close
			oldItem.Vol += origItem.Vol
			oldItem.Time = origItem.Time
			continue
		}

		if oldItem == nil {
			oldItem = &origItem
			continue
		}

		oldItem.Index = int64(l.Len())
		l.PushBack(oldItem)
		oldItem = nil
	}

	if oldItem != nil {
		oldItem.Index = int64(l.Len())
		l.PushBack(oldItem)
		oldItem = nil
	}

	ll := l.Len()
	retList = make([]K_MACDV, ll)

	if ll == 0 {
		return
	}

	b := l.Front()
	for i := 0; i < ll; i++ {
		switch v := b.Value.(type) {
		case *K_MACDV:
			retList[i] = *v
		}
		b = b.Next()
	}

	return
}
