package common

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	DATABASE_URL   = ""
	WsAddr         = "0.0.0.0:9998"
	CurrentVersion = ""
	IsDll          = false
)

type EnvConfigSpec struct {
	Debug            bool   `envconfig:"DEBUG" default:"false"`
	DebugRequests    bool   `envconfig:"DEBUG_REQUESTS" default:"false"`
	DebugSQL         bool   `envconfig:"DEBUG_SQL" default:"false"`
	DebugWS          bool   `envconfig:"DEBUG_WS" default:"false"`
	UIUsername       string `envconfig:"UI_USERNAME" required:"false"`
	UIPassword       string `envconfig:"UI_PASSWORD" required:"false"`
	DisableAnalytics bool   `envconfig:"DISABLE_ANALYTICS" default:"false"`
	DatabaseURL      string `envconfig:"DATABASE_URL" required:"false" default:""`
}

var EnvConfig EnvConfigSpec

func init() {
	godotenv.Load()

	err := envconfig.Process("", &EnvConfig)
	if err != nil {
		Log.Fatal(err.Error())
	}
}

func FileNameAsInFS(fileName string) string {
	if IsDll {
		fileName = strings.ReplaceAll(fileName, "_MKX200_FB360.mp4", "fefb.dll")
		fileName = strings.ReplaceAll(fileName, "_MKX200.mp4", "fe.dll")
		fileName = strings.ReplaceAll(fileName, "_sbs_180.mp4", ".dll")
	}
	return fileName
}

func FileNameAsInDB(fileName string) string {
	if IsDll {
		fileName = strings.ReplaceAll(fileName, "fefb.dll", "_MKX200_FB360.mp4")
		fileName = strings.ReplaceAll(fileName, "fe.dll", "_MKX200.mp4")
		fileName = strings.ReplaceAll(fileName, ".dll", "_sbs_180.mp4")
	}
	return fileName
}
