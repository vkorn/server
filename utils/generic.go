package utils

import (
	"os"
	"strings"
	"time"

	"github.com/go-home-io/server/plugins/device/enums"
)

// SliceContainsString slice.contains implementation for strings.
func SliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// TimeNow returns epoch UTC.
func TimeNow() int64 {
	return time.Now().UTC().Unix()
}

// VerifyDeviceProvider transforms device provider from yaml config into actual type.
func VerifyDeviceProvider(configType string) enums.DeviceType {
	parts := strings.SplitN(configType, "/", 2)
	if len(parts) < 2 {
		return enums.DevUnknown
	}

	t, err := enums.DeviceTypeString(parts[0])
	if err != nil {
		return enums.DevUnknown
	}

	return t
}

// NormalizeDeviceName validates that final device name is correct.
func NormalizeDeviceName(raw string) string {
	raw = strings.ToLower(raw)
	replacer := strings.NewReplacer("%", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		";", "_",
		".", "_",
		"-", "_",
		" ", "_")
	return replacer.Replace(raw)
}

// GetCurrentWorkingDir returns application working directory.
func GetCurrentWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Failed to get current working dir")
	}

	return cwd
}
