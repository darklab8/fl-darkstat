package router

var fixture_app_data *AppData

func GetAppDataFixture() *AppData {
	if fixture_app_data == nil {
		fixture_app_data = NewAppData()
	}

	return fixture_app_data
}
