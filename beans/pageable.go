package beans

type Pageable struct {
	PageNumber int `json:"pageNumber"` // 第几页，1开始

	PageSize int `json:"pageSize"` // 一页限制几条

	CondArr map[string]string `json:"condArr"` // 条件

	Sort string `json:"sort"` // 排序，ex: "-id"

	Columns []string `json:"columns"` // 只请求指定列

	RelatedSel bool `json:"relatedSel"` // 是否多表查询

	RspCodeFormat bool `json:"rspCodeFormat"` // 返回旧的带code json格式

	OffsetOldField int `json:"-"` // 内部
}

func (this *Pageable) CalcOffsetAndLimit(total int) (ok bool, begin int, end int) {
	ok = false

	if total == 0 {
		return
	}

	limit := this.PageSize

	if this.OffsetOldField != 0 {
		begin = this.OffsetOldField
	} else if this.PageNumber == 0 {
		begin = 0
	} else {
		begin = (this.PageNumber - 1) * this.PageSize
	}

	if limit == 0 {
		return
	}

	if limit < 0 {
		limit = total
	}

	if begin < 0 {
		begin = total + begin
	}

	if begin < 0 {
		begin = 0
	}

	end = begin + limit

	if end > total {
		end = total
	}

	if begin >= end {
		return
	}

	ok = true

	return
}
