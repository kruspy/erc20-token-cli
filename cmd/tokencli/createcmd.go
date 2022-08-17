package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/kruspy/erc20-token-cli/cmd/utils"
	token "github.com/kruspy/erc20-token-cli/contracts"

	"github.com/holiman/uint256"
	"github.com/urfave/cli/v2"
)

var (
	createCommand = &cli.Command{
		Action:    createToken,
		Name:      "create",
		Usage:     "Create an ERC20 token",
		ArgsUsage: "<tokenName> <tokenSymbol> <tokenSupply>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "account",
				Usage:    "Account which will deploy the contract",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "ethclient",
				Value: "http://127.0.0.1:8545",
				Usage: "Ethereum client to connect",
			},
		},
		Description: `
Deploys an ERC20 contract to an EVM based chain. The token contract 
can be configured using the first three mandatory arguments. The flag
argument determines the account which will deploy the contract to the network.
An optional flag can be passed to specify the network the contract will be
deployed to.
`,
	}
)

// createToken is the entrypoint func for the create command. It deploys a custom
// ERC20 smart contract and prints to standard output the contract address and the
// transaction in which it was included.
func createToken(ctx *cli.Context) error {
	if ctx.Args().Len() < 3 {
		utils.Fatalf("There are missing arguments.")
	}

	var tokenSupply *big.Int // Smallest representation (supply * 10^decimals)
	// The amount of whole tokens is restricted to 10^18 as it is represented as int64
	wholeSupply, err := strconv.ParseInt(ctx.Args().Get(2), 10, 64)
	if err == nil { // uint256 is used to multiply
		supply, _ := uint256.FromBig(big.NewInt(int64(wholeSupply)))
		decimals, _ := uint256.FromBig(big.NewInt(1000000000)) // Decimals are hardcoded in the contract
		tokenSupply = supply.Mul(supply, decimals).ToBig()
	} else {
		utils.Fatalf(err.Error())
	}

	client := utils.Client(ctx.String("ethclient"))
	tokenName := ctx.Args().Get(0)
	tokenSymbol := ctx.Args().Get(1)

	// Deploy the token to the selected network using the provided account
	address, tx, _, err := token.DeployToken(
		utils.Auth(client, ctx.String("account")),
		client,
		tokenName,
		tokenSymbol,
		tokenSupply,
	)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	fmt.Printf("%v (%v) has been successfully deployed.\nContract address: %v\nTransaction: %v\n", tokenName, tokenSymbol, address, tx.Hash())
	return nil
}
