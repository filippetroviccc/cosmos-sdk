package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

// NewQuerier creates a querier for the IBC module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		case client.SubModuleName:
			switch path[1] {
			case client.QueryClientState:
				res, err = client.QuerierClientState(ctx, req, k.ClientKeeper)
			case client.QueryAllClients:
				res, err = client.QuerierClients(ctx, req, k.ClientKeeper)
			case client.QueryConsensusState:
				res, err = client.QuerierConsensusState(ctx, req, k.ClientKeeper)
			case client.QueryVerifiedRoot:
				res, err = client.QuerierVerifiedRoot(ctx, req, k.ClientKeeper)
			default:
				err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown IBC %s query endpoint", client.SubModuleName)
			}
		case connection.SubModuleName:
			switch path[1] {
			case connection.QueryConnection:
				res, err = connection.QuerierConnection(ctx, req, k.ConnectionKeeper)
			case connection.QueryAllConnections:
				res, err = connection.QuerierConnections(ctx, req, k.ConnectionKeeper)
			case connection.QueryClientConnections:
				res, err = connection.QuerierClientConnections(ctx, req, k.ConnectionKeeper)
			default:
				err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown IBC %s query endpoint", connection.SubModuleName)
			}
		case channel.SubModuleName:
			switch path[1] {
			case channel.QueryChannel:
				res, err = channel.QuerierChannel(ctx, req, k.ChannelKeeper)
			case channel.QueryAllChannels:
				res, err = channel.QuerierChannels(ctx, req, k.ChannelKeeper)
			default:
				err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown IBC %s query endpoint", channel.SubModuleName)
			}
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown IBC query endpoint")
		}

		return res, sdk.ConvertError(err)
	}
}
