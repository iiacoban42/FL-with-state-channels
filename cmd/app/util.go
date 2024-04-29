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
	"math/big"

	"perun.network/go-perun/channel"
)

const numParts = 2

type FieldValue uint8

// type IPFSData struct {
// 	Hash string
// 	len uint8
// }

const (
	notSet FieldValue = iota
	player1
	player2
	maxFieldValue = player2
)

func (v FieldValue) String() string {
	switch v {
	case notSet:
		return " "
	case player1:
		return "x"
	case player2:
		return "o"
	default:
		panic(fmt.Sprintf("unsupported value: %d", v))
	}
}

func makeFieldValueFromPlayerIdx(idx channel.Index) FieldValue {
	switch idx {
	case 0:
		return player1
	case 1:
		return player2
	default:
		panic("invalid")
	}
}

func (v FieldValue) PlayerIndex() channel.Index {
	switch v {
	case player1:
		return 0
	case player2:
		return 1
	default:
		panic("invalid")
	}
}

func (d FLAppData) checkClientReward() bool {
	if d.RoundPhase == uint8(2) && d.Weight[d.Round] != uint8(0) {
		return true
	}
	return false
}

func (d FLAppData) CheckFinal() (isFinal bool, winner *channel.Index) {

	if d.NumberOfRounds == d.Round && d.RoundPhase == uint8(3) {
		if d.Accuracy[0] >= 60 {
			index := makeFieldValueFromPlayerIdx(1)
			playerIndex := index.PlayerIndex()
			return true, &playerIndex
		}
		index := makeFieldValueFromPlayerIdx(0)
		playerIndex := index.PlayerIndex()
		return true, &playerIndex
	}

	return false, nil
}

func uint8safe(a uint16) uint8 {
	b := uint8(a)
	if uint16(b) != a {
		panic("unsafe")
	}
	return b
}

func readUInt8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	return buf[0], err
}

func writeUInt8(w io.Writer, v uint8) error {
	_, err := w.Write([]byte{v})
	return err
}

func readUInt8Array(r io.Reader, n int) ([]uint8, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func writeUInt8Array(w io.Writer, v []uint8) error {
	_, err := w.Write(v)
	return err
}

func writeString(w io.Writer, s string, slen uint8) error {
	n := uint8(len([]byte(s)))
	if n == slen {
		_, err := w.Write([]byte(s))
		return err
	}
	buf := make([]byte, slen)
	_, err := w.Write(buf)
	return err
}


func readString(r io.Reader, n uint8) (string, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return string(buf), err
}

func makeFieldValueArray(a []uint8) []FieldValue {
	b := make([]FieldValue, len(a))
	for i := range b {
		b[i] = FieldValue(a[i])
	}
	return b
}


// func makeUInt8Array(a []FieldValue) []uint8 {
// 	b := make([]uint8, len(a))
// 	for i := range b {
// 		b[i] = uint8(a[i])
// 	}
// 	return b
// }

// func makeStringArray(a []FieldValue) []string {
// 	b := make([]string, len(a))
// 	for i := range b {
// 		b[i] = a[i].String()
// 	}
// 	return b
// }

func makeUInt8Array(a []uint8) []uint8 {
	b := make([]uint8, len(a))
	for i := range b {
		b[i] = uint8(a[i])
	}
	return b
}

// func makeStringArray(a []uint8) []string {
// 	b := make([]string, len(a))
// 	for i := range b {
// 		b[i] = a[i].String()
// 	}
// 	return b
// }

func computeFinalBalances(bals channel.Balances, winner channel.Index) channel.Balances {
	loser := 1 - winner
	finalBals := bals.Clone()
	for i := range finalBals {
		finalBals[i][winner] = new(big.Int).Add(bals[i][0], bals[i][1])
		finalBals[i][loser] = big.NewInt(0)
	}
	return finalBals
}

func payCLientForContrib(bals channel.Balances) channel.Balances {

	contribFee := int64(1000000000000000000)
	client := 1
	server := 0
	newBals := bals.Clone()
	for i := range newBals {
		fmt.Println("newBals[i][client]: ", newBals[i][client])
		newBals[i][client] = new(big.Int).Add(bals[i][client], big.NewInt(contribFee))
		newBals[i][server] = new(big.Int).Sub(bals[i][server], big.NewInt(contribFee))
	}
	return newBals
}

// check if 2 arrays are equal except for one element at position idx
func equalExcept(arr1, arr2 []uint8, idx int) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := range arr1 {
		if i == idx {
			continue
		}
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
