//go:build wasip1

// Package handlers provides HTTP handlers for the Chainlink CRE workflow.
package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	evmcap "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"RWA-Houses/backend/cre/config"
	"RWA-Houses/backend/cre/pkg/encryption"
	"RWA-Houses/backend/cre/pkg/validation"
)

const houseRWAABIJSON = `[
  {"type":"function","name":"setKYCVerification","stateMutability":"nonpayable","inputs":[{"name":"user","type":"address"},{"name":"level","type":"uint8"},{"name":"verificationHash","type":"bytes32"},{"name":"expiryDate","type":"uint48"}],"outputs":[]},
  {"type":"function","name":"mint","stateMutability":"nonpayable","inputs":[{"name":"to","type":"address"},{"name":"houseId","type":"string"},{"name":"documentHash","type":"bytes32"},{"name":"documentURI","type":"string"},{"name":"storageType","type":"uint8"},{"name":"verificationData","type":"string"}],"outputs":[{"name":"","type":"uint256"}]},
  {"type":"function","name":"createListingFromWorkflow","stateMutability":"nonpayable","inputs":[{"name":"tokenId","type":"uint256"},{"name":"owner","type":"address"},{"name":"listingType","type":"uint8"},{"name":"price","type":"uint96"},{"name":"preferredToken","type":"address"},{"name":"isPrivateSale","type":"bool"},{"name":"allowedBuyer","type":"address"},{"name":"durationDays","type":"uint48"}],"outputs":[]},
  {"type":"function","name":"completeSale","stateMutability":"nonpayable","inputs":[{"name":"tokenId","type":"uint256"},{"name":"buyer","type":"address"},{"name":"keyHash","type":"bytes32"},{"name":"encryptedKey","type":"bytes"}],"outputs":[]},
  {"type":"function","name":"startRental","stateMutability":"nonpayable","inputs":[{"name":"tokenId","type":"uint256"},{"name":"renter","type":"address"},{"name":"durationDays","type":"uint48"},{"name":"depositAmount","type":"uint96"},{"name":"monthlyRent","type":"uint96"},{"name":"encryptedAccessKey","type":"bytes"}],"outputs":[]},
  {"type":"function","name":"createBill","stateMutability":"nonpayable","inputs":[{"name":"tokenId","type":"uint256"},{"name":"billType","type":"string"},{"name":"amount","type":"uint96"},{"name":"dueDate","type":"uint48"},{"name":"provider","type":"address"},{"name":"isRecurring","type":"bool"},{"name":"recurrenceInterval","type":"uint8"}],"outputs":[{"name":"billIndex","type":"uint256"}]},
  {"type":"function","name":"recordBillPayment","stateMutability":"nonpayable","inputs":[{"name":"tokenId","type":"uint256"},{"name":"billIndex","type":"uint256"},{"name":"paymentMethod","type":"string"},{"name":"paymentReference","type":"bytes32"}],"outputs":[]},
  {"type":"function","name":"nextTokenId","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"uint256"}]},
  {"type":"function","name":"getTotalBillsCount","stateMutability":"view","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"name":"","type":"uint256"}]},
  {"type":"function","name":"keyExchanges","stateMutability":"view","inputs":[{"name":"","type":"bytes32"}],"outputs":[{"name":"keyHash","type":"bytes32"},{"name":"encryptedKey","type":"bytes"},{"name":"intendedRecipient","type":"address"},{"name":"createdAt","type":"uint48"},{"name":"expiresAt","type":"uint48"},{"name":"isClaimed","type":"bool"},{"name":"exchangeType","type":"uint8"}]}
]`

var houseRWAABI = mustParseABI(houseRWAABIJSON)

// Handler handles HTTP requests.
type Handler struct {
	config    *config.Config
	validator *validation.Validator
}

// NewHandler creates a new HTTP handler.
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config:    cfg,
		validator: validation.NewValidatorWithOptions(cfg.MaxDocumentSize),
	}
}

// Response represents an HTTP response.
type Response struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	TxHash       string      `json:"txHash,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	EncryptedKey string      `json:"encryptedKey,omitempty"`
}

// BaseRequest represents the base request structure.
type BaseRequest struct {
	Action string `json:"action"`
}

// MintRequest represents a house minting request.
// Supports both houseID and houseId to align web/mobile payload variants.
type MintRequest struct {
	Action         string          `json:"action"`
	OwnerAddress   string          `json:"ownerAddress"`
	HouseID        string          `json:"houseID,omitempty"`
	HouseId        string          `json:"houseId,omitempty"`
	Location       string          `json:"location,omitempty"`
	Value          string          `json:"value,omitempty"`
	DocumentsB64   string          `json:"documentsB64"`
	StorageType    string          `json:"storageType"`
	OwnerPublicKey string          `json:"ownerPublicKey"`
	KYCProvider    string          `json:"kycProvider,omitempty"`
	KYCProof       json.RawMessage `json:"kycProof,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}

// SaleRequest represents a house sale request.
type SaleRequest struct {
	Action         string          `json:"action"`
	SellerAddress  string          `json:"sellerAddress,omitempty"`
	BuyerAddress   string          `json:"buyerAddress"`
	TokenID        uint64          `json:"tokenID,omitempty"`
	TokenId        string          `json:"tokenId,omitempty"`
	Price          string          `json:"price"`
	BuyerPublicKey string          `json:"buyerPublicKey,omitempty"`
	IsPrivateSale  bool            `json:"isPrivateSale"`
	KYCProvider    string          `json:"kycProvider,omitempty"`
	KYCProof       json.RawMessage `json:"kycProof,omitempty"`
}

