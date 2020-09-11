package pkg

import (
	"os"
	"strconv"
)

type DBConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func dbConfiguration(prefix string) DBConfiguration {
	s_port := os.Getenv(prefix + "DB__PORT")
	port, err := strconv.Atoi(s_port)
	if err != nil {
		port = 5432
	}
	return DBConfiguration{
		Host:     os.Getenv(prefix + "DB__HOST"),
		Port:     port,
		User:     os.Getenv(prefix + "DB__USER"),
		Password: os.Getenv(prefix + "DB__PASSWORD"),
		DbName:   os.Getenv(prefix + "DB__DBNAME"),
	}
}

func DBConfigurationFromEnv() DBConfiguration                        { return dbConfiguration("") }
func DBConfigurationFromEnvWithPrefix(prefix string) DBConfiguration { return dbConfiguration(prefix) }
