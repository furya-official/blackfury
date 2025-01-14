package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	mgravity "github.com/Gravity-Bridge/Gravity-Bridge/module/x/multigravity"
	mgravitykeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/multigravity/keeper"
	mgravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/multigravity/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	nft "github.com/cosmos/cosmos-sdk/x/nft/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/osmosis-labs/bech32-ibc/x/bech32ibc"
	bech32ibckeeper "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/keeper"
	bech32ibctypes "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/types"
	statikfs "github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	"github.com/tendermint/starport/starport/pkg/openapiconsole"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"
	srvflags "github.com/tharsis/ethermint/server/flags"
	ethermint "github.com/tharsis/ethermint/types"
	"github.com/tharsis/ethermint/x/evm"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/ethermint/x/feemarket"
	feemarketkeeper "github.com/tharsis/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	"github.com/tharsis/evmos/v4/app/ante"

	_ "github.com/furya-official/blackfury/client/docs/statik"
	custombank "github.com/furya-official/blackfury/x/bank"
	custombankclient "github.com/furya-official/blackfury/x/bank/client"
	custombankkeeper "github.com/furya-official/blackfury/x/bank/keeper"
	custombanktypes "github.com/furya-official/blackfury/x/bank/types"
	"github.com/furya-official/blackfury/x/erc20"
	erc20keeper "github.com/furya-official/blackfury/x/erc20/keeper"
	erc20types "github.com/furya-official/blackfury/x/erc20/types"
	"github.com/furya-official/blackfury/x/gauge"
	gaugekeeper "github.com/furya-official/blackfury/x/gauge/keeper"
	gaugetypes "github.com/furya-official/blackfury/x/gauge/types"
	"github.com/furya-official/blackfury/x/maker"
	makerclient "github.com/furya-official/blackfury/x/maker/client"
	makerkeeper "github.com/furya-official/blackfury/x/maker/keeper"
	makertypes "github.com/furya-official/blackfury/x/maker/types"
	"github.com/furya-official/blackfury/x/oracle"
	oracleclient "github.com/furya-official/blackfury/x/oracle/client"
	oraclekeeper "github.com/furya-official/blackfury/x/oracle/keeper"
	oracletypes "github.com/furya-official/blackfury/x/oracle/types"
	customstaking "github.com/furya-official/blackfury/x/staking"
	customstakingkeeper "github.com/furya-official/blackfury/x/staking/keeper"
	"github.com/furya-official/blackfury/x/ve"
	vekeeper "github.com/furya-official/blackfury/x/ve/keeper"
	vetypes "github.com/furya-official/blackfury/x/ve/types"
	customvesting "github.com/furya-official/blackfury/x/vesting"
	customvestingkeeper "github.com/furya-official/blackfury/x/vesting/keeper"
	customvestingtypes "github.com/furya-official/blackfury/x/vesting/types"
	"github.com/furya-official/blackfury/x/voter"
	voterkeeper "github.com/furya-official/blackfury/x/voter/keeper"
	votertypes "github.com/furya-official/blackfury/x/voter/types"
)

// App represents a Cosmos SDK application that can be run as a server and with an exportable state
type App interface {
	servertypes.Application
	ExportableApp
}

// ExportableApp represents an app with an exportable state
type ExportableApp interface {
	ExportAppStateAndValidators(
		forZeroHeight bool,
		jailAllowedAddrs []string,
	) (servertypes.ExportedApp, error)
	LoadHeight(height int64) error
}

