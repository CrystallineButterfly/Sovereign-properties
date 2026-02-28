//go:build wasip1

// Package evm provides Ethereum Virtual Machine interaction capabilities
// for the Chainlink CRE workflow.
package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client represents an EVM client for interacting with smart contracts
type Client struct {
	client      *ethclient.Client
	chainID     *big.Int
	gasLimit    uint64
	gasPrice    *big.Int
	maxGasPrice *big.Int
}

// ClientConfig holds EVM client configuration
type ClientConfig struct {
	RPCURL      string
	ChainID     int64
	GasLimit    uint64
	GasPrice    string // in wei
	MaxGasPrice string // in wei
	Timeout     time.Duration
}

// NewClient creates a new EVM client
func NewClient(config ClientConfig) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Verify connection
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	if chainID.Int64() != config.ChainID {
		return nil, fmt.Errorf("chain ID mismatch: expected %d, got %d", config.ChainID, chainID.Int64())
	}

	gasPrice := big.NewInt(0)
	if config.GasPrice != "" {
		gasPrice, _ = new(big.Int).SetString(config.GasPrice, 10)
	}

	maxGasPrice := big.NewInt(0)
	if config.MaxGasPrice != "" {
		maxGasPrice, _ = new(big.Int).SetString(config.MaxGasPrice, 10)
	}

	return &Client{
		client:      client,
		chainID:     chainID,
		gasLimit:    config.GasLimit,
		gasPrice:    gasPrice,
		maxGasPrice: maxGasPrice,
	}, nil
}

// Close closes the EVM client connection
func (c *Client) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// CallContract performs a read-only contract call
func (c *Client) CallContract(ctx context.Context, contract common.Address, data []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := c.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	return result, nil
}

// SendTransaction sends a transaction to the blockchain
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.client.SendTransaction(ctx, tx)
}

// WaitForReceipt waits for a transaction receipt
func (c *Client) WaitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
			receipt, err := c.client.TransactionReceipt(ctx, txHash)
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				return nil, err
			}
			return receipt, nil
		}
	}
}

// GetBalance gets the balance of an address
func (c *Client) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.client.BalanceAt(ctx, address, nil)
}

// GetNonce gets the transaction nonce for an address
func (c *Client) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	return c.client.PendingNonceAt(ctx, address)
}

// EstimateGas estimates gas for a transaction
func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.client.EstimateGas(ctx, msg)
}

// HouseRWAContract represents the HouseRWA smart contract
type HouseRWAContract struct {
	client   *Client
	address  common.Address
	abi      []byte
	signerFn SignerFunc
}

// SignerFunc is a function type for signing transactions
type SignerFunc func(address common.Address, tx *types.Transaction) (*types.Transaction, error)

// NewHouseRWAContract creates a new HouseRWA contract instance
func NewHouseRWAContract(client *Client, address string, abi []byte, signer SignerFunc) (*HouseRWAContract, error) {
	if !common.IsHexAddress(address) {
		return nil, errors.New("invalid contract address")
	}

	return &HouseRWAContract{
		client:   client,
		address:  common.HexToAddress(address),
		abi:      abi,
		signerFn: signer,
	}, nil
}

// MintHouse calls the mintHouse function on the contract
func (h *HouseRWAContract) MintHouse(ctx context.Context, owner common.Address, location string, value *big.Int, documentURI string) (*types.Transaction, error) {
	// Build transaction data
	data, err := buildMintHouseData(owner, location, value, documentURI)
	if err != nil {
		return nil, fmt.Errorf("failed to build mint data: %w", err)
	}

	return h.sendTransaction(ctx, data)
}

// TransferFrom transfers a token
func (h *HouseRWAContract) TransferFrom(ctx context.Context, from, to common.Address, tokenID *big.Int) (*types.Transaction, error) {
	data, err := buildTransferFromData(from, to, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to build transfer data: %w", err)
	}

	return h.sendTransaction(ctx, data)
}

// CreateRental creates a rental agreement
func (h *HouseRWAContract) CreateRental(ctx context.Context, tokenID *big.Int, renter common.Address, duration uint64, rentAmount *big.Int) (*types.Transaction, error) {
	data, err := buildCreateRentalData(tokenID, renter, duration, rentAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to build rental data: %w", err)
	}

	return h.sendTransaction(ctx, data)
}

// PayBill pays a bill for a property
func (h *HouseRWAContract) PayBill(ctx context.Context, tokenID *big.Int, billIndex uint64, amount *big.Int) (*types.Transaction, error) {
	data, err := buildPayBillData(tokenID, billIndex, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to build payment data: %w", err)
	}

	return h.sendTransaction(ctx, data)
}

// GetOwnerOf gets the owner of a token
func (h *HouseRWAContract) GetOwnerOf(ctx context.Context, tokenID *big.Int) (common.Address, error) {
	data, err := buildOwnerOfData(tokenID)
	if err != nil {
		return common.Address{}, err
	}

	result, err := h.client.CallContract(ctx, h.address, data)
	if err != nil {
		return common.Address{}, err
	}

	// Parse address from result
	if len(result) != 32 {
		return common.Address{}, errors.New("invalid response length")
	}

	return common.BytesToAddress(result[12:]), nil
}