// RentalRequest represents a rental request.
type RentalRequest struct {
	Action          string          `json:"action"`
	TokenID         uint64          `json:"tokenID,omitempty"`
	TokenId         string          `json:"tokenId,omitempty"`
	RenterAddress   string          `json:"renterAddress"`
	DurationDays    uint64          `json:"durationDays"`
	MonthlyRent     string          `json:"monthlyRent"`
	DepositAmount   string          `json:"depositAmount,omitempty"`
	RenterPublicKey string          `json:"renterPublicKey,omitempty"`
	KYCProvider     string          `json:"kycProvider,omitempty"`
	KYCProof        json.RawMessage `json:"kycProof,omitempty"`
}

// CreateListingRequest represents a listing creation request.
type CreateListingRequest struct {
	Action         string          `json:"action"`
	TokenID        uint64          `json:"tokenID,omitempty"`
	TokenId        string          `json:"tokenId,omitempty"`
	OwnerAddress   string          `json:"ownerAddress"`
	ListingType    string          `json:"listingType"`
	Price          string          `json:"price"`
	PreferredToken string          `json:"preferredToken,omitempty"`
	IsPrivateSale  bool            `json:"isPrivateSale"`
	AllowedBuyer   string          `json:"allowedBuyer,omitempty"`
	DurationDays   uint64          `json:"durationDays,omitempty"`
	KYCProvider    string          `json:"kycProvider,omitempty"`
	KYCProof       json.RawMessage `json:"kycProof,omitempty"`
}

// PaymentRequest represents a bill payment request.
type PaymentRequest struct {
	Action        string `json:"action"`
	TokenID       uint64 `json:"tokenID,omitempty"`
	TokenId       string `json:"tokenId,omitempty"`
	BillIndex     uint64 `json:"billIndex"`
	OwnerAddress  string `json:"ownerAddress,omitempty"`
	PaymentMethod string `json:"paymentMethod"`
	StripeToken   string `json:"stripeToken,omitempty"`
}

// CreateBillRequest represents a bill creation request.
type CreateBillRequest struct {
	Action             string  `json:"action"`
	TokenID            uint64  `json:"tokenID,omitempty"`
	TokenId            string  `json:"tokenId,omitempty"`
	BillType           string  `json:"billType"`
	Amount             float64 `json:"amount"`
	DueDate            string  `json:"dueDate"`
	Provider           string  `json:"provider"`
	IsRecurring        bool    `json:"isRecurring"`
	RecurrenceInterval uint8   `json:"recurrenceInterval,omitempty"`
}

// ClaimKeyRequest represents a key claim request.
type ClaimKeyRequest struct {
	Action   string `json:"action"`
	KeyHash  string `json:"keyHash"`
	Claimant string `json:"claimant,omitempty"`
}

type keyExchangeView struct {
	KeyHash           [32]byte
	EncryptedKey      []byte
	IntendedRecipient common.Address
	CreatedAt         *big.Int
	ExpiresAt         *big.Int
	IsClaimed         bool
	ExchangeType      uint8
}

// HandleHTTPAction routes HTTP requests to appropriate handlers.
func (h *Handler) HandleHTTPAction(cfg *config.Config, runtime cre.Runtime, trigger *http.Payload) ([]byte, error) {
	logger := runtime.Logger()

	var baseReq BaseRequest
	if err := json.Unmarshal(trigger.Input, &baseReq); err != nil {
		logger.Error("failed to parse request", "error", err)
		return createErrorResponse("invalid request format")
	}

	action := strings.ToLower(strings.TrimSpace(baseReq.Action))
	logger.Info("processing action", "action", action)

	var (
		resp Response
		err  error
	)

	switch action {
	case "mint":
		resp, err = h.handleMint(cfg, runtime, trigger.Input)
	case "create_listing":
		resp, err = h.handleCreateListing(cfg, runtime, trigger.Input)
	case "sell":
		resp, err = h.handleSell(cfg, runtime, trigger.Input)
	case "rent":
		resp, err = h.handleRent(cfg, runtime, trigger.Input)
	case "pay_bill":
		resp, err = h.handleBillPayment(cfg, runtime, trigger.Input)
	case "create_bill":
		resp, err = h.handleCreateBill(cfg, runtime, trigger.Input)
	case "claim_key":
		resp, err = h.handleClaimKey(cfg, runtime, trigger.Input)
	default:
		return createErrorResponse(fmt.Sprintf("unknown action: %s", baseReq.Action))
	}

	if err != nil {
		logger.Error("action failed", "action", action, "error", err)
		return createErrorResponse(err.Error())
	}

	return json.Marshal(resp)
}

