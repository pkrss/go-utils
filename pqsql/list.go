package pqsql

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkrss/go-utils/beans"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type ValueType int

// iota 初始化后会自动递增
const (
	String ValueType = iota
	Int64
	Int
	InSqlStrVar
)

type ListRawHelper struct {
	DbQuery           *orm.Query
	Db                *pg.DB
	Pageable          *beans.Pageable
	Where             string
	WhereArgs         []interface{}
	ResultListPointer interface{}
	TotalResult       int64
	ObjModel          BaseModelInterface
	UserData          interface{}
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

		this.Where += " " + in_name + " in ? "
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

		this.Where += " " + like_name + " like ? "
		this.WhereArgs = append(this.WhereArgs, like_value)
	}
}

func (this *ListRawHelper) getQueryPageablePostfix(sql *string) {
	pageable := this.Pageable

	if pageable == nil {
		return
	}

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
		if sql == nil {
			this.DbQuery = this.DbQuery.OrderExpr(orderBy)
		} else {
			*sql = *sql + " " + orderBy
		}
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
	if sql == nil {
		this.DbQuery = this.DbQuery.Limit(pageable.PageSize).Offset(offset)
	} else {
		*sql = *sql + fmt.Sprintf(" limit %d offset %v", pageable.PageSize, offset)
	}
}

func (this *ListRawHelper) SetCondArrLike(condKey string, dbKeys ...string) {
	this.SetCondArrParam(condKey, true, String, dbKeys...)
}

func (this *ListRawHelper) SetCondArrEqu(condKey string, valueType ValueType, dbKeys ...string) {
	this.SetCondArrParam(condKey, false, valueType, dbKeys...)
}

func (this *ListRawHelper) SetCondArrParam(condKey string, trueLikeFalseEqual bool, valueType ValueType, dbKeys ...string) {
	pageable := this.Pageable

	c := len(dbKeys)
	if pageable == nil || c == 0 {
		return
	}

	s, ok := pageable.CondArr[condKey]
	if ok {
		if this.Where != "" {
			this.Where += " AND"
		}
		this.Where += " ("

		if valueType == InSqlStrVar {
			this.Where += s
			for i := 0; i < c; i++ {
				this.WhereArgs = append(this.WhereArgs, dbKeys[i])
			}
		} else {

			for i := 0; i < c; i++ {
				v := dbKeys[i]

				var v2 interface{}
				switch valueType {
				case String:
					if trueLikeFalseEqual {
						s = "%" + s + "%"
					}
					v2 = s
				case Int:
					tmp, err := strconv.Atoi(s)
					if err != nil {
						log.Printf("SetCondArrParam Atoi %s=%v error: %s\n", v, v2, err.Error())
						return
					}
					v2 = tmp
				case Int64:
					tmp, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						log.Printf("SetCondArrParam praseInt %s=%v error: %s\n", v, v2, err.Error())
						return
					}
					v2 = tmp
				}
				if v2 == nil {
					log.Printf("SetCondArrParam prase %s=nil\n", v)
					break
				}

				this.WhereArgs = append(this.WhereArgs, v2)
				if trueLikeFalseEqual {
					this.Where += v + " like ?"
				} else {
					this.Where += v + " = ?"
				}
				if i != c-1 {
					this.Where += " OR "
				}
			}
		}
		this.Where += " ) "
	}
}

func (this *ListRawHelper) Query() (int64, error) {

	if this.DbQuery == nil {
		return 0, errors.New("DbQuery is nil")
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

	this.getQueryPageablePostfix(nil)

	cnt, err := this.DbQuery.SelectAndCount(this.ResultListPointer)
	if err != nil {
		return 0, err
	}

	this.TotalResult = int64(cnt)

	return this.TotalResult, nil
}

func (this *ListRawHelper) SelSqlListQuery(selSql string) (total int64, e error) {

	pageable := this.Pageable

	if pageable == nil {
		e = errors.New("pageable is nil")
		return
	}

	this.appendNormalWhereConds()

	sql := "SELECT COUNT(*) FROM " + this.ObjModel.TableName() + " " + this.Where
	_, e = this.Db.QueryOne(pg.Scan(&total), sql, this.WhereArgs...)
	if e != nil {
		return
	}
	this.TotalResult = total

	if selSql == "" {
		selSql = `SELECT `

		c := len(pageable.Columns)
		if c == 0 {
			selSql += `*`
		} else {
			selSql += strings.Join(pageable.Columns, ",")
		}
		selSql += ` FROM ` + this.ObjModel.TableName()
	}
	sql = selSql + " " + this.Where

	this.getQueryPageablePostfix(&sql)

	_, e = this.Db.Query(this.ResultListPointer, sql, this.WhereArgs...)

	return

}
