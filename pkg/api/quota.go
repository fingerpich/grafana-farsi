package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
	"github.com/fingerpich/grafana-farsi/pkg/setting"
)

func GetOrgQuotas(c *middleware.Context) Response {
	if !setting.Quota.Enabled {
		return ApiError(404, "سهمیه فعال نیست", nil)
	}
	query := m.GetOrgQuotasQuery{OrgId: c.ParamsInt64(":orgId")}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "خطا در دریافت سهمیه کاربر", err)
	}

	return Json(200, query.Result)
}

func UpdateOrgQuota(c *middleware.Context, cmd m.UpdateOrgQuotaCmd) Response {
	if !setting.Quota.Enabled {
		return ApiError(404, "سهمیه فعال نیست", nil)
	}
	cmd.OrgId = c.ParamsInt64(":orgId")
	cmd.Target = c.Params(":target")

	if _, ok := setting.Quota.Org.ToMap()[cmd.Target]; !ok {
		return ApiError(404, "سهمیه مقصد نامعتبر است", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "خطا در بروزرسانی سهمیه سازمان", err)
	}
	return ApiSuccess("سهمیه سازمان بروزرسانی شد")
}

func GetUserQuotas(c *middleware.Context) Response {
	if !setting.Quota.Enabled {
		return ApiError(404, "سهمیه فعال نیست", nil)
	}
	query := m.GetUserQuotasQuery{UserId: c.ParamsInt64(":id")}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "خطا در دریافت سهمیه سازمان", err)
	}

	return Json(200, query.Result)
}

func UpdateUserQuota(c *middleware.Context, cmd m.UpdateUserQuotaCmd) Response {
	if !setting.Quota.Enabled {
		return ApiError(404, "سهمیه فعال نیست", nil)
	}
	cmd.UserId = c.ParamsInt64(":id")
	cmd.Target = c.Params(":target")

	if _, ok := setting.Quota.User.ToMap()[cmd.Target]; !ok {
		return ApiError(404, "سهمیه مقصد نامعتبر است", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "خطا در بروزرسانی سهمیه سازمان", err)
	}
	return ApiSuccess("سهمیه سازمان بروزرسانی شد")
}
