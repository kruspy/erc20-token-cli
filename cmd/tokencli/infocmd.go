package main

import (
	"fmt"

	"github.com/kruspy/erc20-token-cli/cmd/utils"
	token "github.com/kruspy/erc20-token-cli/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

var (
	infoCommand = &cli.Command{
		Action:    tokenInfo,
		Name:      "info",
		Usage:     "Retrieve basic information about an ERC20 token",
		ArgsUsage: "<tokenAddress>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ethclient",
				Value: "http://127.0.0.1:8545",
				Usage: "Ethereum client to connect",
			},
		},
		Description: `
Displays basic information about an ERC20 token. The information displayed contains
of the token name, its symbol, the decimals and the total supply.
An optional argument can be passed to specify the network the token contract
is deployed to.
`,
	}
)

// tokenInfo is the entrypoint func for the info command. It quesries a given
// deployed contract and displays the token name, symbol, decimals and totalSupply.
func tokenInfo(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		utils.Fatalf("The token address must be passed as a mandatory argument")
	}

	address := common.HexToAddress(ctx.Args().Get(0))
	// Get an instance of the token to interact with the contract
	instance, err := token.NewToken(address, utils.Client(ctx.String("ethclient")))
	if err != nil {
		utils.Fatalf(err.Error())
	}

	fmt.Printf("NAME: %v\nSYMBOL: %v\nDECIMALS: %v\nSUPPLY: %v\n",
		utils.TokenName(instance), utils.TokenSymbol(instance),
		utils.TokenDecimals(instance), utils.TokenDecimals(instance))
	return nil
}
