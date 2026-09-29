package main

import (
	"log/slog"
	"time"

	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/db"
	"github.com/williamf6894/VB-Events/internal/models"
)

func main() {
	cfg := config.Load()

	database, err := db.InitDB(cfg)
	if err != nil {
		panic("failed to connect to database")
	}

	if err := database.AutoMigrate(&models.Event{}, &models.Participant{}); err != nil {
		panic("failed to migrate database")
	}

	if err := database.Exec("TRUNCATE events, participants, event_participants CASCADE").Error; err != nil {
		panic("failed to clear existing data")
	}

	participants := seedParticipants()
	if err := database.Create(&participants).Error; err != nil {
		panic("failed to seed participants")
	}

	events := seedEvents()
	if err := database.Create(&events).Error; err != nil {
		panic("failed to seed events")
	}

	for i := range events {
		count := 1 + i%4
		attendees := make([]*models.Participant, 0, count)
		for j := 0; j < count; j++ {
			attendees = append(attendees, &participants[(i*2+j)%len(participants)])
		}
		if err := database.Model(&events[i]).Association("Participants").Append(&attendees); err != nil {
			panic("failed to associate participants with event " + events[i].Name)
		}
	}

	var eventCount, participantCount int64
	database.Model(&models.Event{}).Count(&eventCount)
	database.Model(&models.Participant{}).Count(&participantCount)

	slog.Info("seed complete", "events", eventCount, "participants", participantCount)
}

func seedParticipants() []models.Participant {
	return []models.Participant{
		{Name: "Alice Nguyen", Email: "alice@example.com"},
		{Name: "Bob Torres", Email: "bob@example.com"},
		{Name: "Carol Ito", Email: "carol@example.com"},
		{Name: "Dave Okafor", Email: "dave@example.com"},
		{Name: "Eve Kowalski", Email: "eve@example.com"},
		{Name: "Frank Delgado", Email: "frank@example.com"},
		{Name: "Grace Kim", Email: "grace@example.com"},
		{Name: "Heather Blake", Email: "heather@example.com"},
	}
}

func seedEvents() []models.Event {
	now := time.Now()
	day := 24 * time.Hour

	return []models.Event{
		{
			Name:           "North End Open",
			Description:    "Season-opening 2v2 tournament on the north end courts",
			Location:       "VB Oceanfront North",
			Capacity:       48,
			Duration:       6,
			StartTimestamp: now.Add(7 * day),
		},
		{
			Name:           "Wednesday Indoor League",
			Description:    "Weekly 6-team indoor league night",
			Location:       "VB Sports Center",
			Capacity:       36,
			Duration:       3,
			StartTimestamp: now.Add(14 * day),
		},
		{
			Name:           "Beach Volleyball Clinic",
			Description:    "Beginner-friendly coaching clinic with pro players",
			Location:       "VB Oceanfront South",
			Capacity:       24,
			Duration:       2.5,
			StartTimestamp: now.Add(21 * day),
		},
		{
			Name:           "King of the Beach",
			Description:    "Individual-format tournament, rotating partners each round",
			Location:       "First Landing State Park",
			Capacity:       60,
			Duration:       8,
			StartTimestamp: now.Add(30 * day),
		},
		{
			Name:           "Sunset Social Games",
			Description:    "Casual mixed games as the sun goes down",
			Location:       "Sandbridge Beach",
			Capacity:       32,
			Duration:       2,
			StartTimestamp: now.Add(45 * day),
		},
		{
			Name:           "Halloween Spooky Smash",
			Description:    "Costume contest followed by a themed tournament",
			Location:       "VB Oceanfront North",
			Capacity:       40,
			Duration:       5,
			StartTimestamp: now.Add(-3 * day),
		},
		{
			Name:           "Autumn Round Robin",
			Description:    "Relaxed round-robin to close out the outdoor season",
			Location:       "Mount Trashmore Park",
			Capacity:       28,
			Duration:       4,
			StartTimestamp: now.Add(-10 * day),
		},
		{
			Name:           "Referee Certification",
			Description:    "Official certification course for new referees",
			Location:       "VB Sports Center",
			Capacity:       16,
			Duration:       7.5,
			StartTimestamp: now.Add(-21 * day),
		},
	}
}
