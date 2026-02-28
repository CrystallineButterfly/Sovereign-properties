// Code generated — DO NOT EDIT.

//go:build !wasip1

package house_rwa

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
)

var (
	_ = errors.New
	_ = fmt.Errorf
	_ = big.NewInt
	_ = common.Big1
)

// HouseRWAMock is a mock implementation of HouseRWA for testing.
type HouseRWAMock struct {
	BPSDENOMINATOR          func() (*big.Int, error)
	KEYEXCHANGEEXPIRY       func() (*big.Int, error)
	MAXBILLSPERHOUSE        func() (*big.Int, error)
	MAXDOCUMENTSPERHOUSE    func() (*big.Int, error)
	MAXRENTALDURATION       func() (*big.Int, error)
	MINVALIDATORSTAKE       func() (*big.Int, error)
	PROTOCOLFEEBPS          func() (*big.Int, error)
	UPGRADEINTERFACEVERSION func() (string, error)
	Arbitrators             func(ArbitratorsInput) (bool, error)
	AuthorizedCREWorkflows  func(AuthorizedCREWorkflowsInput) (bool, error)
	BalanceOf               func(BalanceOfInput) (*big.Int, error)
	FeeRecipient            func() (common.Address, error)
	GetActiveRental         func(GetActiveRentalInput) (RentalAgreement, error)
	GetApproved             func(GetApprovedInput) (common.Address, error)
	GetBills                func(GetBillsInput) ([]Bill, error)
	GetHouseDetails         func(GetHouseDetailsInput) (House, error)
	GetListing              func(GetListingInput) (Listing, error)
	GetTotalBillsCount      func(GetTotalBillsCountInput) (*big.Int, error)
	HasKYC                  func(HasKYCInput) (bool, error)
	HighValueThresholdUSD   func() (*big.Int, error)
	HouseBills              func(HouseBillsInput) (HouseBillsOutput, error)
	Houses                  func(HousesInput) (HousesOutput, error)
	IsApprovedForAll        func(IsApprovedForAllInput) (bool, error)
	IsRented                func(IsRentedInput) (bool, error)
	KeyExchanges            func(KeyExchangesInput) (KeyExchangesOutput, error)
	KycInfo                 func(KycInfoInput) (KycInfoOutput, error)
	Listings                func(ListingsInput) (ListingsOutput, error)
	MinKYCLevelForHighValue func() (uint8, error)
	MinKYCLevelForMint      func() (uint8, error)
	MintingPaused           func() (bool, error)
	Name                    func() (string, error)
	NextTokenId             func() (*big.Int, error)
	Owner                   func() (common.Address, error)
	OwnerOf                 func(OwnerOfInput) (common.Address, error)
	Paused                  func() (bool, error)
	PaymentsPaused          func() (bool, error)
	PendingOwner            func() (common.Address, error)
	PendingRentalDeposits   func(PendingRentalDepositsInput) (*big.Int, error)
	PriceFeed               func() (common.Address, error)
	ProxiableUUID           func() ([32]byte, error)
	Rentals                 func(RentalsInput) (RentalsOutput, error)
	RentalsPaused           func() (bool, error)
	SalesPaused             func() (bool, error)
	SupportsInterface       func(SupportsInterfaceInput) (bool, error)
	Symbol                  func() (string, error)
	TokenByIndex            func(TokenByIndexInput) (*big.Int, error)
	TokenOfOwnerByIndex     func(TokenOfOwnerByIndexInput) (*big.Int, error)
	TokenURI                func(TokenURIInput) (string, error)
	TotalFeesCollected      func() (*big.Int, error)
	TotalSupply             func() (*big.Int, error)
	TrustedBillProviders    func(TrustedBillProvidersInput) (bool, error)
	Validators              func(ValidatorsInput) (ValidatorsOutput, error)
}

