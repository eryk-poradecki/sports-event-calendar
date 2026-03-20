package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const dbDriverName = "postgres"
const migrationsPath = "/migrations"

func newMigrator(db *sql.DB, migrationsPath string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsPath), migrationsPath, driver)
	if err != nil {
		return nil, fmt.Errorf("could not create migration instance: %w", err)
	}

	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	sourceErr, dbErr := m.Close()
	if sourceErr != nil || dbErr != nil {
		log.Printf("warning: failed to close migration instance: source=%v db=%v", sourceErr, dbErr)
	}
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	m, err := newMigrator(db, migrationsPath)
	if err != nil {
		return err
	}

	defer closeMigrator(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations ran successfully")
	return nil
}

func rollbackMigrations(db *sql.DB, migrationsPath string, steps int, downAll bool) error {
	m, err := newMigrator(db, migrationsPath)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if downAll {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not rollback migrations: %w", err)
		}
		log.Println("Rolled back all migrations")
		return nil
	}

	mVersion, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("no migrations have been applied yet")
		}
		return fmt.Errorf("could not get migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is dirty at version %d; use --force-version to fix it first", mVersion)
	}

	if steps <= 0 {
		return fmt.Errorf("migration steps must be greater than 0")
	}

	if steps > int(mVersion) {
		return fmt.Errorf("migration steps cannot be greater than the current migration version")
	}

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("could not rollback migrations: %w", err)
	}

	log.Printf("Rolled back %d migrations\n\n", steps)
	return nil
}

func forceVersion(db *sql.DB, migrationsPath string, version int) error {
	m, err := newMigrator(db, migrationsPath)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	mVersion, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get migration version: %w", err)
	}

	if !errors.Is(err, migrate.ErrNilVersion) && version > int(mVersion) {
		return fmt.Errorf("forced version cannot be greater than the current migration version")
	}

	if version < 0 {
		return fmt.Errorf("forced version cannot be negative")
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("could not rollback migrations: %w", err)
	}

	log.Printf("Forced migration version to %d\n\n", version)
	return nil
}

func runSQLFile(db *sql.DB, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read seed file %s: %w", path, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction for %s: %w", path, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(contents)); err != nil {
		return fmt.Errorf("could not execute seed file %s: %w", path, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit seed file %s: %w", path, err)
	}

	return nil
}

func seedDatabase(db *sql.DB, seedsPath, seedName string) error {
	if seedsPath == "" {
		return fmt.Errorf("seed path cannot be empty")
	}
	if seedName != "" {
		return runSQLFile(db, filepath.Join(seedsPath, seedName))
	}

	files, err := os.ReadDir(seedsPath)
	if err != nil {
		return fmt.Errorf("could not read seed directory %s: %w", seedsPath, err)
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".sql" {
			continue
		}

		if err := runSQLFile(db, filepath.Join(seedsPath, f.Name())); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open(dbDriverName, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	up := flag.Bool("up", false, "run up migrations")
	down := flag.Bool("down", false, "run down migrations")
	steps := flag.Int("steps", 0, "number of migration steps to rollback")
	downAll := flag.Bool("down-all", false, "rollback all migrations")
	fVersion := flag.Int("force-version", -1, "force migration version and clear dirty state")
	seed := flag.Bool("seed", false, "seed the database")
	seedsPath := flag.String("seeds-path", "/seeds", "path to seed files")
	seedFileName := flag.String("seed-file-name", "", "name of the seed file")

	flag.Parse()

	selectedActions := 0
	if *up {
		selectedActions++
	}
	if *down {
		selectedActions++
	}
	if *fVersion >= 0 {
		selectedActions++
	}
	if *seed {
		selectedActions++
	}

	if selectedActions == 0 {
		log.Fatalf("one of --up, -down, --force-version, or --seed must be provided")
	}

	if selectedActions > 1 {
		log.Fatal("--up, --down, --force-version, and --seed are mutually exclusive")
	}

	if *down && *downAll && *steps > 0 {
		log.Fatal("--down-all and --steps cannot be used together")
	}

	if *up {
		if err := runMigrations(db, migrationsPath); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}

	if *down {
		if !*downAll && *steps <= 0 {
			log.Fatal("--down requires either --down-all or --steps > 0")
		}

		if err := rollbackMigrations(db, migrationsPath, *steps, *downAll); err != nil {
			log.Fatalf("failed to rollback migrations: %v", err)
		}
	}

	if *fVersion >= 0 {
		if err := forceVersion(db, migrationsPath, *fVersion); err != nil {
			log.Fatalf("failed to force migration version: %v", err)
		}
	}

	if *seed == true {
		err := seedDatabase(db, *seedsPath, *seedFileName)
		if err != nil {
			log.Fatalf("failed to seed database: %v", err)
		}
	}

	log.Println("Database command completed successfully")
}