// handleMint processes house minting requests and writes HouseRWA.mint via CRE report.
func (h *Handler) handleMint(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req MintRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid mint request"}, err
	}

	houseID := strings.TrimSpace(req.HouseID)
	if houseID == "" {
		houseID = strings.TrimSpace(req.HouseId)
	}

	if err := h.validator.ValidateEthereumAddress(req.OwnerAddress); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateHouseID(houseID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if req.Location != "" {
		if err := h.validator.ValidateLocation(req.Location); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}
	if err := h.validator.ValidatePublicKey(req.OwnerPublicKey); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateStorageType(req.StorageType); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	docs, err := h.validator.ValidateDocument(req.DocumentsB64)
	if err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	logger.Info("processing mint", "owner", req.OwnerAddress, "houseID", houseID)

	ownerKYC, err := h.verifyKYCForAddress(cfg, runtime, req.OwnerAddress, req.KYCProvider, req.KYCProof)
	if err != nil {
		return Response{Success: false, Message: "owner KYC verification failed"}, err
	}

	if err := h.writeKYCVerification(cfg, runtime, req.OwnerAddress, ownerKYC); err != nil {
		return Response{Success: false, Message: "owner KYC write failed"}, err
	}

	masterKey, err := encryption.GenerateRandomKey(32)
	if err != nil {
		return Response{Success: false, Message: "key generation failed"}, err
	}

	_, shares, err := encryption.EncryptWithThreshold(
		docs,
		hex.EncodeToString(masterKey),
		cfg.ThresholdKeyThreshold,
		cfg.ThresholdKeyTotal,
	)
	if err != nil {
		return Response{Success: false, Message: "encryption failed"}, err
	}

	var privacySalt [32]byte
	if _, err := rand.Read(privacySalt[:]); err != nil {
		return Response{Success: false, Message: "privacy salt generation failed"}, err
	}

	docDigestBytes := sha256.Sum256(append(privacySalt[:], docs...))
	docDigest := docDigestBytes
	var uriNonce [18]byte
	if _, err := rand.Read(uriNonce[:]); err != nil {
		return Response{Success: false, Message: "private pointer generation failed"}, err
	}
	documentURI := fmt.Sprintf("cre://private/%x", uriNonce[:])
	storageEnum := mapStorageType(req.StorageType)
	metadataCommitment := buildMintMetadataCommitment(
		req.OwnerAddress,
		houseID,
		req.Metadata,
		docDigest,
		documentURI,
		storageEnum,
		privacySalt,
	)
	verificationData := buildMintVerificationData(
		req,
		docDigest,
		encryption.HashShares(shares),
		metadataCommitment,
	)

	nextTokenID, nextTokenErr := h.readUint256Method(cfg, runtime, "nextTokenId")
	if nextTokenErr != nil {
		logger.Warn("nextTokenId read failed before mint (tokenId may be empty in simulation)", "error", nextTokenErr)
	}

	calldata, err := houseRWAABI.Pack(
		"mint",
		common.HexToAddress(req.OwnerAddress),
		metadataCommitment,
		docDigest,
		documentURI,
		storageEnum,
		verificationData,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode mint calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain mint failed"}, err
	}

	txHash := txHashHex(writeReply.GetTxHash())
	tokenID := ""
	if nextTokenID != nil {
		tokenID = nextTokenID.String()
	}

	return Response{
		Success:      true,
		Message:      "house minted successfully with private onchain commitment",
		TxHash:       txHash,
		EncryptedKey: base64.StdEncoding.EncodeToString(masterKey),
		Data: map[string]interface{}{
			"tokenId":            tokenID,
			"tokenID":            tokenID,
			"houseID":            houseID,
			"metadataCommitment": metadataCommitment,
			"kycProvider":        ownerKYC.Provider,
			"documentHash":       "0x" + hex.EncodeToString(docDigest[:]),
			"storageType":        strings.ToLower(strings.TrimSpace(req.StorageType)),
			"documentURI":        documentURI,
			"sharesCount":        len(shares),
			"threshold":          cfg.ThresholdKeyThreshold,
		},
	}, nil
}

