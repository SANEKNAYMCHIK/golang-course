package domain

import "time"

type Repo struct {
	Name        string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
}
