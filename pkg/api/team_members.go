package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
	"github.com/fingerpich/grafana-farsi/pkg/util"
)

// GET /api/teams/:teamId/members
func GetTeamMembers(c *middleware.Context) Response {
	query := m.GetTeamMembersQuery{TeamId: c.ParamsInt64(":teamId")}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "خطا در دریافت اعضای تیم", err)
	}

	for _, member := range query.Result {
		member.AvatarUrl = dtos.GetGravatarUrl(member.Email)
	}

	return Json(200, query.Result)
}

// POST /api/teams/:teamId/members
func AddTeamMember(c *middleware.Context, cmd m.AddTeamMemberCommand) Response {
	cmd.TeamId = c.ParamsInt64(":teamId")
	cmd.OrgId = c.OrgId

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrTeamMemberAlreadyAdded {
			return ApiError(400, "کاربر قبلا به این تیم ملحق شده است", err)
		}
		return ApiError(500, "خطا در افزودن کاربر به تیم", err)
	}

	return Json(200, &util.DynMap{
		"message": "عضوی به تیم اضافه گردید",
	})
}

// DELETE /api/teams/:teamId/members/:userId
func RemoveTeamMember(c *middleware.Context) Response {
	if err := bus.Dispatch(&m.RemoveTeamMemberCommand{TeamId: c.ParamsInt64(":teamId"), UserId: c.ParamsInt64(":userId")}); err != nil {
		return ApiError(500, "خطا در حذف عضو از تیم", err)
	}
	return ApiSuccess("اعضای تیم حذف شدند")
}
