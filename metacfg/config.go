package metacfg

import (
	"fmt"
	"strings"
)

// Config represents the configuration loader
type Config struct {
	options *Options
	data    map[string]any
}

// New creates a new Config instance with the given options
func New(opts ...Option) (*Config, error) {
	c := &Config{
		options: &Options{
			IgnoreCase:  true,
			TagName:     "meta",
			UseDefaults: true,
			UseEnv:      true,
		},
		data: make(map[string]any),
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			if opt != nil {
				opt(c)
			}
		}
	}

	return c, nil
}

// Load is a global function to create a new Config and load configuration from a file
func Load(filename string, v any, opts ...Option) error {
	c, err := New(opts...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	if err := c.LoadFile(filename); err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	if c.options.IgnoreCase {
		c.buildCaseInsensitiveCache()
	}

	if err := c.Parse(v); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// buildCaseInsensitiveCache rebuilds the case-insensitive cache
func (c *Config) buildCaseInsensitiveCache() {
	c.data = buildCaseInsensitiveMap(c.data)
}

func buildCaseInsensitiveMap(data map[string]any) map[string]any {
	caseInsensitive := make(map[string]any)
	for key, val := range data {
		lowerKey := strings.ToLower(key)
		switch v := val.(type) {
		case map[string]any:
			caseInsensitive[lowerKey] = buildCaseInsensitiveMap(v)
		default:
			caseInsensitive[lowerKey] = v
		}
	}
	return caseInsensitive
}
