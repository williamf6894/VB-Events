package repository

import (
	"errors"
	"testing"
	"time"
	"uuid"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
)

var (
	tsBase = time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	tsPast = tsBase.Add(-24 * time.Hour)
	tsMid  = tsBase
	tsFar  = tsBase.Add(72 * time.Hour)
)

func TestEventRepository_Create(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	event := newEvent("Beach Volleyball", tsBase)
	if err := repo.Create(event); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	if event.ID == uuid.Nil() {
		t.Fatal("expected ID to be set by BeforeCreate hook")
	}

	var stored models.Event
	if err := testDB.First(&stored, "id = ?", event.ID).Error; err != nil {
		t.Fatalf("event not persisted: %s", err)
	}
	if stored.Name != "Beach Volleyball" || stored.Capacity != 10 || stored.Duration != 1.5 {
		t.Fatalf("unexpected stored values: %+v", stored)
	}
	if !stored.StartTimestamp.Equal(tsBase) {
		t.Fatalf("expected start timestamp %v, got %v", tsBase, stored.StartTimestamp)
	}
}

func TestEventRepository_FindByID(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	created := newEvent("Indoor League", tsBase)
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %s", err)
	}
	if found.Name != "Indoor League" {
		t.Fatalf("unexpected event: %+v", found)
	}

	if _, err := repo.FindByID(uuid.Nil()); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got: %v", err)
	}
}

func TestEventRepository_FindByName(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	if err := repo.Create(newEvent("Beach Volleyball", tsBase)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	found, err := repo.FindByName("Beach Volleyball")
	if err != nil {
		t.Fatalf("FindByName returned error: %s", err)
	}
	if found.Location != "Beach Volleyball location" {
		t.Fatalf("unexpected event: %+v", found)
	}

	if _, err := repo.FindByName("No Such Event"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got: %v", err)
	}
}

func TestEventRepository_Update(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	created := newEvent("Old Name", tsBase)
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	created.Name = "New Name"
	created.Capacity = 99
	if err := repo.Update(created); err != nil {
		t.Fatalf("Update returned error: %s", err)
	}

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %s", err)
	}
	if found.Name != "New Name" || found.Capacity != 99 {
		t.Fatalf("expected updated values, got: %+v", found)
	}
}

func TestEventRepository_DeleteByID(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	created := newEvent("Doomed Event", tsBase)
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	if err := repo.DeleteByID(created.ID); err != nil {
		t.Fatalf("DeleteByID returned error: %s", err)
	}

	if _, err := repo.FindByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected event to be deleted, got: %v", err)
	}

	if err := repo.DeleteByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound on second delete, got: %v", err)
	}
}

func TestEventRepository_ListAll(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	if err := repo.Create(newEvent("Far", tsFar)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Past", tsPast)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Mid", tsMid)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.ListAll()
	if err != nil {
		t.Fatalf("ListAll returned error: %s", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	expectedOrder := []string{"Past", "Mid", "Far"}
	for i, name := range expectedOrder {
		if events[i].Name != name {
			t.Fatalf("expected %q at position %d, got %q", name, i, events[i].Name)
		}
	}
}

func TestEventRepository_FindByPartialNameDescriptionLocation(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	nameMatch := newEvent("Beach Volleyball", tsBase)
	nameMatch.Description = "Casual games"
	nameMatch.Location = "Oceanfront"
	if err := repo.Create(nameMatch); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	descriptionMatch := newEvent("Chess Meetup", tsBase)
	descriptionMatch.Description = "beachside strategy games"
	descriptionMatch.Location = "Community centre"
	if err := repo.Create(descriptionMatch); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	locationMatch := newEvent("Trivia Night", tsBase)
	locationMatch.Description = "Pub quiz"
	locationMatch.Location = "North Beach Bar"
	if err := repo.Create(locationMatch); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.FindByPartialNameDescriptionLocation("beach")
	if err != nil {
		t.Fatalf("FindByPartialNameDescriptionLocation returned error: %s", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 matches (name, description, location), got %d", len(events))
	}

	events, err = repo.FindByPartialNameDescriptionLocation("volley")
	if err != nil {
		t.Fatalf("FindByPartialNameDescriptionLocation returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Beach Volleyball" {
		t.Fatalf("expected only Beach Volleyball, got %+v", events)
	}

	events, err = repo.FindByPartialNameDescriptionLocation("VOLLEYBALL")
	if err != nil {
		t.Fatalf("FindByPartialNameDescriptionLocation returned error: %s", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected case-insensitive match, got %d", len(events))
	}

	events, err = repo.FindByPartialNameDescriptionLocation("nonexistent")
	if err != nil {
		t.Fatalf("FindByPartialNameDescriptionLocation returned error: %s", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no matches, got %d", len(events))
	}
}

func TestEventRepository_FindAllBefore(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	if err := repo.Create(newEvent("Past", tsPast)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Boundary", tsBase)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Future", tsFar)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.FindAllBefore(tsBase)
	if err != nil {
		t.Fatalf("FindAllBefore returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Past" {
		t.Fatalf("expected only strictly-earlier event, got %+v", events)
	}
}

func TestEventRepository_FindAllAfter(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	if err := repo.Create(newEvent("Past", tsPast)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Boundary", tsBase)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Future", tsFar)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.FindAllAfter(tsBase)
	if err != nil {
		t.Fatalf("FindAllAfter returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Future" {
		t.Fatalf("expected only strictly-later event, got %+v", events)
	}
}

func TestEventRepository_FindAllBetween(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	if err := repo.Create(newEvent("Before", tsPast)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Start Boundary", tsMid)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("End Boundary", tsFar)); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.Create(newEvent("Middle", tsMid.Add(24*time.Hour))); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.FindAllBetween(tsMid, tsFar)
	if err != nil {
		t.Fatalf("FindAllBetween returned error: %s", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events with inclusive boundaries, got %d", len(events))
	}
	for _, e := range events {
		if e.Name == "Before" {
			t.Fatal("event outside window should not be returned")
		}
	}
}
