package controllers

import (
	"reflect"

	"github.com/pkrss/go-utils/examples/restAuth/auth"
	"github.com/pkrss/go-utils/examples/restAuth/models"
	base "github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/mvc/controllers/simple"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

func CreateOAuthAppListRestController() base.ControllerInterface {
	postStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.ListRestParams{
		RecordModel: &models.OAuthApp{}, SelectListCbFun: oauthAppGetList,
		SelectPrivilege: auth.Admin, PostPrivilege: auth.Admin, PostStructColsParams: postStructColsParams,
	}

	return simple.CreateListRestController(&params)
}

func CreateOAuthAppRestController() base.ControllerInterface {
	putStructColsParams := &pkReflect.StructSelCols{ExcludeCols: []string{"Id", "CreateTime"}}
	params := simple.ItemRestParams{
		RecordModel: &models.OAuthApp{}, IdUrlParam: ":id", IdType: reflect.Int64,
		SelectPrivilege: auth.Admin, PutPrivilege: auth.Admin, PutStructColsParams: putStructColsParams, DeletePrivilege: auth.Admin,
	}

	return simple.CreateItemRestController(&params)
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
		userContext := ac.LoadUserContext().(*auth.UserContext)
		if "me" == userId {
			userId = userContext.UserId
		}
		pageable.CondArr["userId"] = userId

		role := userContext.Role

		if role < auth.Admin {
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
