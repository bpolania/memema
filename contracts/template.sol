pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract <TOKEN_NAME> is ERC20 {
    constructor()
        ERC20("<TOKEN_NAME>", "<TOKEN_SYMBOL>")
    {
    }
}