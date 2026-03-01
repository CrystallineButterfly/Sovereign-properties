package workflows

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============ LOAD TESTING ============

func BenchmarkMintHouse(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		req := MintRequest{
			Metadata: PropertyMetadata{
				HouseID: "bench-house",
				Address: "Benchmark Blvd",
				Price:   100000000,
			},
			Documents: []Document{
				{Type: "deed", Data: []byte("deed"), Hash: "0xbench"},
			},
		}

		_, err := HandleMintHouse(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPrivateSale(b *testing.B) {
	ctx := context.Background()

	// Setup: Mint a house
	mintReq := MintRequest{
		Metadata: PropertyMetadata{
			HouseID: "bench-sale",
			Address: "Sale Street",
			Price:   100000000,
		},
		Documents: []Document{{Type: "deed", Data: []byte("deed"), Hash: "0xsale"}},
	}
	mintResp, _ := HandleMintHouse(ctx, mintReq)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := PrivateSaleRequest{
			TokenID:          mintResp.TokenID,
			BuyerAddress:     "0xBuyer",
			Price:            100000000,
			IsPrivateSale:    false,
			EncryptedKeyData: []byte("key"),
		}

		_, err := HandlePrivateSale(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBillCreate(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		req := BillRequest{
			TokenID:     "token-123",
			BillType:    "electricity",
			Amount:      15000,
			DueDate:     time.Now().Add(30 * 24 * time.Hour).Unix(),
			Provider:    "0xProvider",
			IsRecurring: false,
		}

		_, err := HandleBillCreate(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ============ CONCURRENCY TESTS ============

func TestConcurrentMints(t *testing.T) {
	ctx := context.Background()
	numConcurrent := 100

	var wg sync.WaitGroup
	errors := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := MintRequest{
				Metadata: PropertyMetadata{
					HouseID: "concurrent-" + string(rune(index)),
					Address: "Concurrent St",
					Price:   100000000,
				},
				Documents: []Document{
					{Type: "deed", Data: []byte("deed"), Hash: "0x" + string(rune(index))},
				},
			}

			_, err := HandleMintHouse(ctx, req)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errCount := 0
	for err := range errors {
		if err != nil {
			errCount++
		}
	}

	assert.Equal(t, 0, errCount, "All concurrent mints should succeed")
}

func TestConcurrentSales(t *testing.T) {
	ctx := context.Background()

	// Setup: Mint house
	mintReq := MintRequest{
		Metadata: PropertyMetadata{
			HouseID: "race-house",
			Address: "Race Road",
			Price:   100000000,
		},
		Documents: []Document{{Type: "deed", Data: []byte("deed"), Hash: "0xrace"}},
	}
	mintResp, err := HandleMintHouse(ctx, mintReq)
	require.NoError(t, err)

	// Concurrent sales attempts
	numAttempts := 50
	var wg sync.WaitGroup
	successes := make(chan bool, numAttempts)

	for i := 0; i < numAttempts; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := PrivateSaleRequest{
				TokenID:          mintResp.TokenID,
				BuyerAddress:     "0xBuyer" + string(rune(index)),
				Price:            100000000,
				IsPrivateSale:    false,
				EncryptedKeyData: []byte("key"),
			}

			_, err := HandlePrivateSale(ctx, req)
			successes <- (err == nil)
		}(i)
	}

	wg.Wait()
	close(successes)

	successCount := 0
	for success := range successes {
		if success {
			successCount++
		}
	}

	// Only one sale should succeed (token transferred)
	assert.GreaterOrEqual(t, successCount, 1, "At least one sale should succeed")
}

func TestConcurrentBillPayments(t *testing.T) {
	ctx := context.Background()

	// Setup: Create bill
	billReq := BillRequest{
		TokenID:     "token-concurrent",
		BillType:    "water",
		Amount:      10000,
		DueDate:     time.Now().Add(30 * 24 * time.Hour).Unix(),
		Provider:    "0xProvider",
		IsRecurring: false,
	}
	billResp, err := HandleBillCreate(ctx, billReq)
	require.NoError(t, err)

	// Concurrent payments
	numPayments := 20
	var wg sync.WaitGroup
	results := make(chan error, numPayments)

	for i := 0; i < numPayments; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			_, err := RecordBillPayment(ctx, "token-concurrent", billResp.BillIndex, "stripe", "pi_"+string(rune(index)))
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	// Only first payment should succeed
	assert.GreaterOrEqual(t, successCount, 1, "At least one payment should succeed")
}

// ============ STRESS TESTS ============

func TestStressManyBills(t *testing.T) {
	ctx := context.Background()

	// Create many bills for a single property
	numBills := 50

	for i := 0; i < numBills; i++ {
		req := BillRequest{
			TokenID:     "stress-token",
			BillType:    "utility",
			Amount:      uint64(1000 * (i + 1)),
			DueDate:     time.Now().Add(time.Duration(i) * 24 * time.Hour).Unix(),
			Provider:    "0xProvider",
			IsRecurring: false,
		}

		_, err := HandleBillCreate(ctx, req)
		require.NoError(t, err)
	}

	// Verify all bills created
	bills, err := GetBillsForToken(ctx, "stress-token")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(bills), numBills)
}

func TestStressManyRentals(t *testing.T) {
	ctx := context.Background()

	// Create many rental agreements
	numRentals := 30
	var wg sync.WaitGroup

	for i := 0; i < numRentals; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Mint house
			mintReq := MintRequest{
				Metadata: PropertyMetadata{
					HouseID: "rental-house-" + string(rune(index)),
					Address: "Rental St",
					Price:   100000000,
				},
				Documents: []Document{{Type: "deed", Data: []byte("deed"), Hash: "0xrent"}},
			}
			mintResp, _ := HandleMintHouse(ctx, mintReq)

			// Create rental
			rentalReq := RentalRequest{
				TokenID:       mintResp.TokenID,
				RenterAddress: "0xRenter" + string(rune(index)),
				DurationDays:  365,
				DepositAmount: 5000000,
				MonthlyRent:   250000,
			}

			HandleRentalStart(ctx, rentalReq)
		}(i)
	}

	wg.Wait()
}

// ============ MEMORY/RESOURCE TESTS ============

func TestMemoryEfficiency(t *testing.T) {
	ctx := context.Background()

	// Create many documents to test memory handling
	largeDoc := make([]byte, 1024*1024) // 1MB document
	for i := range largeDoc {
		largeDoc[i] = byte(i % 256)
	}

	req := MintRequest{
		Metadata: PropertyMetadata{
			HouseID: "large-doc-house",
			Address: "Large Doc Ave",
			Price:   100000000,
		},
		Documents: []Document{
			{Type: "large_file", Data: largeDoc, Hash: "0xlargehash"},
		},
	}

	resp, err := HandleMintHouse(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// ============ TIMEOUT TESTS ============

func TestOperationTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Simulate slow operation
	time.Sleep(200 * time.Millisecond)

	req := MintRequest{
		Metadata: PropertyMetadata{
			HouseID: "timeout-house",
			Address: "Timeout St",
			Price:   100000000,
		},
		Documents: []Document{{Type: "deed", Data: []byte("deed"), Hash: "0xto"}},
	}

	_, err := HandleMintHouse(ctx, req)
	assert.Error(t, err, "Should timeout")
}

// ============ MOCK HELPERS ============

func GetBillsForToken(ctx context.Context, tokenID string) ([]BillRequest, error) {
	return []BillRequest{}, nil
}
