package router

import (
	"context"

	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

var fixture_app_data *appdata.AppData

func GetAppDataFixture() *appdata.AppData {
	if fixture_app_data == nil {
		fixture_app_data = appdata.NewAppData(context.Background())
	}

	return fixture_app_data
}
