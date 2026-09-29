package repository

import (
	"errors"
	"testing"
	"uuid"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
)

func TestParticipantRepository_Create(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	participant := newParticipant("Alice", "alice@example.com")
	if err := repo.Create(participant); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	if participant.ID == uuid.Nil() {
		t.Fatal("expected ID to be set by BeforeCreate hook")
	}

	var stored models.Participant
	if err := testDB.First(&stored, "id = ?", participant.ID).Error; err != nil {
		t.Fatalf("participant not persisted: %s", err)
	}
	if stored.Name != "Alice" || stored.Email != "alice@example.com" {
		t.Fatalf("unexpected stored values: %+v", stored)
	}
}

func TestParticipantRepository_FindByID(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	created := newParticipant("Bob", "bob@example.com")
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %s", err)
	}
	if found.Email != "bob@example.com" {
		t.Fatalf("unexpected participant: %+v", found)
	}

	if _, err := repo.FindByID(uuid.Nil()); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got: %v", err)
	}
}

func TestParticipantRepository_FindByEmail(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	if err := repo.Create(newParticipant("Carol", "carol@example.com")); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	found, err := repo.FindByEmail("carol@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %s", err)
	}
	if found.Name != "Carol" {
		t.Fatalf("unexpected participant: %+v", found)
	}

	if _, err := repo.FindByEmail("missing@example.com"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got: %v", err)
	}
}

func TestParticipantRepository_Update(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	created := newParticipant("Dave", "dave@example.com")
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	created.Name = "David"
	created.Email = "david@example.com"
	if err := repo.Update(created); err != nil {
		t.Fatalf("Update returned error: %s", err)
	}

	found, err := repo.FindByEmail("david@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %s", err)
	}
	if found.Name != "David" {
		t.Fatalf("expected updated name David, got %q", found.Name)
	}
}

func TestParticipantRepository_DeleteByID(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	created := newParticipant("Eve", "eve@example.com")
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	if err := repo.DeleteByID(created.ID); err != nil {
		t.Fatalf("DeleteByID returned error: %s", err)
	}

	if _, err := repo.FindByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected participant to be deleted, got: %v", err)
	}

	if err := repo.DeleteByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound on second delete, got: %v", err)
	}
}

func TestParticipantRepository_ListAll(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewParticipantRepository(testDB)

	if err := repo.Create(newParticipant("Frank", "frank@example.com")); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newParticipant("Grace", "grace@example.com")); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	participants, err := repo.ListAll()
	if err != nil {
		t.Fatalf("ListAll returned error: %s", err)
	}
	if len(participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(participants))
	}
}
