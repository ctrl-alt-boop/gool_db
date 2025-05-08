package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const bindsFileName = "binds.yaml"
const appName = "gooldb"

// FindConfigFile searches for the configuration file in common locations.
// It returns the path of the first found file, or an error if not found.
func FindBindsConfigFile(cmdLinePath, envVarPath string) (string, error) {
	// 1. Check path from command line flag
	if cmdLinePath != "" {
		if _, err := os.Stat(cmdLinePath); err == nil {
			return cmdLinePath, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("cli config path %q: %w", cmdLinePath, err)
		}
	}

	// 2. Check environment variable path
	if envVarPath != "" {
		if _, err := os.Stat(envVarPath); err == nil {
			return envVarPath, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("env var config path %q: %w", envVarPath, err)
		}
	}

	// 3. Potential user config paths using a library like adrg/xdg
	// This handles cross-platform differences automatically
	userConfigPaths := []string{
		filepath.Join(appName, bindsFileName),
	}

	// A more manual list incorporating different OS conventions:
	manualUserConfigPaths := []string{}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		// Linux/macOS ~/.config/
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(homeDir, ".config", appName, bindsFileName))
		// Linux/macOS ~/.appname/
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(homeDir, "."+appName, bindsFileName))
		// macOS Library/Application Support/
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(homeDir, "Library", "Application Support", appName, bindsFileName))
		// macOS Library/Preferences/
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(homeDir, "Library", "Preferences", appName, bindsFileName))
		// Windows AppData/Roaming
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(os.Getenv("APPDATA"), appName, bindsFileName))
		// Windows AppData/Local
		manualUserConfigPaths = append(manualUserConfigPaths, filepath.Join(os.Getenv("LOCALAPPDATA"), appName, bindsFileName))

	}

	allUserConfigPaths := userConfigPaths
	allUserConfigPaths = append(allUserConfigPaths, manualUserConfigPaths...)

	// 4. Check user configuration directories
	for _, configPath := range allUserConfigPaths {
		if configPath != "" {
			if _, err := os.Stat(configPath); err == nil {
				return configPath, nil
			} else if !os.IsNotExist(err) {
				return "", fmt.Errorf("user config path %q: %w", configPath, err)
			}
		}
	}

	// 5. Potential system config paths
	manualSystemConfigPaths := []string{
		// Linux /etc/
		filepath.Join("/etc", appName, bindsFileName),
		// macOS /Library/Application Support/
		filepath.Join("/Library", "Application Support", appName, bindsFileName),
		// Windows ProgramData
		filepath.Join(os.Getenv("ProgramData"), appName, bindsFileName),
	}

	allSystemConfigPaths := manualSystemConfigPaths

	// 6. Check system configuration directories
	for _, configPath := range allSystemConfigPaths {
		if configPath != "" {
			if _, err := os.Stat(configPath); err == nil {
				return configPath, nil
			} else if !os.IsNotExist(err) {
				return "", fmt.Errorf("system config path %q: %w", configPath, err)
			}
		}
	}

	// If no config file was found in any location
	return "", os.ErrNotExist
}
