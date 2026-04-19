// Package main provides a CLI tool for seeding the Super Admin account.
// Usage: go run cmd/seed/main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO: Implement Super Admin seeding (Plan 01-04)
	// Will read ADMIN_EMAIL and ADMIN_PASSWORD from environment,
	// hash the password with bcrypt cost 12, and insert into users table.
	fmt.Fprintln(os.Stderr, "seed: not yet implemented — see Plan 01-04")
	os.Exit(1)
}
