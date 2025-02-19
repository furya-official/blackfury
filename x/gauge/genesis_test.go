package gauge_test

import (
	"testing"

	keepertest "github.com/furya-official/blackfury/testutil/keeper"
	"github.com/furya-official/blackfury/testutil/nullify"
	"github.com/furya-official/blackfury/x/gauge"
	"github.com/furya-official/blackfury/x/gauge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.GaugeKeeper(t)
	gauge.InitGenesis(ctx, *k, genesisState)
	got := gauge.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
