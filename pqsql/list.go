package pqsql

import (
	"errors"
	"fmt"
	"hx98/base/beans"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type ValueType int

// iota 初始化后会自动递增
const (
	String ValueType = iota
	Int64
	Int
)

type ListRawHelper struct {
	Pageable          *beans.Pageable
	TableName         string
	SelectSql         string
	Where             string
	WhereArgs         []interface{}
	ResultListPointer interface{}
	TotalResult       int64
}

func MakeListRawHelper(resultListPointer interface{}, pageable *beans.Pageable, tableName string) *ListRawHelper {
	ret := ListRawHelper{}
	ret.ResultListPointer = resultListPointer
	ret.Pageable = pageable
	ret.TableName = tableName
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
		this.WhereArgs = append(this.WhereArgs, in_values)
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

func (this *ListRawHelper) getQueryPageablePostfix() string {
	pageable := this.Pageable

	if pageable == nil {
		return ""
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

	cond := ""

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
		cond += " ORDER BY " + orderBy
	}

	cond += fmt.Sprintf(" LIMIT %d OFFSET %d", pageable.PageSize, offset)

	return cond
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
	pageable := this.Pageable

	if pageable == nil {
		return 0, errors.New("pageable is nil")
	}

	this.appendNormalWhereConds()

	where := this.Where
	sql := this.SelectSql

	if sql == "" {
		if len(pageable.Columns) > 0 {
			sql = "SELECT " + strings.Join(pageable.Columns, ",") + " FROM " + this.TableName
		} else {
			sql = "SELECT * FROM " + this.TableName
		}
	}

	if where != "" {
		sql += " WHERE " + where
	}

	o := orm.NewOrm()

	err := o.Raw("SELECT COUNT(*) FROM "+this.TableName+" "+where, this.WhereArgs).QueryRow(&this.TotalResult)
	if err != nil {
		return this.TotalResult, err
	}

	sql += " " + this.getQueryPageablePostfix()

	num, err := o.Raw(sql, this.WhereArgs).QueryRows(this.ResultListPointer)
	if num > 0 {

	}
	if err != nil {
		return this.TotalResult, err
	}

	return this.TotalResult, nil
}
