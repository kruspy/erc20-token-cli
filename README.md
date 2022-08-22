# ERC20 Token CLI Tool

Learning CLI tool to create, deploy and interact with ERC20 tokens using Go.


## Installation

1. Install [Go](https://go.dev/doc/install). The project currently uses version 1.18.
2. Install [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).
3. Install [solc](https://docs.soliditylang.org/en/v0.8.9/installing-solidity.html) compiler and [abigen](https://geth.ethereum.org/docs/install-and-build/installing-geth).
4. Run `make` in order to compile the contract and install the CLI tool.
5. Run `tokencli --help` and the tool help message will appear.

## Contract Implementation

The main contract is already defined [here](./contracts/ERC20Token.sol) and just implements the ERC20 interface.

The contract implementation overrides the `decimal()` method of the ERC20 interface to allow a 6 decimal representation. This because the generated bindings only accept a [big.Int](https://pkg.go.dev/math/big#Int) which can represent values up to 64-bit signed integrers (9*10^18 aprox.). This would mean that only a totalSupply of 9 full tokens would be possible with 18 decimals. Since this is just a learning tool, the default decimals have been lowered to 6, allowing more possibilities.


# License

The MIT License (MIT) 2022 - [Marc Puig](https://github.com/kruspy). Please have a look at the [LICENSE](LICENSE) for more details.