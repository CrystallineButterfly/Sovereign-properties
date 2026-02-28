//go:build wasip1

// Package encryption provides threshold encryption for distributed key management
// using Shamir's Secret Sharing scheme.
package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"

	"golang.org/x/crypto/sha3"
)

// ThresholdScheme implements Shamir's Secret Sharing
type ThresholdScheme struct {
	Threshold int
	Total     int
}

// NewThresholdScheme creates a new threshold encryption scheme
func NewThresholdScheme(threshold, total int) (*ThresholdScheme, error) {
	if threshold < 2 {
		return nil, errors.New("threshold must be at least 2")
	}
	if total < threshold {
		return nil, errors.New("total must be >= threshold")
	}
	if total > 255 {
		return nil, errors.New("total must be <= 255")
	}

	return &ThresholdScheme{
		Threshold: threshold,
		Total:     total,
	}, nil
}

// GenerateThresholdShares splits a key into shares using Shamir's Secret Sharing
func (ts *ThresholdScheme) GenerateThresholdShares(secret []byte) ([]Share, error) {
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}

	// Convert secret to big.Int
	secretInt := new(big.Int).SetBytes(secret)

	// Generate polynomial coefficients
	// f(x) = a_0 + a_1*x + a_2*x^2 + ... + a_{t-1}*x^{t-1}
	// where a_0 = secret
	coefficients := make([]*big.Int, ts.Threshold)
	coefficients[0] = secretInt

	for i := 1; i < ts.Threshold; i++ {
		coeff, err := rand.Int(rand.Reader, prime)
		if err != nil {
			return nil, fmt.Errorf("failed to generate coefficient: %w", err)
		}
		coefficients[i] = coeff
	}

	// Generate shares
	shares := make([]Share, ts.Total)
	for i := 1; i <= ts.Total; i++ {
		x := big.NewInt(int64(i))
		y := evaluatePolynomial(coefficients, x)

		shareValue := y.Bytes()
		checksum := calculateShareChecksum(shareValue, i)

		shares[i-1] = Share{
			Index:    i,
			Value:    shareValue,
			Checksum: checksum,
		}
	}

	return shares, nil
}

// ReconstructSecret reconstructs the original secret from shares using Lagrange interpolation
func (ts *ThresholdScheme) ReconstructSecret(shares []Share) ([]byte, error) {
	if len(shares) < ts.Threshold {
		return nil, fmt.Errorf("insufficient shares: got %d, need %d", len(shares), ts.Threshold)
	}

	// Verify share integrity
	for _, share := range shares {
		if !verifyShareChecksum(share) {
			return nil, fmt.Errorf("share %d failed integrity check", share.Index)
		}
	}

	// Use only threshold number of shares
	if len(shares) > ts.Threshold {
		shares = shares[:ts.Threshold]
	}

	// Sort shares by index for consistent reconstruction
	sort.Slice(shares, func(i, j int) bool {
		return shares[i].Index < shares[j].Index
	})

	// Lagrange interpolation to reconstruct secret
	secret := big.NewInt(0)

	for i, shareI := range shares {
		yi := new(big.Int).SetBytes(shareI.Value)
		xi := big.NewInt(int64(shareI.Index))

		// Calculate Lagrange basis polynomial l_i(0)
		li := big.NewInt(1)

		for j, shareJ := range shares {
			if i == j {
				continue
			}

			xj := big.NewInt(int64(shareJ.Index))

			// numerator: (0 - xj) = -xj
			numerator := new(big.Int).Neg(xj)
			numerator.Mod(numerator, prime)

			// denominator: (xi - xj)
			denominator := new(big.Int).Sub(xi, xj)
			denominator.Mod(denominator, prime)

			// denominator^-1 mod prime
			invDenominator := new(big.Int).ModInverse(denominator, prime)
			if invDenominator == nil {
				return nil, errors.New("modular inverse does not exist")
			}

			// li *= numerator * invDenominator
			li.Mul(li, numerator)
			li.Mul(li, invDenominator)
			li.Mod(li, prime)
		}

		// secret += yi * li
		term := new(big.Int).Mul(yi, li)
		secret.Add(secret, term)
		secret.Mod(secret, prime)
	}

	return secret.Bytes(), nil
}

