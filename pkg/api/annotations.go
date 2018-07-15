package api

import (
	"strings"
	"time"

	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/components/simplejson"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	"github.com/fingerpich/grafana-farsi/pkg/services/annotations"
	"github.com/fingerpich/grafana-farsi/pkg/util"
)

func GetAnnotations(c *middleware.Context) Response {

	query := &annotations.ItemQuery{
		From:        c.QueryInt64("from") / 1000,
		To:          c.QueryInt64("to") / 1000,
		OrgId:       c.OrgId,
		AlertId:     c.QueryInt64("alertId"),
		DashboardId: c.QueryInt64("dashboardId"),
		PanelId:     c.QueryInt64("panelId"),
		Limit:       c.QueryInt64("limit"),
		Tags:        c.QueryStrings("tags"),
		Type:        c.Query("type"),
	}

	repo := annotations.GetRepository()

	items, err := repo.Find(query)
	if err != nil {
		return ApiError(500, "شکست در دریافت یادداشت", err)
	}

	for _, item := range items {
		if item.Email != "" {
			item.AvatarUrl = dtos.GetGravatarUrl(item.Email)
		}
		item.Time = item.Time * 1000
	}

	return Json(200, items)
}

type CreateAnnotationError struct {
	message string
}

func (e *CreateAnnotationError) Error() string {
	return e.message
}

func PostAnnotation(c *middleware.Context, cmd dtos.PostAnnotationsCmd) Response {
	repo := annotations.GetRepository()

	if cmd.Text == "" {
		err := &CreateAnnotationError{"text field should not be empty"}
		return ApiError(500, "شکست در ذخیره یادداشت", err)
	}

	item := annotations.Item{
		OrgId:       c.OrgId,
		UserId:      c.UserId,
		DashboardId: cmd.DashboardId,
		PanelId:     cmd.PanelId,
		Epoch:       cmd.Time / 1000,
		Text:        cmd.Text,
		Data:        cmd.Data,
		Tags:        cmd.Tags,
	}

	if item.Epoch == 0 {
		item.Epoch = time.Now().Unix()
	}

	if err := repo.Save(&item); err != nil {
		return ApiError(500, "شکست در ذخیره یادداشت", err)
	}

	startID := item.Id

	// handle regions
	if cmd.IsRegion {
		item.RegionId = startID

		if item.Data == nil {
			item.Data = simplejson.New()
		}

		if err := repo.Update(&item); err != nil {
			return ApiError(500, "شکست در ذخیره مشخصه منطقه یادداشت", err)
		}

		item.Id = 0
		item.Epoch = cmd.TimeEnd / 1000

		if err := repo.Save(&item); err != nil {
			return ApiError(500, "شکست در ذخیره یادداشت برای این منطقه و زمان", err)
		}

		return Json(200, util.DynMap{
			"message": "یادداشت اضافه شد",
			"id":      startID,
			"endId":   item.Id,
		})
	}

	return Json(200, util.DynMap{
		"message": "یادداشت اضافه شد",
		"id":      startID,
	})
}

func formatGraphiteAnnotation(what string, data string) string {
	text := what
	if data != "" {
		text = text + "\n" + data
	}
	return text
}

func PostGraphiteAnnotation(c *middleware.Context, cmd dtos.PostGraphiteAnnotationsCmd) Response {
	repo := annotations.GetRepository()

	if cmd.What == "" {
		err := &CreateAnnotationError{"what field should not be empty"}
		return ApiError(500, "شکست در ذخیره یادداشت Graphite", err)
	}

	if cmd.When == 0 {
		cmd.When = time.Now().Unix()
	}
	text := formatGraphiteAnnotation(cmd.What, cmd.Data)

	// Support tags in prior to Graphite 0.10.0 format (string of tags separated by space)
	var tagsArray []string
	switch tags := cmd.Tags.(type) {
	case string:
		if tags != "" {
			tagsArray = strings.Split(tags, " ")
		} else {
			tagsArray = []string{}
		}
	case []interface{}:
		for _, t := range tags {
			if tagStr, ok := t.(string); ok {
				tagsArray = append(tagsArray, tagStr)
			} else {
				err := &CreateAnnotationError{"tag should be a string"}
				return ApiError(500, "شکست در ذخیره یادداشت Graphite", err)
			}
		}
	default:
		err := &CreateAnnotationError{"unsupported tags format"}
		return ApiError(500, "شکست در ذخیره یادداشت Graphite", err)
	}

	item := annotations.Item{
		OrgId:  c.OrgId,
		UserId: c.UserId,
		Epoch:  cmd.When,
		Text:   text,
		Tags:   tagsArray,
	}

	if err := repo.Save(&item); err != nil {
		return ApiError(500, "شکست در ذخیره یادداشت Graphite", err)
	}

	return Json(200, util.DynMap{
		"message": "Graphite annotation added",
		"id":      item.Id,
	})
}

func UpdateAnnotation(c *middleware.Context, cmd dtos.UpdateAnnotationsCmd) Response {
	annotationId := c.ParamsInt64(":annotationId")

	repo := annotations.GetRepository()

	item := annotations.Item{
		OrgId:  c.OrgId,
		UserId: c.UserId,
		Id:     annotationId,
		Epoch:  cmd.Time / 1000,
		Text:   cmd.Text,
		Tags:   cmd.Tags,
	}

	if err := repo.Update(&item); err != nil {
		return ApiError(500, "شکست در بروزرسانی یادداشت", err)
	}

	if cmd.IsRegion {
		itemRight := item
		itemRight.RegionId = item.Id
		itemRight.Epoch = cmd.TimeEnd / 1000

		// We don't know id of region right event, so set it to 0 and find then using query like
		// ... WHERE region_id = <item.RegionId> AND id != <item.RegionId> ...
		itemRight.Id = 0

		if err := repo.Update(&itemRight); err != nil {
			return ApiError(500, "شکست در بروزرسانی یادداشت برای این منطقه و زمان", err)
		}
	}

	return ApiSuccess("یادداشت به روزرسانی شد")
}

func DeleteAnnotations(c *middleware.Context, cmd dtos.DeleteAnnotationsCmd) Response {
	repo := annotations.GetRepository()

	err := repo.Delete(&annotations.DeleteParams{
		AlertId:     cmd.PanelId,
		DashboardId: cmd.DashboardId,
		PanelId:     cmd.PanelId,
	})

	if err != nil {
		return ApiError(500, "شکست در حذف یادداشت", err)
	}

	return ApiSuccess("یادداشت ها حذف شد")
}

func DeleteAnnotationById(c *middleware.Context) Response {
	repo := annotations.GetRepository()
	annotationId := c.ParamsInt64(":annotationId")

	err := repo.Delete(&annotations.DeleteParams{
		Id: annotationId,
	})

	if err != nil {
		return ApiError(500, "شکست در حذف یادداشت", err)
	}

	return ApiSuccess("یادداشت حذف شد")
}

func DeleteAnnotationRegion(c *middleware.Context) Response {
	repo := annotations.GetRepository()
	regionId := c.ParamsInt64(":regionId")

	err := repo.Delete(&annotations.DeleteParams{
		RegionId: regionId,
	})

	if err != nil {
		return ApiError(500, "شکست در ذخیره منطقه یادداشت", err)
	}

	return ApiSuccess("یادداشت منطقه حذف شد")
}
