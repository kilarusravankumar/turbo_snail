package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	TCP_PORT  = "7777"
	HTTP_PORT = "7000"
	WAL_DIR   = "/data/turboSnail/wal_logs"
)

func Init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error occured while loading Environment Variables \n %s", err.Error())
		log.Printf("Defaulting to Default Configs.")
	}
	if val, ok := isEnvVarPresent("START_LINE_PORT"); ok {
		TCP_PORT = val
	}

	if val, ok := isEnvVarPresent("FINISH_LINE_PORT"); ok {
		HTTP_PORT = val
	}

	if val, ok := isEnvVarPresent("WAL_DIR"); ok {
		WAL_DIR = val
	}
	log.Printf("Write ahead logs directory is set to : %s \n", WAL_DIR)

}

func isEnvVarPresent(str string) (string, bool) {
	val := os.Getenv(str)
	if strings.Trim(val, " ") == "" {
		return "", false
	}
	return val, true
}
