/*
Copyright © 2024 Patrick X. Gray pxgray@proton.me
*/
package cmd

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gray-patrick/circle-cli/internal/utils"
	"github.com/metachris/eth-go-bindings/erc20"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	amountFlag    string
	toAddressFlag string
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer testnet USDC",
	Long: `Transfer USDC tokens from an origin wallet to a recipient wallet on
	the Ethereum Sepolia testnet.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the API endpoint from the config file and connect to it.
		conn, err := ethclient.Dial(viper.GetString("api_endpoint"))
		if err != nil {
			log.Fatal(err)
		}

		// Define the Sepolia address for the USDC contract and connect to it.
		tokenAddress := common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")
		token, err := erc20.NewErc20(tokenAddress, conn)
		if err != nil {
			log.Fatal(err)
		}

		// Get the number of decimals from the contract for converting between
		// human-readable token numbers and wei.
		decimals, err := token.Decimals(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Get the private key for the sending wallet from the config file.
		privateKey, err := crypto.HexToECDSA(viper.GetString("private_key"))
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("publicKey is not of type *ecdsa.PublicKey")
		}

		// Derive the address from the public key, and get the nonce for the
		// wallet.
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := conn.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		// Get a suggested gas price from the API endpoint.
		gasPrice, err := conn.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the amount to send in wei from the provided token number.
		amount := utils.ToWei(amountFlag, int(decimals))
		toAddress := common.HexToAddress(toAddressFlag)

		// Define the function used to sign the transaction.
		signer := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			chainId, err := conn.NetworkID(context.Background())
			if err != nil {
				return nil, err
			}

			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
			if err != nil {
				return nil, err
			}

			return signedTx, nil
		}

		// Define transfer transaction options. Note that Value is 0, this is
		// the value in native tokens (SepoliaETH), since we are sending USDC
		// this value is expected to be 0.
		opts := &bind.TransactOpts{
			From:   fromAddress,
			Nonce:  big.NewInt(int64(nonce)),
			Signer: signer,

			Value:    big.NewInt(0),
			GasPrice: gasPrice,
			GasLimit: uint64(300000),

			Context: context.Background(),
			NoSend:  viper.GetBool("dry-run"),
		}

		// Send the transaction to the chain
		tx, err := token.Transfer(opts, toAddress, amount)

		// Show the transaction hash to the user.
		fmt.Println(tx.Hash())
	},
}

// init defines flags for the transfer command.
func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVarP(&amountFlag, "amount", "a", "0", "Amount of USDC to send")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.Flags().StringVarP(&toAddressFlag, "recipient", "t", "", "Recipient address for USDC")
	transferCmd.MarkFlagRequired("recipient")
}
