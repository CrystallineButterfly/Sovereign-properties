//go:build wasip1

// Package validation provides input validation and sanitization
// for the Chainlink CRE workflow handlers.
package validation

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ethereum/go-ethereum/common"
)

// Validator provides validation functions
type Validator struct {
	maxDocumentSize int64
	allowedPatterns map[string]*regexp.Regexp
}

// NewValidator creates a new validator with default settings
func NewValidator() *Validator {
	return &Validator{
		maxDocumentSize: 10 * 1024 * 1024, // 10MB
		allowedPatterns: map[string]*regexp.Regexp{
			"houseID":  regexp.MustCompile(`^[a-zA-Z0-9-_]{1,64}$`),
			"location": regexp.MustCompile(`^[a-zA-Z0-9\s,.-]{1,256}$`),
			"billType": regexp.MustCompile(`^(utilities|tax|maintenance|insurance|hoa|other|electricity|water|gas|internet|phone|property_tax)$`),
		},
	}
}

// NewValidatorWithOptions creates a validator with custom options
func NewValidatorWithOptions(maxDocSize int64) *Validator {
	v := NewValidator()
	v.maxDocumentSize = maxDocSize
	return v
}

// SetMaxDocumentSize sets the maximum document size
func (v *Validator) SetMaxDocumentSize(size int64) {
	v.maxDocumentSize = size
}

// ValidateEthereumAddress validates an Ethereum address
func (v *Validator) ValidateEthereumAddress(address string) error {
	if address == "" {
		return errors.New("address is required")
	}

	if !common.IsHexAddress(address) {
		return errors.New("invalid Ethereum address format")
	}

	// Validate checksum if mixed case
	if strings.ToLower(address) != address && strings.ToUpper(address) != address {
		if common.HexToAddress(address).Hex() != address {
			return errors.New("invalid address checksum")
		}
	}

	// Check for zero address
	if common.HexToAddress(address) == common.HexToAddress("0x0") {
		return errors.New("zero address not allowed")
	}

	return nil
}

// ValidateHouseID validates a house identifier
func (v *Validator) ValidateHouseID(houseID string) error {
	if houseID == "" {
		return errors.New("house ID is required")
	}

	if len(houseID) > 64 {
		return errors.New("house ID too long (max 64 characters)")
	}

	if !v.allowedPatterns["houseID"].MatchString(houseID) {
		return errors.New("invalid house ID format")
	}

	return nil
}

// ValidateDocument validates document data
func (v *Validator) ValidateDocument(docB64 string) ([]byte, error) {
	if docB64 == "" {
		return nil, errors.New("document is required")
	}

	// Decode base64
	doc, err := base64.StdEncoding.DecodeString(docB64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encoding: %w", err)
	}

	// Check size
	if int64(len(doc)) > v.maxDocumentSize {
		return nil, fmt.Errorf("document too large: %d bytes (max %d)", len(doc), v.maxDocumentSize)
	}

	// Check for empty content
	if len(doc) == 0 {
		return nil, errors.New("empty document")
	}

	return doc, nil
}

// ValidatePublicKey validates a public key
func (v *Validator) ValidatePublicKey(pubKey string) error {
	if pubKey == "" {
		return errors.New("public key is required")
	}

	// Check PEM format
	if strings.Contains(pubKey, "BEGIN PUBLIC KEY") {
		return nil // Valid PEM
	}

	// Check hex format
	if strings.HasPrefix(pubKey, "0x") {
		pubKey = pubKey[2:]
	}

	decoded, err := hex.DecodeString(pubKey)
	if err != nil {
		return errors.New("invalid public key format")
	}

	// Uncompressed secp256k1 public key is 65 bytes
	// Compressed is 33 bytes
	if len(decoded) != 65 && len(decoded) != 33 {
		return errors.New("invalid public key length")
	}

	return nil
}

// ValidateTokenID validates a token ID
func (v *Validator) ValidateTokenID(tokenID uint64) error {
	// HouseRWA uses zero-based token IDs, so tokenID=0 is valid.
	return nil
}

// ValidateAmount validates a monetary amount
func (v *Validator) ValidateAmount(amount string) error {
	if amount == "" {
		return errors.New("amount is required")
	}

	// Remove any decimal point for validation
	amount = strings.ReplaceAll(amount, ".", "")

	// Check if valid positive integer
	val, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return errors.New("invalid amount format")
	}

	if val == 0 {
		return errors.New("amount must be greater than 0")
	}

	return nil
}

// ValidatePrice validates a price string
func (v *Validator) ValidatePrice(price string) error {
	if price == "" {
		return errors.New("price is required")
	}

	// Allow decimal prices
	val, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return errors.New("invalid price format")
	}

	if val <= 0 {
		return errors.New("price must be greater than 0")
	}

	return nil
}

