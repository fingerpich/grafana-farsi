package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
)

func StarDashboard(c *middleware.Context) Response {
	if !c.IsSignedIn {
		return ApiError(412, "برای ستاره کردن داشبورد نیاز است تا وارد سیستم شوید", nil)
	}

	cmd := m.StarDashboardCommand{UserId: c.UserId, DashboardId: c.ParamsInt64(":id")}

	if cmd.DashboardId <= 0 {
		return ApiError(400, "آی دی داشبورد گم شده", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "خطا در ستاره کردن داشبورد", err)
	}

	return ApiSuccess("داشبورد ستاره شد!")
}

func UnstarDashboard(c *middleware.Context) Response {

	cmd := m.UnstarDashboardCommand{UserId: c.UserId, DashboardId: c.ParamsInt64(":id")}

	if cmd.DashboardId <= 0 {
		return ApiError(400, "آی دی داشبورد گم شده", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "خطا در غیر ستاره کردن داشبورد", err)
	}

	return ApiSuccess("داشبورد از حالت ستاره خارج شد")
}
