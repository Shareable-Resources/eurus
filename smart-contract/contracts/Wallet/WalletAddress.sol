pragma solidity >=0.6.0 <0.8.0;
import "../Ownable.sol";

contract WalletAddress is Ownable{
    struct WalletInfo{
        string userName;
        bool isMerchant;
        bool isMetaMask;
        bool isExist;
    }
    mapping(address=>WalletInfo) walletInfoMap;
    address[] walletInfoList;

    constructor()public{}
    
    modifier OnlyNotAdded(address addr) {
        require(walletInfoMap[addr].isExist!=true, "Invalid to add an added address!");
        _;
    }
    
    function addWalletInfo(address _WalletAddress, string memory UserName,bool IsMerchant, bool IsMetaMask)public OnlyNotAdded(_WalletAddress) onlyOwner{
        walletInfoMap[_WalletAddress]=WalletInfo(UserName,IsMerchant,IsMetaMask, true);
        walletInfoList.push(_WalletAddress);
    }
    
    function removeWalletInfo(address _WalletAddress) public onlyOwner{
        walletInfoMap[_WalletAddress].isExist=false;
        for (uint i=0; i<walletInfoList.length - 1; i++)
            if (walletInfoList[i] == _WalletAddress) {
                walletInfoList[i] = walletInfoList[walletInfoList.length - 1];
                break;
            }
        walletInfoList.pop();
    }
    
    function getWalletInfoList()public view returns(address[] memory) {
        return walletInfoList;
    }
    
    function isMerchantAddress(address addr)public returns(bool){
        return walletInfoMap[addr].isMerchant;
    }
    
    
}