// handleCreateListing processes listing requests and writes HouseRWA.createListingFromWorkflow via CRE report.
func (h *Handler) handleCreateListing(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req CreateListingRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid listing request"}, err
	}

	tokenID, err := parseTokenID(req.TokenID, req.TokenId)
	if err != nil {
		return Response{Success: false, Message: "invalid token ID"}, nil
	}

	if err := h.validator.ValidateTokenID(tokenID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateEthereumAddress(req.OwnerAddress); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	listingTypeValue, listingTypeLabel, err := parseListingType(req.ListingType)
	if err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	price, err := parsePositiveBigInt(req.Price, "price")
	if err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if _, err := asUint96(price, "price"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	preferredToken := common.Address{}
	if strings.TrimSpace(req.PreferredToken) != "" {
		if !common.IsHexAddress(req.PreferredToken) {
			return Response{Success: false, Message: "invalid preferredToken address format"}, nil
		}
		preferredToken = common.HexToAddress(req.PreferredToken)
	}

	allowedBuyer := common.Address{}
	if req.IsPrivateSale {
		if err := h.validator.ValidateEthereumAddress(req.AllowedBuyer); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
		allowedBuyer = common.HexToAddress(req.AllowedBuyer)
	}

	if req.DurationDays > 0 {
		if err := h.validator.ValidateDuration(req.DurationDays); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}

	durationBI := new(big.Int).SetUint64(req.DurationDays)
	if _, err := asUint48(durationBI, "durationDays"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	logger.Info(
		"processing listing creation",
		"tokenID", tokenID,
		"owner", req.OwnerAddress,
		"listingType", listingTypeLabel,
		"price", req.Price,
		"isPrivateSale", req.IsPrivateSale,
	)

	ownerKYC, err := h.verifyKYCForAddress(cfg, runtime, req.OwnerAddress, req.KYCProvider, req.KYCProof)
	if err != nil {
		return Response{Success: false, Message: "owner KYC verification failed"}, err
	}
	if err := h.writeKYCVerification(cfg, runtime, req.OwnerAddress, ownerKYC); err != nil {
		return Response{Success: false, Message: "owner KYC write failed"}, err
	}

	calldata, err := houseRWAABI.Pack(
		"createListingFromWorkflow",
		new(big.Int).SetUint64(tokenID),
		common.HexToAddress(req.OwnerAddress),
		listingTypeValue,
		price,
		preferredToken,
		req.IsPrivateSale,
		allowedBuyer,
		durationBI,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode listing calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain listing creation failed"}, err
	}

	allowedBuyerValue := ""
	if req.IsPrivateSale {
		allowedBuyerValue = allowedBuyer.Hex()
	}

	return Response{
		Success: true,
		Message: "listing created successfully",
		TxHash:  txHashHex(writeReply.GetTxHash()),
		Data: map[string]interface{}{
			"tokenID":        tokenID,
			"ownerAddress":   req.OwnerAddress,
			"listingType":    listingTypeLabel,
			"price":          price.String(),
			"preferredToken": preferredToken.Hex(),
			"isPrivateSale":  req.IsPrivateSale,
			"allowedBuyer":   allowedBuyerValue,
			"durationDays":   req.DurationDays,
			"kycProvider":    ownerKYC.Provider,
		},
	}, nil
}

// handleSell processes house sale requests and writes HouseRWA.completeSale via CRE report.
func (h *Handler) handleSell(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req SaleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid sale request"}, err
	}

	tokenID, err := parseTokenID(req.TokenID, req.TokenId)
	if err != nil {
		return Response{Success: false, Message: "invalid token ID"}, nil
	}

	if req.SellerAddress != "" {
		if err := h.validator.ValidateEthereumAddress(req.SellerAddress); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}
	if err := h.validator.ValidateEthereumAddress(req.BuyerAddress); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateTokenID(tokenID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidatePrice(req.Price); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if req.BuyerPublicKey != "" {
		if err := h.validator.ValidatePublicKey(req.BuyerPublicKey); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}

	logger.Info("processing sale", "tokenID", tokenID, "buyer", req.BuyerAddress, "price", req.Price)

	buyerKYC, err := h.verifyKYCForAddress(cfg, runtime, req.BuyerAddress, req.KYCProvider, req.KYCProof)
	if err != nil {
		return Response{Success: false, Message: "buyer KYC verification failed"}, err
	}

	if err := h.writeKYCVerification(cfg, runtime, req.BuyerAddress, buyerKYC); err != nil {
		return Response{Success: false, Message: "buyer KYC write failed"}, err
	}

	transferKey, err := encryption.GenerateRandomKey(32)
	if err != nil {
		return Response{Success: false, Message: "failed to generate transfer key"}, err
	}

	encryptedKey := transferKey
	keyHash := crypto.Keccak256Hash(encryptedKey)

	calldata, err := houseRWAABI.Pack(
		"completeSale",
		new(big.Int).SetUint64(tokenID),
		common.HexToAddress(req.BuyerAddress),
		keyHash,
		encryptedKey,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode sale calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain sale failed"}, err
	}

	txHash := txHashHex(writeReply.GetTxHash())
	keyHashHex := keyHash.Hex()
	encryptedKeyB64 := base64.StdEncoding.EncodeToString(encryptedKey)

	return Response{
		Success:      true,
		Message:      "sale completed successfully",
		TxHash:       txHash,
		EncryptedKey: encryptedKeyB64,
		Data: map[string]interface{}{
			"tokenID":       tokenID,
			"buyer":         req.BuyerAddress,
			"price":         req.Price,
			"isPrivateSale": req.IsPrivateSale,
			"kycProvider":   buyerKYC.Provider,
			"keyHash":       keyHashHex,
			"encryptedKey":  encryptedKeyB64,
		},
	}, nil
}

// handleRent processes rental requests and writes HouseRWA.startRental via CRE report.
func (h *Handler) handleRent(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req RentalRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid rental request"}, err
	}

	tokenID, err := parseTokenID(req.TokenID, req.TokenId)
	if err != nil {
		return Response{Success: false, Message: "invalid token ID"}, nil
	}

	if err := h.validator.ValidateTokenID(tokenID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateEthereumAddress(req.RenterAddress); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateDuration(req.DurationDays); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateAmount(req.MonthlyRent); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if req.RenterPublicKey != "" {
		if err := h.validator.ValidatePublicKey(req.RenterPublicKey); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}

	monthlyRent, err := parsePositiveBigInt(req.MonthlyRent, "monthlyRent")
	if err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if _, err := asUint96(monthlyRent, "monthlyRent"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	depositAmount := new(big.Int).Set(monthlyRent)
	if strings.TrimSpace(req.DepositAmount) != "" {
		depositAmount, err = parsePositiveBigInt(req.DepositAmount, "depositAmount")
		if err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}
	if _, err := asUint96(depositAmount, "depositAmount"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	durationBI := new(big.Int).SetUint64(req.DurationDays)
	if _, err := asUint48(durationBI, "durationDays"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	logger.Info("processing rental", "tokenID", tokenID, "renter", req.RenterAddress, "durationDays", req.DurationDays)

	renterKYC, err := h.verifyKYCForAddress(cfg, runtime, req.RenterAddress, req.KYCProvider, req.KYCProof)
	if err != nil {
		return Response{Success: false, Message: "renter KYC verification failed"}, err
	}

	if err := h.writeKYCVerification(cfg, runtime, req.RenterAddress, renterKYC); err != nil {
		return Response{Success: false, Message: "renter KYC write failed"}, err
	}

	accessKey, err := encryption.GenerateRandomKey(32)
	if err != nil {
		return Response{Success: false, Message: "failed to generate access key"}, err
	}

	encryptedAccessKey := accessKey
	accessKeyHash := crypto.Keccak256Hash(encryptedAccessKey)

	calldata, err := houseRWAABI.Pack(
		"startRental",
		new(big.Int).SetUint64(tokenID),
		common.HexToAddress(req.RenterAddress),
		durationBI,
		depositAmount,
		monthlyRent,
		encryptedAccessKey,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode rental calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain rental failed"}, err
	}

	return Response{
		Success:      true,
		Message:      "rental started successfully",
		TxHash:       txHashHex(writeReply.GetTxHash()),
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedAccessKey),
		Data: map[string]interface{}{
			"tokenID":       tokenID,
			"renter":        req.RenterAddress,
			"durationDays":  req.DurationDays,
			"monthlyRent":   monthlyRent.String(),
			"depositAmount": depositAmount.String(),
			"kycProvider":   renterKYC.Provider,
			"accessKeyHash": accessKeyHash.Hex(),
		},
	}, nil
}

// handleBillPayment processes bill payment requests and writes HouseRWA.recordBillPayment via CRE report.
func (h *Handler) handleBillPayment(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req PaymentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid payment request"}, err
	}

	tokenID, err := parseTokenID(req.TokenID, req.TokenId)
	if err != nil {
		return Response{Success: false, Message: "invalid token ID"}, nil
	}

	if err := h.validator.ValidateTokenID(tokenID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if req.OwnerAddress != "" {
		if err := h.validator.ValidateEthereumAddress(req.OwnerAddress); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}
	if err := h.validator.ValidatePaymentMethod(req.PaymentMethod); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if strings.EqualFold(req.PaymentMethod, "stripe") {
		if err := h.validator.ValidateStripeToken(req.StripeToken); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}

	logger.Info("processing bill payment", "tokenID", tokenID, "billIndex", req.BillIndex, "method", req.PaymentMethod)

	paymentRef := crypto.Keccak256Hash([]byte(fmt.Sprintf(
		"bill:%d:%d:%s",
		tokenID,
		req.BillIndex,
		strings.ToLower(req.PaymentMethod),
	)))

	calldata, err := houseRWAABI.Pack(
		"recordBillPayment",
		new(big.Int).SetUint64(tokenID),
		new(big.Int).SetUint64(req.BillIndex),
		strings.ToLower(strings.TrimSpace(req.PaymentMethod)),
		paymentRef,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode payment calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain payment recording failed"}, err
	}

	return Response{
		Success: true,
		Message: fmt.Sprintf("payment processed via %s", strings.ToLower(req.PaymentMethod)),
		TxHash:  txHashHex(writeReply.GetTxHash()),
		Data: map[string]interface{}{
			"tokenID":          tokenID,
			"billIndex":        req.BillIndex,
			"paymentMethod":    strings.ToLower(strings.TrimSpace(req.PaymentMethod)),
			"paymentReference": paymentRef.Hex(),
		},
	}, nil
}

// handleCreateBill processes bill creation requests and writes HouseRWA.createBill via CRE report.
func (h *Handler) handleCreateBill(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req CreateBillRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid bill request"}, err
	}

	tokenID, err := parseTokenID(req.TokenID, req.TokenId)
	if err != nil {
		return Response{Success: false, Message: "invalid token ID"}, nil
	}

	billType := strings.ToLower(strings.TrimSpace(req.BillType))
	if err := h.validator.ValidateTokenID(tokenID); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateBillType(billType); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	dueDate, err := h.validator.ValidateDueDate(req.DueDate)
	if err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}
	if err := h.validator.ValidateEthereumAddress(req.Provider); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	if req.Amount <= 0 {
		return Response{Success: false, Message: "amount must be greater than 0"}, nil
	}

	amountCents := uint64(math.Round(req.Amount * 100))
	if amountCents == 0 {
		return Response{Success: false, Message: "amount too small"}, nil
	}
	amountBI := new(big.Int).SetUint64(amountCents)
	if _, err := asUint96(amountBI, "amount"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	recurrenceInterval := req.RecurrenceInterval
	if req.IsRecurring && recurrenceInterval == 0 {
		recurrenceInterval = 30
	}

	logger.Info("creating bill", "tokenID", tokenID, "billType", billType, "amountCents", amountCents)

	preCount, preCountErr := h.readUint256Method(cfg, runtime, "getTotalBillsCount", new(big.Int).SetUint64(tokenID))
	if preCountErr != nil {
		logger.Warn("getTotalBillsCount read failed before bill creation (billIndex may be unknown)", "error", preCountErr)
	}

	dueUnix := dueDate.Unix()
	if dueUnix < 0 {
		return Response{Success: false, Message: "due date must be after unix epoch"}, nil
	}
	dueDateBI := new(big.Int).SetUint64(uint64(dueUnix))
	if _, err := asUint48(dueDateBI, "dueDate"); err != nil {
		return Response{Success: false, Message: err.Error()}, nil
	}

	calldata, err := houseRWAABI.Pack(
		"createBill",
		new(big.Int).SetUint64(tokenID),
		billType,
		amountBI,
		dueDateBI,
		common.HexToAddress(req.Provider),
		req.IsRecurring,
		recurrenceInterval,
	)
	if err != nil {
		return Response{Success: false, Message: "failed to encode createBill calldata"}, err
	}

	writeReply, err := h.writeCalldata(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "onchain bill creation failed"}, err
	}

	billIndex := -1
	if preCount != nil {
		billIndex = int(preCount.Uint64())
	}

	return Response{
		Success: true,
		Message: "bill created successfully",
		TxHash:  txHashHex(writeReply.GetTxHash()),
		Data: map[string]interface{}{
			"tokenID":            tokenID,
			"billType":           billType,
			"amount":             amountCents,
			"dueDate":            dueDate.Format(time.RFC3339),
			"isRecurring":        req.IsRecurring,
			"recurrenceInterval": recurrenceInterval,
			"billIndex":          billIndex,
		},
	}, nil
}

// handleClaimKey fetches the encrypted key data from the onchain keyExchanges mapping.
func (h *Handler) handleClaimKey(cfg *config.Config, runtime cre.Runtime, body []byte) (Response, error) {
	logger := runtime.Logger()

	var req ClaimKeyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return Response{Success: false, Message: "invalid claim request"}, err
	}

	if strings.TrimSpace(req.KeyHash) == "" {
		return Response{Success: false, Message: "key hash is required"}, nil
	}
	if req.Claimant != "" {
		if err := h.validator.ValidateEthereumAddress(req.Claimant); err != nil {
			return Response{Success: false, Message: err.Error()}, nil
		}
	}

	keyHash, err := parseBytes32Hex(req.KeyHash)
	if err != nil {
		return Response{Success: false, Message: "invalid key hash format"}, nil
	}

	calldata, err := houseRWAABI.Pack("keyExchanges", keyHash)
	if err != nil {
		return Response{Success: false, Message: "failed to encode keyExchanges calldata"}, err
	}

	raw, err := h.callContract(cfg, runtime, calldata)
	if err != nil {
		return Response{Success: false, Message: "failed to read key exchange"}, err
	}
	if len(raw) == 0 {
		return Response{Success: false, Message: "key exchange not found"}, nil
	}

	method := houseRWAABI.Methods["keyExchanges"]
	decoded, err := method.Outputs.Unpack(raw)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "empty string") {
			return Response{Success: false, Message: "key exchange not found"}, nil
		}
		return Response{Success: false, Message: "failed to decode key exchange"}, err
	}

	var exchange keyExchangeView
	if err := method.Outputs.Copy(&exchange, decoded); err != nil {
		return Response{Success: false, Message: "failed to map key exchange"}, err
	}

	if len(exchange.EncryptedKey) == 0 {
		return Response{Success: false, Message: "key exchange not found"}, nil
	}

	if req.Claimant != "" {
		claimant := common.HexToAddress(req.Claimant)
		if claimant != exchange.IntendedRecipient {
			return Response{Success: false, Message: "claimant is not intended recipient"}, nil
		}
	}

	encryptedKeyB64 := base64.StdEncoding.EncodeToString(exchange.EncryptedKey)
	createdAt := int64(0)
	expiresAt := int64(0)
	if exchange.CreatedAt != nil {
		createdAt = exchange.CreatedAt.Int64()
	}
	if exchange.ExpiresAt != nil {
		expiresAt = exchange.ExpiresAt.Int64()
	}

	logger.Info("key exchange fetched", "keyHash", strings.ToLower(req.KeyHash), "intendedRecipient", exchange.IntendedRecipient.Hex())

	return Response{
		Success:      true,
		Message:      "key fetched successfully",
		EncryptedKey: encryptedKeyB64,
		Data: map[string]interface{}{
			"keyHash":           "0x" + hex.EncodeToString(exchange.KeyHash[:]),
			"encryptedKey":      encryptedKeyB64,
			"intendedRecipient": exchange.IntendedRecipient.Hex(),
			"isClaimed":         exchange.IsClaimed,
			"createdAt":         createdAt,
			"expiresAt":         expiresAt,
			"exchangeType":      exchange.ExchangeType,
		},
	}, nil
}

