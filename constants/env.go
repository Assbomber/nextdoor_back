package constants

var Environments = struct {
	PRODUCTION  string
	DEVELOPMENT string
}{
	PRODUCTION:  "production",
	DEVELOPMENT: "development",
}

var CONFIG_PATH_MAP = map[string]string{
	"development": "./configs/development.json",
	"production":  "./configs/production.json",
}