// Large prime for finite field operations (larger than 2^256)
var prime, _ = new(big.Int).SetString(
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141",
	16,
)

// evaluatePolynomial evaluates the polynomial at point x
func evaluatePolynomial(coefficients []*big.Int, x *big.Int) *big.Int {
	result := big.NewInt(0)
	xPower := big.NewInt(1)

	for _, coeff := range coefficients {
		term := new(big.Int).Mul(coeff, xPower)
		result.Add(result, term)
		result.Mod(result, prime)

		xPower.Mul(xPower, x)
		xPower.Mod(xPower, prime)
	}

	return result
}

// calculateShareChecksum computes a checksum for share verification
func calculateShareChecksum(value []byte, index int) string {
	h := sha3.New256()
	h.Write(value)
	h.Write([]byte(fmt.Sprintf("share-%d", index)))
	return hex.EncodeToString(h.Sum(nil)[:16])
}

// verifyShareChecksum verifies a share's checksum
func verifyShareChecksum(share Share) bool {
	expected := calculateShareChecksum(share.Value, share.Index)
	return share.Checksum == expected
}

// DistributeShares distributes shares to different nodes/servers
func (ts *ThresholdScheme) DistributeShares(shares []Share, nodeIDs []string) (map[string]Share, error) {
	if len(nodeIDs) != len(shares) {
		return nil, errors.New("number of nodes must match number of shares")
	}

	distribution := make(map[string]Share)
	for i, nodeID := range nodeIDs {
		distribution[nodeID] = shares[i]
	}

	return distribution, nil
}

// VerifyShares verifies all shares are valid
func (ts *ThresholdScheme) VerifyShares(shares []Share) ([]Share, []Share) {
	valid := make([]Share, 0)
	invalid := make([]Share, 0)

	for _, share := range shares {
		if verifyShareChecksum(share) {
			valid = append(valid, share)
		} else {
			invalid = append(invalid, share)
		}
	}

	return valid, invalid
}

// GenerateRandomShares generates random shares for testing
func GenerateRandomShares(count int) ([]Share, error) {
	shares := make([]Share, count)
	for i := 0; i < count; i++ {
		value := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, value); err != nil {
			return nil, err
		}

		shares[i] = Share{
			Index:    i + 1,
			Value:    value,
			Checksum: calculateShareChecksum(value, i+1),
		}
	}
	return shares, nil
}

// EncryptWithThreshold performs encryption and splits the key into shares
func EncryptWithThreshold(data []byte, passphrase string, threshold, total int) (*EncryptedDoc, []Share, error) {
	// Encrypt document
	doc, err := EncryptDocument(data, passphrase)
	if err != nil {
		return nil, nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Generate master key for threshold scheme
	masterKey, err := GenerateRandomKey(32)
	if err != nil {
		return nil, nil, fmt.Errorf("master key generation failed: %w", err)
	}

	// Create threshold scheme
	scheme, err := NewThresholdScheme(threshold, total)
	if err != nil {
		return nil, nil, fmt.Errorf("threshold scheme creation failed: %w", err)
	}

	// Generate shares
	shares, err := scheme.GenerateThresholdShares(masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("share generation failed: %w", err)
	}

	doc.KeyShares = shares
	return doc, shares, nil
}

// ReconstructKeyFromShares reconstructs a key from threshold shares
func ReconstructKeyFromShares(shares []Share, threshold int) ([]byte, error) {
	scheme, err := NewThresholdScheme(threshold, len(shares))
	if err != nil {
		return nil, err
	}

	return scheme.ReconstructSecret(shares)
}

// HashShares creates a hash commitment of shares for verification
func HashShares(shares []Share) string {
	h := sha256.New()
	for _, share := range shares {
		h.Write(share.Value)
		h.Write([]byte(fmt.Sprintf("%d", share.Index)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
