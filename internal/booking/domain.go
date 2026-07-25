package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("Seat is already taken")
)

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
	ExpiresAt time.Time
}

// Any type that wants to behave as a BookingStore must have these two methods.
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}