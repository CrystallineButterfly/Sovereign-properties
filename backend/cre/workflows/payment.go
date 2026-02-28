//go:build wasip1

// Package workflows implements the bill payment workflow for Chainlink CRE.
// This workflow supports both on-chain (crypto) and off-chain (Stripe) payments.
package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"time"

	crehttp "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/pkg/validation"
)

// PaymentRequest represents a bill payment request
type PaymentRequest struct {
	Action         string  `json:"action"`
	TokenID        uint64  `json:"tokenID"`
	BillIndex      uint64  `json:"billIndex"`
	OwnerAddress   string  `json:"ownerAddress"`
	PaymentMethod  string  `json:"paymentMethod"`
	StripeToken    string  `json:"stripeToken,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	IdempotencyKey string  `json:"idempotencyKey,omitempty"`
}

// PaymentResponse represents a bill payment response
type PaymentResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	TxHash           string `json:"txHash,omitempty"`
	PaymentReference string `json:"paymentReference,omitempty"`
	Amount           string `json:"amount,omitempty"`
	Currency         string `json:"currency,omitempty"`
	Timestamp        string `json:"timestamp,omitempty"`
	Method           string `json:"method,omitempty"`
}

// PaymentWorkflow handles the complete bill payment process
func PaymentWorkflow(cfg *config.Config, logger *slog.Logger) func(ctx context.Context, req *crehttp.Payload) (*PaymentResponse, error) {
	validator := validation.NewValidator()

	return func(ctx context.Context, req *crehttp.Payload) (*PaymentResponse, error) {
		logger.Info("Starting bill payment workflow")

		// 1. Parse request
		var paymentReq PaymentRequest
		if err := json.Unmarshal(req.Input, &paymentReq); err != nil {
			logger.Error("Failed to parse payment request", "error", err)
			return &PaymentResponse{
				Success: false,
				Message: "Invalid request format",
			}, nil
		}

		// 2. Validate inputs
		if err := validator.ValidateTokenID(paymentReq.TokenID); err != nil {
			return &PaymentResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateEthereumAddress(paymentReq.OwnerAddress); err != nil {
			return &PaymentResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidatePaymentMethod(paymentReq.PaymentMethod); err != nil {
			return &PaymentResponse{Success: false, Message: err.Error()}, nil
		}

		logger.Info("Processing payment",
			"tokenID", paymentReq.TokenID,
			"billIndex", paymentReq.BillIndex,
			"method", paymentReq.PaymentMethod)

		// 3. Check idempotency
		if paymentReq.IdempotencyKey != "" {
			if isDuplicatePayment(paymentReq.IdempotencyKey) {
				logger.Info("Duplicate payment detected", "idempotencyKey", paymentReq.IdempotencyKey)
				return &PaymentResponse{
					Success: false,
					Message: "Duplicate payment request",
				}, nil
			}
		}

		// 4. Process payment based on method
		var paymentRef string
		var txHash string
		var err error

		switch paymentReq.PaymentMethod {
		case "crypto":
			paymentRef, txHash, err = processCryptoPayment(ctx, cfg, paymentReq, logger)
		case "stripe":
			if err := validator.ValidateStripeToken(paymentReq.StripeToken); err != nil {
				return &PaymentResponse{Success: false, Message: err.Error()}, nil
			}
			paymentRef, txHash, err = processStripePayment(ctx, cfg, paymentReq, logger)
		case "bank":
			paymentRef, txHash, err = processBankPayment(ctx, cfg, paymentReq, logger)
		default:
			return &PaymentResponse{
				Success: false,
				Message: "Invalid payment method",
			}, nil
		}

		if err != nil {
			logger.Error("Payment processing failed", "error", err, "method", paymentReq.PaymentMethod)
			return &PaymentResponse{
				Success: false,
				Message: fmt.Sprintf("Payment failed: %s", err.Error()),
			}, nil
		}

		// 5. Record payment on-chain
		if err := recordPaymentOnChain(ctx, cfg, paymentReq, paymentRef, logger); err != nil {
			logger.Error("Failed to record payment on-chain", "error", err)
			// Non-fatal: payment was successful but recording failed
		}

		// 6. Store idempotency key
		if paymentReq.IdempotencyKey != "" {
			storeIdempotencyKey(paymentReq.IdempotencyKey)
		}

		logger.Info("Payment completed successfully",
			"tokenID", paymentReq.TokenID,
			"method", paymentReq.PaymentMethod,
			"reference", paymentRef)

		return &PaymentResponse{
			Success:          true,
			Message:          fmt.Sprintf("Payment processed via %s", paymentReq.PaymentMethod),
			TxHash:           txHash,
			PaymentReference: paymentRef,
			Amount:           fmt.Sprintf("%.2f", paymentReq.Amount),
			Currency:         paymentReq.Currency,
			Timestamp:        time.Now().Format(time.RFC3339),
			Method:           paymentReq.PaymentMethod,
		}, nil
	}
}

// HandleBillPayment is the main entry point for bill payment
func HandleBillPayment(ctx context.Context, request PaymentRequest) (*PaymentResponse, error) {
	validator := validation.NewValidator()

	// Validate inputs
	if err := validator.ValidateTokenID(request.TokenID); err != nil {
		return nil, err
	}
	if err := validator.ValidateEthereumAddress(request.OwnerAddress); err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Success:   true,
		Message:   fmt.Sprintf("Payment request validated for bill %d", request.BillIndex),
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// processCryptoPayment processes a cryptocurrency payment
func processCryptoPayment(ctx context.Context, cfg *config.Config, req PaymentRequest, logger *slog.Logger) (string, string, error) {
	logger.Info("Processing crypto payment",
		"tokenID", req.TokenID,
		"billIndex", req.BillIndex)

	// 1. Use Chainlink Price Feeds for amount verification
	amount, err := verifyAmountWithPriceFeed(ctx, req.Amount, req.Currency)
	if err != nil {
		return "", "", fmt.Errorf("price feed verification failed: %w", err)
	}

	// 2. Create payment reference
	paymentRef := fmt.Sprintf("CRYPTO_%d_%d_%d", req.TokenID, req.BillIndex, time.Now().Unix())

	// 3. In a real implementation, this would:
	//    - Create and sign transaction
	//    - Transfer tokens to property owner
	//    - Update bill status on-chain
	//    - Emit payment event

	logger.Info("Crypto payment processed",
		"reference", paymentRef,
		"amount", amount)

	// Simulate transaction hash
	txHash := fmt.Sprintf("0x%s%d", generateRandomHex(64), time.Now().Unix())

	return paymentRef, txHash, nil
}

// processStripePayment processes a payment via Stripe
func processStripePayment(ctx context.Context, cfg *config.Config, req PaymentRequest, logger *slog.Logger) (string, string, error) {
	logger.Info("Processing Stripe payment",
		"tokenID", req.TokenID,
		"billIndex", req.BillIndex)

	// 1. Create idempotency key for Stripe if not provided
	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = fmt.Sprintf("rwa_%d_%d_%d", req.TokenID, req.BillIndex, time.Now().Unix())
	}

	// 2. In a real implementation, this would call Stripe API:
	// stripeReq := map[string]interface{}{
	//     "amount":   int64(req.Amount * 100), // Convert to cents
	//     "currency": req.Currency,
	//     "source":   req.StripeToken,
	// }
	// response, err := callStripeAPI(cfg, stripeReq, idempotencyKey)

	// Simulate Stripe API call
	stripeRef := fmt.Sprintf("stripe_%d_%d_%d", req.TokenID, req.BillIndex, time.Now().Unix())

	logger.Info("Stripe payment processed",
		"reference", stripeRef,
		"idempotencyKey", idempotencyKey)

	// Simulate transaction hash for on-chain recording
	txHash := fmt.Sprintf("0x%s%d", generateRandomHex(64), time.Now().Unix())

	return stripeRef, txHash, nil
}

// processBankPayment processes a traditional bank transfer
func processBankPayment(ctx context.Context, cfg *config.Config, req PaymentRequest, logger *slog.Logger) (string, string, error) {
	logger.Info("Processing bank payment",
		"tokenID", req.TokenID,
		"billIndex", req.BillIndex)

	// 1. Generate payment reference
	paymentRef := fmt.Sprintf("BANK_%d_%d_%d", req.TokenID, req.BillIndex, time.Now().Unix())

	// 2. In a real implementation, this would:
	//    - Generate bank transfer instructions
	//    - Create payment order
	//    - Notify payment processor
	//    - Wait for confirmation webhook

	logger.Info("Bank payment initiated",
		"reference", paymentRef)

	// Simulate transaction hash
	txHash := fmt.Sprintf("0x%s%d", generateRandomHex(64), time.Now().Unix())

	return paymentRef, txHash, nil
}

// verifyAmountWithPriceFeed verifies the payment amount using Chainlink Price Feeds
func verifyAmountWithPriceFeed(ctx context.Context, amount float64, currency string) (float64, error) {
	// In a real implementation, this would:
	// 1. Query Chainlink Price Feed for currency pair
	// 2. Verify amount is within acceptable range
	// 3. Return verified amount

	// For now, return the amount as-is
	return amount, nil
}

// recordPaymentOnChain records the payment on the HouseRWA contract
func recordPaymentOnChain(ctx context.Context, cfg *config.Config, req PaymentRequest, paymentRef string, logger *slog.Logger) error {
	// In a real implementation, this would:
	// 1. Call HouseRWA.recordPayment()
	// 2. Update bill status
	// 3. Emit payment event

	logger.Info("Recording payment on-chain",
		"tokenID", req.TokenID,
		"billIndex", req.BillIndex,
		"reference", paymentRef)

	return nil
}

// callStripeAPI makes a request to the Stripe API
func callStripeAPI(cfg *config.Config, payload map[string]interface{}, idempotencyKey string) (*nethttp.Response, error) {
	// In a real implementation, this would:
	// 1. Load Stripe API key from secrets
	// 2. Create HTTP request with proper authentication
	// 3. Add idempotency key header
	// 4. Send request to Stripe
	// 5. Parse response

	jsonData, _ := json.Marshal(payload)
	req, err := nethttp.NewRequest("POST", cfg.StripeAPIBaseURL+"/charges", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)
	// req.Header.Set("Authorization", "Bearer "+cfg.StripeAPIKey)

	client := &nethttp.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

// isDuplicatePayment checks if a payment with this idempotency key was already processed
func isDuplicatePayment(idempotencyKey string) bool {
	// In a real implementation, this would:
	// 1. Check cache/database for the key
	// 2. Return true if found (within expiry window)
	// 3. Return false otherwise

	return false
}

// storeIdempotencyKey stores the idempotency key for deduplication
func storeIdempotencyKey(idempotencyKey string) {
	// In a real implementation, this would:
	// 1. Store the key in cache/database
	// 2. Set expiry (e.g., 24 hours)
}

// PriceFeedResponse represents a Chainlink Price Feed response
type PriceFeedResponse struct {
	RoundID         uint64 `json:"roundId"`
	Answer          int64  `json:"answer"`
	StartedAt       uint64 `json:"startedAt"`
	UpdatedAt       uint64 `json:"updatedAt"`
	AnsweredInRound uint64 `json:"answeredInRound"`
}
