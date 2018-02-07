package common

import (
	baseOrm "github.com/pkrss/go-utils/orm"
	custom "github.com/pkrss/go-utils/orm/fields"
)

type OAuthApp struct {
	baseOrm.BaseModel
	Id                    int             `json:"id" orm:"ro"`
	Code                  string          `json:"code"`
	Title                 string          `json:"title"`
	AppId                 string          `json:"appId"`
	AppSecurityKey        string          `json:"appSecurityKey"`
	AppOauthScope         string          `json:"appOauthScope"`
	AppBaseUrl            string          `json:"appBaseUrl" orm:"null"`                // may be null
	AccessToken           string          `json:"accessToken" orm:"ro"`                 // read only
	AccessTokenExpireTime custom.JsonTime `json:"accessTokenExpireTime" orm:"auto_now"` // used update time now
	CreateTime            custom.JsonTime `json:"createTime" orm:"auto_now_add"`        // used create time now
}

func (this *OAuthApp) TableName() string {
	return "myzc_oauth_app"
}

// complex multile table query sample
// func (this *OAuthApp) IdColumn() string {
// 	return "a.id"
// }

// func (this *OAuthApp) SelSql() string {
// 	return "SELECT a.*,b.nick_name AS user_nick_name,c.title AS role_title FROM hx_admin_user_role AS a" +
// 		" LEFT JOIN myzc_oauth_app AS b ON a.xxx = b.id LEFT JOIN hx_admin_role AS c ON c.id = a.xxxx"
// }

func CreateSampleOAuthApp() *OAuthApp {
	return &OAuthApp{Code: "TestCode", Title: "TestTitle", AppId: "TestAppId", AppSecurityKey: "TestAppSecurityKey", AppOauthScope: "TestAppOauthScope", AppBaseUrl: "TestAppBaseUrl",
		AccessToken: "TestAccessToken"}
}
