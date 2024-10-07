# Circle CLI

A demo implementation for interacting with the USDC smart contract on the
Ethereum Sepolia testnet.

## Running the Circle CLI

To run this demo, you must have a working Go environment. You will also need the
private key for your Ethereum Sepolia testnet wallet, and an Infura API key with
the appropriate endpoint. You will also need USDC and Sepolia ETH in the sending
wallet, which can be obtained from public faucets:

- [USDC testnet faucet](https://faucet.circle.com)
- [Sepolia ETH testnet faucet](https://cloud.google.com/application/web3/faucet/ethereum/sepolia)

This repository was developed with Go version `1.23.2`. To run this demo, clone
this repository and install dependencies:

```shell
git clone https://github.com/gray-patrick/circle-cli
cd circle-cli
go mod tidy
```

Create a configuration file in your home directory named `.circle-cli`. This is
a YAML-formatted file that contains the private key for the sending wallet, and
the Infura Sepolia endpoint to use:

```yaml
private_key: <YOUR_PRIVATE_KEY>
api_endpoint: <INFURA_API_ENDPOINT>
```

**Note:** Anyone with your wallet's private key can perform transactions from
your wallet. You should not use the private key from your primary wallet for
this demonstration. Create a new wallet.

The `main.go` file is the entrypoint for the CLI:

```shell
go run main.go --help
```

### The transfer command

You can transfer testnet USDC from your wallet with the transfer command. It is
recommended that you do a dry run before sending the transaction:

```shell
go run main.go transfer -t <RECIPIENT_WALLET_ADDRESS> -a <AMOUNT> --dry-run
```

To send the transaction, remove the `--dry-run` flag from the command.
