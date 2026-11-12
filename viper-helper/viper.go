package viperhelper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFileType      = "yaml"
	defaultConfigFileName      = "config"
	defaultConfigFileLocations = "./,./config/"
	defaultEnvironmentPrefix   = "VIPERDEMO"
	defaultKeyReplacer         = ".,_,-,_"
)

var (
	mu       = sync.Mutex{}
	instance *viperHelper
)

type viperHelper struct {
	instance          *viper.Viper
	fileType          string
	fileName          string
	fileLocations     []string
	keyReplacer       []string
	environmentPrefix string
}

func getViperHelper() *viperHelper {
	//ensure that the viperHelper is initialized before returning it
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		newViperHelper("", "", "", "", "")
	}
	return instance
}

// InitViperHelper allows the user to override the default settings
// calling this is optional, if not called, the default settings will be used
func InitViperHelper(fileType, fileName, fileLocations, keyReplacer, environmentPrefix string) error {
	if instance != nil {
		return fmt.Errorf("viperHelper already initialized")
	} else {
		mu.Lock()
		defer mu.Unlock()
		instance = newViperHelper(fileType, fileName, fileLocations, keyReplacer, environmentPrefix)
	}
	return nil
}

func newViperHelper(fileType, fileName, fileLocations, keyReplacer, environmentPrefix string) *viperHelper {

	if instance == nil {
		fl := stringOrDefault(fileLocations, defaultConfigFileLocations)
		kr := stringOrDefault(keyReplacer, defaultKeyReplacer)

		instance = &viperHelper{
			instance:          viper.New(),
			fileType:          stringOrDefault(fileType, defaultConfigFileType),
			fileName:          stringOrDefault(fileName, defaultConfigFileName),
			fileLocations:     strings.Split(fl, ","),
			environmentPrefix: stringOrDefault(environmentPrefix, defaultEnvironmentPrefix),
			keyReplacer:       strings.Split(kr, ","),
		}
	}
	return nil
}

func stringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// BindFlag binds a flag to the specified flag set as well as viper and sets the default value
func BindFlag(flagset *pflag.FlagSet, name, group, shorthand, def, usage string) error {
	//add the flag to the flagset
	flagset.StringP(name, shorthand, def, usage)

	flagName := name
	if group != "" {
		flagName = fmt.Sprintf("%s.%s", group, name)
	}

	// bind the flag to viper
	if err := getViperHelper().instance.BindPFlag(flagName, flagset.Lookup(name)); err != nil {
		return fmt.Errorf("Failed to bind flag: %s", err.Error())
	}
	// also set the default value in viper, so that if the flag is not set, viper will use the default value
	getViperHelper().instance.SetDefault(flagName, def)

	return nil
}

// GetSettings returns a settings struct populated with values from config file, env vars and flags.
// Precedence is: flags > env vars > config file > defaults
func GetSettings[k interface{}](config *k) (*k, error) {

	vh := getViperHelper()
	keyReplacer := strings.NewReplacer(vh.keyReplacer...)
	vh.instance.SetEnvPrefix(vh.environmentPrefix)
	vh.instance.SetEnvKeyReplacer(keyReplacer)
	vh.instance.AutomaticEnv()

	for _, loc := range vh.fileLocations {
		vh.instance.AddConfigPath(loc)
	}
	vh.instance.SetConfigName(vh.fileName)
	vh.instance.SetConfigType(vh.fileType)

	err := vh.instance.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// do nothing here as the config file is optional
		} else {
			return nil, fmt.Errorf("Warning: failed to read config file: [%+v] [%T]", err, err)
		}

	}

	if err := vh.instance.Unmarshal(config); err != nil {
		// reduce to error once testing is done
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
