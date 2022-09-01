package taroscript

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taro/address"
	"github.com/lightninglabs/taro/asset"
	"github.com/lightninglabs/taro/commitment"
	"github.com/lightninglabs/taro/mssmt"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/stretchr/testify/require"
)

// spendData represents the collection of structs needed to begin a spend.
type spendData struct {
	collectAmt                     uint64
	normalAmt1                     uint64
	normalAmt2                     uint64
	genesis1                       asset.Genesis
	genesis1collect                asset.Genesis
	spenderPrivKey                 btcec.PrivateKey
	spenderPubKey                  btcec.PublicKey
	spenderScriptKey               btcec.PublicKey
	spenderDescriptor              keychain.KeyDescriptor
	receiverPrivKey                btcec.PrivateKey
	receiverPubKey                 btcec.PublicKey
	familyKey                      asset.FamilyKey
	address1                       address.Taro
	address1CollectFamily          address.Taro
	address2                       address.Taro
	address1StateKey               [32]byte
	address1CollectFamilyStateKey  [32]byte
	address2StateKey               [32]byte
	asset1                         asset.Asset
	asset1CollectFamily            asset.Asset
	asset2                         asset.Asset
	asset1PrevID                   asset.PrevID
	asset1CollectFamilyPrevID      asset.PrevID
	asset2PrevID                   asset.PrevID
	asset1InputAssets              commitment.InputSet
	asset1CollectFamilyInputAssets commitment.InputSet
	asset2InputAssets              commitment.InputSet
	asset1TaroTree                 commitment.TaroCommitment
	asset1CollectFamilyTaroTree    commitment.TaroCommitment
	asset2TaroTree                 commitment.TaroCommitment
}

var (
	key1Bytes, _ = hex.DecodeString(
		"a0afeb165f0ec36880b68e0baabd9ad9c62fd1a69aa998bc30e9a346202e" +
			"078e",
	)
	key2Bytes, _ = hex.DecodeString(
		"a0afeb165f0ec36880b68e0baabd9ad9c62fd1a69aa998bc30e9a346202e" +
			"078d",
	)
)

func randKey(t *testing.T) *btcec.PrivateKey {
	t.Helper()
	key, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	return key
}

func randGenesis(t *testing.T, assetType asset.Type) asset.Genesis {
	t.Helper()

	return asset.Genesis{
		FirstPrevOut: wire.OutPoint{},
		Tag:          "",
		Metadata:     []byte{},
		OutputIndex:  rand.Uint32(),
		Type:         assetType,
	}
}

func randFamilyKey(t *testing.T, genesis asset.Genesis) *asset.FamilyKey {
	t.Helper()
	privKey := randKey(t)
	genSigner := asset.NewRawKeyGenesisSigner(privKey)
	fakeKeyDesc := keychain.KeyDescriptor{
		PubKey: privKey.PubKey(),
	}
	familyKey, err := asset.DeriveFamilyKey(genSigner, fakeKeyDesc, genesis)
	require.NoError(t, err)

	return familyKey
}

