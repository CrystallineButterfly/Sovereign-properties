//go:build wasip1

// Package encryption provides secure document encryption using AES-256-GCM
// and threshold encryption for distributed key management in Chainlink CRE.
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

const (
	// Key sizes
	AESKeySize = 32 // 256 bits
	NonceSize  = 12 // 96 bits for GCM
	SaltSize   = 32 // 256 bits
	TagSize    = 16 // 128 bits for GCM authentication tag

	// Argon2 parameters for key derivation
	Argon2Time    = 3
	Argon2Memory  = 64 * 1024 // 64 MB
	Argon2Threads = 4
	Argon2KeyLen  = 32
)

// EncryptedDoc represents an encrypted document with all necessary metadata
type EncryptedDoc struct {
	Ciphertext []byte  `json:"ciphertext"`
	Nonce      []byte  `json:"nonce"`
	Salt       []byte  `json:"salt"`
	Tag        []byte  `json:"tag"`
	KeyShares  []Share `json:"key_shares,omitempty"`
	DocumentID string  `json:"document_id"`
	Algorithm  string  `json:"algorithm"`
	Version    int     `json:"version"`
}

// Share represents a single key share in threshold encryption
type Share struct {
	Index    int    `json:"index"`
	Value    []byte `json:"value"`
	Checksum string `json:"checksum"`
}

// EncryptionResult contains the encryption output and key shares
type EncryptionResult struct {
	Document    *EncryptedDoc
	MasterKey   []byte
	Shares      []Share
	DocumentURI string
}

// TEEContext represents the Trusted Execution Environment context
type TEEContext struct {
	EnclaveID   string
	Attestation []byte
	IsValid     bool
}

// EncryptDocument encrypts data using AES-256-GCM with secure key derivation
func EncryptDocument(data []byte, passphrase string) (*EncryptedDoc, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data provided")
	}

	// Generate random salt
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key using Argon2id
	key := deriveKey(passphrase, salt)

	// Generate nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt data with additional authenticated data (AAD)
	aad := generateAAD(salt, nonce)
	ciphertext := aead.Seal(nil, nonce, data, aad)

	// Separate ciphertext and authentication tag
	if len(ciphertext) < TagSize {
		return nil, errors.New("ciphertext too short")
	}

	actualCiphertext := ciphertext[:len(ciphertext)-TagSize]
	tag := ciphertext[len(ciphertext)-TagSize:]

	// Generate document ID
	docHash := sha256.Sum256(data)
	documentID := hex.EncodeToString(docHash[:16])

	return &EncryptedDoc{
		Ciphertext: actualCiphertext,
		Nonce:      nonce,
		Salt:       salt,
		Tag:        tag,
		DocumentID: documentID,
		Algorithm:  "AES-256-GCM",
		Version:    1,
	}, nil
}

