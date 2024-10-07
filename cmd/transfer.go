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
		conn, err := ethclient.Dial("https://sepolia.infura.io/v3/398ccc0480ae443891c8768995332342")
		if err != nil {
			log.Fatal(err)
		}

		tokenAddress := common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")
		token, err := erc20.NewErc20(tokenAddress, conn)
		if err != nil {
			log.Fatal(err)
		}

		decimals, err := token.Decimals(nil)
		if err != nil {
			log.Fatal(err)
		}

		privateKey, err := crypto.HexToECDSA(viper.GetString("privatekey"))
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("publicKey is not of type *ecdsa.PublicKey")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := conn.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := conn.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		amount := utils.ToWei(amountFlag, int(decimals))
		toAddress := common.HexToAddress(toAddressFlag)

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

		tx, err := token.Transfer(opts, toAddress, amount)

		fmt.Println(tx.Hash())
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVarP(&amountFlag, "amount", "a", "0", "Amount of USDC to send")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.Flags().StringVarP(&toAddressFlag, "recipient", "t", "", "Recipient address for USDC")
	transferCmd.MarkFlagRequired("recipient")
}