func (h *Handler) writeKYCVerification(cfg *config.Config, runtime cre.Runtime, user string, record kycVerificationRecord) error {
	if strings.EqualFold(strings.TrimSpace(record.Provider), kycProviderNone) {
		runtime.Logger().Info("skipping onchain kyc write for anonymous mode", "user", strings.ToLower(user))
		return nil
	}

	level := record.Level
	if level == 0 {
		level = 2
	}

	verificationHash := record.VerificationHash
	if verificationHash == (common.Hash{}) {
		verificationHash = crypto.Keccak256Hash([]byte(fmt.Sprintf(
			"kyc:%s:%s:%d",
			strings.ToLower(user),
			strings.ToLower(strings.TrimSpace(record.Provider)),
			level,
		)))
	}

	expiryTime := record.Expiry
	if expiryTime.IsZero() {
		expiryTime = time.Now().Add(180 * 24 * time.Hour)
	}
	expiry := uint64(expiryTime.Unix())
	expiryBI := new(big.Int).SetUint64(expiry)
	if _, err := asUint48(expiryBI, "expiryDate"); err != nil {
		return err
	}

	calldata, err := houseRWAABI.Pack(
		"setKYCVerification",
		common.HexToAddress(user),
		level,
		verificationHash,
		expiryBI,
	)
	if err != nil {
		return err
	}

	_, err = h.writeCalldata(cfg, runtime, calldata)
	return err
}

