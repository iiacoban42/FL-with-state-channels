"""Create config files for the network.
Make sure ganache is already running on localhost:8545.
Run it with  ganache-cli -m "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic" --account_keys_path keys.json --accounts 55 --block-time 5 --gasPrice 100000000000.
With the flag --accounts the number of accounts to be created can be specified. The config files will be created for these accounts."""
import json
import yaml
import copy
import ipaddress
from eth_account import Account
from web3 import Web3, exceptions

DOCKERIZED = False
MNEMONIC = "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic"
Account.enable_unaudited_hdwallet_features()


def get_private_key(address):
    address = address.lower()
    # Load the JSON file
    with open('keys.json', 'r') as file:
        data = json.load(file)
        private_key = data['private_keys'][address]
        return data['private_keys'][address]


def remove_quotes(file_path):
    """Remove the "'" characters from the config files"""
    # Read the contents of the file
    with open(file_path, 'r', encoding='utf-8') as file:
        content = file.read()

    # Remove the desired character
    modified_content = content.replace("'", "")

    # Write the modified contents back to the file
    with open(file_path, 'w', encoding="utf-8") as file:
        file.write(modified_content)


def generate_peer_and_network_files():
    """Generate the config files for the network and the peers."""
    # Load config template
    config_file_template_path = "config_template.json"

    with open(config_file_template_path, "r", encoding="utf-8") as f:
        config_template = json.load(f)

    # Connect to Ganache (assuming it is running on localhost)
    ganache_url = "http://localhost:8545"

    web3 = Web3(Web3.HTTPProvider(ganache_url))

    # Retrieve the available accounts
    accounts = web3.eth.accounts

    # Create a dictionary to store the network config
    network = {'peers': {}}

    # Create config files for each peer
    for i, account in enumerate(accounts):
        if i == 0:
            config_template['chain']['contractsetup'] = 'validateordeploy'
        else:
            config_template['chain']['contractsetup'] = 'validate'

        config_template['alias'] = f"peer_{i}"
        config_template['accountIndex'] = i
        config_template["walletPath"] = f"/tmp/peer_{i}_wallet"
        config_template['secretKey'] = get_private_key(account)

        config_template['node']['port'] = 5750 + i
        config_template['node']['apiPort'] = 8081 + i
        config_template['node']['persistencePath'] = f"/tmp/peer_{i}"
        config_template['log']['file'] = f"logs/peer_{i}.log"
        if not DOCKERIZED:
            config_template['chain']['url'] = "ws://127.0.0.1:8545"

        with open(f"config/peer_{i}.yaml", "w", encoding="utf-8") as f:
            yaml.safe_dump(config_template, f, default_style=None, default_flow_style=False, sort_keys=False)

        # print(account)
        remove_quotes(f"config/peer_{i}.yaml")
        if DOCKERIZED:
            host = f"peer_{i}"
        else:
            host = "127.0.0.1"
        peer = {
            f"peer_{i}": {
                "perunID": account,
                "hostname": host,
                "port": 5750 + i,
                "apiport": 8081 + i,
            }
        }
        network["peers"].update(peer)

    # Create network config file
    with open("config/network.yaml", "w", encoding="utf-8") as f:
        yaml.safe_dump(network, f, default_style=None, default_flow_style=False, sort_keys=False)

    remove_quotes("config/network.yaml")

    return network

def generate_docker_compose_file(network, subnet, min_ip, max_ip):
    """Generate the docker-compose file."""
    # Load config template
    config_file_template_path = "docker_compose_template.json"
    with open(config_file_template_path, "r", encoding="utf-8") as f:
        config_template = json.load(f)

    # Configure network subnet
    config_template["networks"]["perun-net"]["ipam"]["config"]= [
            {
              "subnet": subnet,
            }
          ]
    # Set ganache IP
    config_template["services"]["ganache"]["networks"]["perun-net"]["ipv4_address"] = min_ip

    # Configure the peer nodes
    peer_config = copy.deepcopy(config_template["services"]["peer_0"])
    # Init IP address counter
    peer_ip = ipaddress.ip_address(min_ip)
    for peer, config in network["peers"].items():
        # Increment the IP address
        peer_ip = peer_ip + 1
        # Check if the IP address is in the subnet
        if peer_ip > ipaddress.ip_address(max_ip):
            raise ValueError("The IP address range is too small for the network.")
        # Configure the peer node
        peer_config["container_name"] = config["hostname"]
        peer_config["hostname"] = config["hostname"]
        peer_config["ports"] = [f"{config['port']}:{config['port']}", f"{config['apiport']}:{config['apiport']}"]
        peer_config["depends_on"] = ["ganache"]
        peer_config["command"] = f'sh -c "./wait-for-it.sh 127.0.0.1:8545 &&./app-channel demo --config config/{peer}.yaml --log-level trace --log-file logs/{peer}.log"'
        # Add dependency on peer_0 if it is not peer_0
        # if peer != "peer_0":
        #     peer_config["depends_on"].append("peer_0")
        #     peer_0_port = network["peers"]["peer_0"]["port"]
        #     peer_config["command"] = f'sh -c "./wait-for-it.sh peer_0:{peer_0_port} &&./app-channel demo --config config/{peer}.yaml"'

        peer_config["networks"]["perun-net"]["ipv4_address"] = str(peer_ip)
        # Add the peer to the config file
        if peer in config_template["services"]:
            config_template["services"][peer] = copy.deepcopy(peer_config)
        else:
            config_template["services"].update({copy.deepcopy(peer): copy.deepcopy(peer_config)})

    # Create the docker-compose file
    with open("docker-compose.yaml", "w", encoding="utf-8") as f:
        yaml.safe_dump(config_template, f, default_style=None, default_flow_style=False, sort_keys=False)

# Generate the config files
net = generate_peer_and_network_files()

generate_docker_compose_file(net, "172.18.0.0/16", "172.18.0.2", "172.18.225.225")
