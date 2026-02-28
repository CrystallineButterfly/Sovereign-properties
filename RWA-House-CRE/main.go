//go:build wasip1

package main

import (
	"context"
	"log/slog"
	"os"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/handlers"
	"RWA-Houses/backend/cre/workflows"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

func initWorkflow(
	cfg *config.Config,
	logger *slog.Logger,
	_ cre.SecretsProvider,
) (cre.Workflow[*config.Config], error) {
	logger.Info(
		"initializing RWA House workflow",
		"chain", cfg.EVMChain,
		"chainId", cfg.EVMChainID,
		"houseRWA", cfg.HouseRWAContractAddr,
		"receiver", cfg.HouseRWAReceiverAddr,
	)

	handler := handlers.NewHandler(cfg)
	workflow := cre.Workflow[*config.Config]{}

	workflow = append(workflow, cre.Handler(
		http.Trigger(&http.Config{}),
		func(c *config.Config, runtime cre.Runtime, trigger *http.Payload) ([]byte, error) {
			return handler.HandleHTTPAction(c, runtime, trigger)
		},
	))

	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 0 * * *"}),
		func(c *config.Config, runtime cre.Runtime, trigger *cron.Payload) (string, error) {
			return workflows.AutomatedRentalPaymentWorkflow(c, runtime.Logger())(
				context.Background(),
				trigger,
			)
		},
	))

	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 */6 * * *"}),
		func(*config.Config, cre.Runtime, *cron.Payload) (string, error) {
			return "key cleanup tick completed", nil
		},
	))

	workflow = append(workflow, cre.Handler(
		cron.Trigger(&cron.Config{Schedule: "0 9 * * *"}),
		func(*config.Config, cre.Runtime, *cron.Payload) (string, error) {
			return "bill reminder tick completed", nil
		},
	))

	return workflow, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger.Info("starting RWA House CRE workflow")

	runner := wasm.NewRunner(cre.ParseJSON[config.Config])
	runner.Run(initWorkflow)
}
