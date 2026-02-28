//go:build wasip1

// Package main implements the RWA House Chainlink CRE workflow.
// This is the entry point for the CRE workflow that handles house tokenization,
// listings, sales, rentals, and bill payments.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/handlers"
	"RWA-Houses/backend/cre/workflows"
)

// WorkflowConfig holds the parsed configuration
type WorkflowConfig struct {
	Config *config.Config
}

// InitWorkflow initializes the CRE workflow with all handlers
func InitWorkflow(cfg *config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*config.Config], error) {
	logger.Info("Initializing RWA House Workflow",
		"chainID", cfg.EVMChainID,
		"contract", cfg.HouseRWAContractAddr)

	// Create main handler
	handler := handlers.NewHandler(cfg)

	// Initialize workflow
	workflow := cre.Workflow[*config.Config]{}

	// 1. HTTP Trigger for user actions (mint, create_listing, sell, rent, pay_bill, create_bill, claim_key)
	workflow = append(workflow, cre.Handler(
		http.Trigger(&http.Config{}),
		func(config *config.Config, runtime cre.Runtime, trigger *http.Payload) ([]byte, error) {
			return handler.HandleHTTPAction(config, runtime, trigger)
		},
	))

	// 2. Cron Trigger for automated rental payments (daily at midnight UTC)
	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 0 * * *"}),
		func(config *config.Config, runtime cre.Runtime, trigger *cron.Payload) (string, error) {
			return workflows.AutomatedRentalPaymentWorkflow(config, runtime.Logger())(context.Background(), trigger)
		},
	))

	// 3. Cron Trigger for key expiry cleanup (every 6 hours)
	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 */6 * * *"}),
		handleKeyExpiryCleanup,
	))

	// 4. Cron Trigger for bill payment reminders (daily at 9 AM UTC)
	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 9 * * *"}),
		handleBillPaymentReminders,
	))

	logger.Info("Workflow initialized successfully",
		"handlers", len(workflow))

	return workflow, nil
}

// handleKeyExpiryCleanup handles cleanup of expired keys
func handleKeyExpiryCleanup(cfg *config.Config, runtime cre.Runtime, trigger *cron.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running key expiry cleanup")

	// In a real implementation, this would:
	// 1. Query all temporary keys
	// 2. Check expiry timestamps
	// 3. Remove expired keys
	// 4. Update on-chain state

	return "Key expiry cleanup completed", nil
}

// handleBillPaymentReminders handles sending payment reminders
func handleBillPaymentReminders(cfg *config.Config, runtime cre.Runtime, trigger *cron.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running bill payment reminders")

	// In a real implementation, this would:
	// 1. Query bills due soon
	// 2. Send notifications to owners
	// 3. Prepare auto-pay for enabled accounts

	return "Bill payment reminders sent", nil
}

// main is the entry point for the CRE workflow
func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Starting RWA House CRE Workflow")

	// Create and run the workflow runner
	// The wasm.NewRunner will:
	// 1. Parse the config from CRE
	// 2. Initialize the workflow
	// 3. Start handling triggers
	runner := wasm.NewRunner(cre.ParseJSON[config.Config])
	runner.Run(InitWorkflow)
}

// Response types for internal use
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// createHealthResponse creates a health check response
func createHealthResponse() []byte {
	resp := HealthResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: "", // Will be set at runtime
	}
	data, _ := json.Marshal(resp)
	return data
}
