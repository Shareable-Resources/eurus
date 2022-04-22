
pragma solidity >=0.6.0 <0.8.0;


contract StoreTest {
  event ResultOffset(uint256 indexed off, bytes content);
    //event TransferRequestFailed(address indexed dest, uint256 indexed userGasUsed, uint256 indexed amount, string assetName, bytes revertReason);

    fallback () external payable{
    //   uint256 gasBegin = gasleft();
    //   address _impl = getUserWalletImplementation();
    //   require(_impl != address(0), "getUserWalletImplementation is 0");

      bytes memory ptr; 
      uint256 offset; 
      assembly {
          ptr := mload(0x40)
          calldatacopy(ptr, 0, calldatasize())
        //   let result := delegatecall(gas(), _impl, ptr, calldatasize(), 0, 0)
        //   let size := returndatasize()
        //   returndatacopy(ptr, 0, size)
          let size:= 100
          let result := 0
          switch result
          case 0{
              
              if gt (size, 0) {
                let x := 0
                if gt (mod (size, 0x20) , 0) {
                    x := 1
                }
                offset := add ( mul ( add(div(size, 0x20), x) , 0x20 ) , ptr)
                mstore(0x40, offset)
                let ptr2 := mload(0x40)
                calldatacopy(ptr2, 0, calldatasize())
                ptr := add(ptr, 4)
              }
          }
          default{
            return(ptr, size)
          }
      }
      emit ResultOffset(offset, ptr);

    }
}