//go:build wasip1

// Package workflows implements the house minting workflow for Chainlink CRE.
// This workflow handles the tokenization of real estate assets as RWA NFTs.
package workflows

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/pkg/encryption"
	"RWA-Houses/backend/cre/pkg/validation"
)

// MintRequest represents the house minting request
type MintRequest struct {
	Action         string `json:"action"`
	OwnerAddress   string `json:"ownerAddress"`
	HouseID        string `json:"houseID"`
	Location       string `json:"location"`
	Value          string `json:"value"`
	DocumentsB64   string `json:"documentsB64"`
	StorageType    string `json:"storageType"`
	OwnerPublicKey string `json:"ownerPublicKey"`
}

// MintResponse represents the house minting response
type MintResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	TokenID      uint64   `json:"tokenID,omitempty"`
	TxHash       string   `json:"txHash,omitempty"`
	DocumentURI  string   `json:"documentURI,omitempty"`
	DocumentHash string   `json:"documentHash,omitempty"`
	EncryptedKey string   `json:"encryptedKey,omitempty"`
	KeyShares    []string `json:"keyShares,omitempty"`
	Threshold    int      `json:"threshold,omitempty"`
	TotalShares  int      `json:"totalShares,omitempty"`
}

// MintHouseWorkflow handles the complete house minting process
func MintHouseWorkflow(cfg *config.Config, logger *slog.Logger) func(ctx context.Context, req *http.Payload) (*MintResponse, error) {
	validator := validation.NewValidatorWithOptions(cfg.MaxDocumentSize)

	return func(ctx context.Context, req *http.Payload) (*MintResponse, error) {
		logger.Info("Starting house minting workflow")

		// 1. Parse request
		var mintReq MintRequest
		if err := json.Unmarshal(req.Input, &mintReq); err != nil {
			logger.Error("Failed to parse mint request", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Invalid request format",
			}, nil
		}

		// 2. Validate KYC status (if enforced)
		if cfg.EnforceKYC {
			if err := validateKYC(ctx, mintReq.OwnerAddress, logger); err != nil {
				logger.Error("KYC validation failed", "error", err, "address", mintReq.OwnerAddress)
				return &MintResponse{
					Success: false,
					Message: "KYC verification required",
				}, nil
			}
		}

		// 3. Validate inputs
		if err := validator.ValidateEthereumAddress(mintReq.OwnerAddress); err != nil {
			return &MintResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateHouseID(mintReq.HouseID); err != nil {
			return &MintResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidateLocation(mintReq.Location); err != nil {
			return &MintResponse{Success: false, Message: err.Error()}, nil
		}
		if err := validator.ValidatePublicKey(mintReq.OwnerPublicKey); err != nil {
			return &MintResponse{Success: false, Message: err.Error()}, nil
		}

		docs, err := validator.ValidateDocument(mintReq.DocumentsB64)
		if err != nil {
			return &MintResponse{Success: false, Message: err.Error()}, nil
		}

		logger.Info("Processing house mint",
			"owner", mintReq.OwnerAddress,
			"houseID", mintReq.HouseID,
			"documentSize", len(docs))

		// 4. Encrypt documents with AES-256-GCM + threshold encryption
		masterKey, err := encryption.GenerateRandomKey(32)
		if err != nil {
			logger.Error("Failed to generate master key", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Key generation failed",
			}, nil
		}

		// Encrypt the document
		passphrase := hex.EncodeToString(masterKey)
		encryptedDoc, err := encryption.EncryptDocument(docs, passphrase)
		if err != nil {
			logger.Error("Document encryption failed", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Encryption failed",
			}, nil
		}

		// Generate threshold shares
		scheme, err := encryption.NewThresholdScheme(cfg.ThresholdKeyThreshold, cfg.ThresholdKeyTotal)
		if err != nil {
			logger.Error("Failed to create threshold scheme", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Threshold scheme creation failed",
			}, nil
		}

		shares, err := scheme.GenerateThresholdShares(masterKey)
		if err != nil {
			logger.Error("Failed to generate key shares", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Key share generation failed",
			}, nil
		}

		logger.Info("Generated threshold shares",
			"threshold", cfg.ThresholdKeyThreshold,
			"total", cfg.ThresholdKeyTotal)

		// 5. Upload to IPFS or store off-chain
		documentURI, err := uploadEncryptedDocument(ctx, cfg, encryptedDoc, logger)
		if err != nil {
			logger.Error("Document upload failed", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Document upload failed",
			}, nil
		}

		// 6. Generate document hash
		docHash := encryption.HashShares(shares)

		// 7. Encrypt master key for owner (ECIES)
		ownerPubKey, err := parseOwnerPublicKey(mintReq.OwnerPublicKey)
		if err != nil {
			logger.Error("Failed to parse owner public key", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Invalid owner public key",
			}, nil
		}

		encryptedKeyForOwner, err := encryption.EncryptWithPublicKey(masterKey, ownerPubKey)
		if err != nil {
			logger.Error("Failed to encrypt key for owner", "error", err)
			return &MintResponse{
				Success: false,
				Message: "Key encryption failed",
			}, nil
		}

		// 8. Call HouseRWA.sol mintHouse() - simulated
		tokenID, txHash, err := mintHouseOnChain(ctx, cfg, mintReq, documentURI, docHash, logger)
		if err != nil {
			logger.Error("On-chain minting failed", "error", err)
			return &MintResponse{
				Success: false,
				Message: "On-chain minting failed",
			}, nil
		}

		// Serialize shares for response
		shareStrings := make([]string, len(shares))
		for i, share := range shares {
			shareData, _ := json.Marshal(share)
			shareStrings[i] = base64.StdEncoding.EncodeToString(shareData)
		}

		logger.Info("House minting completed successfully",
			"tokenID", tokenID,
			"houseID", mintReq.HouseID)

		return &MintResponse{
			Success:      true,
			Message:      "House minted successfully",
			TokenID:      tokenID,
			TxHash:       txHash,
			DocumentURI:  documentURI,
			DocumentHash: docHash,
			EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKeyForOwner),
			KeyShares:    shareStrings,
			Threshold:    cfg.ThresholdKeyThreshold,
			TotalShares:  cfg.ThresholdKeyTotal,
		}, nil
	}
}

// HandleMintHouse is the main entry point for the mint workflow
func HandleMintHouse(ctx context.Context, request MintRequest) (*MintResponse, error) {
	// This is a simplified version for direct invocation
	// In CRE, the workflow is invoked through the handler

	validator := validation.NewValidator()

	// Validate inputs
	if err := validator.ValidateEthereumAddress(request.OwnerAddress); err != nil {
		return nil, err
	}
	if err := validator.ValidateHouseID(request.HouseID); err != nil {
		return nil, err
	}

	return &MintResponse{
		Success: true,
		Message: fmt.Sprintf("Mint request validated for house %s", request.HouseID),
	}, nil
}

// validateKYC validates the KYC status of an address
func validateKYC(ctx context.Context, address string, logger *slog.Logger) error {
	// In a real implementation, this would:
	// 1. Query a KYC provider API
	// 2. Verify KYC status on-chain
	// 3. Check against whitelist

	logger.Info("Validating KYC", "address", address)

	// Simulate KYC check - always pass for demo
	return nil
}

// uploadEncryptedDocument uploads the encrypted document to storage
func uploadEncryptedDocument(ctx context.Context, cfg *config.Config, doc *encryption.EncryptedDoc, logger *slog.Logger) (string, error) {
	// Serialize encrypted document
	docData, err := encryption.SerializeEncryptedDoc(doc)
	if err != nil {
		return "", fmt.Errorf("failed to serialize document: %w", err)
	}

	// In a real implementation, this would:
	// 1. Upload to IPFS using the configured gateway
	// 2. Or store in off-chain encrypted storage
	// 3. Return the URI

	docHash := encryption.HashShares([]encryption.Share{
		{Index: 1, Value: docData},
	})

	// Return IPFS URI (simulated)
	return fmt.Sprintf("ipfs://Qm%s", docHash[:46]), nil
}

// mintHouseOnChain simulates minting the house on-chain
func mintHouseOnChain(ctx context.Context, cfg *config.Config, req MintRequest, documentURI, docHash string, logger *slog.Logger) (uint64, string, error) {
	// In a real implementation, this would:
	// 1. Create EVM client connection
	// 2. Build and sign transaction
	// 3. Call HouseRWA.mintHouse()
	// 4. Wait for confirmation
	// 5. Return token ID and tx hash

	logger.Info("Minting house on-chain",
		"owner", req.OwnerAddress,
		"houseID", req.HouseID)

	// Simulate token ID generation
	tokenID := uint64(time.Now().Unix())
	txHash := fmt.Sprintf("0x%s", docHash[:64])

	return tokenID, txHash, nil
}

// parseOwnerPublicKey parses the owner's public key
func parseOwnerPublicKey(pubKeyStr string) (*ecdsa.PublicKey, error) {
	// In a real implementation, this would parse PEM or hex-encoded public key
	// For now, return a placeholder
	return nil, nil
}