func initSpendScenario(t *testing.T) spendData {
	t.Helper()

	// Amounts and genesises, needed for addresses and assets. We need both
	// a normal and collectible asset, and three amounts to test splits.
	state := spendData{
		collectAmt:      1,
		normalAmt1:      2,
		normalAmt2:      5,
		genesis1:        randGenesis(t, asset.Normal),
		genesis1collect: randGenesis(t, asset.Collectible),
	}

	// Keys for sender, receiver, and family. Default to keypath spend
	// for the spender ScriptKey.
	spenderPrivKey, spenderPubKey := btcec.PrivKeyFromBytes(key1Bytes)
	state.spenderPrivKey = *spenderPrivKey
	state.spenderPubKey = *spenderPubKey
	spenderScriptKey := *txscript.ComputeTaprootKeyNoScript(
		&state.spenderPubKey,
	)
	state.spenderScriptKey = spenderScriptKey
	state.spenderDescriptor = keychain.KeyDescriptor{
		PubKey: &state.spenderScriptKey,
	}
	receiverPrivKey, receiverPubKey := btcec.PrivKeyFromBytes(key2Bytes)
	state.receiverPrivKey = *receiverPrivKey
	state.receiverPubKey = *receiverPubKey
	familyKey := randFamilyKey(t, state.genesis1collect)
	state.familyKey = *familyKey

	// Addesses to cover both asset types and all three asset values.
	// Store the receiver StateKeys as well.
	address1, err := address.New(
		state.genesis1.ID(), nil, state.receiverPubKey,
		state.receiverPubKey, state.normalAmt1,
		asset.Normal, &address.MainNetTaro,
	)
	require.NoError(t, err)
	state.address1 = *address1
	state.address1StateKey = state.address1.AssetCommitmentKey()

	address1CollectFamily, err := address.New(
		state.genesis1collect.ID(), &state.familyKey.FamKey,
		state.receiverPubKey, state.receiverPubKey, state.collectAmt,
		asset.Collectible, &address.TestNet3Taro,
	)
	require.NoError(t, err)
	state.address1CollectFamily = *address1CollectFamily
	state.address1CollectFamilyStateKey = state.address1CollectFamily.
		AssetCommitmentKey()

	address2, err := address.New(
		state.genesis1.ID(), nil, state.receiverPubKey,
		state.receiverPubKey, state.normalAmt2,
		asset.Normal, &address.MainNetTaro,
	)
	require.NoError(t, err)
	state.address2 = *address2
	state.address2StateKey = state.address2.AssetCommitmentKey()

	// Generate matching assets and PrevIDs.
	updateScenarioAssets(t, &state)

	// Generate matching TaroCommitments.
	updateScenarioCommitments(t, &state)

	return state
}

func updateScenarioAssets(t *testing.T, state *spendData) {
	t.Helper()

	require.NotNil(t, state)

	locktime := uint64(1)
	relLocktime := uint64(1)

	// Assets to cover both asset types and all three asset values.
	asset1, err := asset.New(
		state.genesis1, state.normalAmt1, locktime,
		relLocktime, state.spenderDescriptor, nil,
	)
	require.NoError(t, err)
	state.asset1 = *asset1

	asset1CollectFamily, err := asset.New(
		state.genesis1collect, state.collectAmt, locktime,
		relLocktime, state.spenderDescriptor, &state.familyKey,
	)
	require.NoError(t, err)
	state.asset1CollectFamily = *asset1CollectFamily

	asset2, err := asset.New(
		state.genesis1, state.normalAmt2, locktime,
		relLocktime, state.spenderDescriptor, nil,
	)
	require.NoError(t, err)
	state.asset2 = *asset2

	// Asset PrevIDs, required to represent an input asset for a spend.
	state.asset1PrevID = asset.PrevID{
		OutPoint:  wire.OutPoint{},
		ID:        state.asset1.ID(),
		ScriptKey: state.spenderScriptKey,
	}
	state.asset1CollectFamilyPrevID = asset.PrevID{
		OutPoint:  wire.OutPoint{},
		ID:        state.asset1CollectFamily.ID(),
		ScriptKey: state.spenderScriptKey,
	}
	state.asset2PrevID = asset.PrevID{
		OutPoint:  wire.OutPoint{},
		ID:        state.asset2.ID(),
		ScriptKey: state.spenderScriptKey,
	}

	state.asset1InputAssets = commitment.InputSet{
		state.asset1PrevID: &state.asset1,
	}
	state.asset1CollectFamilyInputAssets = commitment.InputSet{
		state.asset1CollectFamilyPrevID: &state.asset1CollectFamily,
	}
	state.asset2InputAssets = commitment.InputSet{
		state.asset2PrevID: &state.asset2,
	}
}

func updateScenarioCommitments(t *testing.T, state *spendData) {
	t.Helper()

	require.NotNil(t, state)

	// TaroCommitments for each asset.
	asset1AssetTree, err := commitment.NewAssetCommitment(&state.asset1)
	require.NoError(t, err)
	asset1TaroTree, err := commitment.NewTaroCommitment(asset1AssetTree)
	require.NoError(t, err)
	state.asset1TaroTree = *asset1TaroTree

	asset1CollectFamilyAssetTree, err := commitment.NewAssetCommitment(
		&state.asset1CollectFamily,
	)
	require.NoError(t, err)
	asset1CollectFamilyTaroTree, err := commitment.NewTaroCommitment(
		asset1CollectFamilyAssetTree,
	)
	require.NoError(t, err)
	state.asset1CollectFamilyTaroTree = *asset1CollectFamilyTaroTree

	asset2AssetTree, err := commitment.NewAssetCommitment(&state.asset2)
	require.NoError(t, err)
	asset2TaroTree, err := commitment.NewTaroCommitment(asset2AssetTree)
	require.NoError(t, err)
	state.asset2TaroTree = *asset2TaroTree
	require.NoError(t, err)
}

