package evm

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"

	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

func TestUnpackTransmitTxInput(t *testing.T) {
	registryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)

	packer := &evmRegistryPackerV2_1{abi: registryABI}
	decodedReport, err := packer.UnpackTransmitTxInput(hexutil.MustDecode("0x00011a04d404e571ead64b2f08cfae623a0d96b9beb326c20e322001cbbd344700000000000000000000000000000000000000000000000000000000000d580e35681c68a0426c30f4686e837c0cd7864200f48dbfe48c80c51f92aa5ac607b300000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002e00000000000000000000000000000000000000000000000000000000000000360000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000773594000000000000000000000000000000000000000000000000000010fb9cd2f34a00000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001de1256139081c6b165a3aee0432f605d3dee0e6087ea53b46ca9478c253ea9c8000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000827075c4bacd41884f60c2ca7af3630400bedd92ad7ad0ba4e1f000e70297de0573e180000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a8c0000000000000000000000000000000000000000000000000000000000086470000000000000000000000000000000000000000000000000000000000000000326e2b521089d44f1457ae51b3f8d76e8577e08c4af9374bdc62aebbfad081a78a13941ab209ad44a905ee0fd704a46b2ebc022dcb60659bed87342fd94dadb70827af523f59c7c9bb8dcc77e959b0476869612e8cf84e63a2e9a5617290633f70000000000000000000000000000000000000000000000000000000000000003723d77998618c5959396115fc61380215e0395f68c18a6cf0647c3e759ee013040c2967fdd369aac59b464f931dacd7b8863498757eda53a9f6a4b6150f2dbe640771f3c242c297265c36c5e78f4c660ae74dcd1f5bda8687b6afed3d3f27e0d"))
	assert.Nil(t, err)

	// We expect one upkeep ID in the report at block number
	var expectedBlock uint32 = 8548469
	expectedID, _ := new(big.Int).SetString("100445849710294316610676143149039812931260394722330855891004881602834541226440", 10)

	assert.Equal(t, len(decodedReport), 1)

	rpt, ok := decodedReport[0].(EVMAutomationUpkeepResult21)
	assert.True(t, ok)

	assert.Equal(t, rpt.Block, expectedBlock)
	assert.Equal(t, rpt.ID.String(), expectedID.String())
}

func TestUnpackTransmitTxInputErrors(t *testing.T) {

	tests := []struct {
		Name    string
		RawData string
	}{
		{
			Name:    "Empty Data",
			RawData: "0x",
		},
		{
			Name:    "Random Data",
			RawData: "0x2f08cfae623a0d96b9beb326c20e322001cbbd344700",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			abi, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
			assert.Nil(t, err)

			packer := &evmRegistryPackerV2_1{abi: abi}
			_, err = packer.UnpackTransmitTxInput(hexutil.MustDecode(test.RawData))
			assert.NotNil(t, err)
		})
	}
}

