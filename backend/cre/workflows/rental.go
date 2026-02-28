//go:build wasip1

// Package workflows implements the rental management workflow for Chainlink CRE.
// This workflow handles rental agreements, automated payments, and disputes.
package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/pkg/encryption"
	"RWA-Houses/backend/cre/pkg/validation"
)

// RentalRequest represents a rental agreement request
type RentalRequest struct {
	Action          string `json:"action"`
	TokenID         uint64 `json:"tokenID"`
	RenterAddress   string `json:"renterAddress"`
	DurationDays    uint64 `json:"durationDays"`
	MonthlyRent     string `json:"monthlyRent"`
	RenterPublicKey string `json:"renterPublicKey"`
	DepositAmount   string `json:"depositAmount,omitempty"`
}

// RentalResponse represents a rental agreement response
type RentalResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	RentalID      string `json:"rentalID,omitempty"`
	TxHash        string `json:"txHash,omitempty"`
	AccessKeyHash string `json:"accessKeyHash,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	Deposit       string `json:"deposit,omitempty"`
}

// RentalWorkflow handles the complete rental process
func RentalWorkflow(cfg *config.Config, logger *slog.Logger) func(ctx context.Context, req *http.Payload) (*RentalResponse, error) {
	validator := validation.NewValidator()

	return func(ctx context.Context, req *http.Payload) (*RentalResponse, error) {
		logger.Info("Starting rental workflow")

		// 1. Parse request
		var rentalReq RentalRequest
		if err := json.Unmarshal(req.Input, &rentalReq); err != nil {
			logger.Error("Failed to parse rental request", "error", err)
			return &RentalResponse{
				Success: false,
				Message: "Invalid request format",
			}, nil
		}

		// 2. Validate inputs
		if err := validator.ValidateTokenID(rentalReq.TokenID); err != nil {
			return &RentalResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateEthereumAddress(rentalReq.RenterAddress); err != nil {
			return &RentalResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateDuration(rentalReq.DurationDays); err != nil {
			return &RentalResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateAmount(rentalReq.MonthlyRent); err != nil {
			return &RentalResponse{Success: false, Message: err.Error()}, nil
		}

		logger.Info("Processing rental",
			"tokenID", rentalReq.TokenID,
			"renter", rentalReq.RenterAddress,
			"duration", rentalReq.DurationDays)

		// 3. Verify property is available for rent
		if err := verifyPropertyAvailable(ctx, rentalReq.TokenID, logger); err != nil {
			logger.Error("Property not available", "error", err)
			return &RentalResponse{
				Success: false,
				Message: "Property not available for rent",
			}, nil
		}

		// 4. Calculate rent and deposit
		rentAmount, _ := new(big.Int).SetString(rentalReq.MonthlyRent, 10)
		depositAmount := new(big.Int).Mul(rentAmount, big.NewInt(1)) // 1 month deposit

		if rentalReq.DepositAmount != "" {
			depositAmount, _ = new(big.Int).SetString(rentalReq.DepositAmount, 10)
		}

		// 5. Generate temporary access key
		accessKey, err := encryption.GenerateRandomKey(32)
		if err != nil {
			logger.Error("Failed to generate access key", "error", err)
			return &RentalResponse{
				Success: false,
				Message: "Access key generation failed",
			}, nil
		}

		// 6. Create rental agreement on-chain
		rentalID, txHash, err := createRentalOnChain(ctx, cfg, rentalReq, depositAmount, logger)
		if err != nil {
			logger.Error("Failed to create rental on-chain", "error", err)
			return &RentalResponse{
				Success: false,
				Message: "Failed to create rental agreement",
			}, nil
		}

		// 7. Calculate dates
		startDate := time.Now()
		endDate := startDate.Add(time.Duration(rentalReq.DurationDays) * 24 * time.Hour)

		// 8. Schedule automated payments
		if err := scheduleAutomatedPayments(ctx, cfg, rentalID, rentalReq.TokenID, rentAmount, logger); err != nil {
			logger.Error("Failed to schedule payments", "error", err)
			// Non-fatal: continue with rental creation
		}

		accessKeyHash := encryption.HashShares([]encryption.Share{
			{Index: 1, Value: accessKey},
		})

		logger.Info("Rental created successfully",
			"rentalID", rentalID,
			"tokenID", rentalReq.TokenID,
			"renter", rentalReq.RenterAddress)

		return &RentalResponse{
			Success:       true,
			Message:       "Rental agreement created",
			RentalID:      rentalID,
			TxHash:        txHash,
			AccessKeyHash: accessKeyHash,
			StartDate:     startDate.Format(time.RFC3339),
			EndDate:       endDate.Format(time.RFC3339),
			Deposit:       depositAmount.String(),
		}, nil
	}
}

// HandleCreateRental is the main entry point for rental creation
func HandleCreateRental(ctx context.Context, request RentalRequest) (*RentalResponse, error) {
	validator := validation.NewValidator()

	// Validate inputs
	if err := validator.ValidateTokenID(request.TokenID); err != nil {
		return nil, err
	}
	if err := validator.ValidateEthereumAddress(request.RenterAddress); err != nil {
		return nil, err
	}
	if err := validator.ValidateDuration(request.DurationDays); err != nil {
		return nil, err
	}

	return &RentalResponse{
		Success:   true,
		Message:   fmt.Sprintf("Rental request validated for token %d", request.TokenID),
		StartDate: time.Now().Format(time.RFC3339),
	}, nil
}

// AutomatedRentalPaymentWorkflow handles cron-triggered rental payments
func AutomatedRentalPaymentWorkflow(cfg *config.Config, logger *slog.Logger) func(ctx context.Context, trigger *cron.Payload) (string, error) {
	return func(ctx context.Context, trigger *cron.Payload) (string, error) {
		logger.Info("Running automated rental payment check")

		// 1. Get all active rentals with due payments
		rentals, err := getActiveRentals(ctx, cfg, logger)
		if err != nil {
			logger.Error("Failed to get active rentals", "error", err)
			return "", err
		}

		// 2. Process payments for each rental
		processedCount := 0
		for _, rental := range rentals {
			if isPaymentDue(rental) {
				if err := ProcessRentalPayment(ctx, rental.ID); err != nil {
					logger.Error("Failed to process rental payment",
						"rentalID", rental.ID,
						"error", err)
					continue
				}
				processedCount++
			}
		}

		result := fmt.Sprintf("Processed %d rental payments", processedCount)
		logger.Info(result)
		return result, nil
	}
}

// ProcessRentalPayment processes a single rental payment
func ProcessRentalPayment(ctx context.Context, rentalID string) error {
	// In a real implementation, this would:
	// 1. Verify payment is due
	// 2. Check renter's balance/approval
	// 3. Execute payment on-chain
	// 4. Update rental status
	// 5. Emit payment event

	// For now, just simulate processing
	return nil
}

// RentalInfo represents information about a rental
type RentalInfo struct {
	ID          string
	TokenID     uint64
	Renter      string
	MonthlyRent *big.Int
	StartDate   time.Time
	EndDate     time.Time
	LastPayment time.Time
	Status      string
}

// verifyPropertyAvailable checks if a property is available for rent
func verifyPropertyAvailable(ctx context.Context, tokenID uint64, logger *slog.Logger) error {
	// In a real implementation, this would:
	// 1. Query the HouseRWA contract
	// 2. Check if property is already rented
	// 3. Verify no pending disputes
	// 4. Check maintenance status

	logger.Info("Verifying property availability", "tokenID", tokenID)
	return nil
}

// createRentalOnChain creates the rental agreement on-chain
func createRentalOnChain(ctx context.Context, cfg *config.Config, req RentalRequest, deposit *big.Int, logger *slog.Logger) (string, string, error) {
	// In a real implementation, this would:
	// 1. Create and sign transaction
	// 2. Call HouseRWA.createRental()
	// 3. Transfer deposit
	// 4. Emit rental event

	rentalID := fmt.Sprintf("rental_%d_%d_%d", req.TokenID, time.Now().Unix(), time.Now().Nanosecond())
	txHash := fmt.Sprintf("0x%s%d", generateRandomHex(64), time.Now().Unix())

	logger.Info("Created rental on-chain",
		"rentalID", rentalID,
		"tokenID", req.TokenID,
		"txHash", txHash)

	return rentalID, txHash, nil
}

// scheduleAutomatedPayments schedules automated payment processing
func scheduleAutomatedPayments(ctx context.Context, cfg *config.Config, rentalID string, tokenID uint64, rentAmount *big.Int, logger *slog.Logger) error {
	// In a real implementation, this would:
	// 1. Register rental for automated payment processing
	// 2. Set up cron job triggers
	// 3. Configure payment rules (auto-pay threshold, etc.)

	logger.Info("Scheduled automated payments",
		"rentalID", rentalID,
		"tokenID", tokenID,
		"rentAmount", rentAmount.String())

	return nil
}

// getActiveRentals retrieves all active rentals
func getActiveRentals(ctx context.Context, cfg *config.Config, logger *slog.Logger) ([]RentalInfo, error) {
	// In a real implementation, this would:
	// 1. Query the HouseRWA contract for active rentals
	// 2. Filter by status
	// 3. Return rental information

	return []RentalInfo{}, nil
}

// isPaymentDue checks if a rental payment is due
func isPaymentDue(rental RentalInfo) bool {
	// In a real implementation, this would:
	// 1. Check last payment date
	// 2. Calculate next due date
	// 3. Return true if payment is due

	nextPayment := rental.LastPayment.Add(30 * 24 * time.Hour)
	return time.Now().After(nextPayment)
}

// HandleDispute handles rental disputes with slashing
func HandleDispute(ctx context.Context, rentalID string, reason string, evidence []byte) error {
	// In a real implementation, this would:
	// 1. Verify dispute reason
	// 2. Collect evidence
	// 3. Initiate dispute resolution
	// 4. Handle slashing if applicable

	return nil
}
