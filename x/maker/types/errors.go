package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/maker module sentinel errors
var (
	ErrMerPriceTooLow  = sdkerrors.Register(ModuleName, 2, "black stablecoin price too low")
	ErrMerPriceTooHigh = sdkerrors.Register(ModuleName, 3, "black stablecoin price too high")

	ErrBackingParamsInvalid    = sdkerrors.Register(ModuleName, 4, "backing params invalid")
	ErrCollateralParamsInvalid = sdkerrors.Register(ModuleName, 5, "collateral params invalid")

	ErrBackingCoinDisabled    = sdkerrors.Register(ModuleName, 6, "backing coin disabled")
	ErrCollateralCoinDisabled = sdkerrors.Register(ModuleName, 7, "collateral coin disabled")

	ErrBackingCoinAlreadyExists    = sdkerrors.Register(ModuleName, 8, "backing coin already exists")
	ErrCollateralCoinAlreadyExists = sdkerrors.Register(ModuleName, 9, "collateral coin already exists")
	ErrBackingCoinNotFound         = sdkerrors.Register(ModuleName, 10, "backing coin not found")
	ErrCollateralCoinNotFound      = sdkerrors.Register(ModuleName, 11, "collateral coin not found")

	ErrMerSlippage         = sdkerrors.Register(ModuleName, 12, "black over slippage")
	ErrBackingCoinSlippage = sdkerrors.Register(ModuleName, 13, "backing coin over slippage")
	ErrFuryCoinSlippage    = sdkerrors.Register(ModuleName, 14, "fury coin over slippage")

	ErrBackingCeiling    = sdkerrors.Register(ModuleName, 15, "total backing coin over ceiling")
	ErrCollateralCeiling = sdkerrors.Register(ModuleName, 16, "total collateral coin over ceiling")
	ErrMerCeiling        = sdkerrors.Register(ModuleName, 17, "total black coin over ceiling")

	ErrBackingCoinInsufficient    = sdkerrors.Register(ModuleName, 18, "backing coin balance insufficient")
	ErrCollateralCoinInsufficient = sdkerrors.Register(ModuleName, 19, "collateral coin balance insufficient")
	ErrFuryCoinInsufficient       = sdkerrors.Register(ModuleName, 20, "insufficient available fury coin")

	ErrAccountNoCollateral           = sdkerrors.Register(ModuleName, 21, "account has no collateral")
	ErrAccountInsufficientCollateral = sdkerrors.Register(ModuleName, 22, "account collateral insufficient")
	ErrAccountNoDebt                 = sdkerrors.Register(ModuleName, 23, "account has no debt")
	ErrNotUndercollateralized        = sdkerrors.Register(ModuleName, 24, "position is not undercollateralized")

	ErrLTVOutOfRange = sdkerrors.Register(ModuleName, 25, "LTV is out of range")
	ErrOverSlippage  = sdkerrors.Register(ModuleName, 26, "over slippage")
)