// ValidateDuration validates a duration in days
func (v *Validator) ValidateDuration(days uint64) error {
	if days == 0 {
		return errors.New("duration must be greater than 0")
	}
	if days > 3650 { // Max 10 years
		return errors.New("duration too long (max 3650 days)")
	}
	return nil
}

// ValidateBillType validates a bill type
func (v *Validator) ValidateBillType(billType string) error {
	if billType == "" {
		return errors.New("bill type is required")
	}

	billType = strings.ToLower(billType)
	if !v.allowedPatterns["billType"].MatchString(billType) {
		return errors.New("invalid bill type")
	}

	return nil
}

// ValidateDueDate validates a due date
func (v *Validator) ValidateDueDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("due date is required")
	}

	// Try multiple formats
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// Check if date is in the future
			if t.Before(time.Now()) {
				return time.Time{}, errors.New("due date must be in the future")
			}
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}

// ValidateStorageType validates a storage type
func (v *Validator) ValidateStorageType(storageType string) error {
	if storageType == "" {
		return errors.New("storage type is required")
	}

	storageType = strings.ToLower(storageType)
	if storageType != "ipfs" && storageType != "offchain" {
		return errors.New("invalid storage type (must be 'ipfs' or 'offchain')")
	}

	return nil
}

// ValidatePaymentMethod validates a payment method
func (v *Validator) ValidatePaymentMethod(method string) error {
	if method == "" {
		return errors.New("payment method is required")
	}

	method = strings.ToLower(method)
	if method != "crypto" && method != "stripe" && method != "bank" {
		return errors.New("invalid payment method (must be 'crypto', 'stripe', or 'bank')")
	}

	return nil
}

// ValidateStripeToken validates a Stripe token
func (v *Validator) ValidateStripeToken(token string) error {
	if token == "" {
		return errors.New("stripe token is required")
	}

	// Stripe tokens typically start with 'tok_' or 'pi_'
	if !strings.HasPrefix(token, "tok_") && !strings.HasPrefix(token, "pi_") {
		return errors.New("invalid stripe token format")
	}

	return nil
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	// Additional validation
	if !strings.Contains(addr.Address, "@") {
		return errors.New("invalid email format")
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateLocation validates a location string
func (v *Validator) ValidateLocation(location string) error {
	if location == "" {
		return errors.New("location is required")
	}

	if len(location) > 256 {
		return errors.New("location too long (max 256 characters)")
	}

	if !v.allowedPatterns["location"].MatchString(location) {
		return errors.New("invalid location format")
	}

	return nil
}

// SanitizeString sanitizes a string input
func (v *Validator) SanitizeString(input string) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Limit length
	if len(sanitized) > 4096 {
		sanitized = sanitized[:4096]
	}

	return sanitized
}

// SanitizeJSON sanitizes a JSON string
func (v *Validator) SanitizeJSON(input string) (string, error) {
	// Remove potentially dangerous characters
	sanitized := strings.Map(func(r rune) rune {
		if r == 0x00 || r == 0x1F || r == 0x7F {
			return -1
		}
		return r
	}, input)

	// Check for valid JSON structure
	if !isValidJSON(sanitized) {
		return "", errors.New("invalid JSON structure")
	}

	return sanitized, nil
}

// isValidJSON checks if a string is valid JSON
func isValidJSON(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}

	firstChar := s[0]
	lastChar := s[len(s)-1]

	// Object or Array
	return (firstChar == '{' && lastChar == '}') ||
		(firstChar == '[' && lastChar == ']')
}

// ContainsSQLInjection checks for SQL injection patterns
func (v *Validator) ContainsSQLInjection(input string) bool {
	sqlPatterns := []string{
		`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|EXEC|EXECUTE)`,
		`(?i)(UNION|UNION ALL|JOIN|INNER JOIN|OUTER JOIN)`,
		`(?i)(--|;|/\*|\*/)`,
		`(?i)('|")\s*(OR|AND)\s*('|")`,
	}

	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}

	return false
}

// ContainsXSS checks for XSS patterns
func (v *Validator) ContainsXSS(input string) bool {
	xssPatterns := []string{
		`(?i)<script[^>]*>`,
		`(?i)javascript:`,
		`(?i)on\w+\s*=`,
		`(?i)<iframe`,
		`(?i)<object`,
		`(?i)<embed`,
	}

	for _, pattern := range xssPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}

	return false
}

// IsValidID checks if a string is a valid identifier
func (v *Validator) IsValidID(id string) bool {
	if id == "" {
		return false
	}

	// Check length
	if len(id) > 128 {
		return false
	}

	// Check characters
	for _, r := range id {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}

// IsSafeString checks if a string is safe (no control characters)
func (v *Validator) IsSafeString(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return false
		}
	}
	return true
}

// ValidationResult holds validation results
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// AddError adds an error to the result
func (r *ValidationResult) AddError(err string) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// AddWarning adds a warning to the result
func (r *ValidationResult) AddWarning(warning string) {
	r.Warnings = append(r.Warnings, warning)
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:    true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}
}
