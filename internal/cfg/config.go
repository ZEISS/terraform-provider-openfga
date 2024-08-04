package cfg

// Flags contains the command line flags.
type Flags struct {
	Debug bool
}

// Config contains the configuration.
type Config struct {
	Flags *Flags
}

// New returns a new instance of Config.
func New() *Config {
	return &Config{
		Flags: &Flags{},
	}
}
