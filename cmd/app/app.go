// Copyright 2021 PolyCrypt GmbH, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/pkg/errors"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

// FLApp is a channel app.
type FLApp struct {
	Addr wallet.Address
}

func NewFLApp(addr wallet.Address) *FLApp {
	return &FLApp{
		Addr: addr,
	}
}

// Def returns the app address.
func (a *FLApp) Def() wallet.Address {
	return a.Addr
}

func (a *FLApp) InitData(firstActor channel.Index) *FLAppData {
	return &FLAppData{
		NextActor: uint8(firstActor),
		// Model: make([]byte, 64),
	}
}

// DecodeData decodes the channel data.
func (a *FLApp) DecodeData(r io.Reader) (channel.Data, error) {
	d := FLAppData{}

	var err error
	d.NextActor, err = readUInt8(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading actor")
	}

	d.CIDLen, err = readUInt8(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading CIDLen")
	}

	d.NumberOfRounds, err = readUInt8(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading numberOfRounds")
	}

	d.Round, err = readUInt8(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading round")
	}

	d.RoundPhase, err = readUInt8(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading roundPhase")
	}

	d.Model, err = readString(r, d.CIDLen)
	if err != nil {
		return nil, errors.WithMessage(err, "reading model")
	}

	d.Weight, err = readString(r, d.CIDLen)
	if err != nil {
		return nil, errors.WithMessage(err, "reading weight")
	}

	accuracy, err := readUInt8Array(r, len(d.Accuracy))
	if err != nil {
		return nil, errors.WithMessage(err, "reading accuracy")
	}
	copy(d.Accuracy[:], accuracy)

	loss, err := readUInt8Array(r, len(d.Loss))
	if err != nil {
		return nil, errors.WithMessage(err, "reading loss")
	}
	copy(d.Loss[:], loss)


	// grid, err := readUInt8Array(r, len(d.Grid))
	// if err != nil {
	// 	return nil, errors.WithMessage(err, "reading grid")
	// }
	// copy(d.Grid[:], makeFieldValueArray(grid))
	return &d, nil
}

// ValidInit checks that the initial state is valid.
func (a *FLApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != numParts {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", numParts, len(p.Parts))
	}

	appData, ok := s.Data.(*FLAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", s.Data)
	}

	zero := FLAppData{}
	// if appData.Grid != zero.Grid {
	// 	return fmt.Errorf("invalid starting grid: %v", appData.Grid)
	// }

	if appData.Model != zero.Model {
		return fmt.Errorf("invalid starting model: %v", appData.Model)
	}

	if appData.NumberOfRounds != zero.NumberOfRounds {
		return fmt.Errorf("invalid starting numberOfRounds: %v", appData.NumberOfRounds)
	}

	if appData.Round != zero.Round {
		return fmt.Errorf("invalid starting round: %v", appData.Round)
	}

	if appData.RoundPhase != zero.RoundPhase {
		return fmt.Errorf("invalid starting roundPhase: %v", appData.RoundPhase)
	}

	if appData.Weight != zero.Weight {
		return fmt.Errorf("invalid starting weight: %v", appData.Weight)
	}

	if appData.Accuracy != zero.Accuracy {
		return fmt.Errorf("invalid starting accuracy: %v", appData.Accuracy)
	}

	if appData.Loss != zero.Loss {
		return fmt.Errorf("invalid starting loss: %v", appData.Loss)
	}

	if s.IsFinal {
		return fmt.Errorf("must not be final")
	}

	if appData.NextActor >= numParts {
		return fmt.Errorf("invalid next actor: got %d, expected < %d", appData.NextActor, numParts)
	}
	return nil
}

func checkServerTransitionConstraints(fromData, toData *FLAppData) error {
	//print fromData
	// fmt.Println("fromData: %v", fromData.String())
	// //print toData
	// fmt.Println("toData: %v", toData.String())

	if fromData.NextActor == uint8(0) { // Server conditions

		if fromData.Model == toData.Model {
			return fmt.Errorf("actor: %v cannot skip model", fromData.NextActor)
		}

		if fromData.RoundPhase != 0 && toData.Round != fromData.Round+1 {
			return fmt.Errorf("actor: %v must increment round: expected %v, got %v", fromData.NextActor, fromData.Round+1, toData.Round)
		}

		if fromData.RoundPhase != 0 && !reflect.DeepEqual(fromData.Weight, toData.Weight){
			return fmt.Errorf("actor: %v cannot override weight: expected %v, got %v", fromData.NextActor, fromData.Weight, toData.Weight)
		}

		if !equalExcept(fromData.Accuracy[:], toData.Accuracy[:], int(fromData.Round)) {
			return fmt.Errorf("actor: %v cannot override accuracy outside current round: expected %v, got %v", fromData.NextActor, fromData.Accuracy, toData.Accuracy)
		}

		if fromData.RoundPhase != 0 && toData.Accuracy[fromData.Round] == 0 && toData.Loss[fromData.Round] == 0 { //accuracy and loss are not set
			return fmt.Errorf("actor: %v cannot skip accuracy and loss", fromData.NextActor)
		}

		if !equalExcept(fromData.Loss[:], toData.Loss[:], int(fromData.Round)) {
			return fmt.Errorf("actor: %v cannot override loss outside current round: expected %v, got %v", fromData.NextActor, fromData.Loss, toData.Loss)
		}

	}
	return nil
}