func (h *Handler) writeCalldata(cfg *config.Config, runtime cre.Runtime, calldata []byte) (*evmcap.WriteReportReply, error) {
	logger := runtime.Logger()

	if err := h.validator.ValidateEthereumAddress(cfg.HouseRWAReceiverAddr); err != nil {
		return nil, fmt.Errorf("invalid houseRWAReceiverAddr: %w", err)
	}

	selector, err := cfg.ResolveChainSelector()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve chain selector: %w", err)
	}

	receiverAddr := common.HexToAddress(cfg.HouseRWAReceiverAddr)
	logger.Info("submitting CRE EVM report",
		"chainSelector", selector,
		"evmChain", cfg.EVMChain,
		"receiver", receiverAddr.Hex(),
		"gasLimit", cfg.GasLimit,
		"calldataBytes", len(calldata),
	)

	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: calldata,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to generate CRE report: %w", err)
	}

	evClient := &evmcap.Client{ChainSelector: selector}
	reply, err := evClient.WriteReport(runtime, &evmcap.WriteCreReportRequest{
		Receiver: receiverAddr.Bytes(),
		Report:   report,
		GasConfig: &evmcap.GasConfig{
			GasLimit: cfg.GasLimit,
		},
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to submit CRE report: %w", err)
	}

	txHash := txHashHex(reply.GetTxHash())
	txStatus := reply.GetTxStatus()
	receiverStatus := reply.GetReceiverContractExecutionStatus()
	errMsg := strings.TrimSpace(reply.GetErrorMessage())

	logger.Info("CRE EVM write reply",
		"chainSelector", selector,
		"receiver", receiverAddr.Hex(),
		"txStatus", txStatus.String(),
		"receiverStatus", receiverStatus.String(),
		"txHash", txHash,
		"errorMessage", errMsg,
	)

	if txStatus != evmcap.TxStatus_TX_STATUS_SUCCESS {
		if errMsg == "" {
			errMsg = txStatus.String()
		}
		return nil, fmt.Errorf(
			"write report failed (selector=%d receiver=%s txStatus=%s receiverStatus=%s txHash=%s): %s",
			selector,
			receiverAddr.Hex(),
			txStatus.String(),
			receiverStatus.String(),
			txHash,
			errMsg,
		)
	}

	if receiverStatus != evmcap.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS {
		if errMsg == "" {
			errMsg = receiverStatus.String()
		}
		return nil, fmt.Errorf(
			"receiver execution failed (selector=%d receiver=%s txStatus=%s receiverStatus=%s txHash=%s): %s",
			selector,
			receiverAddr.Hex(),
			txStatus.String(),
			receiverStatus.String(),
			txHash,
			errMsg,
		)
	}

	if txHash == "" {
		logger.Warn("CRE write reply returned empty tx hash (common during local simulation)",
			"chainSelector", selector,
			"receiver", receiverAddr.Hex(),
		)
	}

	return reply, nil
}

