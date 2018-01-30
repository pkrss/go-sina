package kline

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	pkBeans "github.com/pkrss/go-utils/beans"
	pkFile "github.com/pkrss/go-utils/file"
	"github.com/pkrss/go-utils/profile"
	pkTime "github.com/pkrss/go-utils/time"
)

var kline_peroids_list []string
var kline_indciates_list []string

func KlineDataQuery(stk string, period string, indicate string, pageable *pkBeans.Pageable) (retList []K_MACDV, retListSize int64, retE error) {

	if stk == "" || period == "" {
		retE = errors.New("参数[stk、period]不能为空")
		return
	}

	// stk = strings.ToLower(stk)
	period = strings.ToLower(period)
	indicate = strings.ToLower(indicate)

	if kline_peroids_list == nil {
		kline_peroids_list = strings.Split(profile.ProfileReadString("kline_peroids_list"), ",")
	}

	if kline_indciates_list == nil {
		kline_indciates_list = strings.Split(profile.ProfileReadString("kline_indciates_list"), ",")
	}

	saveFileName := profile.ProfileReadString("kline_save_path_fmt")
	saveFileName = strings.Replace(saveFileName, "{stk}", stk, -1)
	saveFileName = strings.Replace(saveFileName, "{period}", period, -1)
	if indicate != "" {
		saveFileName = strings.Replace(saveFileName, "{indciate}", indicate, -1)
	} else {
		saveFileName = strings.Replace(saveFileName, "_{indciate}", indicate, -1)
	}

	lastAccessTime, fileExist := pkFile.FileLastWriteTime(saveFileName)

	if !fileExist {
		pkFile.CreateDir(pkFile.FileDir(saveFileName))
	}

	needFetch := !fileExist

	if !needFetch {
		needFetch = !pkTime.CheckSamePeriod(period, time.Unix(lastAccessTime, 0))
	}

	var itemList []K_MACDV

	if needFetch {
		itemList, retE = queryKline(stk, period, indicate)

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

		if indicate != "" && retE == nil {
			// itemList, retE = handleKIndciate(saveFileName, itemList)
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

	if retE == nil && itemList != nil {
		retListSize = int64(len(itemList))
		// itemList = pkContainer.ListSubPage(itemList, pageable)
	}

	retList = itemList

	return
}
