package keeper_test

import (
	"github.com/atomone-hub/atomone/testutil/testdata"
	sdk "github.com/atomone-hub/atomone/types"
	stakingtypes "github.com/atomone-hub/atomone/x/staking/types"
)

func (s *KeeperTestSuite) TestAfterValidatorBonded() {
	ctx, keeper := s.ctx, s.slashingKeeper
	require := s.Require()

	valAddr := sdk.ValAddress(consAddr.Bytes())
	keeper.Hooks().AfterValidatorBonded(ctx, consAddr, valAddr)

	_, ok := keeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(ok)
}

func (s *KeeperTestSuite) TestAfterValidatorCreatedOrRemoved() {
	ctx, keeper := s.ctx, s.slashingKeeper
	require := s.Require()

	_, pubKey, addr := testdata.KeyTestPubAddr()
	valAddr := sdk.ValAddress(addr)

	validator, err := stakingtypes.NewValidator(sdk.ValAddress(addr), pubKey, stakingtypes.Description{})
	require.NoError(err)

	s.stakingKeeper.EXPECT().Validator(ctx, valAddr).Return(validator)
	err = keeper.Hooks().AfterValidatorCreated(ctx, valAddr)
	require.NoError(err)

	ePubKey, err := keeper.GetPubkey(ctx, addr.Bytes())
	require.NoError(err)
	require.Equal(ePubKey, pubKey)

	err = keeper.Hooks().AfterValidatorRemoved(ctx, sdk.ConsAddress(addr), nil)
	require.NoError(err)

	_, err = keeper.GetPubkey(ctx, addr.Bytes())
	require.Error(err)
}
