package config

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	Version           = "undefined"
	Commit            = "undefined"
	Date              = "undefined"
	CheckForUpdates   = "false"
	License           = "Apache-2.0"
	Discord           = "https://discord.gg/S4xgmRSHdC"
	ProjectName       = "Process Compose 🔥"
	RemoteProjectName = "Process Compose ⚡"

	scFiles = []string{
		"shortcuts.yaml",
		"shortcuts.yml",
	}
)

const (
	pcConfigEnv       = "PROC_COMP_CONFIG"
	LogPathEnvVarName = "PC_LOG_FILE"
	LogFileFlags      = os.O_CREATE | os.O_APPEND | os.O_WRONLY | os.O_TRUNC
	LogFileMode       = os.FileMode(0600)
	themeFileName     = "theme.yaml"
	settingsFileName  = "settings.yaml"
)

var (
	clientCommands = []string{
		"down",
		"attach",
		"process",
		"project",
	}
)

func GetLogFilePath() string {
	val, found := os.LookupEnv(LogPathEnvVarName)
	if found {
		return val
	}
	userName := getUser()
	if len(userName) != 0 {
		userName = "-" + userName
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("process-compose%s%s.log", userName, mode()))
}

func getTuiDefault() bool {
	_, found := os.LookupEnv(EnvVarNameTui)
	return !found
}

func getNoServerDefault() bool {
	_, found := os.LookupEnv(EnvVarNameNoServer)
	return found
}

func getPortDefault() int {
	val, found := os.LookupEnv(EnvVarNamePort)
	if found {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal().Err(err).Msgf("Invalid port number: %s", val)
			return DefaultPortNum
		}
		return port
	}
	return DefaultPortNum
}

func GetConfigDefault() []string {
	val, found := os.LookupEnv(EnvVarNameConfig)
	if found {
		return strings.Split(val, ",")
	}
	return []string{}
}

func procCompHome() string {
	if env := os.Getenv(pcConfigEnv); env != "" {
		return env
	}
	xdgPcHome, err := xdg.ConfigFile("process-compose")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create configuration directory for process compose")
	}
	return xdgPcHome
}

func GetShortCutsPath() string {
	for _, path := range scFiles {
		scPath := filepath.Join(procCompHome(), path)
		if _, err := os.Stat(scPath); err == nil {
			return scPath
		}
	}
	return ""
}

func GetThemesPath() string {
	themePath := filepath.Join(procCompHome(), themeFileName)
	return themePath
}

func GetSettingsPath() string {
	settingsPath := filepath.Join(procCompHome(), settingsFileName)
	return settingsPath
}

func getUser() string {
	usr, err := user.Current()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to retrieve user info.")
		return ""
	}
	username := usr.Username
	if strings.Contains(username, "\\") {
		parts := strings.Split(username, "\\")
		username = parts[len(parts)-1]
	}
	return username
}

func mode() string {
	if isClient() {
		return "-client"
	}
	return ""
}

func isClient() bool {
	for _, proc := range os.Args {
		for _, cmd := range clientCommands {
			if proc == cmd {
				return true
			}
		}
	}
	return false
}

func IsLogSelectionOn() bool {
	_, found := os.LookupEnv("WAYLAND_DISPLAY")
	return !found
}

func GetUnixSocketPath() string {
	val, found := os.LookupEnv(EnvVarUnixSocketPath)
	if found {
		return val
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("process-compose-%d.sock", os.Getpid()))
}
