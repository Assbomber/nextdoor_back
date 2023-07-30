package configs

import (
	"fmt"
	"os"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/spf13/viper"
)

// merging .env properties & config file properties
func Init() {

	envProperties := getConfigs(".env", "env")

	environment := envProperties.GetString("RUNTIME_ENV")
	if environment == "" {
		panic("'RUNTIME_ENV' property not set in environment variables")
	}

	// getting config file path based on environment
	configPath, ok := constants.CONFIG_PATH_MAP[environment]
	if !ok {
		panic("Invalid value for 'RUNTIME_ENV' in environment variables")
	}

	configProperties := getConfigs(configPath, "json")

	//config Properties will overwrite env properties here
	viper.GetViper().MergeConfigMap(envProperties.AllSettings())
	viper.GetViper().MergeConfigMap(configProperties.AllSettings())
	viper.GetViper().AutomaticEnv()
}

// retrieves config properties from files
func getConfigs(path string, fileType string) *viper.Viper {
	xViper := viper.New()
	xViper.SetConfigType(fileType)

	reader, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintln(constants.FAILURE, "Error opening "+path+" file", err))
	}
	err = xViper.ReadConfig(reader)
	if err != nil {
		panic(fmt.Sprintln(constants.FAILURE, "Error reading "+path+" file", err))
	}
	xViper.AutomaticEnv()
	return xViper
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBoolean(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringMap(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}
