// Package utils contains internal helper functions for erc20-token-cli commands.
package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"runtime"

	token "github.com/kruspy/erc20-token-cli/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}

// Client opens a connection to an Ethereum RPC node and returns an
// object representing the connection.
func Client(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		Fatalf(err.Error())
	}
	return client
}

// Auth generates a transaction signer from a single private key.
func Auth(client *ethclient.Client, privateKeyAddress string) *bind.TransactOpts {

	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
	if err != nil {
		Fatalf(err.Error())
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		Fatalf("The provided private key is not valid.")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		Fatalf(err.Error())
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		Fatalf(err.Error())
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		Fatalf(err.Error())
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		Fatalf(err.Error())
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}

func TokenName(token *token.Token) string {
	name, err := token.Name(nil)
	if err != nil {
		Fatalf(err.Error())
	}
	return name
}

func TokenSymbol(token *token.Token) string {
	symbol, err := token.Symbol(nil)
	if err != nil {
		Fatalf(err.Error())
	}
	return symbol
}

func TokenDecimals(token *token.Token) uint8 {
	decimals, err := token.Decimals(nil)
	if err != nil {
		Fatalf(err.Error())
	}
	return decimals
}

func TokenTotalSupply(token *token.Token) *big.Float {
	minRepSupply, err := token.TotalSupply(nil)
	if err != nil {
		Fatalf(err.Error())
	}
	return new(big.Float).Quo(new(big.Float).SetInt(minRepSupply),
		big.NewFloat(math.Pow10(int(TokenDecimals(token)))))
}
