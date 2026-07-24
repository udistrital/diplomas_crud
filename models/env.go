package models

import (
	"bufio"
	"os"
	"strings"

	"github.com/astaxie/beego"
)

// LoadDotEnv loads simple KEY=VALUE pairs from .env into process env
// without overriding variables already exported in the shell.
func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}

func SyncEnvToAppConfig() {
	envToConf := map[string]string{
		"DIPLOMAS_CRUD_HTTP_PORT":  "httpport",
		"DIPLOMAS_CRUD_RUNMODE":    "runmode",
		"DIPLOMAS_CRUD_PGUSER":     "PGuser",
		"DIPLOMAS_CRUD_PGPASS":     "PGpass",
		"DIPLOMAS_CRUD_PGHOST":     "PGhost",
		"DIPLOMAS_CRUD_PGPORT":     "PGport",
		"DIPLOMAS_CRUD_PGDB":       "PGdb",
		"DIPLOMAS_CRUD_PGSCHEMA":   "PGschema",
		"DIPLOMAS_CRUD_PGSSLMODE":  "PGsslmode",
		"DIPLOMAS_CRUD_PGTIMEZONE": "PGtimezone",
		"PARAMETER_STORE":          "parameterStore",
	}

	for envKey, confKey := range envToConf {
		if value := os.Getenv(envKey); value != "" {
			_ = beego.AppConfig.Set(confKey, value)
		}
	}
}
