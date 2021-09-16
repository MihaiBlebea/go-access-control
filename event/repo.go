package event

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type EventRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *EventRepo {
	return &EventRepo{conn}
}

func (r *EventRepo) eventsWithUserID(userID int) ([]Event, error) {
	events := []Event{}
	err := r.conn.Where("user_id = ?", userID).Find(&events).Error
	if err != nil {
		return events, err
	}

	return events, nil
}

func (r *EventRepo) store(event *Event) error {
	return r.conn.Create(event).Error
}
