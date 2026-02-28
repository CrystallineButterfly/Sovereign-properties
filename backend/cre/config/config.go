//go:build wasip1

// Package config provides configuration management for the Chainlink CRE workflow.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	evmcap "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
)

// Config holds the workflow configuration
type Config struct {
	// Chain configuration
	EVMChain             string `json:"evmChain,omitempty"`
	EVMChainID           int64  `json:"evmChainID"`
	HouseRWAContractAddr string `json:"houseRWAContractAddr"`
	HouseRWAReceiverAddr string `json:"houseRWAReceiverAddr"`
	RPCURL               string `json:"rpcURL"`

	// Gas settings
	GasLimit    uint64 `json:"gasLimit"`
	GasPrice    string `json:"gasPrice,omitempty"`
	MaxGasPrice string `json:"maxGasPrice,omitempty"`

	// External services
	StripeAPIBaseURL string `json:"stripeAPIBaseURL"`
	StripeAPIKey     string `json:"-"` // Loaded from secrets
	IPFSGateway      string `json:"ipfsGateway"`
	IPFSAPIKey       string `json:"-"` // Loaded from secrets

	// Trusted entities
	TrustedProviders []string `json:"trustedProviders"`
	ValidatorAddress string   `json:"validatorAddress"`

	// Business logic settings
	AutoPayThreshold   float64 `json:"autoPayThreshold"`
	KeyExpiryDuration  string  `json:"keyExpiryDuration"`
	MaxDocumentSize    int64   `json:"maxDocumentSize"`
	RateLimitPerMinute int     `json:"rateLimitPerMinute"`

	// Threshold encryption settings
	ThresholdKeyThreshold int `json:"thresholdKeyThreshold"`
	ThresholdKeyTotal     int `json:"thresholdKeyTotal"`

	// Security settings
	EnableTEE           bool   `json:"enableTEE"`
	EnforceKYC          bool   `json:"enforceKYC"`
	KYCProvider         string `json:"kycProvider"`
	KYCVerifierURL      string `json:"kycVerifierURL,omitempty"`
	KYCProviderKey      string `json:"-"` // Loaded from secrets
	RequireDONConsensus bool   `json:"requireDONConsensus"`

	// Logging
	LogLevel string `json:"logLevel"`

	// Parsed duration
	KeyExpiry time.Duration `json:"-"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		EVMChain:              "ethereum-testnet-sepolia",
		EVMChainID:            11155111, // Sepolia
		GasLimit:              500000,
		StripeAPIBaseURL:      "https://api.stripe.com/v1",
		IPFSGateway:           "https://ipfs.io",
		TrustedProviders:      make([]string, 0),
		AutoPayThreshold:      1000.0,
		KeyExpiryDuration:     "24h",
		MaxDocumentSize:       10 * 1024 * 1024, // 10MB
		RateLimitPerMinute:    100,
		ThresholdKeyThreshold: 3,
		ThresholdKeyTotal:     5,
		EnableTEE:             true,
		EnforceKYC:            true,
		KYCProvider:           "mock",
		RequireDONConsensus:   true,
		LogLevel:              "info",
		KeyExpiry:             24 * time.Hour,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Parse duration
	if config.KeyExpiryDuration != "" {
		duration, err := time.ParseDuration(config.KeyExpiryDuration)
		if err != nil {
			return nil, fmt.Errorf("invalid key expiry duration: %w", err)
		}
		config.KeyExpiry = duration
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// Chain settings
	if chainName := os.Getenv("CRE_CHAIN_NAME"); chainName != "" {
		config.EVMChain = chainName
	}

	if chainID := os.Getenv("CRE_CHAIN_ID"); chainID != "" {
		if id, err := strconv.ParseInt(chainID, 10, 64); err == nil {
			config.EVMChainID = id
		}
	}

	if contractAddr := os.Getenv("CRE_CONTRACT_ADDRESS"); contractAddr != "" {
		config.HouseRWAContractAddr = contractAddr
	}

	if receiverAddr := os.Getenv("CRE_RECEIVER_ADDRESS"); receiverAddr != "" {
		config.HouseRWAReceiverAddr = receiverAddr
	}

	if rpcURL := os.Getenv("CRE_RPC_URL"); rpcURL != "" {
		config.RPCURL = rpcURL
	}

	// Gas settings
	if gasLimit := os.Getenv("CRE_GAS_LIMIT"); gasLimit != "" {
		if limit, err := strconv.ParseUint(gasLimit, 10, 64); err == nil {
			config.GasLimit = limit
		}
	}

	// External services
	if stripeURL := os.Getenv("CRE_STRIPE_API_URL"); stripeURL != "" {
		config.StripeAPIBaseURL = stripeURL
	}

	if ipfsGateway := os.Getenv("CRE_IPFS_GATEWAY"); ipfsGateway != "" {
		config.IPFSGateway = ipfsGateway
	}

	// Security settings
	if enforceKYC := os.Getenv("CRE_ENFORCE_KYC"); enforceKYC != "" {
		config.EnforceKYC = strings.ToLower(enforceKYC) == "true"
	}

	if kycProvider := os.Getenv("CRE_KYC_PROVIDER"); kycProvider != "" {
		config.KYCProvider = strings.ToLower(strings.TrimSpace(kycProvider))
	}

	if verifierURL := os.Getenv("CRE_KYC_VERIFIER_URL"); verifierURL != "" {
		config.KYCVerifierURL = strings.TrimSpace(verifierURL)
	}

	if enableTEE := os.Getenv("CRE_ENABLE_TEE"); enableTEE != "" {
		config.EnableTEE = strings.ToLower(enableTEE) == "true"
	}

	// Business logic
	if threshold := os.Getenv("CRE_AUTO_PAY_THRESHOLD"); threshold != "" {
		if t, err := strconv.ParseFloat(threshold, 64); err == nil {
			config.AutoPayThreshold = t
		}
	}

	if docSize := os.Getenv("CRE_MAX_DOC_SIZE"); docSize != "" {
		if size, err := strconv.ParseInt(docSize, 10, 64); err == nil {
			config.MaxDocumentSize = size
		}
	}

	// Parse duration
	if expiry := os.Getenv("CRE_KEY_EXPIRY"); expiry != "" {
		if duration, err := time.ParseDuration(expiry); err == nil {
			config.KeyExpiry = duration
			config.KeyExpiryDuration = expiry
		}
	}

	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.EVMChainID == 0 {
		return errors.New("evmChainID is required")
	}

	if c.HouseRWAContractAddr == "" {
		return errors.New("houseRWAContractAddr is required")
	}

	if c.HouseRWAReceiverAddr == "" {
		return errors.New("houseRWAReceiverAddr is required")
	}

	if c.EVMChain == "" {
		if chain, ok := chainNameFromID(c.EVMChainID); ok {
			c.EVMChain = chain
		}
	}

	if c.EVMChain != "" {
		if _, err := evmcap.ChainSelectorFromName(c.EVMChain); err != nil {
			return fmt.Errorf("invalid evmChain %q: %w", c.EVMChain, err)
		}
	}

	if c.GasLimit == 0 {
		c.GasLimit = 500000
	}

	if c.AutoPayThreshold < 0 {
		return errors.New("autoPayThreshold must be non-negative")
	}

	if c.ThresholdKeyThreshold < 2 {
		return errors.New("thresholdKeyThreshold must be at least 2")
	}

	if c.ThresholdKeyTotal < c.ThresholdKeyThreshold {
		return errors.New("thresholdKeyTotal must be >= thresholdKeyThreshold")
	}

	if c.MaxDocumentSize <= 0 {
		c.MaxDocumentSize = 10 * 1024 * 1024
	}

	if c.RateLimitPerMinute <= 0 {
		c.RateLimitPerMinute = 100
	}

	provider := strings.ToLower(strings.TrimSpace(c.KYCProvider))
	if provider == "" {
		provider = "mock"
	}

	switch provider {
	case "mock", "none", "zkpassport":
		c.KYCProvider = provider
	default:
		return fmt.Errorf("unsupported kycProvider %q (expected mock, none, or zkpassport)", c.KYCProvider)
	}

	if c.KYCProvider == "zkpassport" && strings.TrimSpace(c.KYCVerifierURL) == "" {
		return errors.New("kycVerifierURL is required when kycProvider is zkpassport")
	}

	return nil
}

// IsProduction returns true if this is a production configuration
func (c *Config) IsProduction() bool {
	return c.EVMChainID == 1 // Mainnet
}

// IsStaging returns true if this is a staging configuration
func (c *Config) IsStaging() bool {
	return c.EVMChain == "ethereum-testnet-sepolia" || c.EVMChainID == 11155111 // Sepolia
}

// GetRPCURL returns the RPC URL, with fallback
func (c *Config) GetRPCURL() string {
	if c.RPCURL != "" {
		return c.RPCURL
	}

	// Fallback to default RPCs based on chain ID
	switch c.EVMChainID {
	case 1:
		return "https://ethereum-rpc.publicnode.com"
	case 11155111:
		return "https://sepolia.drpc.org"
	case 137:
		return "https://polygon-rpc.com"
	case 42161:
		return "https://arb1.arbitrum.io/rpc"
	default:
		return ""
	}
}

// ResolveChainSelector returns the configured CRE chain selector.
func (c *Config) ResolveChainSelector() (uint64, error) {
	if c.EVMChain != "" {
		return evmcap.ChainSelectorFromName(c.EVMChain)
	}

	if chainName, ok := chainNameFromID(c.EVMChainID); ok {
		return evmcap.ChainSelectorFromName(chainName)
	}

	return 0, fmt.Errorf("unsupported chain configuration (evmChain=%q evmChainID=%d)", c.EVMChain, c.EVMChainID)
}

func chainNameFromID(chainID int64) (string, bool) {
	switch chainID {
	case 1:
		return "ethereum-mainnet", true
	case 11155111:
		return "ethereum-testnet-sepolia", true
	case 137:
		return "polygon-mainnet", true
	case 80002:
		return "polygon-testnet-amoy", true
	default:
		return "", false
	}
}

// LoadSecrets loads secrets from environment or secrets provider
func (c *Config) LoadSecrets(getSecret func(key string) (string, error)) error {
	var err error

	// Try to load from secrets provider first
	if getSecret != nil {
		c.StripeAPIKey, err = getSecret("STRIPE_API_KEY")
		if err != nil && c.StripeAPIKey == "" {
			// Try environment variable
			c.StripeAPIKey = os.Getenv("STRIPE_API_KEY")
		}

		c.IPFSAPIKey, err = getSecret("IPFS_API_KEY")
		if err != nil && c.IPFSAPIKey == "" {
			c.IPFSAPIKey = os.Getenv("IPFS_API_KEY")
		}

		c.KYCProviderKey, err = getSecret("KYC_PROVIDER_KEY")
		if err != nil && c.KYCProviderKey == "" {
			c.KYCProviderKey = os.Getenv("KYC_PROVIDER_KEY")
		}
	} else {
		// Load from environment
		c.StripeAPIKey = os.Getenv("STRIPE_API_KEY")
		c.IPFSAPIKey = os.Getenv("IPFS_API_KEY")
		c.KYCProviderKey = os.Getenv("KYC_PROVIDER_KEY")
	}

	return nil
}

// String returns a string representation of the config (excluding secrets)
func (c *Config) String() string {
	data, _ := json.MarshalIndent(struct {
		EVMChain              string  `json:"evmChain,omitempty"`
		EVMChainID            int64   `json:"evmChainID"`
		HouseRWAContractAddr  string  `json:"houseRWAContractAddr"`
		HouseRWAReceiverAddr  string  `json:"houseRWAReceiverAddr"`
		GasLimit              uint64  `json:"gasLimit"`
		IPFSGateway           string  `json:"ipfsGateway"`
		AutoPayThreshold      float64 `json:"autoPayThreshold"`
		MaxDocumentSize       int64   `json:"maxDocumentSize"`
		ThresholdKeyThreshold int     `json:"thresholdKeyThreshold"`
		ThresholdKeyTotal     int     `json:"thresholdKeyTotal"`
		EnableTEE             bool    `json:"enableTEE"`
		EnforceKYC            bool    `json:"enforceKYC"`
		KYCProvider           string  `json:"kycProvider"`
		KYCVerifierURL        string  `json:"kycVerifierURL,omitempty"`
	}{
		EVMChain:              c.EVMChain,
		EVMChainID:            c.EVMChainID,
		HouseRWAContractAddr:  c.HouseRWAContractAddr,
		HouseRWAReceiverAddr:  c.HouseRWAReceiverAddr,
		GasLimit:              c.GasLimit,
		IPFSGateway:           c.IPFSGateway,
		AutoPayThreshold:      c.AutoPayThreshold,
		MaxDocumentSize:       c.MaxDocumentSize,
		ThresholdKeyThreshold: c.ThresholdKeyThreshold,
		ThresholdKeyTotal:     c.ThresholdKeyTotal,
		EnableTEE:             c.EnableTEE,
		EnforceKYC:            c.EnforceKYC,
		KYCProvider:           c.KYCProvider,
		KYCVerifierURL:        c.KYCVerifierURL,
	}, "", "  ")

	return string(data)
}

// SecretsConfig holds secret configuration
type SecretsConfig struct {
	StripeAPIKey   string `json:"stripe_api_key"`
	IPFSAPIKey     string `json:"ipfs_api_key"`
	EncryptionKey  string `json:"encryption_key"`
	ValidatorKey   string `json:"validator_key"`
	KYCProviderKey string `json:"kyc_provider_key"`
}

// LoadSecretsConfig loads secrets from a JSON file
func LoadSecretsConfig(path string) (*SecretsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %w", err)
	}

	var secrets SecretsConfig
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets: %w", err)
	}

	return &secrets, nil
}
