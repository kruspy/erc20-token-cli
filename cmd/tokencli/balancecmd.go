package main

import (
	"fmt"
	"math"
	"math/big"

	"github.com/kruspy/erc20-token-cli/cmd/utils"
	token "github.com/kruspy/erc20-token-cli/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

var (
	balanceCommand = &cli.Command{
		Action: tokenBalance,
		Name:   "balance",
		Usage:  "Check the balance of an ERC20 token for a given wallet address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Token contract address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "wallet",
				Usage:    "Wallet address",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "ethclient",
				Value: "http://127.0.0.1:8545",
				Usage: "Ethereum client to connect",
			},
		},
		Description: `
Retrieves the token balance of an address. Both the address of the token contract
and the address of the wallet are required flags.
An optional flag can be passed to specify the network where the query is going
to take place.
`,
	}
)

func tokenBalance(ctx *cli.Context) error {
	tokenAddress := common.HexToAddress(ctx.String("token"))
	instance, err := token.NewToken(tokenAddress, utils.Client(ctx.String("ethclient")))
	if err != nil {
		utils.Fatalf(err.Error())
	}

	walletAddress := common.HexToAddress(ctx.String("wallet"))
	minimalRepBalance, err := instance.BalanceOf(&bind.CallOpts{}, walletAddress)
	if err != nil {
		utils.Fatalf(err.Error())
	}

	fbal := new(big.Float)
	fbal.SetString(minimalRepBalance.String())
	balance := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(utils.TokenDecimals(instance)))))

	fmt.Printf("Wallet address: %v\nCurrent balance: %v%v\n", ctx.String("wallet"), balance, utils.TokenSymbol(instance))

	return nil
}
