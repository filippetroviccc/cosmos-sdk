package utils

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
)

// QueryAllConnections returns all the connections. It _does not_ return
// any merkle proof.
func QueryAllConnections(cliCtx context.CLIContext, page, limit int) ([]types.ConnectionEnd, int64, error) {
	params := types.NewQueryAllConnectionsParams(page, limit)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal query params: %w", err)
	}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllConnections)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, 0, err
	}

	var connections []types.ConnectionEnd
	err = cliCtx.Codec.UnmarshalJSON(res, &connections)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal connections: %w", err)
	}
	return connections, height, nil
}

// QueryConnection queries the store to get a connection end and a merkle
// proof.
func QueryConnection(
	cliCtx context.CLIContext, connectionID string, prove bool,
) (types.ConnectionResponse, error) {
	req := abci.RequestQuery{
		Path:  "store/ibc/key",
		Data:  types.KeyConnection(connectionID),
		Prove: prove,
	}

	res, err := cliCtx.QueryABCI(req)
	if err != nil {
		return types.ConnectionResponse{}, err
	}

	var connection types.ConnectionEnd
	if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res.Value, &connection); err != nil {
		return types.ConnectionResponse{}, err
	}

	connRes := types.NewConnectionResponse(connectionID, connection, res.Proof, res.Height)

	return connRes, nil
}

// QueryClientConnections queries the store to get the registered connection paths
// registered for a particular client and a merkle proof.
func QueryClientConnections(
	cliCtx context.CLIContext, clientID string, prove bool,
) (types.ClientConnectionsResponse, error) {
	req := abci.RequestQuery{
		Path:  "store/ibc/key",
		Data:  types.KeyClientConnections(clientID),
		Prove: prove,
	}

	res, err := cliCtx.QueryABCI(req)
	if err != nil {
		return types.ClientConnectionsResponse{}, err
	}

	var paths []string
	if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res.Value, &paths); err != nil {
		return types.ClientConnectionsResponse{}, err
	}

	connPathsRes := types.NewClientConnectionsResponse(clientID, paths, res.Proof, res.Height)
	return connPathsRes, nil
}
