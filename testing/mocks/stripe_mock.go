package mocks

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// StripeMock implements a mock Stripe API for testing
type StripeMock struct {
	mu             sync.RWMutex
	Payments       map[string]*PaymentIntent
	Customers      map[string]*Customer
	Refunds        map[string]*Refund
	PaymentMethods map[string]*PaymentMethod
}

type PaymentIntent struct {
	ID            string            `json:"id"`
	Amount        int64             `json:"amount"`
	Currency      string            `json:"currency"`
	Status        string            `json:"status"`
	CustomerID    string            `json:"customer,omitempty"`
	PaymentMethod string            `json:"payment_method,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Created       int64             `json:"created"`
	Charges       []Charge          `json:"charges,omitempty"`
}

type Customer struct {
	ID       string            `json:"id"`
	Email    string            `json:"email"`
	Name     string            `json:"name"`
	Created  int64             `json:"created"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Refund struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Status   string `json:"status"`
	ChargeID string `json:"charge"`
	Created  int64  `json:"created"`
}

type PaymentMethod struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Card    *Card  `json:"card,omitempty"`
	Created int64  `json:"created"`
}

type Card struct {
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
}

type Charge struct {
	ID      string `json:"id"`
	Amount  int64  `json:"amount"`
	Status  string `json:"status"`
	Created int64  `json:"created"`
}

// NewStripeMock creates a new Stripe mock instance
func NewStripeMock() *StripeMock {
	return &StripeMock{
		Payments:       make(map[string]*PaymentIntent),
		Customers:      make(map[string]*Customer),
		Refunds:        make(map[string]*Refund),
		PaymentMethods: make(map[string]*PaymentMethod),
	}
}

// CreatePaymentIntent creates a mock payment intent
func (m *StripeMock) CreatePaymentIntent(ctx context.Context, params *PaymentIntentParams) (*PaymentIntent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate required fields
	if params.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if params.Currency == "" {
		return nil, errors.New("currency is required")
	}

	// Create payment intent
	pi := &PaymentIntent{
		ID:            "pi_" + generateID(),
		Amount:        params.Amount,
		Currency:      params.Currency,
		Status:        "requires_confirmation",
		CustomerID:    params.CustomerID,
		PaymentMethod: params.PaymentMethod,
		Metadata:      params.Metadata,
		Created:       time.Now().Unix(),
	}

	// Auto-confirm for test cards
	if params.PaymentMethod != "" {
		pm, exists := m.PaymentMethods[params.PaymentMethod]
		if exists && pm.Card != nil && pm.Card.Last4 == "4242" {
			pi.Status = "succeeded"
			pi.Charges = []Charge{
				{
					ID:      "ch_" + generateID(),
					Amount:  params.Amount,
					Status:  "succeeded",
					Created: time.Now().Unix(),
				},
			}
		}
	}

	m.Payments[pi.ID] = pi
	return pi, nil
}

// GetPaymentIntent retrieves a payment intent
func (m *StripeMock) GetPaymentIntent(ctx context.Context, id string) (*PaymentIntent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pi, exists := m.Payments[id]
	if !exists {
		return nil, errors.New("payment intent not found")
	}

	return pi, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (m *StripeMock) ConfirmPaymentIntent(ctx context.Context, id string, params *ConfirmParams) (*PaymentIntent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pi, exists := m.Payments[id]
	if !exists {
		return nil, errors.New("payment intent not found")
	}

	if pi.Status == "succeeded" {
		return nil, errors.New("payment intent already succeeded")
	}

	// Simulate confirmation
	pi.Status = "succeeded"
	pi.Charges = []Charge{
		{
			ID:      "ch_" + generateID(),
			Amount:  pi.Amount,
			Status:  "succeeded",
			Created: time.Now().Unix(),
		},
	}

	return pi, nil
}

// CreateRefund creates a refund
func (m *StripeMock) CreateRefund(ctx context.Context, params *RefundParams) (*Refund, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the payment intent
	pi, exists := m.Payments[params.PaymentIntentID]
	if !exists {
		return nil, errors.New("payment intent not found")
	}

	if pi.Status != "succeeded" {
		return nil, errors.New("cannot refund non-succeeded payment")
	}

	refundAmount := params.Amount
	if refundAmount == 0 {
		refundAmount = pi.Amount
	}

	refund := &Refund{
		ID:       "re_" + generateID(),
		Amount:   refundAmount,
		Status:   "succeeded",
		ChargeID: pi.Charges[0].ID,
		Created:  time.Now().Unix(),
	}

	m.Refunds[refund.ID] = refund
	return refund, nil
}

// CreateCustomer creates a customer
func (m *StripeMock) CreateCustomer(ctx context.Context, params *CustomerParams) (*Customer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	customer := &Customer{
		ID:       "cus_" + generateID(),
		Email:    params.Email,
		Name:     params.Name,
		Created:  time.Now().Unix(),
		Metadata: params.Metadata,
	}

	m.Customers[customer.ID] = customer
	return customer, nil
}

// CreatePaymentMethod creates a payment method
func (m *StripeMock) CreatePaymentMethod(ctx context.Context, params *PaymentMethodParams) (*PaymentMethod, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pm := &PaymentMethod{
		ID:      "pm_" + generateID(),
		Type:    params.Type,
		Created: time.Now().Unix(),
	}

	if params.Card != nil {
		pm.Card = &Card{
			Brand:    params.Card.Brand,
			Last4:    params.Card.Last4,
			ExpMonth: params.Card.ExpMonth,
			ExpYear:  params.Card.ExpYear,
		}
	}

	m.PaymentMethods[pm.ID] = pm
	return pm, nil
}

// Helper types

type PaymentIntentParams struct {
	Amount        int64
	Currency      string
	CustomerID    string
	PaymentMethod string
	Metadata      map[string]string
}

type ConfirmParams struct {
	PaymentMethod string
}

type RefundParams struct {
	PaymentIntentID string
	Amount          int64
}

type CustomerParams struct {
	Email    string
	Name     string
	Metadata map[string]string
}

type PaymentMethodParams struct {
	Type string
	Card *CardParams
}

type CardParams struct {
	Number   string
	ExpMonth int
	ExpYear  int
	CVC      string
	Brand    string
	Last4    string
}

// generateID generates a random ID for testing
func generateID() string {
	return time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// Reset clears all mock data
func (m *StripeMock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Payments = make(map[string]*PaymentIntent)
	m.Customers = make(map[string]*Customer)
	m.Refunds = make(map[string]*Refund)
	m.PaymentMethods = make(map[string]*PaymentMethod)
}

// GetTotalVolume returns total payment volume for testing
func (m *StripeMock) GetTotalVolume() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int64
	for _, pi := range m.Payments {
		if pi.Status == "succeeded" {
			total += pi.Amount
		}
	}
	return total
}
