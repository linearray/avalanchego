// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snow

import (
	"testing"

	"github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/proto/pb/p2p"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/prometheus/client_golang/prometheus"
)

func DefaultContextTest() *Context {
	sk, err := bls.NewSecretKey()
	if err != nil {
		panic(err)
	}
	pk := bls.PublicFromSecretKey(sk)
	return &Context{
		NetworkID:    0,
		SubnetID:     ids.Empty,
		ChainID:      ids.Empty,
		NodeID:       ids.EmptyNodeID,
		PublicKey:    pk,
		Log:          logging.NoLog{},
		BCLookup:     ids.NewAliaser(),
		Metrics:      metrics.NewOptionalGatherer(),
		ChainDataDir: "",
	}
}

func DefaultConsensusContextTest(t *testing.T) *ConsensusContext {
	var (
		startedState State = Initializing
		stoppedState State = Initializing
		engineType         = p2p.EngineType_ENGINE_TYPE_UNSPECIFIED
	)

	return &ConsensusContext{
		Context:             DefaultContextTest(),
		Registerer:          prometheus.NewRegistry(),
		AvalancheRegisterer: prometheus.NewRegistry(),
		BlockAcceptor:       noOpAcceptor{},
		TxAcceptor:          noOpAcceptor{},
		VertexAcceptor:      noOpAcceptor{},
		SubnetStateTracker: &SubnetStateTrackerTest{
			T: t,
			IsSyncedF: func() bool {
				return stoppedState == Bootstrapping || stoppedState == StateSyncing
			},
			IsChainBootstrappedF: func(ids.ID) bool {
				return stoppedState == Bootstrapping || stoppedState == StateSyncing
			},
			StartStateF: func(chainID ids.ID, state State, currentEngineType p2p.EngineType) {
				startedState = state
				engineType = currentEngineType
			},
			StopStateF: func(chainID ids.ID, state State) {
				stoppedState = state
			},
			GetStateF: func(chainID ids.ID) (State, p2p.EngineType) {
				return startedState, engineType
			},
		},
	}
}