package minute

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	pkFile "github.com/pkrss/go-utils/file"
	"github.com/pkrss/go-utils/profile"
	pkTime "github.com/pkrss/go-utils/time"
)

func MinuteDataQuery(stk string) (retList []MINUTEDATA, retE error) {

	if stk == "" {
		retE = errors.New("参数[stk]不能为空")
		return
	}

	saveFileName := profile.ProfileReadString("minute_save_path_fmt")
	saveFileName = strings.Replace(saveFileName, "{stk}", stk, -1)

	lastAccessTime, fileExist := pkFile.FileLastWriteTime(saveFileName)

	if !fileExist {
		pkFile.CreateDir(pkFile.FileDir(saveFileName))
	}

	needFetch := !fileExist

	if !needFetch {
		needFetch = !pkTime.CheckSamePeriod("1m", time.Unix(lastAccessTime, 0))
	}

	var itemList []MINUTEDATA

	if needFetch {
		itemList, retE = queryMinute(stk)

		if retE == nil {
			jsonData, err := json.Marshal(itemList)

			if err == nil {
				err = json.Unmarshal(jsonData, &itemList)
				if err == nil && (len(itemList) > 1) && (itemList[0].Time > itemList[1].Time) {
					sort.Slice(itemList, func(i, j int) bool {
						return itemList[i].Time < itemList[j].Time
					})
					jsonData, err = json.Marshal(&itemList)
				}
			}

			if err == nil {
				retE = ioutil.WriteFile(saveFileName, jsonData, 0666)
			} else {
				retE = err
			}
		}
	} else {
		fileData, e := ioutil.ReadFile(saveFileName)
		if e != nil {
			retE = e
			return
		}
		retE = json.Unmarshal(fileData, &itemList)
		if retE != nil {
			return
		}
	}

	retList = itemList

	return
}