func checkClientTransitionConstraints(fromData, toData *FLAppData) error {
	if fromData.NextActor == uint8(1) { //Client conditions
		if fromData.Model != toData.Model {
			return fmt.Errorf("actor: %v cannot override model: expected %v, got %v", fromData.NextActor, fromData.Model, toData.Model)
		}
		if fromData.Round != toData.Round {
			return fmt.Errorf("actor: %v cannot override round: expected %v, got %v", fromData.NextActor, fromData.Round, toData.Round)
		}
		if !reflect.DeepEqual(fromData.Accuracy, toData.Accuracy) {
			return fmt.Errorf("actor: %v cannot override accuracy: expected %v, got %v", fromData.NextActor, fromData.Accuracy, toData.Accuracy)
		}
		if !reflect.DeepEqual(fromData.Loss, toData.Loss) {
			return fmt.Errorf("actor: %v cannot override loss: expected %v, got %v", fromData.NextActor, fromData.Loss, toData.Loss)
		}

		// if !equalExcept(fromData.Weight[:], toData.Weight[:], int(toData.Round)){
		// 	return fmt.Errorf("actor: %v cannot override weights outside current round: expected %v, got %v", fromData.NextActor, fromData.Weight, toData.Weight)
		// }

		if fromData.RoundPhase == 1 && toData.Weight == fromData.Weight { //weight is not set
			return fmt.Errorf("actor: %v cannot skip weight %v -> %v", fromData.NextActor, fromData.Weight, toData.Weight)
		}
	}
	return nil
}

func checkFLRoundTransitionConstraints(fromData, toData *FLAppData) error {

	if fromData.RoundPhase != 0 && fromData.NumberOfRounds != toData.NumberOfRounds {
		return fmt.Errorf("cannot override number of rounds: expected %v, got %v", fromData.NumberOfRounds, toData.NumberOfRounds)
	}

	if toData.Round > toData.NumberOfRounds {
		return fmt.Errorf("round out of bounds: %v", toData.Round)
	}
	return nil
}

// ValidTransition is called whenever the channel state transitions.
func (a *FLApp) ValidTransition(params *channel.Params, from, to *channel.State, idx channel.Index) error {
	err := channel.AssetsAssertEqual(from.Assets, to.Assets)
	if err != nil {
		return fmt.Errorf("invalid assets: %v", to.Assets)
	}

	fromData, ok := from.Data.(*FLAppData)
	if !ok {
		return fmt.Errorf("from state: invalid data type: %T", from.Data)
	}

	toData, ok := to.Data.(*FLAppData)
	if !ok {
		return fmt.Errorf("to state: invalid data type: %T", to.Data)
	}

	// Check actor.
	if fromData.NextActor != uint8safe(uint16(idx)) {
		return fmt.Errorf("invalid actor: expected %v, got %v", fromData.NextActor, idx)
	}

	// Check next actor.
	if len(params.Parts) != numParts {
		return fmt.Errorf("invalid number of participants expected %v, got %v", numParts, len(params.Parts))
	}

	expectedToNextActor := calcNextActor(fromData.NextActor)
	if toData.NextActor != expectedToNextActor {
		return fmt.Errorf("invalid next actor: expected %v, got %v", expectedToNextActor, toData.NextActor)
	}

	var roundCheck = checkFLRoundTransitionConstraints(fromData, toData)
	if roundCheck != nil {
		return roundCheck
	}

	var serverCheck = checkServerTransitionConstraints(fromData, toData)
	if serverCheck != nil {
		return serverCheck
	}

	var clientCheck = checkClientTransitionConstraints(fromData, toData)
	if clientCheck != nil {
		return clientCheck
	}

	// Check final and allocation.
	isFinal, _ := toData.CheckFinal()
	if to.IsFinal != isFinal {
		return fmt.Errorf("final flag: expected %v, got %v", isFinal, to.IsFinal)
	}
	// expectedAllocation := from.Allocation.Clone()
	// if winner != nil {
	// 	expectedAllocation.Balances = computeFinalBalances(from.Allocation.Balances, *winner)
	// }
	// if err := expectedAllocation.Equal(&to.Allocation); err != nil {
	// 	return errors.WithMessagef(err, "wrong allocation: expected %v, got %v", expectedAllocation, to.Allocation)
	// }
	return nil
}

func (a *FLApp) Set(s *channel.State, model string, numberOfRounds int, weight string, accuracy, loss int, actorIdx channel.Index) error {
	d, ok := s.Data.(*FLAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	fmt.Println("🔁 Setting state")

	err := d.Set(model, numberOfRounds, weight, accuracy, loss, actorIdx)
	if err != nil {
		return err
	}
	log.Println("\n" + d.String())

	fmt.Println(d.String())

	if checkClientReward := d.checkClientReward(); checkClientReward {
		fmt.Println("Client rewarded")
		s.Balances = payCLientForContrib(s.Balances)
		fmt.Printf("💰 Sent payment. New balance: [My: %v Ξ, Peer: %v Ξ]\n", s.Balances[0][1], s.Balances[0][0])

	}

	if isFinal, _ := d.CheckFinal(); isFinal {
		s.IsFinal = true
		// if winner != nil {
		// 	s.Balances = computeFinalBalances(s.Balances, *winner)
		// }
	}
	return nil
}
