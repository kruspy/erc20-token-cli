// erc20-token-cli is a learning tool to interact with ERC20 smart contracts in the Ethereum chain.
package main

import (
	"os"

	"github.com/kruspy/erc20-token-cli/cmd/utils"

	"github.com/urfave/cli/v2"
)

var (
	app = NewApp("the erc20 token command line interface")
)

func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Usage = usage
	return app
}

func init() {
	app.Commands = []*cli.Command{
		createCommand,
		infoCommand,
		balanceCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		utils.Fatalf(err.Error())
	}
}