func (h *Handler) callContract(cfg *config.Config, runtime cre.Runtime, calldata []byte) ([]byte, error) {
	logger := runtime.Logger()

	if err := h.validator.ValidateEthereumAddress(cfg.HouseRWAContractAddr); err != nil {
		return nil, fmt.Errorf("invalid houseRWAContractAddr: %w", err)
	}

	selector, err := cfg.ResolveChainSelector()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve chain selector: %w", err)
	}

	contractAddr := common.HexToAddress(cfg.HouseRWAContractAddr)
	logger.Info("executing CRE EVM call",
		"chainSelector", selector,
		"contract", contractAddr.Hex(),
		"calldataBytes", len(calldata),
	)

	evClient := &evmcap.Client{ChainSelector: selector}
	reply, err := evClient.CallContract(runtime, &evmcap.CallContractRequest{
		Call: &evmcap.CallMsg{
			To:   contractAddr.Bytes(),
			Data: calldata,
		},
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	logger.Info("CRE EVM call completed",
		"chainSelector", selector,
		"contract", contractAddr.Hex(),
		"resultBytes", len(reply.GetData()),
	)

	return reply.GetData(), nil
}

func (h *Handler) readUint256Method(cfg *config.Config, runtime cre.Runtime, method string, args ...interface{}) (*big.Int, error) {
	calldata, err := houseRWAABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	raw, err := h.callContract(cfg, runtime, calldata)
	if err != nil {
		return nil, err
	}

	decoded, err := houseRWAABI.Methods[method].Outputs.Unpack(raw)
	if err != nil {
		return nil, err
	}
	if len(decoded) == 0 {
		return nil, fmt.Errorf("%s returned no values", method)
	}

	value, ok := decoded[0].(*big.Int)
	if ok {
		return value, nil
	}

	if v, ok := decoded[0].(big.Int); ok {
		return &v, nil
	}

	return nil, fmt.Errorf("%s returned unsupported type %T", method, decoded[0])
}

func parseTokenID(tokenID uint64, tokenIDStr string) (uint64, error) {
	if strings.TrimSpace(tokenIDStr) == "" {
		return tokenID, nil
	}

	parsed, err := strconv.ParseUint(strings.TrimSpace(tokenIDStr), 10, 64)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func parsePositiveBigInt(raw string, field string) (*big.Int, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, fmt.Errorf("%s is required", field)
	}

	bi, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid %s amount", field)
	}
	if bi.Sign() <= 0 {
		return nil, fmt.Errorf("%s must be greater than 0", field)
	}

	return bi, nil
}

