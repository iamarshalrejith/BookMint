package booking

import (
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/iamarshalrejith/BookMint/internal/adapters/redis"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClient("localhost:6379"))
	svc := NewService(store)

	const numGoroutines = 100_000 // Pretend 100,000 users are trying to book the same seat at exactly the same time.

	var (
		successes atomic.Int64   // Counts how many bookings succeeded.
		failures  atomic.Int64   // Counts how many bookings failed.
		wg        sync.WaitGroup // Makes the main test wait until all 100,000 goroutines finish.
	)

	// Starting the Goroutines
	wg.Add(numGoroutines) // Telling the WaitGroup abt starting 100,000 goroutines
	for i := range numGoroutines{
		go func(userNum int){ // Each iteration starts one goroutine
			defer wg.Done()
			_,err := svc.Book(Booking{
				MovieID: "abc-1",
				SeatID: "A1",
				UserID: uuid.New().String(),
			})

			if err == nil {
				successes.Add(1)
			}else{
				failures.Add(1)
			}
		}(i)
	}
		wg.Wait() // The test waits until all 100,000 goroutines have completed


		// Verifying the result
		if got := successes.Load(); got != 1 {
			t.Errorf("Expected exactly 1 success, but got %d",got)
		}

		if got := failures.Load(); got != int64(numGoroutines - 1){
			t.Errorf("expected %d failures, got %d", numGoroutines-1, got)
		}
	}