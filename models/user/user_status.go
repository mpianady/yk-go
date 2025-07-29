package user

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
	StatusBanned   Status = "BANNED"
	StatusPending  Status = "PENDING"
)
