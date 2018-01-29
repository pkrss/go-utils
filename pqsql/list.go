package pqsql

import (
	"errors"
	"hx98/base/beans"
	"strconv"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type ValueType int

// iota 初始化后会自动递增
const (
	String ValueType = iota
	Int64
	Int
)

type ListRawHelper struct {
	DbQuery           *orm.Query
	Pageable          *beans.Pageable
	Where             string
	WhereArgs         []interface{}
	ResultListPointer interface{}
	TotalResult       int64
}

func MakeListRawHelper(resultListPointer interface{}, pageable *beans.Pageable) *ListRawHelper {
	ret := ListRawHelper{}
	ret.ResultListPointer = resultListPointer
	ret.Pageable = pageable
	ret.WhereArgs = make([]interface{}, 0)

	return &ret
}

func (this *ListRawHelper) appendNormalWhereConds() {
	if this.Pageable == nil {
		return
	}

	condArr := this.Pageable.CondArr

	if condArr == nil {
		return
	}

	for f := true; f; f = false {
		in_name, ok := condArr["in_name"]
		if !ok || in_name == "" {
			break
		}

		in_list, ok := condArr["in_list"]
		if !ok || in_list == "" {
			break
		}

		in_name = strings.Trim(in_name, " ,")
		if in_list == "" {
			break
		}

		in_type, _ := condArr["in_type"]

		in_values := make([]interface{}, 0)
		in_values2 := strings.Split(in_list, ",")
		for _, v := range in_values2 {
			if v == "" {
				continue
			}
			switch in_type {
			case "i":
				i, e := strconv.Atoi(v)
				if e == nil {
					in_values = append(in_values, i)
				}
				break
			case "l":
				i, e := strconv.ParseInt(v, 10, 64)
				if e == nil {
					in_values = append(in_values, i)
				}
				break
			case "u":
				break
			}
		}

		if len(in_values) == 0 {
			break
		}

		this.Where += "in_name in ?"
		this.WhereArgs = append(this.WhereArgs, pg.In(in_values))
	}

	for f := true; f; f = false {
		like_name, ok := condArr["like_name"]
		if !ok || like_name == "" {
			break
		}

		like_value, ok := condArr["like_value"]
		if !ok || like_value == "" {
			break
		}

		this.Where += "like_name like ?"
		this.WhereArgs = append(this.WhereArgs, like_value)
	}
}

func (this *ListRawHelper) getQueryPageablePostfix() {
	pageable := this.Pageable

	if pageable == nil {
		return
	}

	if pageable.PageNumber < 1 {
		pageable.PageNumber = 1
	}

	if pageable.PageSize == 0 {
		pageable.PageSize = 20
	}

	offset := pageable.OffsetOldField
	if offset == 0 {
		offset = (pageable.PageNumber - 1) * pageable.PageSize
	}

	this.DbQuery = this.DbQuery.Limit(pageable.PageSize).Offset(offset)

	// o := orm.NewOrm()
	// var total int64
	// err := o.Raw("SELECT COUNT(*) FROM " + tableName + " " + cond).QueryRow(&total)
	// if err != nil {
	// 	return total
	// }

	if pageable.Sort != "" {
		orderBy := pageable.Sort
		if strings.ContainsAny(orderBy, ` ()'"`) {
			orderBy = ""
		}
		if strings.HasPrefix(orderBy, "-") {
			orderBy = orderBy[1:] + " DESC"
		} else {
			if strings.HasPrefix(orderBy, "+") {
				orderBy = orderBy[1:]
			}
			orderBy += " ASC"
		}
		this.DbQuery = this.DbQuery.OrderExpr(orderBy)
	}
}

func (this *ListRawHelper) SetCondArrLike(condKey string, dbKeys ...string) {
	this.SetCondArrParam(condKey, true, String, dbKeys...)
}

func (this *ListRawHelper) SetCondArrParam(condKey string, trueLikeFalseEqual bool, valueType ValueType, dbKeys ...string) {
	pageable := this.Pageable

	c := len(dbKeys)
	if pageable == nil || c == 0 {
		return
	}

	s, ok := pageable.CondArr[condKey]
	if ok {
		this.Where += "AND ("
		for i := 0; i < c; i++ {
			v := dbKeys[i]
			if trueLikeFalseEqual {
				this.Where += v + " like ?"
			} else {
				this.Where += v + " = ?"
			}
			if i != c-1 {
				this.Where += " OR "
			}

			var v2 interface{}
			switch valueType {
			case String:
				v2 = s
			case Int:
				tmp, err := strconv.Atoi(s)
				if err != nil {
					return
				}
				v2 = tmp
			case Int64:
				tmp, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return
				}
				v2 = tmp
			}
			if v2 == nil {
				return
			}
			this.WhereArgs = append(this.WhereArgs, v2)
		}
		this.Where += ")"
	}
}

func (this *ListRawHelper) Query() (int64, error) {

	if this.DbQuery == nil {
		return 0, errors.New("db is nil")
	}

	pageable := this.Pageable

	if pageable == nil {
		return 0, errors.New("pageable is nil")
	}

	this.appendNormalWhereConds()

	if len(pageable.Columns) > 0 {
		this.DbQuery = this.DbQuery.Column(pageable.Columns...)
	}

	if this.Where != "" {
		this.DbQuery = this.DbQuery.Where(this.Where, this.WhereArgs...)
	}

	this.getQueryPageablePostfix()

	cnt, err := this.DbQuery.SelectAndCount(this.ResultListPointer)
	if err != nil {
		return 0, err
	}

	this.TotalResult = int64(cnt)

	return this.TotalResult, nil

	// v := reflect.ValueOf(this.ResultListPointer)
	// switch v.Kind() {
	// case reflect.Ptr:
	// 	v = v.Elem()
	// }

	// for i := 0; rows.Next(); i++ {
	// 	// json.Unmarshal()

	// 	// Get element of array, growing if necessary.
	// 	if v.Kind() == reflect.Slice {
	// 		// Grow slice if necessary
	// 		if i >= v.Cap() {
	// 			newcap := v.Cap() + v.Cap()/2
	// 			if newcap < 4 {
	// 				newcap = 4
	// 			}
	// 			newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
	// 			reflect.Copy(newv, v)
	// 			v.Set(newv)
	// 		}
	// 		if i >= v.Len() {
	// 			v.SetLen(i + 1)
	// 		}
	// 	}

	// 	if i < v.Len() {
	// 		// Decode into element.
	// 		rows.Scan(v.Index(i))
	// 		// d.value(v.Index(i))
	// 	} else {
	// 		// Ran out of fixed array: skip.
	// 		// d.value(reflect.Value{})
	// 	}

	// }

	// return this.TotalResult, nil
}
