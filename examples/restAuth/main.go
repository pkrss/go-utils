package main

import (
	"hx98/base/constants"
	"log"
	"net/http"
	"reflect"

	"github.com/pkrss/go-utils/mvc/controllers"
	simple "github.com/pkrss/go-utils/mvc/controllers/simple"
	pkRouters "github.com/pkrss/go-utils/mvc/routers"
	orm "github.com/pkrss/go-utils/orm"
	custom "github.com/pkrss/go-utils/orm/fields"
	"github.com/pkrss/go-utils/orm/pqsql"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

////////////////////////////////////////////////////////////////////////////////
// system auth part

const (
	User int = iota
	Admin
)

type UserContext struct {
	Role   int
	UserId string
}

type MyAuthImpl struct {
}

func (this *MyAuthImpl) TokenKey() string {
	return "X-PKRSS-SAMPLE"
}

func (this *MyAuthImpl) LoadTokenObj(token string) interface{} {
	return &UserContext{Role: Admin, UserId: "1"} // json.Unmarshal(redis.Get(token))
}

func (this *MyAuthImpl) SaveTokenObj(token string, obj interface{}) {
	// redis.Set(token, json.Marshal(obj))
}

func (this *MyAuthImpl) CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool {
	return userContext.(*UserContext).Role >= requiredPrivilege.(int)
}

func (this *MyAuthImpl) IsClientManagerOrSelf(userContext interface{}, targetUserId interface{}) bool {
	return userContext.(*UserContext).Role >= Admin || userContext.(*UserContext).UserId == targetUserId.(string)
}

// system auth part
////////////////////////////////////////////////////////////////////////////////

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

func CreateOAuthAppListRestController() *simple.SimpleAuthRestListCreateController {
	postStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.SimpleAuthRestListCreateParams{
		RecordModel: &OAuthApp{}, SelectListCbFun: oauthAppGetList,
		SelectPrivilege: Admin, PostPrivilege: Admin, PostStructColsParams: postStructColsParams,
	}

	return simple.CreateSimpleListRestController(&params)
}

func CreateOAuthAppRestController() *simple.SimpleAuthRestCreateController {
	putStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.SimpleAuthRestCreateParams{
		RecordModel: &OAuthApp{}, IdUrlParam: ":id", IdType: reflect.Int64,
		SelectPrivilege: Admin, PutPrivilege: Admin, PutStructColsParams: putStructColsParams, DeletePrivilege: Admin,
	}

	return simple.CreateSimpleRestController(&params)
}

func oauthAppGetList(listRawHelper *orm.ListRawHelper) error {
	pageable := listRawHelper.Pageable
	if len(pageable.Sort) == 0 {
		pageable.Sort = "-id"
	}

	listRawHelper.SetCondArrLike("q", "title")

	c := listRawHelper.UserData.(controllers.ControllerInterface)
	ac := listRawHelper.UserData.(controllers.AuthControllerInterface)

	parentId := c.GetString("parentId")
	if parentId != "" {
		pageable.CondArr["parentId"] = parentId
		listRawHelper.SetCondArrEqu("parentId", orm.Int64, "parent_id")
	}

	userId := c.GetString("userId")
	if userId != "" {
		userContext := ac.LoadUserContext().(*UserContext)
		if "me" == userId {
			userId = userContext.UserId
		}
		pageable.CondArr["userId"] = userId

		role := userContext.Role

		if role < constants.EAdmin {
			sql := `id in (` +
				`  select menu_id_list from hx_admin_role where user_role = (` +
				`    select role_id from hx_admin_user_role where user_id = ?` +
				`  )` +
				`)`

			pageable.CondArr["id_insql"] = sql
			listRawHelper.SetCondArrEqu("id_insql", orm.InSqlStrVar, userId)
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////

func main() {

	pqsql.Db = pqsql.CreatePgSql()
	orm.DefaultOrmAdapter = &pqsql.PgSqlAdapter{}

	pkRouters.AddRouter("/oauth/apps/{id:\\d+}", CreateOAuthAppListRestController())
	pkRouters.AddRouterOptSlash("/oauth/apps", CreateOAuthAppRestController())

	controllers.DefaultAuthImpl = &MyAuthImpl{}

	localAddr := "127.0.0.1:8080"
	log.Printf("Server bind in %s\n", localAddr)
	log.Fatal(http.ListenAndServe(localAddr, pkRouters.GetApp()))

}
