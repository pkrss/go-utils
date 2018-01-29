package main

import (
	"log"

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
	AccessTokenExpireTime baseOrm.JsonTime `json:"accessTokenExpireTime" pg:",datetime"`
	CreateTime            baseOrm.JsonTime `json:"createTime" pg:",datetime"`
}

func (this *OAuthApp) TableName() string {
	return "myzc_oauth_app"
}

func main() {
	pqsql.Db = pqsql.CreatePgSql()

	id := 1
	var e error
	var oAuthApp OAuthApp
	var r pqsql.BaseModelInterface
	dao := pqsql.CreateBaseDao(&oAuthApp)
	r, e = dao.FindOneById(id)
	if e != nil {
		log.Print(e.Error())
		// return
	}
	dao.FindOneById(id)

	dao.FindOneById(id)

	dao.FindOneById(id)
	log.Printf("FetchById(%v) = %v\n", id, r)
}