func TestUnpackCheckResults(t *testing.T) {
	registryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		assert.Nil(t, err)
	}

	upkeepId, _ := new(big.Int).SetString("1843548457736589226156809205796175506139185429616502850435279853710366065936", 10)

	tests := []struct {
		Name           string
		UpkeepKey      ocr2keepers.UpkeepKey
		RawData        string
		ExpectedResult EVMAutomationUpkeepResult21
	}{
		{
			Name:      "upkeep not needed",
			UpkeepKey: ocr2keepers.UpkeepKey(fmt.Sprintf("19447615|%s", upkeepId)),
			RawData:   "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000421c000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: EVMAutomationUpkeepResult21{
				Block:            19447615,
				ID:               upkeepId,
				Eligible:         false,
				FailureReason:    UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED,
				GasUsed:          big.NewInt(16924),
				PerformData:      nil,
				FastGasWei:       big.NewInt(1000000000),
				LinkNative:       big.NewInt(3532383906411401),
				CheckBlockNumber: 0,
				CheckBlockHash:   [32]byte{},
				ExecuteGas:       5000000,
			},
		},
		{
			Name:      "target check reverted",
			UpkeepKey: ocr2keepers.UpkeepKey(fmt.Sprintf("19448272|%s", upkeepId)),
			RawData:   "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000007531000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000008914039bf676e20aad43a5642485e666575ed0d927a4b5679745e947e7d125ee2687c10000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000024462e8a50d00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000128c1d000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000009666565644944537472000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000184554482d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000000000184254432d5553442d415242495452554d2d544553544e45540000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d6265720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: EVMAutomationUpkeepResult21{
				Block:            19448272,
				ID:               upkeepId,
				Eligible:         false,
				FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
				GasUsed:          big.NewInt(30001),
				PerformData:      []byte{98, 232, 165, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 193, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 83, 116, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				FastGasWei:       big.NewInt(1000000000),
				LinkNative:       big.NewInt(3532383906411401),
				CheckBlockNumber: 8983555,
				CheckBlockHash:   [32]byte{155, 246, 118, 226, 10, 173, 67, 165, 100, 36, 133, 230, 102, 87, 94, 208, 217, 39, 164, 181, 103, 151, 69, 233, 71, 231, 209, 37, 238, 38, 135, 193},
				ExecuteGas:       5000000,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := &evmRegistryPackerV2_1{abi: registryABI}
			rs, err := packer.UnpackCheckResult(test.UpkeepKey, test.RawData)
			assert.Nil(t, err)
			assert.Equal(t, test.ExpectedResult, rs)
		})
	}
}

func TestUnpackPerformResult(t *testing.T) {
	registryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		assert.Nil(t, err)
	}

	tests := []struct {
		Name    string
		RawData string
	}{
		{
			Name:    "unpack success",
			RawData: "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000a52d",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := &evmRegistryPackerV2_1{abi: registryABI}
			rs, err := packer.UnpackSimulatePerformResult(test.RawData)
			assert.Nil(t, err)
			assert.True(t, rs)
		})
	}
}

