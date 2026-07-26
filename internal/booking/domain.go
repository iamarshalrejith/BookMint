package booking

import (
	"context"
	"errors"
	"time"
)

var (
    ErrSeatAlreadyBooked = errors.New("seat already booked")
    ErrSessionExpired    = errors.New("session expired")
    ErrNotSeatOwner      = errors.New("user does not own this seat")
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
    Book(Booking) (Booking, error)
    ListBookings(movieID string) []Booking

    Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
    Release(ctx context.Context, sessionID string, userID string) error
}