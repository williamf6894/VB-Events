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
	participantRepo := NewParticipantRepository(testDB)

	created := newEvent("Doomed Event", tsBase)
	if err := repo.Create(created); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	participant := newParticipant("Lena Duarte", "lena@example.com")
	if err := participantRepo.Create(participant); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := repo.AddParticipant(created.ID, participant.ID); err != nil {
		t.Fatalf("AddParticipant returned error: %s", err)
	}

	if err := repo.DeleteByID(created.ID); err != nil {
		t.Fatalf("DeleteByID returned error: %s", err)
	}

	if _, err := repo.FindByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected event to be deleted, got: %v", err)
	}

	var registrations int64
	testDB.Table("event_participants").Where("event_id = ?", created.ID).Count(&registrations)
	if registrations != 0 {
		t.Fatalf("expected registrations to be removed, got %d", registrations)
	}

	if err := repo.DeleteByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound on second delete, got: %v", err)
	}
}

func TestEventRepository_List(t *testing.T) {
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

	events, err := repo.List(models.EventQuery{})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
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

func TestEventRepository_AddAndRemoveParticipant(t *testing.T) {
	cleanTables(t, testDB)
	eventRepo := NewEventRepository(testDB)
	participantRepo := NewParticipantRepository(testDB)

	event := newEvent("Registration Test", tsBase)
	if err := eventRepo.Create(event); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	participant := newParticipant("Iris Chen", "iris@example.com")
	if err := participantRepo.Create(participant); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	if err := eventRepo.AddParticipant(event.ID, participant.ID); err != nil {
		t.Fatalf("AddParticipant returned error: %s", err)
	}

	var count int64
	testDB.Table("event_participants").Where("event_id = ?", event.ID).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 registration, got %d", count)
	}

	if err := eventRepo.RemoveParticipant(event.ID, participant.ID); err != nil {
		t.Fatalf("RemoveParticipant returned error: %s", err)
	}

	testDB.Table("event_participants").Where("event_id = ?", event.ID).Count(&count)
	if count != 0 {
		t.Fatalf("expected registration removed, got %d", count)
	}

	if err := eventRepo.RemoveParticipant(event.ID, participant.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound on second remove, got: %v", err)
	}

	if err := eventRepo.AddParticipant(uuid.Nil(), participant.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound for missing event, got: %v", err)
	}

	if err := eventRepo.AddParticipant(event.ID, uuid.Nil()); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound for missing participant, got: %v", err)
	}
}

func TestEventRepository_FindByIDPreloadsParticipants(t *testing.T) {
	cleanTables(t, testDB)
	eventRepo := NewEventRepository(testDB)
	participantRepo := NewParticipantRepository(testDB)

	event := newEvent("Preload Test", tsBase)
	if err := eventRepo.Create(event); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	first := newParticipant("Jack Amari", "jack@example.com")
	second := newParticipant("Kara Silva", "kara@example.com")
	if err := participantRepo.Create(first); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := participantRepo.Create(second); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	for _, participant := range []*struct {
		ID uuid.UUID
	}{{first.ID}, {second.ID}} {
		if err := eventRepo.AddParticipant(event.ID, participant.ID); err != nil {
			t.Fatalf("AddParticipant returned error: %s", err)
		}
	}

	found, err := eventRepo.FindByID(event.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %s", err)
	}
	if len(found.Participants) != 2 {
		t.Fatalf("expected 2 preloaded participants, got %d", len(found.Participants))
	}
}

func TestEventRepository_ListSearch(t *testing.T) {
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

	events, err := repo.List(models.EventQuery{Search: "beach"})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 matches (name, description, location), got %d", len(events))
	}

	events, err = repo.List(models.EventQuery{Search: "volley"})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Beach Volleyball" {
		t.Fatalf("expected only Beach Volleyball, got %+v", events)
	}

	events, err = repo.List(models.EventQuery{Search: "VOLLEYBALL"})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected case-insensitive match, got %d", len(events))
	}

	events, err = repo.List(models.EventQuery{Search: "nonexistent"})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no matches, got %d", len(events))
	}
}

func TestEventRepository_ListBeforeAndAfter(t *testing.T) {
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

	events, err := repo.List(models.EventQuery{Before: &tsBase})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Past" {
		t.Fatalf("expected only strictly-earlier event, got %+v", events)
	}

	events, err = repo.List(models.EventQuery{After: &tsBase})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Future" {
		t.Fatalf("expected only strictly-later event, got %+v", events)
	}
}

func TestEventRepository_ListBetween(t *testing.T) {
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

	events, err := repo.List(models.EventQuery{After: &tsMid, Before: &tsFar})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Middle" {
		t.Fatalf("expected only the strictly-inside event, got %+v", events)
	}

	widerBefore := tsFar.Add(time.Hour)
	events, err = repo.List(models.EventQuery{After: &tsMid, Before: &widerBefore})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events with widened window (end boundary in, start boundary out), got %d", len(events))
	}
	for _, e := range events {
		if e.Name == "Before" {
			t.Fatal("event outside window should not be returned")
		}
	}
}

func TestEventRepository_ListFullFilter(t *testing.T) {
	cleanTables(t, testDB)
	eventRepo := NewEventRepository(testDB)
	participantRepo := NewParticipantRepository(testDB)

	fullEvent := newEvent("Full Event", tsFar)
	fullEvent.Capacity = 1
	if err := eventRepo.Create(fullEvent); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	spaciousEvent := newEvent("Spacious Event", tsFar)
	spaciousEvent.Capacity = 10
	if err := eventRepo.Create(spaciousEvent); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	participant := newParticipant("Ola Petrov", "ola@example.com")
	if err := participantRepo.Create(participant); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}
	if err := eventRepo.AddParticipant(fullEvent.ID, participant.ID); err != nil {
		t.Fatalf("AddParticipant returned error: %s", err)
	}

	full := true
	events, err := eventRepo.List(models.EventQuery{Full: &full})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Full Event" {
		t.Fatalf("expected only the full event, got %+v", events)
	}

	notFull := false
	events, err = eventRepo.List(models.EventQuery{Full: &notFull})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Spacious Event" {
		t.Fatalf("expected only the spacious event, got %+v", events)
	}
}

func TestEventRepository_ListCombinesFilters(t *testing.T) {
	cleanTables(t, testDB)
	repo := NewEventRepository(testDB)

	target := newEvent("Beach Target", tsFar)
	if err := repo.Create(target); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	pastMatch := newEvent("Beach Past", tsPast)
	if err := repo.Create(pastMatch); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	otherFuture := newEvent("Indoor Future", tsFar)
	if err := repo.Create(otherFuture); err != nil {
		t.Fatalf("Create returned error: %s", err)
	}

	events, err := repo.List(models.EventQuery{Search: "beach", After: &tsBase})
	if err != nil {
		t.Fatalf("List returned error: %s", err)
	}
	if len(events) != 1 || events[0].Name != "Beach Target" {
		t.Fatalf("expected only future beach event, got %+v", events)
	}
}