// sendTransaction sends a transaction to the contract
func (h *HouseRWAContract) sendTransaction(ctx context.Context, data []byte) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := h.client.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Apply max gas price limit if set
	if h.client.maxGasPrice != nil && h.client.maxGasPrice.Cmp(big.NewInt(0)) > 0 {
		if gasPrice.Cmp(h.client.maxGasPrice) > 0 {
			gasPrice = h.client.maxGasPrice
		}
	}

	// Get nonce
	// In CRE, this would come from the runtime context
	nonce := uint64(0) // Placeholder

	// Create transaction
	tx := types.NewTransaction(
		nonce,
		h.address,
		big.NewInt(0),
		h.client.gasLimit,
		gasPrice,
		data,
	)

	// Sign transaction
	signedTx, err := h.signerFn(h.address, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	if err := h.client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// EventHandler handles contract events
type EventHandler struct {
	client *Client
}

// NewEventHandler creates a new event handler
func NewEventHandler(client *Client) *EventHandler {
	return &EventHandler{client: client}
}

// WatchEvent watches for specific contract events
func (e *EventHandler) WatchEvent(ctx context.Context, contract common.Address, eventSignature string) (<-chan types.Log, error) {
	// Compute event signature hash
	sigHash := crypto.Keccak256Hash([]byte(eventSignature))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
		Topics:    [][]common.Hash{{sigHash}},
	}

	logs := make(chan types.Log)
	sub, err := e.client.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		close(logs)
	}()

	return logs, nil
}

// ParseEventData parses event data from log
func ParseEventData(log types.Log, abi []byte) (map[string]interface{}, error) {
	// In a real implementation, this would use the contract ABI
	// to decode the event data
	return map[string]interface{}{
		"address": log.Address.Hex(),
		"topics":  log.Topics,
		"data":    hex.EncodeToString(log.Data),
	}, nil
}

// Helper functions for building transaction data

func buildMintHouseData(owner common.Address, location string, value *big.Int, documentURI string) ([]byte, error) {
	// Simplified: In production, use proper ABI encoding
	// This would encode: mintHouse(address,string,uint256,string)
	return []byte{}, nil
}

func buildTransferFromData(from, to common.Address, tokenID *big.Int) ([]byte, error) {
	// Encode: transferFrom(address,address,uint256)
	return []byte{}, nil
}

func buildCreateRentalData(tokenID *big.Int, renter common.Address, duration uint64, rentAmount *big.Int) ([]byte, error) {
	// Encode: createRental(uint256,address,uint64,uint256)
	return []byte{}, nil
}

func buildPayBillData(tokenID *big.Int, billIndex uint64, amount *big.Int) ([]byte, error) {
	// Encode: payBill(uint256,uint256,uint256)
	return []byte{}, nil
}

func buildOwnerOfData(tokenID *big.Int) ([]byte, error) {
	// Encode: ownerOf(uint256)
	return []byte{}, nil
}

// ValidateAddress validates an Ethereum address
func ValidateAddress(address string) error {
	if !common.IsHexAddress(address) {
		return errors.New("invalid Ethereum address format")
	}

	// Check checksum if mixed case
	if strings.ToLower(address) != address && strings.ToUpper(address) != address {
		if common.HexToAddress(address).Hex() != address {
			return errors.New("invalid address checksum")
		}
	}

	return nil
}

// FormatAddress formats an address with checksum
func FormatAddress(address string) (string, error) {
	if err := ValidateAddress(address); err != nil {
		return "", err
	}
	return common.HexToAddress(address).Hex(), nil
}

// ParseBigInt parses a string to *big.Int
func ParseBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		return nil, fmt.Errorf("invalid number: %s", s)
	}
	return n, nil
}

// CREIntegration provides integration helpers for Chainlink CRE runtime
type CREIntegration struct {
	LogFunc func(level, msg string, args ...interface{})
}

// NewCREIntegration creates a new CRE integration helper
func NewCREIntegration(logFunc func(level, msg string, args ...interface{})) *CREIntegration {
	return &CREIntegration{LogFunc: logFunc}
}

// LogInfo logs an info message
func (c *CREIntegration) LogInfo(msg string, args ...interface{}) {
	if c.LogFunc != nil {
		c.LogFunc("info", msg, args...)
	}
}

// LogError logs an error message
func (c *CREIntegration) LogError(msg string, args ...interface{}) {
	if c.LogFunc != nil {
		c.LogFunc("error", msg, args...)
	}
}

// JSONResponse creates a standard JSON response
func JSONResponse(success bool, message string, data interface{}) []byte {
	resp := map[string]interface{}{
		"success": success,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}

	bytes, _ := json.Marshal(resp)
	return bytes
}

// Event represents a contract event
type Event struct {
	Name     string
	Contract common.Address
	Data     map[string]interface{}
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// EventParser parses contract events
type EventParser struct {
	abi []byte
}

// NewEventParser creates a new event parser
func NewEventParser(abi []byte) *EventParser {
	return &EventParser{abi: abi}
}

// Parse parses a log into an Event
func (p *EventParser) Parse(log types.Log) (*Event, error) {
	if len(log.Topics) == 0 {
		return nil, errors.New("no topics in log")
	}

	// Get event signature from first topic
	eventSig := log.Topics[0].Hex()

	return &Event{
		Name:     eventSig,
		Contract: log.Address,
		BlockNum: log.BlockNumber,
		TxHash:   log.TxHash,
		LogIndex: log.Index,
		Data: map[string]interface{}{
			"topics": log.Topics,
			"data":   hex.EncodeToString(log.Data),
		},
	}, nil
}
