package configs

import "fmt"

// MigrationDSN returns the database URL in format required by golang-migrate.
func (d DatabaseConfig) MigrationDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia%%2FTehran",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}
