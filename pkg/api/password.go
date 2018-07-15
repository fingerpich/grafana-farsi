package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
	"github.com/fingerpich/grafana-farsi/pkg/util"
)

func SendResetPasswordEmail(c *middleware.Context, form dtos.SendResetPasswordEmailForm) Response {
	userQuery := m.GetUserByLoginQuery{LoginOrEmail: form.UserOrEmail}

	if err := bus.Dispatch(&userQuery); err != nil {
		c.Logger.Info("Requested password reset for user that was not found", "user", userQuery.LoginOrEmail)
		return ApiError(200, "ایمیل ارسال شد", err)
	}

	emailCmd := m.SendResetPasswordEmailCommand{User: userQuery.Result}
	if err := bus.Dispatch(&emailCmd); err != nil {
		return ApiError(500, "خطا در ارسال ایمیل", err)
	}

	return ApiSuccess("ایمیل ارسال شد")
}

func ResetPassword(c *middleware.Context, form dtos.ResetUserPasswordForm) Response {
	query := m.ValidateResetPasswordCodeQuery{Code: form.Code}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrInvalidEmailCode {
			return ApiError(400, "کد ریست رمز عبور نامعتبر است", nil)
		}
		return ApiError(500, "خطا در اعتبار سنجی کد ایمیل", err)
	}

	if form.NewPassword != form.ConfirmPassword {
		return ApiError(400, "رمز های عبور منطبق نیستند", nil)
	}

	cmd := m.ChangeUserPasswordCommand{}
	cmd.UserId = query.Result.Id
	cmd.NewPassword = util.EncodePassword(form.NewPassword, query.Result.Salt)

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "خطا در تغییر رمز عبور", err)
	}

	return ApiSuccess("رمز عبور کاربر تغییر یافت")
}
