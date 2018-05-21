package auth

import "github.com/pkrss/go-utils/examples/mvc/complex/models"

func getUserContextRole(userContext *models.UserContext) int {
	if userContext == nil {
		return EUnAct
	}

	return userContext.Role
}

func CheckUserPrivilege(userContext *models.UserContext, requiredPrivilege int) bool {
	if getUserContextRole(userContext) >= requiredPrivilege {
		return true
	}

	return false
}

func IsGuest(userContext *models.UserContext) bool {
	return CheckUserPrivilege(userContext, EGuest)
}

func IsUser(userContext *models.UserContext) bool {
	return CheckUserPrivilege(userContext, EUser)
}

func IsRobot(userContext *models.UserContext) bool {
	return CheckUserPrivilege(userContext, ERobot)
}

func IsClientManager(userContext *models.UserContext) bool {
	return CheckUserPrivilege(userContext, EOperManagerStart)
}

func IsClientManagerOrSelf(userContext *models.UserContext, targetUserId int64) bool {
	if userContext == nil {
		return false
	}
	if targetUserId == userContext.UserId {
		return true
	}
	return CheckUserPrivilege(userContext, EOperManagerStart)
}

func IsAdmin(userContext *models.UserContext) bool {
	return CheckUserPrivilege(userContext, EAdmin)
}
