# Instructions for creating a network with `n` nodes
1. Make sure configuration files for each node are present in the `config` directory. Additionally, make sure the `network.yaml` contains all nodes desired to be part of the network.

    To generate these files `generate_config_files.py` can be used. The script assumes a ganache network is running and generates configuration files and a network configuration file for the accounts it finds in the network. To also include the account keys in the configuration files, those need to be provided in `keys.json`.

    `start_ganache.sh` can be used to start a ganache network with the desired number of accounts.

2. Start the ganache network with `start_ganache.sh` script. The script will start a ganache network with the desired number of accounts and pre-fund them with the desired amount of ETH.

3. To bring up the Perun network, run `start_network.sh` script which will open a terminal window for and initialize each node. In `start_network.sh` the number of nodes to start can be specified. Note that `start_network.sh` might need to be run twice if the smart contracts have not been deployed yet.

4. To interact with the network, the CLI on the nodes can be used to open channels, send payments, and close channels.

5. The smart contract governing the state channels can be modified in `cmd/contracts/FL.sol` and compiled with `cmd/contracts/generate.sh`