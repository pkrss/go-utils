package main

import (
	"log"

	"github.com/pkrss/go-utils/beans"
	"github.com/pkrss/go-utils/examples/common"
	baseOrm "github.com/pkrss/go-utils/orm"
	"github.com/pkrss/go-utils/orm/inner"

	pqsql "github.com/pkrss/go-utils/orm/pqsql"
)

var dao baseOrm.BaseDaoInterface

func testMakeSlice() {
	l := dao.CreateModelSlice(10, 10)
	log.Printf("testMakeSlice :%v\n", l)
	// testMakeSlice :&[{{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}} {{} 0        {0 0 <nil>} {0 0 <nil>}}]
}

func testGetStructFields() {
	app := common.CreateSampleOAuthApp()
	app.AppBaseUrl = ""
	m := inner.GetStructDbFieldsAndValues(app, true)
	log.Printf("GetStructDbFieldsAndValues(true) = %v\n", m)
	m = inner.GetStructDbFieldsAndValues(app, false)
	log.Printf("GetStructDbFieldsAndValues(false) = %v\n", m)
}

func testFindOne() {
	id := 1
	r, e := dao.FindOneById(id)
	if e != nil {
		log.Printf("FindOneById error:%s\n", e.Error())
	}
	log.Printf("FetchById(%v) = %v\n", id, r)
}

func testUpdate() {
	r, e := dao.FindOneById(7)
	if e != nil {
		log.Printf("update select one error:%s\n", e.Error())
		return
	}
	app := r.(*common.OAuthApp)
	app.AppId = "123456"
	e = dao.UpdateById(app, app.Id)
	if e != nil {
		log.Printf("Update error:%s\n", e.Error())
	}
	log.Printf("Update()=%v\n", app)
}

func testDelete() {
	app := common.CreateSampleOAuthApp()
	e := dao.DeleteByFilter("Code", app.Code)
	log.Printf("Delete()=%v\n", e)
}

func testInsert() {

	app := common.CreateSampleOAuthApp()
	e := dao.Insert(app)
	if e != nil {
		log.Printf("Insert error:%s\n", e.Error())
	}
	log.Printf("Insert(%v)\n", app)
}

func testFindList() {
	pageable := beans.Pageable{}
	pageable.PageSize = 20
	pageable.Sort = "-id"
	pageable.RspCodeFormat = true
	pageable.CondArr = make(map[string]string, 0)
	// pageable.PageNumber = 3
	pageable.CondArr["q"] = "WX"

	l, total, e := dao.SelectSelSqlList("", &pageable, nil, func(listRawHelper *baseOrm.ListRawHelper) error {
		listRawHelper.SetCondArrLike("q", "title", "code")
		return nil
	})
	if e != nil {
		log.Printf("SelectList error:%s\n", e.Error())
	}
	log.Printf("SelectList() total=%v list=%v\n", total, l)
}

func main() {
	pqsql.Db = pqsql.CreatePgSql()
	baseOrm.DefaultOrmAdapter = &pqsql.PgSqlAdapter{}

	dao = baseOrm.CreateBaseDao(&common.OAuthApp{})

	// testMakeSlice()
	// testGetStructFields()
	// testFindOne()
	// testFindList()
	testDelete()
	testInsert()
	// testUpdate()

}
