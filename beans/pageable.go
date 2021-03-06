package beans

type Pageable struct {
	Page int `json:"page"` // from page number, from 0

	Size int `json:"size"` // query limit

	CondArr map[string]string `json:"condArr"` // query conditions

	Sort string `json:"sort"` // order, ex: "-id"

	Columns []string `json:"columns"` // fetch specify columns

	RelatedSel bool `json:"relatedSel"` // is multiple query

	RspCodeFormat bool `json:"rspCodeFormat"` // response format. true:code json format, false: normal rest format

	OffsetOldField int `json:"-"` // inner used
}

func (this *Pageable) CalcOffsetAndLimit(total int) (ok bool, begin int, end int) {
	ok = false

	if total == 0 {
		return
	}

	limit := this.Size

	if this.OffsetOldField != 0 {
		begin = this.OffsetOldField
	} else {
		begin = this.Page * this.Size
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
