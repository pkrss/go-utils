package models

import (
	orm "github.com/pkrss/go-utils/orm"
	custom "github.com/pkrss/go-utils/orm/fields"
)

type OAuthApp struct {
	orm.BaseModel
	Id                    int             `json:"id"`
	Code                  string          `json:"code"`
	Title                 string          `json:"title"`
	AppId                 string          `json:"appId"`
	AppSecurityKey        string          `json:"appSecurityKey"`
	AppOauthScope         string          `json:"appOauthScope"`
	AppBaseUrl            string          `json:"appBaseUrl"`
	AccessToken           string          `json:"accessToken"`
	AccessTokenExpireTime custom.JsonTime `json:"accessTokenExpireTime"`
	CreateTime            custom.JsonTime `json:"createTime"`
}
