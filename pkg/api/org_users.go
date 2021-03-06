package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
)

// POST /api/org/users
func AddOrgUserToCurrentOrg(c *middleware.Context, cmd m.AddOrgUserCommand) Response {
	cmd.OrgId = c.OrgId
	return addOrgUserHelper(cmd)
}

// POST /api/orgs/:orgId/users
func AddOrgUser(c *middleware.Context, cmd m.AddOrgUserCommand) Response {
	cmd.OrgId = c.ParamsInt64(":orgId")
	return addOrgUserHelper(cmd)
}

func addOrgUserHelper(cmd m.AddOrgUserCommand) Response {
	if !cmd.Role.IsValid() {
		return ApiError(400, "نقش نا معتبری تعیین شده است", nil)
	}

	userQuery := m.GetUserByLoginQuery{LoginOrEmail: cmd.LoginOrEmail}
	err := bus.Dispatch(&userQuery)
	if err != nil {
		return ApiError(404, "کاربر یافت نشد", nil)
	}

	userToAdd := userQuery.Result

	cmd.UserId = userToAdd.Id

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrOrgUserAlreadyAdded {
			return ApiError(409, "کاربر قبلا در این سازمان عضو بوده است", nil)
		}
		return ApiError(500, "نمیتوان کاربر را به سازمان افزود", err)
	}

	return ApiSuccess("کاربر به سازمان اضافه شد")
}

// GET /api/org/users
func GetOrgUsersForCurrentOrg(c *middleware.Context) Response {
	return getOrgUsersHelper(c.OrgId)
}

// GET /api/orgs/:orgId/users
func GetOrgUsers(c *middleware.Context) Response {
	return getOrgUsersHelper(c.ParamsInt64(":orgId"))
}

func getOrgUsersHelper(orgId int64) Response {
	query := m.GetOrgUsersQuery{OrgId: orgId}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "خطا در دریافت حساب کاربری", err)
	}

	for _, user := range query.Result {
		user.AvatarUrl = dtos.GetGravatarUrl(user.Email)
	}

	return Json(200, query.Result)
}

// PATCH /api/org/users/:userId
func UpdateOrgUserForCurrentOrg(c *middleware.Context, cmd m.UpdateOrgUserCommand) Response {
	cmd.OrgId = c.OrgId
	cmd.UserId = c.ParamsInt64(":userId")
	return updateOrgUserHelper(cmd)
}

// PATCH /api/orgs/:orgId/users/:userId
func UpdateOrgUser(c *middleware.Context, cmd m.UpdateOrgUserCommand) Response {
	cmd.OrgId = c.ParamsInt64(":orgId")
	cmd.UserId = c.ParamsInt64(":userId")
	return updateOrgUserHelper(cmd)
}

func updateOrgUserHelper(cmd m.UpdateOrgUserCommand) Response {
	if !cmd.Role.IsValid() {
		return ApiError(400, "نقش نا معتبری تعیین شده است", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrLastOrgAdmin {
			return ApiError(400, "Cannot change role so that there is no organization admin left", nil)
		}
		return ApiError(500, "خطا در بروزرسانی کاربر سازمان", err)
	}

	return ApiSuccess("سازمان کاربر به روزرسانی شد")
}

// DELETE /api/org/users/:userId
func RemoveOrgUserForCurrentOrg(c *middleware.Context) Response {
	userId := c.ParamsInt64(":userId")
	return removeOrgUserHelper(c.OrgId, userId)
}

// DELETE /api/orgs/:orgId/users/:userId
func RemoveOrgUser(c *middleware.Context) Response {
	userId := c.ParamsInt64(":userId")
	orgId := c.ParamsInt64(":orgId")
	return removeOrgUserHelper(orgId, userId)
}

func removeOrgUserHelper(orgId int64, userId int64) Response {
	cmd := m.RemoveOrgUserCommand{OrgId: orgId, UserId: userId}

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrLastOrgAdmin {
			return ApiError(400, "آخرین سازمان ادمین حذف نمیشود", nil)
		}
		return ApiError(500, "خطا در حذف سازمان کاربر", err)
	}

	return ApiSuccess("کاربر از سازمان حذف شد")
}
