package api

import (
	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/components/apikeygen"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
)

func GetApiKeys(c *middleware.Context) Response {
	query := m.GetApiKeysQuery{OrgId: c.OrgId}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "شکست در لیست کردن کلید های API", err)
	}

	result := make([]*m.ApiKeyDTO, len(query.Result))
	for i, t := range query.Result {
		result[i] = &m.ApiKeyDTO{
			Id:   t.Id,
			Name: t.Name,
			Role: t.Role,
		}
	}

	return Json(200, result)
}

func DeleteApiKey(c *middleware.Context) Response {
	id := c.ParamsInt64(":id")

	cmd := &m.DeleteApiKeyCommand{Id: id, OrgId: c.OrgId}

	err := bus.Dispatch(cmd)
	if err != nil {
		return ApiError(500, "شکست در حذف کلید API", err)
	}

	return ApiSuccess("کلید API حذف شد.")
}

func AddApiKey(c *middleware.Context, cmd m.AddApiKeyCommand) Response {
	if !cmd.Role.IsValid() {
		return ApiError(400, "نقش مشخص شده مجاز نمی‌باشد", nil)
	}

	cmd.OrgId = c.OrgId

	newKeyInfo := apikeygen.New(cmd.OrgId, cmd.Name)
	cmd.Key = newKeyInfo.HashedKey

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "شکست در اضافه کردن کلید API", err)
	}

	result := &dtos.NewApiKeyResult{
		Name: cmd.Result.Name,
		Key:  newKeyInfo.ClientSecret}

	return Json(200, result)
}
