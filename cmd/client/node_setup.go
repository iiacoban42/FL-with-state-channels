// Copyright (c) 2020 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-eth-demo. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	phd "perun.network/go-perun/backend/ethereum/wallet/hd"
	pkeystore "perun.network/go-perun/backend/ethereum/wallet/keystore"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/persistence/keyvalue"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
	"perun.network/go-perun/pkg/sortedkv/leveldb"
	"perun.network/go-perun/watcher/local"
	wirenet "perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
	"perun.network/perun-examples/app-channel/cmd/app"
	"perun.network/perun-examples/app-channel/cmd/contracts/generated/FLApp"
)

var (
	backend         *node
	ethereumBackend *ethclient.Client
)

// Setup initializes the node, can not be done in init() since it needs the
// configuration from viper.
func Setup() {
	SetConfig()

	var err error
	if ethereumBackend, err = ethclient.Dial(config.Chain.URL); err != nil {
		log.WithError(err).Fatalln("Could not connect to ethereum node.")
	}
	if backend, err = newNode(); err != nil {
		log.WithError(err).Fatalln("Could not initialize node.")
	}
}

func newNode() (*node, error) {
	// wallet, acc, err := setupWallet(config.Mnemonic, config.AccountIndex)

	wallet, acc, err := importAccount(config.SecretKey)
	if err != nil {
		return nil, errors.WithMessage(err, "importing secret key")
	}
	dialer := simple.NewTCPDialer(config.Node.DialTimeout)
	signer := types.NewEIP155Signer(big.NewInt(config.Chain.ID))

	n := &node{
		log:     log.Get(),
		onChain: acc,
		wallet:  wallet,
		dialer:  dialer,
		cb:      ethchannel.NewContractBackend(ethereumBackend, pkeystore.NewTransactor(*wallet, signer), config.Chain.TxFinalityDepth),
		peers:   make(map[string]*peer),
	}
	return n, n.setup()
}

// setup does:
//   - Create a new offChain account.
//   - Create a client with the node's dialer, funder, adjudicator and wallet.
//   - Setup a TCP listener for incoming connections.
//   - Load or create the database and setting up persistence with it.
//   - Set the OnNewChannel, Proposal and Update handler.
//   - Print the configuration.
func (n *node) setup() error {
	if err := n.setupContracts(); err != nil {
		return errors.WithMessage(err, "setting up contracts")
	}

	var err error

	n.offChain = n.wallet.NewAccount()
	n.log.WithField("off-chain", n.offChain.Address()).Info("Generating account")

	n.bus = wirenet.NewBus(n.onChain, n.dialer)

	watcher, err := local.NewWatcher(n.adjudicator)
	if err != nil {
		return errors.WithMessage(err, "creating watcher")
	}
	if n.client, err = client.New(n.onChain.Address(), n.bus, n.funder, n.adjudicator, n.wallet, watcher); err != nil {
		return errors.WithMessage(err, "creating client")
	}

	host := config.Node.IP + ":" + strconv.Itoa(int(config.Node.Port))
	n.log.WithField("host", host).Trace("Listening for connections")
	listener, err := simple.NewTCPListener(host)
	if err != nil {
		return errors.WithMessage(err, "could not start tcp listener")
	}

	n.app = app.NewFLApp(ethwallet.AsWalletAddr(config.Chain.app))
	n.stake = big.NewInt(1)

	n.client.OnNewChannel(n.setupChannel)
	if err := n.setupPersistence(); err != nil {
		return errors.WithMessage(err, "setting up persistence")
	}
	channel.RegisterApp(n.app)

	go n.client.Handle(n, n)
	go n.bus.Listen(listener)
	n.PrintConfig()
	return nil
}

