package main

import (
	"log"

	"github.com/pkrss/go-utils/beans"
	baseOrm "github.com/pkrss/go-utils/pqsql/custom"

	pqsql "github.com/pkrss/go-utils/pqsql"
)

type OAuthApp struct {
	pqsql.BaseModel
	Id                    int              `json:"id"`
	Code                  string           `json:"code"`
	Title                 string           `json:"title"`
	AppId                 string           `json:"appId"`
	AppSecurityKey        string           `json:"appSecurityKey"`
	AppOauthScope         string           `json:"appOauthScope"`
	AppBaseUrl            string           `json:"appBaseUrl"`
	AccessToken           string           `json:"accessToken"`
	AccessTokenExpireTime baseOrm.JsonTime `json:"accessTokenExpireTime"`
	CreateTime            baseOrm.JsonTime `json:"createTime"`
}

var dao pqsql.BaseDaoInterface

func (this *OAuthApp) TableName() string {
	return "myzc_oauth_app"
}

func testMakeSlice() {
	l := dao.CreateModelSlice(10, 10)
	log.Printf("testMakeSlice :%v\n", l)
	// testMakeSlice :&[{{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}}]
}

func testFindOne() {
	id := 1
	r, e := dao.FindOneById(id)
	if e != nil {
		log.Printf("FindOneById error:%s\n", e.Error())
	}
	log.Printf("FetchById(%v) = %v\n", id, r)
}

func testFindList() {
	pageable := beans.Pageable{}
	pageable.PageSize = 20
	pageable.Sort = "-id"
	pageable.RspCodeFormat = true
	pageable.CondArr = make(map[string]string, 0)
	pageable.CondArr["q"] = "WX"

	l, total, e := dao.SelectSelSqlList("", &pageable, nil, func(listRawHelper *pqsql.ListRawHelper) {
		listRawHelper.SetCondArrLike("q", "title", "code")
	})
	if e != nil {
		log.Printf("SelectList error:%s\n", e.Error())
	}
	log.Printf("SelectList() total=%v list=%v\n", total, l)
}

func main() {
	pqsql.Db = pqsql.CreatePgSql()

	var oAuthApp OAuthApp
	dao = pqsql.CreateBaseDao(&oAuthApp)

	// testMakeSlice()
	// testFindOne()
	testFindList()

}
