package multikf

import (
	"time"

	"github.com/spf13/viper"
)

type viperConfigKey string

func (t viperConfigKey) String() string {
	return string(t)
}

func (t viperConfigKey) GetString() string {
	val := viper.Get(t.String())
	_, isStr := val.(string)
	if !isStr {
		return ""
	}
	return val.(string)
}

func (t viperConfigKey) GetDuration() time.Duration {
	val := viper.Get(t.String())
	_, isDur := val.(time.Duration)
	if !isDur {
		return time.Duration(0)
	}
	return val.(time.Duration)
}

func (t viperConfigKey) GetBool() bool {
	val := viper.Get(t.String())
	_, isBool := val.(bool)
	if !isBool {
		return false
	}
	return val.(bool)
}

func (t viperConfigKey) Get() interface{} {
	return viper.Get(t.String())
}

func (t viperConfigKey) Set(val interface{}) {
	viper.Set(t.String(), val)
}

func (t viperConfigKey) setDefault(val interface{}) {
	viper.SetDefault(t.String(), val)
}

var (
	viperConfigKeyRootDir viperConfigKey = "rootDir"
	viperConfigKeyVerbose viperConfigKey = "verbose"
)