func TestUnpackMercuryCallbackResult(t *testing.T) {
	registryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		assert.Nil(t, err)
	}

	tests := []struct {
		Name         string
		CallbackResp []byte
		UpkeepNeeded bool
		PerformData  []byte
		ErrorString  string
	}{
		{
			Name:         "unpack upkeep needed",
			CallbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 117, 97, 54, 121, 44, 145, 48, 168, 145, 51, 172, 64, 131, 204, 54, 198, 186, 169, 18, 232, 136, 106, 134, 39, 115, 38, 154, 170, 6, 227, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 193, 3, 12, 77, 215, 29, 228, 93, 26, 35, 179, 50, 246, 100, 137, 30, 42, 158, 250, 32, 22, 120, 155, 204, 125, 96, 212, 27, 57, 47, 88, 112, 125, 28, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 170, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 131, 144, 161, 34, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 130, 203, 146, 154, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 133, 100, 70, 83, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 15, 68, 151, 27, 218, 234, 151, 253, 29, 138, 59, 147, 216, 86, 82, 202, 156, 216, 128, 252, 22, 203, 82, 143, 236, 163, 169, 159, 50, 109, 203, 132, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 16, 203, 54, 77, 83, 46, 183, 200, 215, 27, 252, 39, 80, 154, 96, 138, 143, 133, 138, 13, 160, 75, 241, 255, 67, 155, 254, 34, 224, 166, 218, 102, 230, 143, 135, 248, 238, 231, 114, 244, 147, 243, 153, 198, 143, 252, 92, 169, 175, 161, 233, 232, 152, 118, 168, 54, 167, 144, 85, 242, 235, 105, 54, 246, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 113, 38, 35, 185, 62, 8, 131, 218, 83, 46, 227, 168, 238, 127, 4, 202, 207, 252, 216, 184, 250, 226, 126, 154, 72, 187, 71, 90, 247, 149, 230, 190, 1, 210, 116, 90, 140, 23, 15, 81, 32, 15, 57, 8, 209, 86, 204, 153, 31, 11, 138, 27, 108, 36, 114, 92, 220, 238, 61, 244, 180, 238, 243, 245, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 100, 83, 124, 126, 209, 163, 161, 185, 65, 44, 219, 214, 164, 118, 54, 102, 85, 244, 245, 247, 70, 199, 10, 201, 214, 103, 241, 47, 5, 211, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 189, 3, 54, 141, 184, 180, 160, 108, 142, 160, 173, 206, 250, 156, 144, 216, 72, 213, 199, 59, 231, 52, 116, 150, 96, 123, 72, 215, 130, 151, 131, 149, 219, 45, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 170, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 246, 229, 108, 136, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 240, 113, 76, 105, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 253, 89, 140, 167, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 15, 68, 151, 27, 218, 234, 151, 253, 29, 138, 59, 147, 216, 86, 82, 202, 156, 216, 128, 252, 22, 203, 82, 143, 236, 163, 169, 159, 50, 109, 203, 132, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 24, 135, 38, 53, 221, 121, 90, 51, 43, 35, 56, 237, 245, 153, 230, 14, 35, 121, 20, 100, 254, 234, 214, 52, 24, 131, 193, 116, 226, 183, 196, 225, 78, 205, 181, 32, 77, 174, 88, 78, 136, 124, 110, 170, 145, 167, 190, 71, 133, 215, 94, 71, 171, 213, 15, 67, 62, 101, 152, 59, 132, 76, 117, 226, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 99, 160, 121, 17, 169, 117, 133, 205, 32, 52, 90, 255, 127, 128, 177, 242, 86, 252, 103, 102, 57, 163, 118, 206, 141, 102, 190, 120, 28, 74, 245, 170, 21, 66, 46, 80, 74, 30, 169, 74, 132, 96, 96, 232, 90, 4, 143, 93, 83, 6, 83, 38, 69, 218, 4, 106, 96, 49, 126, 169, 68, 11, 97, 189, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: true,
			PerformData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 117, 97, 54, 121, 44, 145, 48, 168, 145, 51, 172, 64, 131, 204, 54, 198, 186, 169, 18, 232, 136, 106, 134, 39, 115, 38, 154, 170, 6, 227, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 193, 3, 12, 77, 215, 29, 228, 93, 26, 35, 179, 50, 246, 100, 137, 30, 42, 158, 250, 32, 22, 120, 155, 204, 125, 96, 212, 27, 57, 47, 88, 112, 125, 28, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 170, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 131, 144, 161, 34, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 130, 203, 146, 154, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 133, 100, 70, 83, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 15, 68, 151, 27, 218, 234, 151, 253, 29, 138, 59, 147, 216, 86, 82, 202, 156, 216, 128, 252, 22, 203, 82, 143, 236, 163, 169, 159, 50, 109, 203, 132, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 16, 203, 54, 77, 83, 46, 183, 200, 215, 27, 252, 39, 80, 154, 96, 138, 143, 133, 138, 13, 160, 75, 241, 255, 67, 155, 254, 34, 224, 166, 218, 102, 230, 143, 135, 248, 238, 231, 114, 244, 147, 243, 153, 198, 143, 252, 92, 169, 175, 161, 233, 232, 152, 118, 168, 54, 167, 144, 85, 242, 235, 105, 54, 246, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 113, 38, 35, 185, 62, 8, 131, 218, 83, 46, 227, 168, 238, 127, 4, 202, 207, 252, 216, 184, 250, 226, 126, 154, 72, 187, 71, 90, 247, 149, 230, 190, 1, 210, 116, 90, 140, 23, 15, 81, 32, 15, 57, 8, 209, 86, 204, 153, 31, 11, 138, 27, 108, 36, 114, 92, 220, 238, 61, 244, 180, 238, 243, 245, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 100, 83, 124, 126, 209, 163, 161, 185, 65, 44, 219, 214, 164, 118, 54, 102, 85, 244, 245, 247, 70, 199, 10, 201, 214, 103, 241, 47, 5, 211, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 189, 3, 54, 141, 184, 180, 160, 108, 142, 160, 173, 206, 250, 156, 144, 216, 72, 213, 199, 59, 231, 52, 116, 150, 96, 123, 72, 215, 130, 151, 131, 149, 219, 45, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 170, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 246, 229, 108, 136, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 240, 113, 76, 105, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 111, 253, 89, 140, 167, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 15, 68, 151, 27, 218, 234, 151, 253, 29, 138, 59, 147, 216, 86, 82, 202, 156, 216, 128, 252, 22, 203, 82, 143, 236, 163, 169, 159, 50, 109, 203, 132, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 64, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 24, 135, 38, 53, 221, 121, 90, 51, 43, 35, 56, 237, 245, 153, 230, 14, 35, 121, 20, 100, 254, 234, 214, 52, 24, 131, 193, 116, 226, 183, 196, 225, 78, 205, 181, 32, 77, 174, 88, 78, 136, 124, 110, 170, 145, 167, 190, 71, 133, 215, 94, 71, 171, 213, 15, 67, 62, 101, 152, 59, 132, 76, 117, 226, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 99, 160, 121, 17, 169, 117, 133, 205, 32, 52, 90, 255, 127, 128, 177, 242, 86, 252, 103, 102, 57, 163, 118, 206, 141, 102, 190, 120, 28, 74, 245, 170, 21, 66, 46, 80, 74, 30, 169, 74, 132, 96, 96, 232, 90, 4, 143, 93, 83, 6, 83, 38, 69, 218, 4, 106, 96, 49, 126, 169, 68, 11, 97, 189, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:         "unpack upkeep not needed",
			CallbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 117, 97, 54, 121, 44, 145, 48, 168, 145, 51, 172, 64, 131, 204, 54, 198, 186, 169, 18, 232, 136, 106, 134, 39, 115, 38, 154, 170, 6, 227, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 31, 84, 1, 214, 103, 34, 170, 203, 87, 228, 153, 158, 168, 164, 185, 66, 10, 86, 151, 125, 155, 115, 131, 78, 233, 162, 240, 58, 180, 72, 243, 115, 251, 133, 109, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 180, 80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 53, 188, 231, 171, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 53, 63, 76, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 41, 54, 238, 142, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 83, 91, 231, 79, 120, 88, 2, 27, 207, 212, 17, 12, 48, 161, 221, 17, 228, 209, 58, 231, 146, 100, 61, 76, 90, 51, 29, 163, 98, 244, 113, 129, 31, 99, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 83, 88, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 247, 43, 231, 91, 54, 189, 170, 47, 66, 212, 187, 44, 13, 168, 87, 189, 51, 193, 77, 158, 195, 231, 222, 231, 26, 2, 138, 99, 121, 20, 23, 161, 194, 188, 247, 78, 246, 236, 42, 209, 43, 9, 106, 132, 75, 109, 68, 74, 41, 134, 245, 88, 250, 236, 89, 253, 9, 82, 1, 134, 232, 227, 75, 18, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 113, 233, 89, 125, 216, 137, 43, 17, 125, 242, 163, 135, 161, 170, 224, 52, 10, 101, 123, 3, 6, 157, 159, 247, 154, 80, 14, 121, 1, 17, 149, 68, 58, 166, 145, 229, 29, 229, 82, 0, 52, 209, 53, 38, 63, 203, 205, 180, 79, 150, 48, 94, 41, 7, 122, 210, 219, 24, 199, 179, 27, 245, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 192, 0, 1, 100, 83, 124, 126, 209, 163, 161, 185, 65, 44, 219, 214, 164, 118, 54, 102, 85, 244, 245, 247, 70, 199, 10, 201, 214, 103, 241, 47, 5, 211, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 31, 81, 1, 149, 170, 40, 163, 85, 94, 22, 80, 44, 48, 162, 37, 249, 57, 84, 205, 16, 11, 109, 178, 162, 55, 33, 170, 97, 144, 59, 199, 46, 91, 123, 89, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 93, 180, 80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 109, 47, 161, 93, 137, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 109, 40, 224, 245, 93, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 109, 54, 97, 197, 181, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 83, 91, 231, 79, 120, 88, 2, 27, 207, 212, 17, 12, 48, 161, 221, 17, 228, 209, 58, 231, 146, 100, 61, 76, 90, 51, 29, 163, 98, 244, 113, 129, 31, 99, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 41, 83, 88, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 172, 16, 186, 158, 160, 113, 175, 66, 61, 76, 123, 96, 96, 158, 69, 193, 63, 218, 146, 117, 51, 147, 158, 240, 47, 247, 75, 225, 146, 201, 101, 53, 79, 11, 234, 195, 245, 179, 152, 73, 138, 65, 124, 143, 63, 232, 157, 155, 91, 213, 208, 208, 45, 85, 79, 205, 48, 105, 180, 219, 155, 133, 192, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 104, 23, 211, 145, 155, 142, 212, 94, 72, 76, 57, 224, 167, 225, 128, 173, 132, 32, 162, 1, 33, 116, 254, 254, 101, 104, 163, 66, 91, 228, 102, 200, 84, 144, 32, 33, 238, 108, 79, 183, 172, 159, 133, 96, 243, 184, 102, 44, 180, 174, 92, 2, 28, 233, 218, 44, 168, 192, 191, 253, 237, 13, 183, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: false,
			PerformData:  nil,
		},
		{
			Name:         "unpack malformed data",
			CallbackResp: []byte{0, 0, 0, 23, 4, 163, 66, 91, 228, 102, 200, 84, 144, 233, 218, 44, 168, 192, 191, 253, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: false,
			PerformData:  nil,
			ErrorString:  "abi: improperly encoded boolean value: unpack checkUpkeep return: ",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := &evmRegistryPackerV2_1{abi: registryABI}
			needed, pd, _, _, err := packer.UnpackMercuryLookupResult(test.CallbackResp)

			if test.ErrorString != "" {
				assert.EqualError(t, err, test.ErrorString+hexutil.Encode(test.CallbackResp))
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, test.UpkeepNeeded, needed)
			assert.Equal(t, test.PerformData, pd)
		})
	}
}

