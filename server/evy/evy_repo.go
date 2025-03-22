package evy

type Repository interface {
	CreateEvent(e EventPayload) error
	GetEventsBatch(size int32) ([]EventPayload, error)
	DeleteEvent(id uint) error
}