func assertAssetEqual(t *testing.T, a, b *asset.Asset) {
	t.Helper()

	require.Equal(t, a.Version, b.Version)
	require.Equal(t, a.Genesis, b.Genesis)
	require.Equal(t, a.Type, b.Type)
	require.Equal(t, a.Amount, b.Amount)
	require.Equal(t, a.LockTime, b.LockTime)
	require.Equal(t, a.RelativeLockTime, b.RelativeLockTime)
	require.Equal(t, len(a.PrevWitnesses), len(b.PrevWitnesses))

	for i := range a.PrevWitnesses {
		witA, witB := a.PrevWitnesses[i], b.PrevWitnesses[i]
		require.Equal(t, witA.PrevID, witB.PrevID)
		require.Equal(t, witA.TxWitness, witB.TxWitness)
		splitA, splitB := witA.SplitCommitment, witB.SplitCommitment

		if witA.SplitCommitment != nil && witB.SplitCommitment != nil {
			require.Equal(
				t, len(splitA.Proof.Nodes),
				len(splitB.Proof.Nodes),
			)
			for i := range splitA.Proof.Nodes {
				nodeA := splitA.Proof.Nodes[i]
				nodeB := splitB.Proof.Nodes[i]
				require.True(t, mssmt.IsEqualNode(nodeA, nodeB))
			}
			require.Equal(t, splitA.RootAsset, splitB.RootAsset)
		} else {
			require.Equal(t, splitA, splitB)
		}
	}

	require.Equal(t, a.SplitCommitmentRoot, b.SplitCommitmentRoot)
	require.Equal(t, a.ScriptVersion, b.ScriptVersion)
	require.Equal(t, a.ScriptKey, b.ScriptKey)
	require.Equal(t, a.FamilyKey, b.FamilyKey)
}

func checkPreparedSplitSpend(t *testing.T, spend *SpendDelta, addr address.Taro,
	prevInput asset.PrevID, scriptKey btcec.PublicKey) {

	t.Helper()

	require.NotNil(t, spend.SplitCommitment)
	require.Equal(t, *spend.NewAsset.ScriptKey.PubKey, scriptKey)
	require.Equal(
		t, spend.NewAsset.Amount,
		spend.InputAssets[prevInput].Amount-addr.Amount,
	)

	receiverStateKey := addr.AssetCommitmentKey()
	receiverLocator, ok := spend.Locators[receiverStateKey]
	require.True(t, ok)
	receiverAsset, ok := spend.SplitCommitment.SplitAssets[receiverLocator]
	require.True(t, ok)
	require.Equal(t, receiverAsset.Asset.Amount, addr.Amount)
	require.Equal(t, *receiverAsset.Asset.ScriptKey.PubKey, addr.ScriptKey)
}

func checkPreparedCompleteSpend(t *testing.T, spend *SpendDelta,
	addr address.Taro, prevInput asset.PrevID) {

	t.Helper()

	require.Nil(t, spend.SplitCommitment)
	require.Equal(t, *spend.NewAsset.ScriptKey.PubKey, addr.ScriptKey)
	require.Equal(t, *spend.NewAsset.PrevWitnesses[0].PrevID, prevInput)
	require.Nil(t, spend.NewAsset.PrevWitnesses[0].TxWitness)
	require.Nil(t, spend.NewAsset.PrevWitnesses[0].SplitCommitment)
}

