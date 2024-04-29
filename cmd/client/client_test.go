// Copyright (c) 2021 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-eth-demo. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package client_test

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Cmd struct {
	*exec.Cmd
	stdin io.WriteCloser
}

func nodeCmd(name string) (*Cmd, error) {
	args := []string{
		"run", "../../main.go",
		"demo",
		"--config", fmt.Sprintf("../../config/%v.yaml", name),
		"--network", "../../config/network.yaml",
		"--log-level", "trace",
		"--log-file", fmt.Sprintf("../../logs/%v.log", name),
		"--stdio",
	}
	cmd := exec.Command("go", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return &Cmd{cmd, stdin}, nil
}

const (
	blockTime       = 10 * time.Second
	txFinalityDepth = 3
	numUpdates      = 25
	ethUrl          = "ws://127.0.0.1:8545"
	addresspeer_0    = "0x2EE1ac154435f542ECEc55C5b0367650d8A5343B"
	addresspeer_1      = "0x70765701b79a4e973dAbb4b30A72f5a845f22F9E"
)

func TestNodes(t *testing.T) {
	// Start peer_0.
	peer_0, err := nodeCmd("peer_0")
	require.NoError(t, err)
	require.NoError(t, peer_0.Start())
	defer peer_0.Process.Kill()
	time.Sleep(blockTime*2 + txFinalityDepth*blockTime) // Wait 2 blocks for contract deployment.

	// Start peer_1.
	peer_1, err := nodeCmd("peer_1")
	require.NoError(t, err)
	require.NoError(t, peer_1.Start())
	defer peer_1.Process.Kill()
	time.Sleep(10 * time.Second) // Give peer_1 some time to initialize.

	// Get the initial on-chain balances from peer_0 and peer_1.
	initBals, err := getOnChainBals()
	require.NoError(t, err)
	t.Logf("Initial on-chain balances: peer_0 = %f, peer_1 = %f", initBals[0], initBals[1])

	// peer_0 opens channel with peer_1.
	require.NoError(t, peer_0.sendCommand("open peer_1 500 500\n"), "proposing channel")
	time.Sleep(10 * time.Second) // Ensure that peer_1 really received the proposal.
	// require.NoError(t, peer_1.sendCommand("y\n"), "accepting channel proposal")
	t.Log("Opening channel…")
	time.Sleep(blockTime + txFinalityDepth*blockTime) // Wait 1 block for funding transactions to be confirmed.

	// peer_0 sends to peer_1 and peer_1 to peer_0.

	t.Log("Server: Init FL")
	require.NoError(t, peer_0.sendCommand("set peer_1 qmqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 1 0 0 0\n"))
	time.Sleep(1 * time.Second)
	t.Log("Client: Set Weight")
	require.NoError(t, peer_1.sendCommand("set peer_0 1mqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 1 1mqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 0 0\n"))
	time.Sleep(1 * time.Second)
	t.Log("Server: Aggregate and Evaluate")
	require.NoError(t, peer_0.sendCommand("set peer_1 amqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 1 0 66 34\n"))
	time.Sleep(1 * time.Second)
	// round 2
	t.Log("Client: Set Weight")
	require.NoError(t, peer_1.sendCommand("set peer_0 amqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 3 2mQFKFA74GJ4t1xwQFx25HbpT7ZmuyE3MBpoxRu1hZyuaS 66 34\n"))
	time.Sleep(1 * time.Second)
	t.Log("Server: Aggregate and Evaluate")
	require.NoError(t, peer_0.sendCommand("set peer_1 bmqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 3 2mQFKFA74GJ4t1xwQFx25HbpT7ZmuyE3MBpoxRu1hZyuaS 67 33\n"))
	time.Sleep(1 * time.Second)
	// round 3
	t.Log("Client: Set Weight")
	require.NoError(t, peer_1.sendCommand("set peer_0 bmqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 3 3mQFKFA74GJ4t1xwQFx25HbpT7ZmuyE3MBpoxRu1hZyuaS 67 33\n"))
	time.Sleep(1 * time.Second)
	t.Log("Server: Aggregate and Evaluate")
	require.NoError(t, peer_0.sendCommand("set peer_1 cmqfkfa74gj4t1xwqfx25hbpt7zmuye3mbpoxru1hzyuas 3 3mQFKFA74GJ4t1xwQFx25HbpT7ZmuyE3MBpoxRu1hZyuaS 68 32\n"))
	time.Sleep(1 * time.Second)
	t.Log("Info")
	t.Log(peer_0.sendCommand("info\n"))

	t.Log("Closing channel…")
	require.NoError(t, peer_0.sendCommand("settle peer_1\n"))
	// Wait 2 blocks for the settle and withdrawal transactions plus some additional seconds.
	time.Sleep(2*blockTime + 10*time.Second + txFinalityDepth*blockTime)
	require.NoError(t, peer_1.sendCommand("settle peer_0\n"))

	// Wait 2 blocks for the settle and withdrawal transactions plus some additional seconds.
	time.Sleep(2*blockTime + 10*time.Second + txFinalityDepth*blockTime)

	// Get the final balances from peer_0 and peer_1 after the settlement.
	finalBals, err := getOnChainBals()
	require.NoError(t, err)

	t.Logf("Final on-chain balances: peer_0 = %f, peer_1 = %f", finalBals[0], finalBals[1])

	var diffBals [2]float64

	// Calculate the differences between the initial and final balances.
	diffBals[0], _ = finalBals[0].Sub(finalBals[0], initBals[0]).Float64()
	diffBals[1], _ = finalBals[1].Sub(finalBals[1], initBals[1]).Float64()

	// Check the on-chain balance differences while allowing a higher deviation for peer_0.
	assert.InEpsilon(t, -3, diffBals[0], 1)
	assert.InEpsilon(t, 3, diffBals[1], 1)

	t.Log("Done")
}

func (cmd *Cmd) sendCommand(str string) error {
	_, err := io.WriteString(cmd.stdin, str)
	return err
}

func getOnChainBals() ([2]*big.Float, error) {
	var bals [2]*big.Float

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ethereumBackend, err := ethclient.Dial(ethUrl)
	if err != nil {
		return bals, errors.New("Could not connect to ethereum node.")
	}

	for idx, adr := range [2]string{addresspeer_0, addresspeer_1} {
		wei, err := ethereumBackend.BalanceAt(ctx,
			common.HexToAddress(adr), nil)
		if err != nil {
			return bals, err
		}
		bals[idx] = new(big.Float).Quo(new(big.Float).SetInt(wei),
			new(big.Float).SetFloat64(params.Ether))
	}

	return bals, nil
}
