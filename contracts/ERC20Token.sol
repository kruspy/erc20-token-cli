// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract ERC20Token is ERC20 {
    constructor(string memory name, string memory symbol, uint256 initialSupply) ERC20(name, symbol) {
        _mint(msg.sender, initialSupply);
    }

    function decimals() override public virtual view returns (uint8) {
        // Decimals is overriden to 9 as the Go bindings generated
        // only support passing the initial supply of the contract as
        // int64 which can only hold up to 9*10^18. This would effectevly 
        // leave he whole supply with just 9 full tokens.
        return 6;
    }
}