func parseListingType(raw string) (uint8, string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "for_sale":
		return 1, "for_sale", nil
	case "for_rent":
		return 2, "for_rent", nil
	default:
		return 0, "", fmt.Errorf("listingType must be either `for_sale` or `for_rent`")
	}
}

func asUint96(v *big.Int, field string) (*big.Int, error) {
	if v.Sign() < 0 {
		return nil, fmt.Errorf("%s cannot be negative", field)
	}
	if v.BitLen() > 96 {
		return nil, fmt.Errorf("%s exceeds uint96 range", field)
	}
	return v, nil
}

func asUint48(v *big.Int, field string) (*big.Int, error) {
	if v.Sign() < 0 {
		return nil, fmt.Errorf("%s cannot be negative", field)
	}
	if v.BitLen() > 48 {
		return nil, fmt.Errorf("%s exceeds uint48 range", field)
	}
	return v, nil
}

func parseBytes32Hex(input string) ([32]byte, error) {
	var out [32]byte

	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")

	b, err := hex.DecodeString(s)
	if err != nil {
		return out, err
	}
	if len(b) != 32 {
		return out, fmt.Errorf("expected 32 bytes, got %d", len(b))
	}

	copy(out[:], b)
	return out, nil
}

func mapStorageType(input string) uint8 {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "offchain", "off_chain", "off-chain":
		return 1 // OFF_CHAIN
	case "arweave":
		return 2 // ARWEAVE
	case "encrypted_db", "encrypteddb", "db":
		return 3 // ENCRYPTED_DB
	default:
		return 0 // IPFS
	}
}

func buildMintMetadataCommitment(
	ownerAddress string,
	houseID string,
	metadata json.RawMessage,
	docDigest [32]byte,
	documentURI string,
	storageEnum uint8,
	privacySalt [32]byte,
) string {
	hasher := sha256.New()
	hasher.Write([]byte(strings.ToLower(strings.TrimSpace(ownerAddress))))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(strings.TrimSpace(houseID)))
	hasher.Write([]byte("|"))
	hasher.Write(docDigest[:])
	hasher.Write([]byte("|"))
	hasher.Write([]byte(documentURI))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(strconv.FormatUint(uint64(storageEnum), 10)))
	hasher.Write([]byte("|"))
	if len(metadata) > 0 {
		hasher.Write(metadata)
	}
	hasher.Write([]byte("|"))
	hasher.Write(privacySalt[:])
	return "0x" + hex.EncodeToString(hasher.Sum(nil))
}

func buildMintVerificationData(
	req MintRequest,
	docDigest [32]byte,
	shareCommitment string,
	metadataCommitment string,
) string {
	payload := map[string]interface{}{
		"schemaVersion":      2,
		"documentHash":       "0x" + hex.EncodeToString(docDigest[:]),
		"shareCommitment":    shareCommitment,
		"metadataCommitment": metadataCommitment,
		"storageType":        strings.ToLower(strings.TrimSpace(req.StorageType)),
		"visibility":         "private_onchain_commitment",
	}

	if req.Location != "" {
		payload["location"] = req.Location
	}
	if req.Value != "" {
		payload["declaredValue"] = req.Value
	}
	if len(req.Metadata) > 0 {
		payload["metadata"] = json.RawMessage(req.Metadata)
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return `{"schemaVersion":2}`
	}

	return string(encoded)
}

func txHashHex(txHash []byte) string {
	if len(txHash) == 0 {
		return ""
	}
	if len(txHash) == common.HashLength {
		return common.BytesToHash(txHash).Hex()
	}
	return "0x" + hex.EncodeToString(txHash)
}

func mustParseABI(raw string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		panic(fmt.Sprintf("failed to parse HouseRWA ABI: %v", err))
	}
	return parsed
}

// createErrorResponse creates an error response.
func createErrorResponse(message string) ([]byte, error) {
	resp := Response{
		Success: false,
		Message: message,
	}
	return json.Marshal(resp)
}

// CreateSuccessResponse creates a success response.
func CreateSuccessResponse(message string, data interface{}) []byte {
	resp := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	bytes, _ := json.Marshal(resp)
	return bytes
}

// Logger wraps slog.Logger for CRE compatibility.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new logger.
func NewLogger() *Logger {
	return &Logger{
		Logger: slog.Default(),
	}
}