const (
	Name = "blackfury"
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		custombankclient.SetDenomMetaDataProposalHandler,
		makerclient.RegisterBackingProposalHandler,
		makerclient.RegisterCollateralProposalHandler,
		makerclient.SetBackingProposalHandler,
		makerclient.SetCollateralProposalHandler,
		makerclient.BatchSetBackingProposalHandler,
		makerclient.BatchSetCollateralProposalHandler,
		oracleclient.RegisterTargetProposalHandler,
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		customstaking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()...),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		erc20.AppModuleBasic{},
		oracle.AppModuleBasic{},
		maker.AppModuleBasic{},
		nft.AppModuleBasic{},
		ve.AppModuleBasic{},
		gauge.AppModuleBasic{},
		voter.AppModuleBasic{},
		customvesting.AppModuleBasic{},
		mgravity.AppModuleBasic{},
		bech32ibc.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
		oracletypes.ModuleName:         nil,
		makertypes.ModuleName:          {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		nfttypes.ModuleName:            nil,
		vetypes.ModuleName:             {authtypes.Burner},
		vetypes.EmissionPoolName:       {authtypes.Minter},
		vetypes.DistributionPoolName:   nil,
		gaugetypes.ModuleName:          nil,
		votertypes.ModuleName:          nil,
		customvestingtypes.ModuleName:  {authtypes.Minter},
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		mgravitytypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		oracletypes.ModuleName: true,
	}
)

var (
	_ App                     = (*Blackfury)(nil)
	_ servertypes.Application = (*Blackfury)(nil)
	_ simapp.App              = (*Blackfury)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)

	// Manually update the power reduction by replacing micro (u) -> fur
	sdk.DefaultPowerReduction = ethermint.PowerReduction
}

// Blackfury extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type Blackfury struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    customstakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	Erc20Keeper erc20keeper.Keeper

	OracleKeeper oraclekeeper.Keeper
	MakerKeeper  makerkeeper.Keeper

	NftKeeper   vekeeper.NftKeeper
	VeKeeper    vekeeper.Keeper
	GaugeKeeper gaugekeeper.Keeper
	VoterKeeper voterkeeper.Keeper

	VestingKeeper customvestingkeeper.Keeper

	GravityKeeper   mgravitykeeper.Keeper
	Bech32IbcKeeper bech32ibckeeper.Keeper

	// mm is the module manager
	mm *module.Manager

	// sm is the simulation manager
	sm *module.SimulationManager

	tpsCounter *tpsCounter
}

