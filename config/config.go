package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var global *Store

// Store holds merged configuration as nested maps.
type Store struct {
	data map[string]any
}

// Load reads YAML files in order and merges them; later files override.
// Placeholders ${VAR} and ${VAR:default} are expanded from the environment.
func Load(files ...string) (*Store, error) {
	merged := map[string]any{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		var doc map[string]any
		if err := yaml.Unmarshal(b, &doc); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		expandEnvInPlace(doc)
		mergeMaps(merged, doc)
	}
	s := &Store{data: merged}
	global = s
	return s, nil
}

// SetGlobal sets the process-wide config store.
func SetGlobal(s *Store) {
	global = s
}

// Global returns the process-wide store (may be nil).
func Global() *Store {
	return global
}

var envRe = regexp.MustCompile(`\$\{([^}:]+)(?::([^}]*))?\}`)

func expandEnvInPlace(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case string:
			m[k] = expandEnvString(t)
		case map[string]any:
			expandEnvInPlace(t)
		}
	}
}

func expandEnvString(s string) string {
	return envRe.ReplaceAllStringFunc(s, func(match string) string {
		parts := envRe.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}
		key := parts[1]
		def := ""
		if len(parts) > 2 {
			def = parts[2]
		}
		if v := os.Getenv(key); v != "" {
			return v
		}
		return def
	})
}

func mergeMaps(dst, src map[string]any) {
	for k, v := range src {
		if dm, ok := dst[k].(map[string]any); ok {
			if sm, ok2 := v.(map[string]any); ok2 {
				mergeMaps(dm, sm)
				continue
			}
		}
		dst[k] = v
	}
}

// Get returns a dotted path value as any.
func (s *Store) Get(path string) (any, bool) {
	if s == nil {
		return nil, false
	}
	parts := strings.Split(path, ".")
	var cur any = s.data
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		cur, ok = m[p]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}

// String returns string at path.
func (s *Store) String(path string) string {
	v, ok := s.Get(path)
	if !ok {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

// Int returns int at path.
func (s *Store) Int(path string) int {
	v, ok := s.Get(path)
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(t)
		return n
	default:
		return 0
	}
}

// Bool returns bool at path.
func (s *Store) Bool(path string) bool {
	v, ok := s.Get(path)
	if !ok {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(t, "true") || t == "1"
	default:
		return false
	}
}
