//go:build wasip1

package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"RWA-Houses/backend/cre/config"
)

const (
	kycProviderNone       = "none"
	kycProviderMock       = "mock"
	kycProviderZKPassport = "zkpassport"
)

type kycVerificationRecord struct {
	Provider         string
	Level            uint8
	VerificationHash common.Hash
	Expiry           time.Time
}

type kycVerifierRequest struct {
	Provider      string          `json:"provider"`
	WalletAddress string          `json:"walletAddress"`
	Proof         json.RawMessage `json:"proof"`
}

type kycVerifierResponse struct {
	Verified         bool   `json:"verified"`
	Level            uint8  `json:"level,omitempty"`
	VerificationHash string `json:"verificationHash,omitempty"`
	ExpiresAt        int64  `json:"expiresAt,omitempty"`
	Message          string `json:"message,omitempty"`
}

func (h *Handler) verifyKYCForAddress(
	cfg *config.Config,
	runtime cre.Runtime,
	address string,
	providerOverride string,
	proof json.RawMessage,
) (kycVerificationRecord, error) {
	provider := strings.ToLower(strings.TrimSpace(providerOverride))
	if provider == "" {
		provider = strings.ToLower(strings.TrimSpace(cfg.KYCProvider))
	}
	if provider == "" {
		provider = kycProviderMock
	}

	record := kycVerificationRecord{
		Provider: provider,
		Level:    2,
		Expiry:   time.Now().Add(180 * 24 * time.Hour),
	}

	switch provider {
	case kycProviderNone:
		return kycVerificationRecord{
			Provider: kycProviderNone,
		}, nil
	case kycProviderMock:
		record.VerificationHash = crypto.Keccak256Hash([]byte(fmt.Sprintf(
			"kyc:mock:%s",
			strings.ToLower(address),
		)))
		return record, nil
	case kycProviderZKPassport:
		if len(bytes.TrimSpace(proof)) == 0 || string(bytes.TrimSpace(proof)) == "null" {
			return kycVerificationRecord{}, fmt.Errorf("zkpassport proof required when kycProvider=zkpassport")
		}

		verified, err := h.verifyWithExternalKYCProvider(cfg, runtime, address, provider, proof)
		if err != nil {
			return kycVerificationRecord{}, err
		}

		record.Level = verified.Level
		record.Expiry = verified.Expiry
		record.VerificationHash = verified.VerificationHash
		return record, nil
	default:
		return kycVerificationRecord{}, fmt.Errorf("unsupported kycProvider: %s", provider)
	}
}

func (h *Handler) verifyWithExternalKYCProvider(
	cfg *config.Config,
	runtime cre.Runtime,
	address string,
	provider string,
	proof json.RawMessage,
) (kycVerificationRecord, error) {
	if strings.TrimSpace(cfg.KYCVerifierURL) == "" {
		return kycVerificationRecord{}, fmt.Errorf("kyc verifier URL is not configured")
	}

	payload := kycVerifierRequest{
		Provider:      provider,
		WalletAddress: strings.ToLower(address),
		Proof:         proof,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return kycVerificationRecord{}, fmt.Errorf("failed to marshal kyc verifier request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, cfg.KYCVerifierURL, bytes.NewReader(body))
	if err != nil {
		return kycVerificationRecord{}, fmt.Errorf("failed to create kyc verifier request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.KYCProviderKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.KYCProviderKey)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return kycVerificationRecord{}, fmt.Errorf("failed to call kyc verifier: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return kycVerificationRecord{}, fmt.Errorf("failed to read kyc verifier response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return kycVerificationRecord{}, fmt.Errorf(
			"kyc verifier returned status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(responseBody)),
		)
	}

	var verifierResp kycVerifierResponse
	if err := json.Unmarshal(responseBody, &verifierResp); err != nil {
		return kycVerificationRecord{}, fmt.Errorf("invalid kyc verifier response: %w", err)
	}

	if !verifierResp.Verified {
		msg := strings.TrimSpace(verifierResp.Message)
		if msg == "" {
			msg = "proof not verified"
		}
		return kycVerificationRecord{}, fmt.Errorf("kyc verification failed: %s", msg)
	}

	level := verifierResp.Level
	if level == 0 {
		level = 2
	}

	expiry := time.Now().Add(180 * 24 * time.Hour)
	if verifierResp.ExpiresAt > 0 {
		expiry = time.Unix(verifierResp.ExpiresAt, 0)
	}

	verificationHash := crypto.Keccak256Hash(
		[]byte("kyc:zkpassport:"),
		[]byte(strings.ToLower(address)),
		proof,
	)
	if strings.TrimSpace(verifierResp.VerificationHash) != "" {
		parsedHash, err := parseOptionalHexHash(verifierResp.VerificationHash)
		if err != nil {
			return kycVerificationRecord{}, fmt.Errorf("invalid verificationHash from kyc verifier: %w", err)
		}
		verificationHash = parsedHash
	}

	runtime.Logger().Info("kyc verified by external provider",
		"provider", provider,
		"address", strings.ToLower(address),
		"level", level,
	)

	return kycVerificationRecord{
		Provider:         provider,
		Level:            level,
		VerificationHash: verificationHash,
		Expiry:           expiry,
	}, nil
}

func parseOptionalHexHash(value string) (common.Hash, error) {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimPrefix(trimmed, "0x")
	trimmed = strings.TrimPrefix(trimmed, "0X")
	if len(trimmed) != 64 {
		return common.Hash{}, fmt.Errorf("expected 32-byte hex hash")
	}

	decoded, err := hex.DecodeString(trimmed)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(decoded), nil
}