// NewBlackfury returns a reference to an initialized blockchain app
func NewBlackfury(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) App {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		erc20types.StoreKey,
		oracletypes.StoreKey,
		makertypes.StoreKey,
		nfttypes.StoreKey,
		vetypes.StoreKey,
		gaugetypes.StoreKey,
		votertypes.StoreKey,
		customvestingtypes.StoreKey,
		mgravitytypes.StoreKey,
		bech32ibctypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &Blackfury{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), ethermint.ProtoAccount, maccPerms,
	)
	var erc20Keeper erc20keeper.Keeper
	getErc20Keeper := func() custombanktypes.Erc20Keeper {
		return erc20Keeper
	}
	app.BankKeeper = custombankkeeper.NewKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(), getErc20Keeper,
	)

	var veKeeper vekeeper.Keeper
	getVeKeeper := func() vekeeper.Keeper {
		return veKeeper
	}
	nftKeeper := nftkeeper.NewKeeper(keys[nfttypes.StoreKey], appCodec, app.AccountKeeper, app.BankKeeper)
	app.NftKeeper = vekeeper.NewNftKeeper(nftKeeper, getVeKeeper)

	app.VeKeeper = *vekeeper.NewKeeper(appCodec, keys[vetypes.StoreKey], keys[vetypes.MemStoreKey], app.GetSubspace(vetypes.ModuleName), app.AccountKeeper, app.BankKeeper, app.NftKeeper)
	veKeeper = app.VeKeeper

	stakingKeeper := customstakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.GetSubspace(stakingtypes.ModuleName), app.AccountKeeper, app.BankKeeper, app.NftKeeper, app.VeKeeper,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.BlockedAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
			app.GravityKeeper.Hooks(),
		),
	)

	// ... other modules keepers

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Create ethermint keepers
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, keys[feemarkettypes.StoreKey], app.GetSubspace(feemarkettypes.ModuleName),
	)
	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.FeeMarketKeeper,
		tracer,
	)

	app.Erc20Keeper = *erc20keeper.NewKeeper(
		appCodec,
		keys[erc20types.StoreKey],
		app.GetSubspace(erc20types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.EvmKeeper,
	)
	erc20Keeper = app.Erc20Keeper
	erc20Module := erc20.NewAppModule(appCodec, app.Erc20Keeper, app.AccountKeeper, app.BankKeeper)

	app.OracleKeeper = *oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		app.GetSubspace(oracletypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.StakingKeeper,
		distrtypes.ModuleName,
	)
	oracleModule := oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper)

	app.MakerKeeper = *makerkeeper.NewKeeper(
		appCodec,
		keys[makertypes.StoreKey],
		keys[makertypes.MemStoreKey],
		app.GetSubspace(makertypes.ModuleName),

		app.AccountKeeper,
		app.BankKeeper,
		app.OracleKeeper,
	)
	makerModule := maker.NewAppModule(appCodec, app.MakerKeeper, app.AccountKeeper, app.BankKeeper)

	nftModule := vekeeper.NewNftAppModule(nft.NewAppModule(appCodec, nftKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry), app.NftKeeper)

	veModule := ve.NewAppModule(appCodec, app.VeKeeper, app.AccountKeeper, app.BankKeeper)

	app.GaugeKeeper = *gaugekeeper.NewKeeper(appCodec, keys[gaugetypes.StoreKey], keys[gaugetypes.MemStoreKey],
		app.GetSubspace(gaugetypes.ModuleName), app.AccountKeeper, app.BankKeeper, app.NftKeeper, app.VeKeeper)
	gaugeModule := gauge.NewAppModule(appCodec, app.GaugeKeeper, app.AccountKeeper, app.BankKeeper)

	app.VoterKeeper = *voterkeeper.NewKeeper(appCodec, keys[votertypes.StoreKey], keys[votertypes.MemStoreKey],
		app.GetSubspace(votertypes.ModuleName), app.AccountKeeper, app.BankKeeper, app.VeKeeper, app.GaugeKeeper)
	voterModule := voter.NewAppModule(appCodec, app.VoterKeeper, app.AccountKeeper, app.BankKeeper)

	app.VestingKeeper = *customvestingkeeper.NewKeeper(appCodec, keys[customvestingtypes.StoreKey], app.GetSubspace(customvestingtypes.ModuleName), app.AccountKeeper, app.BankKeeper, app.DistrKeeper, app.VeKeeper, authtypes.FeeCollectorName)
	vestingModule := customvesting.NewAppModule(appCodec, app.VestingKeeper, app.AccountKeeper, app.BankKeeper)

	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			app.Erc20Keeper.EvmHooks(),
		),
	)

	app.Bech32IbcKeeper = *bech32ibckeeper.NewKeeper(
		app.IBCKeeper.ChannelKeeper, appCodec, keys[bech32ibctypes.StoreKey],
		app.TransferKeeper,
	)
	bech32ibcModule := bech32ibc.NewAppModule(
		appCodec,
		app.Bech32IbcKeeper,
	)

	app.GravityKeeper = mgravitykeeper.NewKeeper(
		app.BaseApp,
		keys[mgravitytypes.StoreKey],
		app.GetSubspace(mgravitytypes.ModuleName),
		appCodec,
		app.BankKeeper,
		&stakingKeeper,
		&app.SlashingKeeper,
		&app.DistrKeeper,
		&app.AccountKeeper,
		&app.TransferKeeper,
		&app.Bech32IbcKeeper,
		&app.ParamsKeeper,
	)
	gravityModule := mgravity.NewAppModule(app.GravityKeeper, app.BankKeeper)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(makertypes.RouterKey, maker.NewMakerProposalHandler(app.MakerKeeper)).
		AddRoute(oracletypes.RouterKey, oracle.NewOracleProposalHandler(app.OracleKeeper)).
		AddRoute(banktypes.RouterKey, custombank.NewBankProposalHandler(app.BankKeeper)).
		AddRoute(mgravitytypes.RouterKey, mgravitykeeper.NewGravityProposalHandler(app.GravityKeeper)).
		AddRoute(bech32ibctypes.RouterKey, bech32ibc.NewBech32IBCProposalHandler(app.Bech32IbcKeeper))

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		custombank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		customstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20Module,
		oracleModule,
		makerModule,
		nftModule,
		veModule,
		gaugeModule,
		voterModule,
		vestingModule,
		gravityModule,
		bech32ibcModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		// no-op modules
		genutiltypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		feegrant.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		paramstypes.ModuleName,
		ibctransfertypes.ModuleName,
		erc20types.ModuleName,
		oracletypes.ModuleName,
		makertypes.ModuleName,
		nfttypes.ModuleName,
		vetypes.ModuleName,
		gaugetypes.ModuleName,
		votertypes.ModuleName,
		customvestingtypes.ModuleName,
		mgravitytypes.ModuleName,
		bech32ibctypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		oracletypes.ModuleName,
		makertypes.ModuleName,
		// no-op modules
		genutiltypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		capabilitytypes.ModuleName,
		feegrant.ModuleName,
		slashingtypes.ModuleName,
		distrtypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		ibchost.ModuleName,
		paramstypes.ModuleName,
		ibctransfertypes.ModuleName,
		erc20types.ModuleName,
		nfttypes.ModuleName,
		vetypes.ModuleName,
		gaugetypes.ModuleName,
		votertypes.ModuleName,
		customvestingtypes.ModuleName,
		mgravitytypes.ModuleName,
		bech32ibctypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		nfttypes.ModuleName,
		vetypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		oracletypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		erc20types.ModuleName,
		makertypes.ModuleName,
		gaugetypes.ModuleName,
		votertypes.ModuleName,
		customvestingtypes.ModuleName,
		mgravitytypes.ModuleName,
		bech32ibctypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// create the simulation manager and define the order of the modules for deterministic simulations
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		customstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20Module,
		oracleModule,
		makerModule,
		nftModule,
		veModule,
		gaugeModule,
		voterModule,
		vestingModule,
	)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		StakingKeeper:   app.StakingKeeper,
		FeegrantKeeper:  app.FeeGrantKeeper,
		IBCKeeper:       app.IBCKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  SigVerificationGasConsumer,
		Cdc:             appCodec,
		MaxTxGasWanted:  maxGasWanted,
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(ante.NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic,
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background())
	}()

	mgravitykeeper.RegisterProposalTypes()

	return app
}