// TestPrepareAssetSplitSpend tests the creating of split commitment data with
// different sets of split locators. The validity of locators is assumed to be
// checked earlier via areValidIndexes().
func TestPrepareAssetSplitSpend(t *testing.T) {
	t.Parallel()

	prepareAssetSplitSpendTestCases := []struct {
		name string
		f    func() error
		err  error
	}{
		{
			name: "asset split with custom locators",
			f: func() error {
				state := initSpendScenario(t)
				spend := SpendDelta{
					InputAssets: state.asset2InputAssets,
				}

				spenderStateKey := asset.AssetCommitmentKey(
					state.asset2.ID(),
					&state.spenderScriptKey, true,
				)
				receiverStateKey := state.address1StateKey

				spend.Locators = make(SpendLocators)
				spend.Locators[spenderStateKey] = commitment.
					SplitLocator{OutputIndex: 0}
				spend.Locators[receiverStateKey] = commitment.
					SplitLocator{OutputIndex: 2}
				spendPrepared, err := prepareAssetSplitSpend(
					state.address1, state.asset2PrevID,
					state.spenderScriptKey, spend,
				)
				require.NoError(t, err)

				checkPreparedSplitSpend(
					t, spendPrepared, state.address1,
					state.asset2PrevID,
					state.spenderScriptKey,
				)
				return nil
			},
			err: nil,
		},
		{
			name: "asset split with mock locators",
			f: func() error {
				state := initSpendScenario(t)
				spend := SpendDelta{
					InputAssets: state.asset2InputAssets,
				}
				spendPrepared, err := prepareAssetSplitSpend(
					state.address1, state.asset2PrevID,
					state.spenderScriptKey, spend,
				)
				require.NoError(t, err)

				checkPreparedSplitSpend(
					t, spendPrepared, state.address1,
					state.asset2PrevID,
					state.spenderScriptKey,
				)
				return nil
			},
			err: nil,
		},
	}

	for _, testCase := range prepareAssetSplitSpendTestCases {
		success := t.Run(testCase.name, func(t *testing.T) {
			err := testCase.f()
			require.ErrorIs(t, err, testCase.err)
		})
		if !success {
			return
		}
	}
}

// TestPrepareAssetCompleteSpend tests the two cases where an asset is spent
// completely, asserting that new asset leaves are correctly created.
func TestPrepareAssetCompleteSpend(t *testing.T) {
	t.Parallel()

	prepareAssetCompleteSpendTestCases := []struct {
		name string
		f    func() error
		err  error
	}{
		{
			name: "collectible with family key",
			f: func() error {
				state := initSpendScenario(t)
				spend := SpendDelta{
					InputAssets: state.
						asset1CollectFamilyInputAssets,
				}
				spendPrepared := prepareAssetCompleteSpend(
					state.address1CollectFamily,
					state.asset1CollectFamilyPrevID, spend,
				)
				checkPreparedCompleteSpend(
					t, spendPrepared,
					state.address1CollectFamily,
					state.asset1CollectFamilyPrevID,
				)
				return nil
			},
			err: nil,
		},
		{
			name: "normal asset without split",
			f: func() error {
				state := initSpendScenario(t)
				spend := SpendDelta{
					InputAssets: state.asset1InputAssets,
				}
				spendPrepared := prepareAssetCompleteSpend(
					state.address1, state.asset1PrevID,
					spend,
				)
				checkPreparedCompleteSpend(
					t, spendPrepared, state.address1,
					state.asset1PrevID,
				)
				return nil
			},
			err: nil,
		},
	}

	for _, testCase := range prepareAssetCompleteSpendTestCases {
		success := t.Run(testCase.name, func(t *testing.T) {
			err := testCase.f()
			require.ErrorIs(t, err, testCase.err)
		})
		if !success {
			return
		}
	}
}

