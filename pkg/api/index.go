package api

import (
	"fmt"
	"strings"

	"github.com/fingerpich/grafana-farsi/pkg/api/dtos"
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	"github.com/fingerpich/grafana-farsi/pkg/middleware"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
	"github.com/fingerpich/grafana-farsi/pkg/plugins"
	"github.com/fingerpich/grafana-farsi/pkg/setting"
)

func setIndexViewData(c *middleware.Context) (*dtos.IndexViewData, error) {
	settings, err := getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{OrgId: c.OrgId, UserId: c.UserId}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appUrl := setting.AppUrl
	appSubUrl := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appUrl = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubUrl = ""
		settings["appSubUrl"] = ""
	}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:             c.UserId,
			IsSignedIn:     c.IsSignedIn,
			Login:          c.Login,
			Email:          c.Email,
			Name:           c.Name,
			OrgCount:       c.OrgCount,
			OrgId:          c.OrgId,
			OrgName:        c.OrgName,
			OrgRole:        c.OrgRole,
			GravatarUrl:    dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin: c.IsGrafanaAdmin,
			LightTheme:     prefs.Theme == "light",
			Timezone:       prefs.Timezone,
			Locale:         locale,
			HelpFlags1:     c.HelpFlags1,
		},
		Settings:                settings,
		Theme:                   prefs.Theme,
		AppUrl:                  appUrl,
		AppSubUrl:               appSubUrl,
		GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/transparent.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeUrlParam := c.Query("theme")
	if themeUrlParam == "light" {
		data.User.LightTheme = true
		data.Theme = "light"
	}

	if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text: "ایجاد",
			Id:   "create",
			Icon: "fa fa-fw fa-plus",
			Url:  setting.AppSubUrl + "/dashboard/new",
			Children: []*dtos.NavLink{
				{Text: "داشبورد", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
				{Text: "پوشه", SubTitle: "ایجاد یک پوشه برای مرتب سازی داشبورد", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"},
				{Text: "دریافت", SubTitle: "وارد کردن داشبورد از فایل ذخیره شده", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"},
			},
		})
	}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "خانه", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		{Divider: true, HideFromTabs: true},
		{Text: "مدیریت", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		{Text: "پلی لیست", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		{Text: "اسنپ شات", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "داشبورد ها",
		Id:       "dashboards",
		SubTitle: "مدیریت داشبوردها و فودلدرها",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	if c.IsSignedIn {
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     c.SignedInUser.Login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			Children: []*dtos.NavLink{
				{Text: "ترجیحات", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
				{Text: "تغییر رمز عبور", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "خروج", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		alertChildNavs := []*dtos.NavLink{
			{Text: "قوانین هشدار", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
			{Text: "کانال نوفیکیشن ها", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
		}

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "هشداردهی",
			SubTitle: "قوانین هشداردهی و نوتیفیکیشن ها",
			Id:       "alerting",
			Icon:     "gicon gicon-alert",
			Url:      setting.AppSubUrl + "/alerting/list",
			Children: alertChildNavs,
		})
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}

			for _, include := range plugin.Includes {
				if !c.HasUserRole(include.Role) {
					continue
				}

				if include.Type == "page" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}

				if include.Type == "dashboard" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}
			}

			if len(appLink.Children) > 0 && c.OrgRole == m.ROLE_ADMIN {
				appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
				appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "تنظیمات پلاگین", Icon: "gicon gicon-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
			}

			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}

	if c.OrgRole == m.ROLE_ADMIN {
		cfgNode := &dtos.NavLink{
			Id:       "cfg",
			Text:     "تنظیمات",
			SubTitle: "تنظیمات: " + c.OrgName,
			Icon:     "gicon gicon-cog",
			Url:      setting.AppSubUrl + "/datasources",
			Children: []*dtos.NavLink{
				{
					Text:        "منبع داده",
					Icon:        "gicon gicon-datasources",
					Description: "اضافه کردن و تنظیم منبع داده",
					Id:          "datasources",
					Url:         setting.AppSubUrl + "/datasources",
				},
				{
					Text:        "کاربران",
					Id:          "users",
					Description: "مدیریت کاربران",
					Icon:        "gicon gicon-user",
					Url:         setting.AppSubUrl + "/org/users",
				},
				{
					Text:        "تیم ها",
					Id:          "teams",
					Description: "مدیریت گروه ها",
					Icon:        "gicon gicon-team",
					Url:         setting.AppSubUrl + "/org/teams",
				},
				{
					Text:        "پلاگین ها",
					Id:          "plugins",
					Description: "نمایش و تنظیمات پلاگین ها",
					Icon:        "gicon gicon-plugins",
					Url:         setting.AppSubUrl + "/plugins",
				},
				{
					Text:        "اولویت ها",
					Id:          "org-settings",
					Description: "تنظیم اولویت ها",
					Icon:        "gicon gicon-preferences",
					Url:         setting.AppSubUrl + "/org",
				},

				{
					Text:        "کلید های API",
					Id:          "apikeys",
					Description: "ایجاد و مدیریت کلیدهای ای پی آی",
					Icon:        "gicon gicon-apikeys",
					Url:         setting.AppSubUrl + "/org/apikeys",
				},
			},
		}

		if c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Divider: true, HideFromTabs: true,
			})
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Text:         "مدیر سرور",
				HideFromTabs: true,
				SubTitle:     "مدیریت تمام کاربران و بخش ها",
				Id:           "admin",
				Icon:         "gicon gicon-shield",
				Url:          setting.AppSubUrl + "/admin/users",
				Children: []*dtos.NavLink{
					{Text: "کاربران", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
					{Text: "بخش ها", Id: "global-orgs", Url: setting.AppSubUrl + "/admin/orgs", Icon: "gicon gicon-org"},
					{Text: "تنظیمات", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
					{Text: "وضعیت", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
					{Text: "استایل راهنما", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
				},
			})
		}

		data.NavTree = append(data.NavTree, cfgNode)
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:         "کمک",
		Id:           "help",
		Url:          "#",
		Icon:         "gicon gicon-question",
		HideFromMenu: true,
		Children: []*dtos.NavLink{
			{Text: "میانبرهای صفحه کلید", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
			{Text: "سایت کاربران", Url: "http://community.grafana.com", Icon: "fa fa-fw fa-comment", Target: "_blank"},
			{Text: "مستندات", Url: "http://docs.grafana.org", Icon: "fa fa-fw fa-file", Target: "_blank"},
		},
	})

	return &data, nil
}

func Index(c *middleware.Context) {
	if data, err := setIndexViewData(c); err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	} else {
		c.HTML(200, "index", data)
	}
}

func NotFoundHandler(c *middleware.Context) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "صفحه مورد نظر یافت نشد.", nil)
		return
	}

	if data, err := setIndexViewData(c); err != nil {
		c.Handle(500, "تنظیمات مورد نظر با دریافت نشد.", err)
		return
	} else {
		c.HTML(404, "index", data)
	}
}
