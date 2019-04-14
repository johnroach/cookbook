package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *viper.Viper

type UTCFormatter struct {
	log.Formatter
}

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) error{
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	config.AddConfigPath("/config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Errorf("error on parsing configuration file for %v", env)
		return err
	}

	// We allow override via environment variables
	// The environment variables need to be prefixed with COOKBOOK_
	// so a FIREBASE_HOST config will need to be an environment variable of COOKBOOK_FIREBASE_HOST
	// You should be able to retrieve these keys via
	//  c := config.GetConfig()
	//  c.Get("FIREBASE_HOST")

	config.SetEnvPrefix("cookbook")
	config.AutomaticEnv()

	if config.Get("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	if config.Get("DEPLOYMENT_TYPE") == "release" {
		log.SetFormatter(UTCFormatter{&log.JSONFormatter{
			DisableTimestamp: false,
		}})
	} else {
		log.SetFormatter(UTCFormatter{&log.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp: true,
		}})
	}



	return nil
}

// GetConfig returns the configurations set
func GetConfig() *viper.Viper {
	return config
}

// Format formats the time for UTC
func (u UTCFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}