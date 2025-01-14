package keeper_test

import "github.com/furya-official/blackfury/x/ve/types"

func (suite *KeeperTestSuite) TestKeeper_GetParams() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	params := k.GetParams(suite.ctx)
	suite.Require().Equal("afury", params.LockDenom)
}

func (suite *KeeperTestSuite) TestKeeper_SetParams() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	k.SetParams(suite.ctx, types.Params{LockDenom: "aaa"})
	params := k.GetParams(suite.ctx)
	suite.Require().Equal("aaa", params.LockDenom)
}

func (suite *KeeperTestSuite) TestKeeper_LockDenom() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	res := k.LockDenom(suite.ctx)
	suite.Require().Equal("afury", res)
}