func (n *node) setupContracts() error {
	var adjAddr common.Address
	var assAddr common.Address
	var appAddr common.Address
	var err error

	fmt.Println("💭 Validating contracts...")

	switch contractSetup := config.Chain.contractSetup; contractSetup {
	case contractSetupOptionValidate:
		if adjAddr, err = validateAdjudicator(n.cb); err == nil { // validate adjudicator
			if assAddr, err = validateAssetHolder(n.cb, adjAddr); err == nil { // validate asset holder
				appAddr, err = validateContract(n.cb) // validate app contract
			}
		}
	case contractSetupOptionDeploy:
		if adjAddr, err = deployAdjudicator(n.cb, n.onChain.Account); err == nil { // deploy adjudicator
			if assAddr, err = deployAssetHolder(n.cb, adjAddr, n.onChain.Account); err == nil { // deploy asset holder
				appAddr, err = deployFLApp(n.cb, n.onChain.Account) // deploy app
			}
		}
	case contractSetupOptionValidateOrDeploy:
		if adjAddr, err = validateAdjudicator(n.cb); err != nil { // validate adjudicator
			fmt.Println("❌ Adjudicator invalid")
			adjAddr, err = deployAdjudicator(n.cb, n.onChain.Account) // deploy adjudicator
		}

		if err == nil {
			if assAddr, err = validateAssetHolder(n.cb, adjAddr); err != nil { // validate asset holder
				fmt.Println("❌ Asset holder invalid")
				assAddr, err = deployAssetHolder(n.cb, adjAddr, n.onChain.Account) // deploy asset holder
			}
		}
		if err == nil {
			if appAddr, err = validateContract(n.cb); err != nil { // validate app contract
				fmt.Println("❌ App contract invalid")
				appAddr, err = deployFLApp(n.cb, n.onChain.Account) // deploy app
			}

		}
	default:
		// unsupported setup method
		err = errors.New(fmt.Sprintf("Unsupported contract setup method '%s'.", contractSetup))
	}

	fmt.Println("✅ Contracts validated.")

	if err != nil {
		return errors.WithMessage(err, "contract setup failed")
	}

	n.adjAddr = adjAddr
	n.assetAddr = assAddr
	recvAddr := ethwallet.AsEthAddr(n.onChain.Address())
	n.adjudicator = ethchannel.NewAdjudicator(n.cb, n.adjAddr, recvAddr, n.onChain.Account)
	n.asset = (*ethwallet.Address)(&n.assetAddr)
	n.log.WithField("Adj", n.adjAddr).WithField("Asset", n.assetAddr).Debug("Set contracts")

	funder := ethchannel.NewFunder(n.cb)
	dep := ethchannel.NewETHDepositor()
	asset := ethwallet.Address(n.assetAddr)
	funder.RegisterAsset(asset, dep, n.onChain.Account)

	// funder := ethchannel.NewFunder(n.cb)
	// funder.RegisterAsset(ethwallet.Address(n.assetAddr), new(ethchannel.ETHDepositor), n.onChain.Account)
	n.funder = funder

	n.appAddr = appAddr

	return nil
}

func (n *node) setupPersistence() error {
	if config.Node.PersistenceEnabled {
		n.log.Info("Starting persistence")
		db, err := leveldb.LoadDatabase(config.Node.PersistencePath)
		if err != nil {
			return errors.WithMessage(err, "creating/loading database")
		}
		persister := keyvalue.NewPersistRestorer(db)
		n.client.EnablePersistence(persister)

		ctx, cancel := context.WithTimeout(context.Background(), config.Node.ReconnecTimeout)
		defer cancel()
		if err := n.client.Restore(ctx); err != nil {
			n.log.WithError(err).Warn("Could not restore client")
			// return the error.
		}
	} else {
		n.log.Info("Persistence disabled")
	}
	return nil
}

func validateAdjudicator(cb ethchannel.ContractBackend) (common.Address, error) {
	fmt.Println("📋 Validating adjudicator")
	ctx, cancel := newTransactionContext()
	defer cancel()

	adjAddr := config.Chain.adjudicator
	return adjAddr, ethchannel.ValidateAdjudicator(ctx, cb, adjAddr)
}

func validateAssetHolder(cb ethchannel.ContractBackend, adjAddr common.Address) (common.Address, error) {
	fmt.Println("📋 Validating asset holder")

	ctx, cancel := newTransactionContext()
	defer cancel()

	assAddr := config.Chain.assetholder
	return assAddr, ethchannel.ValidateAssetHolderETH(ctx, cb, assAddr, adjAddr)
}

func validateContract(cb ethchannel.ContractBackend) (common.Address, error) {
	fmt.Println("📋 Validating app contract")
	ctx, cancel := newTransactionContext()
	defer cancel()

	appAddr := config.Chain.app

	emptyAddr := common.Address{}
	if appAddr.String() == emptyAddr.String() {
		return appAddr, errors.New("address is empty")
	}
	// fmt.Println(appAddr.String())
	bin := FLApp.FLAppMetaData.Bin

	bytecode, err := cb.CodeAt(ctx, appAddr, nil) // nil is latest block

	if err != nil {
		return appAddr, errors.WithMessage(err, "getting contract code")
	}
	if hex.EncodeToString(bytecode) != bin {
		return appAddr, errors.WithMessage(err, "incorrect contract code")
	}
	return appAddr, err
}

// deployAdjudicator deploys the Adjudicator to the blockchain and returns its address
// or an error.
func deployAdjudicator(cb ethchannel.ContractBackend, acc accounts.Account) (common.Address, error) {
	fmt.Println("🌐 Deploying adjudicator")
	ctx, cancel := context.WithTimeout(context.Background(), config.Chain.TxTimeout)
	defer cancel()
	adjAddr, err := ethchannel.DeployAdjudicator(ctx, cb, acc)
	fmt.Println("🚀 Adjudicator Deployed")
	return adjAddr, errors.WithMessage(err, "deploying eth adjudicator")
}

