package quote

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/pkrss/go-sina/check"
)

var oldStrList []string = make([]string, 0)

type FuncOnPublishChannel func(string, string)

var OnPublishChannel FuncOnPublishChannel

func SetOnPublishChannel(f FuncOnPublishChannel) {
	OnPublishChannel = f
}

func NotifyStkChanged(stkStr string) {
	stkList := strings.Split(stkStr, ",")

	list := make([]string, 0)
	for _, v := range stkList {
		if check.IsZh(v) || check.IsFx(v) || check.IsCtp(v) {
			list = append(list, v)
		}
	}
	oldStrList = list
}

func Start() {
	for {
		time.Sleep(3 * 1000 * 1000 * 1000)

		RunOnce()
	}
}

func RunOnce() {
	t := time.Now().Unix()
	t = int64(math.Mod(float64(t), float64(24*60*60))) + 8*60*60

	// if !((t > (9 * 60 * 60 + 25 * 60) && t < (11 * 60 * 60 + 32 * 60)) || (t > (13 * 60 * 60) && t < (15 * 60 * 60 + 2 * 60))) {
	// 	utils.Println("startAsyncFetchData is not working time")
	// 	continue
	// }
	// utils.Println("startAsyncFetchData is working time")

	list := oldStrList
	if len(list) == 0 {
		// utils.Println("startAsyncFetchData list len = 0")
		return
	}

	// go
	fetchHqAndNotifyBatch(list)
}

func fetchHqAndNotifyBatch(list []string) {
	l := len(list)
	if l == 0 {
		return
	}
	const limit = 90
	if l <= limit {
		fetchHqAndNotify(list)
	} else {
		b := 0
		for b < l {
			e := b + limit
			if e > l {
				e = l
			}
			fetchHqAndNotify(list[b:e])
			b = e
		}
	}
}

func fetchHqAndNotify(list []string) {

	items := FetchSinaHqDo(list...)
	if items == nil {
		return
	}

	for stk, quote := range items {

		if !checkStkNeedSend(stk, quote) {
			// utils.Printf("fetchSinaHq not need send stk %s", stk)
			continue
		}

		jsonStr, err5 := json.Marshal(quote)
		if err5 != nil {
			continue
		}

		// utils.Printf("fetchSinaHq stk_quote_changed %s", quote["id"])

		if OnPublishChannel != nil {
			OnPublishChannel("stk_quote_changed", string(jsonStr))
		}
	}

}
