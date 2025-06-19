package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/refine-software/afrad-api/config"
	"github.com/refine-software/afrad-api/internal/database"
)

var testService database.Service

func TestMain(m *testing.M) {
	fmt.Println("🚀 Setting up test environment...")

	// Setup environment config
	env := config.NewTestEnv()

	// Connect to test DB

	testService = database.New(env)

	sqlDB, err := sql.Open("pgx", env.DBUrl)
	if err != nil {
		log.Fatalf("❌ Failed to open DB for migrations: %v", err)
	}
	defer sqlDB.Close()

	err = goose.Up(sqlDB, "../database/migrations/")
	if err != nil {
		log.Fatalf("❌ Could not run migrations: %v", err)
	}

	// Clean DB
	err = truncateAllTables(testService)
	if err != nil {
		log.Fatalf("could not clean test DB: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Teardown
	fmt.Println("🧹 Cleaning up after tests...")
	testService.Close()

	os.Exit(code)
}
