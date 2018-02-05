package controllers

import (
	"reflect"

	"github.com/pkrss/go-utils/examples/restAuth/models"
	base "github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/mvc/controllers/simple"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

func CreateOAuthAppListRestController() *simple.SimpleAuthRestListCreateController {
	postStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.SimpleAuthRestListCreateParams{
		RecordModel: &models.OAuthApp{}, SelectListCbFun: oauthAppGetList,
		SelectPrivilege: Admin, PostPrivilege: Admin, PostStructColsParams: postStructColsParams,
	}

	return simple.CreateSimpleListRestController(&params)
}

func CreateOAuthAppRestController() *simple.SimpleAuthRestCreateController {
	putStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.SimpleAuthRestCreateParams{
		RecordModel: &models.OAuthApp{}, IdUrlParam: ":id", IdType: reflect.Int64,
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

	c := listRawHelper.UserData.(base.ControllerInterface)
	ac := listRawHelper.UserData.(base.AuthControllerInterface)

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

		if role < Admin {
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