// deployAssetHolder deploys the Assetholder to the blockchain and returns its address
// or an error. Needs an Adjudicator address as second argument.
func deployAssetHolder(cb ethchannel.ContractBackend, adjudicator common.Address, acc accounts.Account) (common.Address, error) {
	fmt.Println("🌐 Deploying assetholder")
	ctx, cancel := context.WithTimeout(context.Background(), config.Chain.TxTimeout)
	defer cancel()
	asset, err := ethchannel.DeployETHAssetholder(ctx, cb, adjudicator, acc)
	fmt.Println("🚀 Assetholder Deployed")
	return asset, errors.WithMessage(err, "deploying eth assetholder")
}

// deployFLApp deploys the FLApp to the blockchain and returns its address
func deployFLApp(cb ethchannel.ContractBackend, acc accounts.Account) (common.Address, error) {
	fmt.Println("🌐 Deploying FLApp")
	const gasLimit = 30000000 // Must be sufficient for deploying FL.sol.
	ctx, cancel := context.WithTimeout(context.Background(), config.Chain.TxTimeout)
	defer cancel()
	tops, err := cb.NewTransactor(ctx, gasLimit, acc)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "deploying FLApp")
	}
	// Deploy FL App.
	app, tx, _, err := FLApp.DeployFLApp(tops, cb)
	if err != nil {
		return app, errors.WithMessage(err, "deploying FLApp")
	}
	_, err = bind.WaitDeployed(ctx, cb, tx)
	fmt.Println("🚀 FLApp Deployed")
	fmt.Println("📋 FLApp Address: ", app.String())
	return app, errors.WithMessage(err, "wait to deploy FLApp")
}

// setupWallet imports the mnemonic and returns a corresponding wallet and
// the derived account at the given account index.
func setupWallet(mnemonic string, accountIndex uint) (*phd.Wallet, *phd.Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating hdwallet")
	}

	perunWallet, err := phd.NewWallet(wallet, accounts.DefaultBaseDerivationPath.String(), accountIndex)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating perun wallet")
	}
	acc, err := perunWallet.NewAccount()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating account")
	}

	return perunWallet, acc, nil
}

// importAccount is a helper method to import secret keys until we have the ethereum wallet done.
func importAccount(secret string) (*pkeystore.Wallet, *pkeystore.Account, error) {
	ks := keystore.NewKeyStore(config.WalletPath, 2, 1)
	sk, err := crypto.HexToECDSA(secret[2:])
	if err != nil {
		return nil, nil, errors.WithMessage(err, "decoding secret key")
	}
	var ethAcc accounts.Account
	addr := crypto.PubkeyToAddress(sk.PublicKey)
	if ethAcc, err = ks.Find(accounts.Account{Address: addr}); err != nil {
		ethAcc, err = ks.ImportECDSA(sk, "")
		if err != nil && errors.Cause(err).Error() != "account already exists" {
			return nil, nil, errors.WithMessage(err, "importing secret key")
		}
	}

	wallet, err := pkeystore.NewWallet(ks, "")
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating wallet")
	}

	wAcc := pkeystore.NewAccountFromEth(wallet, &ethAcc)
	acc, err := wallet.Unlock(wAcc.Address())
	return wallet, acc.(*pkeystore.Account), err
}

func newTransactionContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.Chain.TxTimeout)
}

func (n *node) PrintConfig() (string, error) {
	res := fmt.Sprintf(
		"Alias: %s\n"+
			"Listening: %s:%d\n"+
			"ETH RPC URL: %s\n"+
			"API port: %d\n"+
			"Perun ID: %s\n"+
			"OffChain: %s\n"+
			"ETHAssetHolder: %s\n"+
			"Adjudicator: %s\n"+
			"App: %s\n"+
			"", config.Alias, config.Node.IP, config.Node.Port, config.Chain.URL, config.Node.APIPort, n.onChain.Address().String(), n.offChain.Address().String(), n.assetAddr.String(), n.adjAddr.String(), n.appAddr.String())
	fmt.Println(res)
	res += "\n Known peers:\n"
	fmt.Println("Known peers:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	for alias, peer := range config.Peers {
		res += fmt.Sprintf("%s\t%v\t%s:%d\n", alias, peer.PerunID, peer.Hostname, peer.Port)
		fmt.Fprintf(w, "%s\t%v\t%s:%d\n", alias, peer.PerunID, peer.Hostname, peer.Port)
	}
	return res, w.Flush()
}
