package userclient

type UserRole string

const (
	Guest UserRole = "guest"
	Host  UserRole = "host"
	Admin UserRole = "admin"
)