func TestUnpackLogTriggerConfig(t *testing.T) {
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.NoError(t, err)
	tests := []struct {
		name    string
		raw     []byte
		res     iregistry21.KeeperRegistryBase21LogTriggerConfig
		errored bool
	}{
		{
			"happy flow",
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000007456fadf415b7c34b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			iregistry21.KeeperRegistryBase21LogTriggerConfig{
				ContractAddress: common.HexToAddress("0x7456FadF415b7c34B1182Bd20B0537977e945e3E"),
				Topic0:          [32]uint8{0x3d, 0x53, 0xa3, 0x95, 0x50, 0xe0, 0x46, 0x88, 0x6, 0x58, 0x27, 0xf3, 0xbb, 0x86, 0x58, 0x4c, 0xb0, 0x7, 0xab, 0x9e, 0xbc, 0xa7, 0xeb, 0xd5, 0x28, 0xe7, 0x30, 0x1c, 0x9c, 0x31, 0xeb, 0x5d},
			},
			false,
		},
		{
			"invalid",
			func() []byte {
				b, _ := hexutil.Decode("0x000000000000000000000000b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			iregistry21.KeeperRegistryBase21LogTriggerConfig{},
			true,
		},
	}

	packer := &evmRegistryPackerV2_1{abi: keeperRegistryABI}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := packer.UnpackLogTriggerConfig(tc.raw)
			if tc.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.res, res)
			}
		})
	}
}
