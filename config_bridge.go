package gostack

import "github.com/rohitdas13595/go-stack/config"

// Config returns a string config value using dotted path (requires config.Load).
func Config(path string) string {
	if config.Global() == nil {
		return ""
	}
	return config.Global().String(path)
}

// ConfigInt returns int at path.
func ConfigInt(path string) int {
	if config.Global() == nil {
		return 0
	}
	return config.Global().Int(path)
}

// ConfigBool returns bool at path.
func ConfigBool(path string) bool {
	if config.Global() == nil {
		return false
	}
	return config.Global().Bool(path)
}
