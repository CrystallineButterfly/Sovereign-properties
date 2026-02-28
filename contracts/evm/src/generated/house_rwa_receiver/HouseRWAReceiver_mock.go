// Code generated — DO NOT EDIT.

//go:build !wasip1

package house_rwa_receiver

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

// HouseRWAReceiverMock is a mock implementation of HouseRWAReceiver for testing.
type HouseRWAReceiverMock struct {
	AllowedSelectors func(AllowedSelectorsInput) (bool, error)
	Forwarder        func() (common.Address, error)
	HouseRWA         func() (common.Address, error)
	Owner            func() (common.Address, error)
	PendingOwner     func() (common.Address, error)
}

// NewHouseRWAReceiverMock creates a new HouseRWAReceiverMock for testing.
func NewHouseRWAReceiverMock(address common.Address, clientMock *evmmock.ClientCapability) *HouseRWAReceiverMock {
	mock := &HouseRWAReceiverMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["allowedSelectors"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AllowedSelectors == nil {
				return nil, errors.New("allowedSelectors method not mocked")
			}
			inputs := abi.Methods["allowedSelectors"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := AllowedSelectorsInput{
				Arg0: values[0].([4]byte),
			}

			result, err := mock.AllowedSelectors(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["allowedSelectors"].Outputs.Pack(result)
		},
		string(abi.Methods["forwarder"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Forwarder == nil {
				return nil, errors.New("forwarder method not mocked")
			}
			result, err := mock.Forwarder()
			if err != nil {
				return nil, err
			}
			return abi.Methods["forwarder"].Outputs.Pack(result)
		},
		string(abi.Methods["houseRWA"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.HouseRWA == nil {
				return nil, errors.New("houseRWA method not mocked")
			}
			result, err := mock.HouseRWA()
			if err != nil {
				return nil, err
			}
			return abi.Methods["houseRWA"].Outputs.Pack(result)
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
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
