package quote

import (
	"sync"
	"time"

	"github.com/pkrss/gosina/check"
)

type HqInfo struct {
	Quote      map[string]interface{}
	UpdateTime time.Time
}

var hqinfo map[string]HqInfo
var hqinfoMutex *sync.Mutex

func GetCacheHqInfo(stk string) (ret HqInfo, ok bool) {
	if stk == "" {
		return
	}

	if hqinfoMutex == nil {
		hqinfoMutex = new(sync.Mutex)
	}

	hqinfoMutex.Lock()
	ret, ok = hqinfo[stk]
	hqinfoMutex.Unlock()

	return
}

func SetCacheHqInfo(stk string, quote map[string]interface{}) {

	if stk == "" || quote == nil {
		return
	}

	if hqinfoMutex == nil {
		hqinfoMutex = new(sync.Mutex)
	}

	now := time.Now()

	v := HqInfo{
		Quote:      quote,
		UpdateTime: now,
	}

	hqinfoMutex.Lock()
	hqinfo[stk] = v
	hqinfoMutex.Unlock()
}

func PruneCacheHqInfo(params ...string) {

	var invalidPeriod string
	if len(params) > 0 {
		invalidPeriod = params[0]
	}

	if hqinfoMutex == nil {
		return
	}

	now := time.Now()
	l := make([]string, 0)
	hqinfoMutex.Lock()

	for stk, v := range hqinfo {
		if !check.CheckSamePeriod(invalidPeriod, v.UpdateTime, now) {
			l = append(l, stk)
		}
	}

	for _, stk := range l {
		delete(hqinfo, stk)
	}

	hqinfoMutex.Unlock()

}
func GetHqInfo(stk string, invalidPeriod string) (ret map[string]interface{}, ok bool) {

	var v HqInfo
	now := time.Now()

	if invalidPeriod == "" {
		return
	}

	v, ok = GetCacheHqInfo(stk)

	if ok {
		ok = check.CheckSamePeriod(invalidPeriod, v.UpdateTime, now)

		if ok {
			ret = v.Quote
			return
		}
	}

	m := FetchSinaHqDo(stk)
	if m == nil {
		ok = false
		return
	}

	ret, ok = m[stk]

	if ok {
		SetCacheHqInfo(stk, ret)
	}

	return
}

func GetHqField(stk string, field string, invalidPeriod string) (ret interface{}, ok bool) {
	if field == "" {
		return
	}

	v, o := GetHqInfo(stk, invalidPeriod)
	if !o {
		return
	}

	ret, ok = v[field]
	if !ok {
		return
	}
	return
}

func GetHqFloatField(stk string, field string, invalidPeriod string) (ret float64, ok bool) {

	v, o := GetHqField(stk, field, invalidPeriod)
	if !o {
		return
	}

	switch v2 := v.(type) {
	case float64:
		ret = v2
	case float32:
		ret = float64(v2)
	case int64:
		ret = float64(v2)
	case int:
		ret = float64(v2)
	}
	return
}

func GetHqIntField(stk string, field string, invalidPeriod string) (ret int64, ok bool) {

	v, o := GetHqField(stk, field, invalidPeriod)
	if !o {
		return
	}

	switch v2 := v.(type) {
	case int64:
		ret = v2
	case int:
		ret = int64(v2)
	case float64:
		ret = int64(v2)
	case float32:
		ret = int64(v2)
	}
	return
}

func GetHqStringField(stk string, field string, invalidPeriod string) (ret string, ok bool) {

	v, o := GetHqField(stk, field, invalidPeriod)
	if !o {
		return
	}

	switch v2 := v.(type) {
	case string:
		ret = v2
	}
	return
}
