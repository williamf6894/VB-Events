package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %s", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Printf("failed to terminate postgres container: %s", err)
		}
	}()

	testDB, err = gorm.Open(gormpostgres.Open(container.MustConnectionString(ctx, "sslmode=disable")), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %s", err)
	}

	if err := testDB.AutoMigrate(&models.Event{}, &models.Participant{}); err != nil {
		log.Fatalf("failed to migrate test database: %s", err)
	}

	os.Exit(m.Run())
}
