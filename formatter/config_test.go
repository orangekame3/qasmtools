package formatter

import (
	"testing"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default_config",
			config: Config{
				Write:    false,
				Check:    false,
				Indent:   2,
				Newline:  true,
				Verbose:  false,
				Diff:     false,
				Unescape: false,
			},
		},
		{
			name: "custom_config",
			config: Config{
				Write:    true,
				Check:    true,
				Indent:   4,
				Newline:  false,
				Verbose:  true,
				Diff:     true,
				Unescape: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the config struct can be created and accessed
			c := tt.config
			
			// Verify all fields are accessible
			_ = c.Write
			_ = c.Check
			_ = c.Indent
			_ = c.Newline
			_ = c.Verbose
			_ = c.Diff
			_ = c.Unescape
			
			// Test specific values
			if c.Unescape != tt.config.Unescape {
				t.Errorf("Config.Unescape = %v, want %v", c.Unescape, tt.config.Unescape)
			}
			
			if c.Indent != tt.config.Indent {
				t.Errorf("Config.Indent = %v, want %v", c.Indent, tt.config.Indent)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	// Test zero value of Config struct
	var config Config
	
	// All bool fields should be false by default
	if config.Write != false ||
		config.Check != false ||
		config.Verbose != false ||
		config.Diff != false ||
		config.Unescape != false {
		t.Errorf("Config zero value has unexpected non-false bool fields: %+v", config)
	}
	
	// Numeric fields should be zero
	if config.Indent != 0 {
		t.Errorf("Config.Indent zero value = %v, want 0", config.Indent)
	}
	
	// Newline defaults to false (Go zero value)
	if config.Newline != false {
		t.Errorf("Config.Newline zero value = %v, want false", config.Newline)
	}
}