// DecryptDocument decrypts an encrypted document
func DecryptDocument(doc *EncryptedDoc, passphrase string) ([]byte, error) {
	if doc == nil {
		return nil, errors.New("nil document")
	}

	// Derive key
	key := deriveKey(passphrase, doc.Salt)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Reconstruct ciphertext with tag
	ciphertext := append(doc.Ciphertext, doc.Tag...)

	// Generate AAD
	aad := generateAAD(doc.Salt, doc.Nonce)

	// Decrypt
	plaintext, err := aead.Open(nil, doc.Nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// DecryptInTEE performs decryption within a Trusted Execution Environment
func DecryptInTEE(encryptedDoc *EncryptedDoc, teeContext *TEEContext) ([]byte, error) {
	if !teeContext.IsValid {
		return nil, errors.New("invalid TEE attestation")
	}

	if encryptedDoc == nil {
		return nil, errors.New("nil document")
	}

	// In a real implementation, this would:
	// 1. Verify TEE attestation with the hardware provider
	// 2. Establish secure channel with enclave
	// 3. Send encrypted data to enclave
	// 4. Enclave decrypts using sealed key
	// 5. Return plaintext through secure channel

	// For now, simulate TEE operation
	// In production, this would interface with Intel SGX, AMD SEV, or similar
	simulatedKey := deriveTEEKey(teeContext.EnclaveID)

	block, err := aes.NewCipher(simulatedKey)
	if err != nil {
		return nil, fmt.Errorf("TEE cipher creation failed: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("TEE GCM creation failed: %w", err)
	}

	ciphertext := append(encryptedDoc.Ciphertext, encryptedDoc.Tag...)
	aad := generateAAD(encryptedDoc.Salt, encryptedDoc.Nonce)

	plaintext, err := aead.Open(nil, encryptedDoc.Nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("TEE decryption failed: %w", err)
	}

	return plaintext, nil
}

// deriveKey derives an encryption key using Argon2id
func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLen,
	)
}

// deriveTEEKey derives a key for TEE operations
func deriveTEEKey(enclaveID string) []byte {
	hkdfReader := hkdf.New(sha256.New, []byte(enclaveID), nil, []byte("tee-encryption"))
	key := make([]byte, AESKeySize)
	io.ReadFull(hkdfReader, key)
	return key
}

// generateAAD generates Additional Authenticated Data
func generateAAD(salt, nonce []byte) []byte {
	h := sha256.New()
	h.Write(salt)
	h.Write(nonce)
	h.Write([]byte("RWA-House-CRE-v1"))
	return h.Sum(nil)
}

// EncryptWithPublicKey encrypts data using ECIES with ephemeral ECDH
func EncryptWithPublicKey(data []byte, pubKey *ecdsa.PublicKey) ([]byte, error) {
	if pubKey == nil {
		return nil, errors.New("nil public key")
	}

	// Generate ephemeral key pair
	ephemeralKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ephemeral key generation failed: %w", err)
	}

	// Compute shared secret using ECDH
	x, _ := ephemeralKey.PublicKey.Curve.ScalarMult(
		pubKey.X, pubKey.Y, ephemeralKey.D.Bytes(),
	)
	sharedSecret := sha256.Sum256(x.Bytes())

	// Encrypt data with shared secret
	block, err := aes.NewCipher(sharedSecret[:])
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %w", err)
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation failed: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation failed: %w", err)
	}

	encrypted := aead.Seal(nil, nonce, data, nil)

	// Package: ephemeral public key (uncompressed) + nonce + ciphertext
	ephemeralPubBytes := elliptic.Marshal(
		ephemeralKey.PublicKey.Curve,
		ephemeralKey.PublicKey.X,
		ephemeralKey.PublicKey.Y,
	)

	result := make([]byte, 0, len(ephemeralPubBytes)+len(nonce)+len(encrypted))
	result = append(result, ephemeralPubBytes...)
	result = append(result, nonce...)
	result = append(result, encrypted...)

	return result, nil
}

// DecryptWithPrivateKey decrypts data using ECIES
func DecryptWithPrivateKey(encryptedData []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	if privKey == nil {
		return nil, errors.New("nil private key")
	}

	if len(encryptedData) < 65+NonceSize+TagSize {
		return nil, errors.New("encrypted data too short")
	}

	// Extract ephemeral public key (65 bytes for uncompressed P-256)
	ephemeralPubBytes := encryptedData[:65]
	nonce := encryptedData[65 : 65+NonceSize]
	ciphertext := encryptedData[65+NonceSize:]

	// Parse ephemeral public key
	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, ephemeralPubBytes)
	if x == nil {
		return nil, errors.New("invalid ephemeral public key")
	}

	ephemeralPub := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// Compute shared secret
	x, _ = curve.ScalarMult(ephemeralPub.X, ephemeralPub.Y, privKey.D.Bytes())
	sharedSecret := sha256.Sum256(x.Bytes())

	// Decrypt
	block, err := aes.NewCipher(sharedSecret[:])
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation failed: %w", err)
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// SerializeEncryptedDoc serializes an encrypted document to JSON
func SerializeEncryptedDoc(doc *EncryptedDoc) ([]byte, error) {
	return json.Marshal(doc)
}

// DeserializeEncryptedDoc deserializes an encrypted document from JSON
func DeserializeEncryptedDoc(data []byte) (*EncryptedDoc, error) {
	var doc EncryptedDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("deserialization failed: %w", err)
	}
	return &doc, nil
}

// SerializeShares serializes key shares to base64-encoded JSON
func SerializeShares(shares []Share) (string, error) {
	data, err := json.Marshal(shares)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DeserializeShares deserializes key shares from base64-encoded JSON
func DeserializeShares(encoded string) ([]Share, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var shares []Share
	if err := json.Unmarshal(data, &shares); err != nil {
		return nil, err
	}
	return shares, nil
}

// VerifyShare verifies a single key share integrity
func VerifyShare(share Share) bool {
	expectedChecksum := calculateChecksum(share.Value, share.Index)
	return subtle.ConstantTimeCompare([]byte(share.Checksum), []byte(expectedChecksum)) == 1
}

// calculateChecksum calculates a checksum for a share
func calculateChecksum(value []byte, index int) string {
	h := sha256.New()
	h.Write(value)
	h.Write([]byte(fmt.Sprintf("%d", index)))
	return hex.EncodeToString(h.Sum(nil)[:8])
}

// SecureCompare performs constant-time comparison
func SecureCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// GenerateRandomKey generates a cryptographically secure random key
func GenerateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
