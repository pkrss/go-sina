package quote

import (
	"sync"

	"github.com/pkrss/go-sina/check"
)

var stk2quotes map[string]map[string]interface{}
var stk2quotesMutex sync.Mutex

func checkStkNeedSend(stk string, quote map[string]interface{}) bool {
	if stk2quotes == nil {
		stk2quotes = make(map[string]map[string]interface{}, 0)
	}

	stk2quotesMutex.Lock()
	defer stk2quotesMutex.Unlock()

	oldQuote, ok := stk2quotes[stk]
	if ok {
		if quote["time"] == oldQuote["time"] && !check.IsCtp(stk) {
			return false
		}
	}

	stk2quotes[stk] = quote

	return true
}
