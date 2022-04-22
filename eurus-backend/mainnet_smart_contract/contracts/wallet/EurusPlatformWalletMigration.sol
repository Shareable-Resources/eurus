pragma solidity >=0.6.0 <0.8.0;

pragma experimental ABIEncoderV2;
import "./EurusPlatformWallet.sol";

contract EurusPlatformWalletMigration is EurusPlatformWallet {
        using SafeMath for uint256;
        using Address for address;
        using Address for address payable;
        using SafeERC20 for ERC20;

    function directTransfer(string memory assetName, address destWallet) public onlyOwner {
        address self = address (this);

        if (keccak256( abi.encodePacked(assetName) ) != keccak256(abi.encodePacked("ETH"))){
            address ercAddr = internalConfig.getErc20SmartContractAddrByAssetName(assetName);
            require(ercAddr != address(0), "Asset not found");
            ERC20 erc20 = ERC20(ercAddr);
            erc20.safeTransfer(destWallet, erc20.balanceOf(self));
        }else {
            uint256 balance = self.balance;
            payable(destWallet).sendValue(balance);
        }
    }
}