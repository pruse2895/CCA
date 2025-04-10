package scheduler

import (
	"context"
	"log"
	"time"

	"cricketApp/db"
	"cricketApp/models"
)

type ReminderScheduler struct {
	db db.Database
}

func NewReminderScheduler(db db.Database) *ReminderScheduler {
	return &ReminderScheduler{db: db}
}

func (s *ReminderScheduler) Start() {
	// Run immediately on start
	go s.checkAndSendReminders()

	// Schedule to run daily at 11 AM
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(next.Sub(now))
		s.checkAndSendReminders()
	}
}

func (s *ReminderScheduler) checkAndSendReminders() {
	ctx := context.Background()

	// Get all cricketers
	cricketers, err := s.db.GetAllCricketers(ctx)
	if err != nil {
		log.Printf("Error fetching cricketers for reminders: %v", err)
		return
	}

	// Get current time
	now := time.Now()

	// Check each cricketer
	for _, cricketer := range cricketers {
		if cricketer.DueDate == nil {
			continue
		}

		// Calculate days until due date
		daysUntilDue := int(cricketer.DueDate.Sub(now).Hours() / 24)

		// If due date is exactly 2 days away
		if daysUntilDue == 2 {
			s.sendReminder(cricketer)
		}
	}
}

func (s *ReminderScheduler) sendReminder(cricketer models.Cricketer) {
	// TODO: Implement actual notification sending
	// This could be email, SMS, or any other notification method
	log.Printf("Sending fee reminder to cricketer: %s (ID: %s). Due date: %s",
		cricketer.Name,
		cricketer.ID.Hex(),
		cricketer.DueDate.Format("2006-01-02"))
}
