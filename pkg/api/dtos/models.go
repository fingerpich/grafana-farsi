package dtos

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"github.com/fingerpich/grafana-farsi/pkg/components/simplejson"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
	"github.com/fingerpich/grafana-farsi/pkg/setting"
)

type AnyId struct {
	Id int64 `json:"id"`
}

type LoginCommand struct {
	User     string `json:"user" binding:"Required"`
	Password string `json:"password" binding:"Required"`
	Remember bool   `json:"remember"`
}

type CurrentUser struct {
	IsSignedIn     bool         `json:"isSignedIn"`
	Id             int64        `json:"id"`
	Login          string       `json:"login"`
	Email          string       `json:"email"`
	Name           string       `json:"name"`
	LightTheme     bool         `json:"lightTheme"`
	OrgCount       int          `json:"orgCount"`
	OrgId          int64        `json:"orgId"`
	OrgName        string       `json:"orgName"`
	OrgRole        m.RoleType   `json:"orgRole"`
	IsGrafanaAdmin bool         `json:"isGrafanaAdmin"`
	GravatarUrl    string       `json:"gravatarUrl"`
	Timezone       string       `json:"timezone"`
	Locale         string       `json:"locale"`
	HelpFlags1     m.HelpFlags1 `json:"helpFlags1"`
}

type MetricRequest struct {
	From    string             `json:"from"`
	To      string             `json:"to"`
	Queries []*simplejson.Json `json:"queries"`
}

type UserStars struct {
	DashboardIds map[string]bool `json:"dashboardIds"`
}

func GetGravatarUrl(text string) string {
	if text == "" {
		return ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.ToLower(text)))
	return fmt.Sprintf(setting.AppSubUrl+"/avatar/%x", hasher.Sum(nil))
}

func GetGravatarUrlWithDefault(text string, defaultText string) string {
	if text != "" {
		return GetGravatarUrl(text)
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return ""
	}

	text = reg.ReplaceAllString(defaultText, "") + "@localhost"

	return GetGravatarUrl(text)
}
