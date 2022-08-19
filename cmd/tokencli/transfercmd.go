package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/kruspy/erc20-token-cli/cmd/utils"
	token "github.com/kruspy/erc20-token-cli/contracts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/sha3"
)

var (
	transferCommand = &cli.Command{
		Action:    transferToken,
		Name:      "transfer",
		Usage:     "Transfer ERC20 tokens",
		ArgsUsage: "<quantity>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "from",
				Usage:    "Sender address private key",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "to",
				Usage:    "Receiver address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Token contract address",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "ethclient",
				Value: "http://127.0.0.1:8545",
				Usage: "Ethereum client to connect",
			},
		},
		Description: `
Transfer ERC20 tokens between two addresses. The sender address passed
to the command needs to be the private key of the sender, not the public one.
Any error while performing the transfer will rollback all the process and
no tokens will be transferred. 
`,
	}
)

func transferToken(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		utils.Fatalf("The transfer quantity hasn't been provided")
	}

	client := utils.Client(ctx.String("ethclient"))
	senderKey := ctx.String("from")
	if strings.Contains(senderKey, "0x") { // The private key is needed in order to sign the tx
		utils.Fatalf("The private key of the sender address needs to be provided, not the public key")
	}

	privateKey, err := crypto.HexToECDSA(senderKey)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	// Generate the public key from the private one to retrieve the tx nonce
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		utils.Fatalf("Could not generate a valid public key from the key provided")
	}

	// Use the public key of the sender to get the tx nonce
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	value := big.NewInt(0) // When not transacting ETH, the value of the tx is 0

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		utils.Fatalf(err.Error())
	}

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4] // Hex signature of the Transfer contract method

	// Left padding to 32 bytes is required since the EVM use 32 byte wide data structures
	paddedAddress := common.LeftPadBytes(common.HexToAddress(ctx.String("to")).Bytes(), 32)

	tokenAddress := common.HexToAddress(ctx.String("token"))
	instance, err := token.NewToken(tokenAddress, utils.Client(ctx.String("ethclient")))
	if err != nil {
		utils.Fatalf(err.Error())
	}

	// Calculate the minimal representation of the tokens to be transfered
	// using the amount provided by the user and the number of decimals defined
	// in the contract, then left padding is also required
	decimals := math.Pow(10, float64(utils.TokenDecimals(instance)))
	wholeAmount, _ := strconv.ParseFloat(ctx.Args().Get(0), 64)
	amount := big.NewInt(int64((wholeAmount * decimals)))
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte // The tx data needs to be passed as a byte slice
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// The EstimateGas method provided by the client is able to esimate the gas
	// for us based on the most recent state of the blockchain, so no hardcoded
	// values are needed
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		utils.Fatalf(err.Error())
	}

	// Create the tx using all the values calculated before
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		utils.Fatalf(err.Error())
	}

	// Sign the tx
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	// Send it to the chain
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	fmt.Printf("Transaction sent.\nHash: %v\n", signedTx.Hash().Hex())

	return nil
}