// Name returns the name of the Blackfury
func (app *Blackfury) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the base app of the application
func (app Blackfury) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *Blackfury) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *Blackfury) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *Blackfury) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		// TODO: Record the count along with the code and or reason so as to display
		// in the transactions per second live dashboards.
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()

	return app.BaseApp.DeliverTx(req)
}

// InitChainer application update at chain initialization
func (app *Blackfury) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *Blackfury) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Blackfury) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive tokens.
func (app *Blackfury) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Blackfury) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Blackfury) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *Blackfury) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Blackfury) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Blackfury) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Blackfury) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Blackfury) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *Blackfury) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	statikFS, err := statikfs.New()
	if err != nil {
		panic(err)
	}
	apiSvr.Router.Handle("/static/openapi.yml", http.FileServer(statikFS))
	apiSvr.Router.HandleFunc("/", openapiconsole.Handler(Name, "/static/openapi.yml"))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *Blackfury) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Blackfury) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(erc20types.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)
	paramsKeeper.Subspace(makertypes.ModuleName)
	paramsKeeper.Subspace(nfttypes.ModuleName)
	paramsKeeper.Subspace(vetypes.ModuleName)
	paramsKeeper.Subspace(gaugetypes.ModuleName)
	paramsKeeper.Subspace(votertypes.ModuleName)
	paramsKeeper.Subspace(customvestingtypes.ModuleName)
	paramsKeeper.Subspace(mgravitytypes.ModuleName)

	return paramsKeeper
}

// SimulationManager implements the SimulationApp interface
func (app *Blackfury) SimulationManager() *module.SimulationManager {
	return app.sm
}
