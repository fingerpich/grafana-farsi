package sqlstore

import (
	"github.com/fingerpich/grafana-farsi/pkg/bus"
	m "github.com/fingerpich/grafana-farsi/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}
