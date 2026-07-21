package booking

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
}

// Any type that wants to behave as a BookingStore must have these two methods.
type BookingStore interface{
	Book(b Booking) (Booking,error)
	ListBookings(movieID string) []Booking
}