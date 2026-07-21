package booking

import "sync"

type ConcurentStore struct {
	// Mapping seats to booking
	bookings map[string]Booking // "A2" -> booking
	sync.RWMutex // By embedding sync.RWMutex, your store gets these methods directly: Lock(), Unlock(), RLock(), RUnlock()
}

func NewConcurentStore() *ConcurentStore {
	return &ConcurentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	// is the seat taken? - return error
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	// if not populate
	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurentStore) ListBookings(movieId string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieId {
			result = append(result, b)
		}
	}
	return result
}