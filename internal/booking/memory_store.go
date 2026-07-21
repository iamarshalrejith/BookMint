package booking

type MemoryStore struct {
	// Mapping seats to booking
	bookings map[string]Booking // "A2" -> booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	// is the seat taken? - return error
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	// if not populate
	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieId string) []Booking {
	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieId {
			result = append(result,b)
		}
	}
	return result
}