// NewHouseRWAMock creates a new HouseRWAMock for testing.
func NewHouseRWAMock(address common.Address, clientMock *evmmock.ClientCapability) *HouseRWAMock {
	mock := &HouseRWAMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["BPS_DENOMINATOR"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.BPSDENOMINATOR == nil {
				return nil, errors.New("BPS_DENOMINATOR method not mocked")
			}
			result, err := mock.BPSDENOMINATOR()
			if err != nil {
				return nil, err
			}
			return abi.Methods["BPS_DENOMINATOR"].Outputs.Pack(result)
		},
		string(abi.Methods["KEY_EXCHANGE_EXPIRY"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.KEYEXCHANGEEXPIRY == nil {
				return nil, errors.New("KEY_EXCHANGE_EXPIRY method not mocked")
			}
			result, err := mock.KEYEXCHANGEEXPIRY()
			if err != nil {
				return nil, err
			}
			return abi.Methods["KEY_EXCHANGE_EXPIRY"].Outputs.Pack(result)
		},
		string(abi.Methods["MAX_BILLS_PER_HOUSE"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MAXBILLSPERHOUSE == nil {
				return nil, errors.New("MAX_BILLS_PER_HOUSE method not mocked")
			}
			result, err := mock.MAXBILLSPERHOUSE()
			if err != nil {
				return nil, err
			}
			return abi.Methods["MAX_BILLS_PER_HOUSE"].Outputs.Pack(result)
		},
		string(abi.Methods["MAX_DOCUMENTS_PER_HOUSE"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MAXDOCUMENTSPERHOUSE == nil {
				return nil, errors.New("MAX_DOCUMENTS_PER_HOUSE method not mocked")
			}
			result, err := mock.MAXDOCUMENTSPERHOUSE()
			if err != nil {
				return nil, err
			}
			return abi.Methods["MAX_DOCUMENTS_PER_HOUSE"].Outputs.Pack(result)
		},
		string(abi.Methods["MAX_RENTAL_DURATION"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MAXRENTALDURATION == nil {
				return nil, errors.New("MAX_RENTAL_DURATION method not mocked")
			}
			result, err := mock.MAXRENTALDURATION()
			if err != nil {
				return nil, err
			}
			return abi.Methods["MAX_RENTAL_DURATION"].Outputs.Pack(result)
		},
		string(abi.Methods["MIN_VALIDATOR_STAKE"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MINVALIDATORSTAKE == nil {
				return nil, errors.New("MIN_VALIDATOR_STAKE method not mocked")
			}
			result, err := mock.MINVALIDATORSTAKE()
			if err != nil {
				return nil, err
			}
			return abi.Methods["MIN_VALIDATOR_STAKE"].Outputs.Pack(result)
		},
		string(abi.Methods["PROTOCOL_FEE_BPS"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PROTOCOLFEEBPS == nil {
				return nil, errors.New("PROTOCOL_FEE_BPS method not mocked")
			}
			result, err := mock.PROTOCOLFEEBPS()
			if err != nil {
				return nil, err
			}
			return abi.Methods["PROTOCOL_FEE_BPS"].Outputs.Pack(result)
		},
		string(abi.Methods["UPGRADE_INTERFACE_VERSION"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.UPGRADEINTERFACEVERSION == nil {
				return nil, errors.New("UPGRADE_INTERFACE_VERSION method not mocked")
			}
			result, err := mock.UPGRADEINTERFACEVERSION()
			if err != nil {
				return nil, err
			}
			return abi.Methods["UPGRADE_INTERFACE_VERSION"].Outputs.Pack(result)
		},
		string(abi.Methods["arbitrators"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Arbitrators == nil {
				return nil, errors.New("arbitrators method not mocked")
			}
			inputs := abi.Methods["arbitrators"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := ArbitratorsInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.Arbitrators(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["arbitrators"].Outputs.Pack(result)
		},
		string(abi.Methods["authorizedCREWorkflows"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AuthorizedCREWorkflows == nil {
				return nil, errors.New("authorizedCREWorkflows method not mocked")
			}
			inputs := abi.Methods["authorizedCREWorkflows"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := AuthorizedCREWorkflowsInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.AuthorizedCREWorkflows(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["authorizedCREWorkflows"].Outputs.Pack(result)
		},
		string(abi.Methods["balanceOf"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.BalanceOf == nil {
				return nil, errors.New("balanceOf method not mocked")
			}
			inputs := abi.Methods["balanceOf"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := BalanceOfInput{
				Owner: values[0].(common.Address),
			}

			result, err := mock.BalanceOf(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["balanceOf"].Outputs.Pack(result)
		},
		string(abi.Methods["feeRecipient"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.FeeRecipient == nil {
				return nil, errors.New("feeRecipient method not mocked")
			}
			result, err := mock.FeeRecipient()
			if err != nil {
				return nil, err
			}
			return abi.Methods["feeRecipient"].Outputs.Pack(result)
		},
		string(abi.Methods["getActiveRental"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetActiveRental == nil {
				return nil, errors.New("getActiveRental method not mocked")
			}
			inputs := abi.Methods["getActiveRental"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetActiveRentalInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetActiveRental(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getActiveRental"].Outputs.Pack(result)
		},
		string(abi.Methods["getApproved"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetApproved == nil {
				return nil, errors.New("getApproved method not mocked")
			}
			inputs := abi.Methods["getApproved"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetApprovedInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetApproved(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getApproved"].Outputs.Pack(result)
		},
		string(abi.Methods["getBills"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetBills == nil {
				return nil, errors.New("getBills method not mocked")
			}
			inputs := abi.Methods["getBills"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetBillsInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetBills(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getBills"].Outputs.Pack(result)
		},
		string(abi.Methods["getHouseDetails"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetHouseDetails == nil {
				return nil, errors.New("getHouseDetails method not mocked")
			}
			inputs := abi.Methods["getHouseDetails"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetHouseDetailsInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetHouseDetails(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getHouseDetails"].Outputs.Pack(result)
		},
		string(abi.Methods["getListing"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetListing == nil {
				return nil, errors.New("getListing method not mocked")
			}
			inputs := abi.Methods["getListing"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetListingInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetListing(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getListing"].Outputs.Pack(result)
		},
		string(abi.Methods["getTotalBillsCount"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetTotalBillsCount == nil {
				return nil, errors.New("getTotalBillsCount method not mocked")
			}
			inputs := abi.Methods["getTotalBillsCount"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetTotalBillsCountInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.GetTotalBillsCount(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getTotalBillsCount"].Outputs.Pack(result)
		},
		string(abi.Methods["hasKYC"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.HasKYC == nil {
				return nil, errors.New("hasKYC method not mocked")
			}
			inputs := abi.Methods["hasKYC"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := HasKYCInput{
				User: values[0].(common.Address),
			}

			result, err := mock.HasKYC(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["hasKYC"].Outputs.Pack(result)
		},
		string(abi.Methods["highValueThresholdUSD"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.HighValueThresholdUSD == nil {
				return nil, errors.New("highValueThresholdUSD method not mocked")
			}
			result, err := mock.HighValueThresholdUSD()
			if err != nil {
				return nil, err
			}
			return abi.Methods["highValueThresholdUSD"].Outputs.Pack(result)
		},
		string(abi.Methods["houseBills"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.HouseBills == nil {
				return nil, errors.New("houseBills method not mocked")
			}
			inputs := abi.Methods["houseBills"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := HouseBillsInput{
				Arg0: values[0].(*big.Int),
				Arg1: values[1].(*big.Int),
			}

			result, err := mock.HouseBills(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["houseBills"].Outputs.Pack(
				result.BillType,
				result.Amount,
				result.DueDate,
				result.PaidAt,
				result.Status,
				result.PaymentReference,
				result.IsRecurring,
				result.Provider,
				result.RecurrenceInterval,
			)
		},
		string(abi.Methods["houses"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Houses == nil {
				return nil, errors.New("houses method not mocked")
			}
			inputs := abi.Methods["houses"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := HousesInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.Houses(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["houses"].Outputs.Pack(
				result.HouseId,
				result.DocumentHash,
				result.DocumentURI,
				result.StorageType,
				result.OriginalOwner,
				result.MintedAt,
				result.IsVerified,
				result.DocumentCount,
			)
		},
		string(abi.Methods["isApprovedForAll"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.IsApprovedForAll == nil {
				return nil, errors.New("isApprovedForAll method not mocked")
			}
			inputs := abi.Methods["isApprovedForAll"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := IsApprovedForAllInput{
				Owner:    values[0].(common.Address),
				Operator: values[1].(common.Address),
			}

			result, err := mock.IsApprovedForAll(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["isApprovedForAll"].Outputs.Pack(result)
		},
		string(abi.Methods["isRented"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.IsRented == nil {
				return nil, errors.New("isRented method not mocked")
			}
			inputs := abi.Methods["isRented"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := IsRentedInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.IsRented(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["isRented"].Outputs.Pack(result)
		},
		string(abi.Methods["keyExchanges"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.KeyExchanges == nil {
				return nil, errors.New("keyExchanges method not mocked")
			}
			inputs := abi.Methods["keyExchanges"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := KeyExchangesInput{
				Arg0: values[0].([32]byte),
			}

			result, err := mock.KeyExchanges(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["keyExchanges"].Outputs.Pack(
				result.KeyHash,
				result.EncryptedKey,
				result.IntendedRecipient,
				result.CreatedAt,
				result.ExpiresAt,
				result.IsClaimed,
				result.ExchangeType,
			)
		},
		string(abi.Methods["kycInfo"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.KycInfo == nil {
				return nil, errors.New("kycInfo method not mocked")
			}
			inputs := abi.Methods["kycInfo"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := KycInfoInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.KycInfo(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["kycInfo"].Outputs.Pack(
				result.IsVerified,
				result.VerifiedAt,
				result.VerificationHash,
				result.VerificationLevel,
				result.ExpiryDate,
			)
		},
		string(abi.Methods["listings"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Listings == nil {
				return nil, errors.New("listings method not mocked")
			}
			inputs := abi.Methods["listings"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := ListingsInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.Listings(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["listings"].Outputs.Pack(
				result.ListingType,
				result.Price,
				result.PreferredToken,
				result.IsPrivateSale,
				result.AllowedBuyer,
				result.CreatedAt,
				result.ExpiresAt,
				result.PlatformFee,
			)
		},
		string(abi.Methods["minKYCLevelForHighValue"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MinKYCLevelForHighValue == nil {
				return nil, errors.New("minKYCLevelForHighValue method not mocked")
			}
			result, err := mock.MinKYCLevelForHighValue()
			if err != nil {
				return nil, err
			}
			return abi.Methods["minKYCLevelForHighValue"].Outputs.Pack(result)
		},
		string(abi.Methods["minKYCLevelForMint"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MinKYCLevelForMint == nil {
				return nil, errors.New("minKYCLevelForMint method not mocked")
			}
			result, err := mock.MinKYCLevelForMint()
			if err != nil {
				return nil, err
			}
			return abi.Methods["minKYCLevelForMint"].Outputs.Pack(result)
		},
		string(abi.Methods["mintingPaused"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.MintingPaused == nil {
				return nil, errors.New("mintingPaused method not mocked")
			}
			result, err := mock.MintingPaused()
			if err != nil {
				return nil, err
			}
			return abi.Methods["mintingPaused"].Outputs.Pack(result)
		},
		string(abi.Methods["name"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Name == nil {
				return nil, errors.New("name method not mocked")
			}
			result, err := mock.Name()
			if err != nil {
				return nil, err
			}
			return abi.Methods["name"].Outputs.Pack(result)
		},
		string(abi.Methods["nextTokenId"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.NextTokenId == nil {
				return nil, errors.New("nextTokenId method not mocked")
			}
			result, err := mock.NextTokenId()
			if err != nil {
				return nil, err
			}
			return abi.Methods["nextTokenId"].Outputs.Pack(result)
		},
		string(abi.Methods["owner"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Owner == nil {
				return nil, errors.New("owner method not mocked")
			}
			result, err := mock.Owner()
			if err != nil {
				return nil, err
			}
			return abi.Methods["owner"].Outputs.Pack(result)
		},
		string(abi.Methods["ownerOf"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.OwnerOf == nil {
				return nil, errors.New("ownerOf method not mocked")
			}
			inputs := abi.Methods["ownerOf"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := OwnerOfInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.OwnerOf(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["ownerOf"].Outputs.Pack(result)
		},
		string(abi.Methods["paused"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Paused == nil {
				return nil, errors.New("paused method not mocked")
			}
			result, err := mock.Paused()
			if err != nil {
				return nil, err
			}
			return abi.Methods["paused"].Outputs.Pack(result)
		},
		string(abi.Methods["paymentsPaused"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PaymentsPaused == nil {
				return nil, errors.New("paymentsPaused method not mocked")
			}
			result, err := mock.PaymentsPaused()
			if err != nil {
				return nil, err
			}
			return abi.Methods["paymentsPaused"].Outputs.Pack(result)
		},
		string(abi.Methods["pendingOwner"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PendingOwner == nil {
				return nil, errors.New("pendingOwner method not mocked")
			}
			result, err := mock.PendingOwner()
			if err != nil {
				return nil, err
			}
			return abi.Methods["pendingOwner"].Outputs.Pack(result)
		},
		string(abi.Methods["pendingRentalDeposits"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PendingRentalDeposits == nil {
				return nil, errors.New("pendingRentalDeposits method not mocked")
			}
			inputs := abi.Methods["pendingRentalDeposits"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := PendingRentalDepositsInput{
				Arg0: values[0].(*big.Int),
				Arg1: values[1].(common.Address),
			}

			result, err := mock.PendingRentalDeposits(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["pendingRentalDeposits"].Outputs.Pack(result)
		},
		string(abi.Methods["priceFeed"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PriceFeed == nil {
				return nil, errors.New("priceFeed method not mocked")
			}
			result, err := mock.PriceFeed()
			if err != nil {
				return nil, err
			}
			return abi.Methods["priceFeed"].Outputs.Pack(result)
		},
		string(abi.Methods["proxiableUUID"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.ProxiableUUID == nil {
				return nil, errors.New("proxiableUUID method not mocked")
			}
			result, err := mock.ProxiableUUID()
			if err != nil {
				return nil, err
			}
			return abi.Methods["proxiableUUID"].Outputs.Pack(result)
		},
		string(abi.Methods["rentals"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Rentals == nil {
				return nil, errors.New("rentals method not mocked")
			}
			inputs := abi.Methods["rentals"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := RentalsInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.Rentals(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["rentals"].Outputs.Pack(
				result.Renter,
				result.StartTime,
				result.EndTime,
				result.DepositAmount,
				result.MonthlyRent,
				result.IsActive,
				result.EncryptedAccessKeyHash,
				result.DisputeStatus,
			)
		},
		string(abi.Methods["rentalsPaused"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.RentalsPaused == nil {
				return nil, errors.New("rentalsPaused method not mocked")
			}
			result, err := mock.RentalsPaused()
			if err != nil {
				return nil, err
			}
			return abi.Methods["rentalsPaused"].Outputs.Pack(result)
		},
		string(abi.Methods["salesPaused"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.SalesPaused == nil {
				return nil, errors.New("salesPaused method not mocked")
			}
			result, err := mock.SalesPaused()
			if err != nil {
				return nil, err
			}
			return abi.Methods["salesPaused"].Outputs.Pack(result)
		},
		string(abi.Methods["supportsInterface"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.SupportsInterface == nil {
				return nil, errors.New("supportsInterface method not mocked")
			}
			inputs := abi.Methods["supportsInterface"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := SupportsInterfaceInput{
				InterfaceId: values[0].([4]byte),
			}

			result, err := mock.SupportsInterface(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["supportsInterface"].Outputs.Pack(result)
		},
		string(abi.Methods["symbol"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Symbol == nil {
				return nil, errors.New("symbol method not mocked")
			}
			result, err := mock.Symbol()
			if err != nil {
				return nil, err
			}
			return abi.Methods["symbol"].Outputs.Pack(result)
		},
		string(abi.Methods["tokenByIndex"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TokenByIndex == nil {
				return nil, errors.New("tokenByIndex method not mocked")
			}
			inputs := abi.Methods["tokenByIndex"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := TokenByIndexInput{
				Index: values[0].(*big.Int),
			}

			result, err := mock.TokenByIndex(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["tokenByIndex"].Outputs.Pack(result)
		},
		string(abi.Methods["tokenOfOwnerByIndex"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TokenOfOwnerByIndex == nil {
				return nil, errors.New("tokenOfOwnerByIndex method not mocked")
			}
			inputs := abi.Methods["tokenOfOwnerByIndex"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := TokenOfOwnerByIndexInput{
				Owner: values[0].(common.Address),
				Index: values[1].(*big.Int),
			}

			result, err := mock.TokenOfOwnerByIndex(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["tokenOfOwnerByIndex"].Outputs.Pack(result)
		},
		string(abi.Methods["tokenURI"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TokenURI == nil {
				return nil, errors.New("tokenURI method not mocked")
			}
			inputs := abi.Methods["tokenURI"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := TokenURIInput{
				TokenId: values[0].(*big.Int),
			}

			result, err := mock.TokenURI(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["tokenURI"].Outputs.Pack(result)
		},
		string(abi.Methods["totalFeesCollected"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalFeesCollected == nil {
				return nil, errors.New("totalFeesCollected method not mocked")
			}
			result, err := mock.TotalFeesCollected()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalFeesCollected"].Outputs.Pack(result)
		},
		string(abi.Methods["totalSupply"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalSupply == nil {
				return nil, errors.New("totalSupply method not mocked")
			}
			result, err := mock.TotalSupply()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalSupply"].Outputs.Pack(result)
		},
		string(abi.Methods["trustedBillProviders"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TrustedBillProviders == nil {
				return nil, errors.New("trustedBillProviders method not mocked")
			}
			inputs := abi.Methods["trustedBillProviders"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := TrustedBillProvidersInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.TrustedBillProviders(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["trustedBillProviders"].Outputs.Pack(result)
		},
		string(abi.Methods["validators"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Validators == nil {
				return nil, errors.New("validators method not mocked")
			}
			inputs := abi.Methods["validators"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := ValidatorsInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.Validators(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["validators"].Outputs.Pack(
				result.StakedAmount,
				result.StakedAt,
				result.Reputation,
				result.IsSlashed,
				result.SuccessfulValidations,
				result.FailedValidations,
			)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
