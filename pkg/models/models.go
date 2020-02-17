package models

import (
	"errors"
	"time"
)

type Environment struct {
	ID      int       `db:"id"`
	Name    string    `db:"name"`
	Image   string    `db:"image"`
	About   string    `db:"about"`
	Claimed bool      `db:"claimed"`
	Claimer string    `db:"claimer"`
	Date    time.Time `db:"date"`
}

var (
	ErrClaimed = errors.New("You already claimed an environment. You need to unclaim it before you can claim another one.")

	ErrUnclaimed = errors.New("You can just unclaim environments that are claimed by you")
)
