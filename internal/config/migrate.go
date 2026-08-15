package config

type Migration func(cfg *Config)

func Migrations() []Migration {
	return []Migration{
		func(cfg *Config) {
			// first one to start with
		},
	}
}

func Migrate(cfg *Config) bool {
	migrations := Migrations()

	if cfg.Version == len(migrations) {
		return false
	}

	toApply := migrations[cfg.Version:]
	cfg.Version = len(migrations)

	for _, migrate := range toApply {
		migrate(cfg)
	}

	return true
}
