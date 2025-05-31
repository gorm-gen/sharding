package list

import (
	"sort"

	"github.com/shopspring/decimal"
)

type List struct {
	offset   int64
	page     int64
	pageSize int64
	desc     bool
	asc      bool
	list     []*St
}

type St struct {
	ShardingValue string
	Total         uint64
}

func New(list []*St, opts ...Option) *List {
	l := &List{
		offset: 100,
		list:   list,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

type Sld struct {
	Page     int64
	PageSize int64
	Start    int64
	End      int64
	Num      int64
}

type Sl struct {
	ShardingValue string
	Total         int64
	Start         int64
	End           int64
	Num           int64
	List          []Sld
}

type SL []Sl

func (sl SL) ToSliceIndex() {
	for k, v := range sl {
		sl[k].Start -= 1
		for kk := range v.List {
			sl[k].List[kk].Start -= 1
		}
	}
}

func (l *List) Analysis() SL {
	if l.list == nil {
		return nil
	}

	// 排序
	if l.desc || l.asc {
		sort.Slice(l.list, func(i, j int) bool {
			if l.desc {
				return l.list[i].ShardingValue > l.list[j].ShardingValue
			} else {
				return l.list[i].ShardingValue < l.list[j].ShardingValue
			}
		})
	}

	start := int64(1)
	end := int64(0)
	for _, v := range l.list {
		end += int64(v.Total)
	}
	// 需要分页
	if l.page > 0 && l.pageSize > 0 {
		start = (l.page-1)*l.pageSize + 1
		end = l.page * l.pageSize
	}

	n := int64(0)
	list := make([]Sl, 0, len(l.list))
	_break := false

	for _, v := range l.list {
		if _break {
			break
		}

		if v.Total == 0 || v.ShardingValue == "" {
			continue
		}

		sl := Sl{
			ShardingValue: v.ShardingValue,
			Total:         int64(v.Total),
		}
		if n >= start {
			sl.Start = 1
		}
		for i := int64(1); i <= int64(v.Total); i++ {
			n++
			if n == start {
				sl.Start = i
			}
			if n == end {
				sl.End = i
				_break = true
				break
			}
			if i == int64(v.Total) {
				sl.End = i
			}
		}

		if sl.Start == 0 || sl.End == 0 {
			continue
		}

		sl.Num = sl.End - sl.Start + 1
		list = append(list, sl)
	}

	for k, v := range list {
		offset := l.offset
		if offset > v.Total {
			offset = v.Total
		}
		totalPage := (decimal.NewFromInt(v.Total).Div(decimal.NewFromInt(offset))).Ceil().BigInt().Int64()
		_break = false
		n = 0
		sldList := make([]Sld, 0, totalPage)

		for i := int64(1); i <= totalPage; i++ {
			if _break {
				break
			}

			sld := Sld{
				Page:     i,
				PageSize: offset,
			}
			if n >= v.Start {
				sld.Start = 1
			}
			for j := int64(1); j <= offset; j++ {
				n++
				if n == v.Start {
					sld.Start = j
				}
				if n == v.End {
					sld.End = j
					_break = true
					break
				}
				if j == offset {
					sld.End = j
				}
			}

			if sld.Start == 0 || sld.End == 0 {
				continue
			}

			sld.Num = sld.End - sld.Start + 1

			sldList = append(sldList, sld)
		}

		list[k].List = sldList
	}

	return list
}
