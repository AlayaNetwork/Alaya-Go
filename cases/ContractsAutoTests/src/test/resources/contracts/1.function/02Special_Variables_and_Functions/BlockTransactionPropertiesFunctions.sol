pragma solidity 0.5.9;
/**
 * 验证区块和交易属性的内置函数
 * @author liweic
 * @dev 2019/12/27 19:10
 */

contract BlockTransactionPropertiesFunctions {
    
    function getBlockhash(uint blockNumber) public view returns (bytes32) {
        // 获取指定区块的区块哈希
        return blockhash(blockNumber);
    }

    function getBlockCoinbase() public view returns(address) {
        // 获取当前块矿工的地址
        return block.coinbase;
    }

    function getBlockDifficulty() public view returns(uint) {
        // 获取当前块的难度
        return block.difficulty;
    }

    function getGaslimit() public view returns(uint) {
        // 获取当前区块 gas 限额
        return block.gaslimit;
    }

    function getBlockNumber() public view returns(uint) {
        // 获取当前区块的块高
        return block.number;
    }

    function getBlockTimestamp() public view returns(uint) {
        // 获取当前块的Unix时间戳
        return block.timestamp;
    }

    function getData() public view returns(bytes memory) {
        // 获取完整的 calldata
        return msg.data;
    }

     function getGasleft() public view returns(uint) {
         // 获取当前还剩的gas
         return gasleft();
     }

    function getSender() public view returns(address) {
        // 获取当前调用发起人的地址
        return msg.sender;
    }

    function getSig() public view returns(bytes4) {
        // 调用数据的前四个字节
        return msg.sig;
    }

    function getValue() public payable returns(uint) {
        // 获取这个消息所附带的以太币，单位为wei
        return msg.value;
    }

    function getNow() public view returns(uint) {
        // 获取当前块的时间戳
        return now;
    }

    function getGasprice() public view returns(uint) {
        // 获取交易的gas价格
        return tx.gasprice;
    }

    function getOrigin() public view returns(address) {
        // 获取交易发起者（完全的调用链）
        return tx.origin;
    }
}