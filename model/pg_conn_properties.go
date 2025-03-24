package model

// PgConnProperties is used for storing connection properties for database
type PgConnProperties struct {
	Url      string
	Username string
	Password string
	RoHost   string
	Role     string
}
