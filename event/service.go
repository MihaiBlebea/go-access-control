package event

type Service interface {
	StoreEvent(userID int, action string) error
	UserEvents(userID int) ([]Event, error)
}

type service struct {
	eventRepo *EventRepo
}

func NewService(eventRepo *EventRepo) Service {
	return &service{eventRepo}
}

func (s *service) StoreEvent(userID int, action string) error {
	return s.eventRepo.store(&Event{
		UserID: userID,
		Action: action,
	})
}

func (s *service) UserEvents(userID int) ([]Event, error) {
	return s.eventRepo.eventsWithUserID(userID)
}
