#!/bin/sh

solcjs FL.sol --abi --base-path ./ -o abi --bin --optimize
abigen --pkg FLApp --abi abi/FL_sol_FLApp.abi --bin abi/FL_sol_FLApp.bin --out ./generated/FLApp/FLApp.go
# generate_bindings ./perun-eth-contracts/contracts/Adjudicator.sol adjudicator
# generate_bindings ./perun-eth-contracts/contracts/AssetHolderETH.sol assetHolderETH