// TestValidIndexes tests various sets of asset locators to assert that we can
// detect an incomplete set of locators, and sets that form a valid Bitcoin
// transaction.
func TestValidIndexes(t *testing.T) {
	t.Parallel()

	state := initSpendScenario(t)

	spenderStateKey := asset.AssetCommitmentKey(
		state.asset1.ID(), &state.spenderScriptKey, true,
	)
	receiverStateKey := state.address1.AssetCommitmentKey()
	receiver2StateKey := state.address2.AssetCommitmentKey()

	locators := make(SpendLocators)

	// Insert a locator for the sender.
	locators[spenderStateKey] = commitment.SplitLocator{
		OutputIndex: 0,
	}

	// Reject groups of locators smaller than 2.
	taroOnlySpend, err := areValidIndexes(locators)
	require.False(t, taroOnlySpend)
	require.ErrorIs(t, err, ErrInvalidOutputIndexes)

	// Insert a locator for the receiver, that would form a Taro-only spend.
	locators[receiverStateKey] = commitment.SplitLocator{
		OutputIndex: 1,
	}

	taroOnlySpend, err = areValidIndexes(locators)
	require.True(t, taroOnlySpend)
	require.NoError(t, err)

	// Modify the receiver locator so the indexes are no longer continuous.
	locators[receiverStateKey] = commitment.SplitLocator{
		OutputIndex: 2,
	}

	taroOnlySpend, err = areValidIndexes(locators)
	require.False(t, taroOnlySpend)
	require.NoError(t, err)

	// Check for correctness with more than 2 locators.
	locators[receiver2StateKey] = commitment.SplitLocator{
		OutputIndex: 1,
	}

	taroOnlySpend, err = areValidIndexes(locators)
	require.True(t, taroOnlySpend)
	require.NoError(t, err)
}

// TestAddressValidInput tests edge cases around validating inputs for asset
// transfers with isValidInput.
func TestAddressValidInput(t *testing.T) {
	t.Parallel()

	state := initSpendScenario(t)

	address1testnet, err := address.New(
		state.genesis1.ID(), nil, state.receiverPubKey,
		state.receiverPubKey, state.normalAmt1, asset.Normal,
		&address.TestNet3Taro,
	)
	require.NoError(t, err)

	testCases := []struct {
		name string
		f    func() (*asset.Asset, *asset.Asset, error)
		err  error
	}{
		{
			name: "valid normal",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1TaroTree, state.address1,
					state.spenderScriptKey,
					address.MainNetTaro,
				)
				require.False(t, needsSplit)
				return &state.asset1, inputAsset, err
			},
			err: nil,
		},
		{
			name: "valid collectible with family key",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1CollectFamilyTaroTree,
					state.address1CollectFamily,
					state.spenderScriptKey,
					address.TestNet3Taro,
				)
				require.False(t, needsSplit)
				return &state.asset1CollectFamily,
					inputAsset, err
			},
			err: nil,
		},
		{
			name: "valid asset split",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset2TaroTree, state.address1,
					state.spenderScriptKey,
					address.MainNetTaro,
				)
				require.True(t, needsSplit)
				return &state.asset2, inputAsset, err
			},
			err: nil,
		},
		{
			name: "normal with insufficient amount",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1TaroTree, state.address2,
					state.spenderScriptKey,
					address.MainNetTaro,
				)
				require.False(t, needsSplit)
				return &state.asset1, inputAsset, err
			},
			err: ErrInsufficientInputAsset,
		},
		{
			name: "collectible with missing input asset",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1TaroTree,
					state.address1CollectFamily,
					state.spenderScriptKey,
					address.TestNet3Taro,
				)
				require.False(t, needsSplit)
				return &state.asset1, inputAsset, err
			},
			err: ErrMissingInputAsset,
		},
		{
			name: "normal with bad sender script key",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1TaroTree,
					*address1testnet,
					state.receiverPubKey,
					address.TestNet3Taro,
				)
				require.False(t, needsSplit)
				return &state.asset1, inputAsset, err
			},
			err: ErrMissingInputAsset,
		},
		{
			name: "normal with mismatched network",
			f: func() (*asset.Asset, *asset.Asset, error) {
				inputAsset, needsSplit, err := isValidInput(
					state.asset1TaroTree,
					*address1testnet,
					state.receiverPubKey,
					address.MainNetTaro,
				)
				require.False(t, needsSplit)
				return &state.asset1, inputAsset, err
			},
			err: address.ErrMismatchedHRP,
		},
	}

	for _, testCase := range testCases {
		success := t.Run(testCase.name, func(t *testing.T) {
			inputAsset, checkedInputAsset, err := testCase.f()
			require.ErrorIs(t, err, testCase.err)
			if testCase.err == nil {
				assertAssetEqual(
					t, inputAsset, checkedInputAsset,
				)
			}
		})
		if !success {
			return
		}
	}
}
