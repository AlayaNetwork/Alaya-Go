package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Bool;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes2;
import com.alaya.abi.solidity.datatypes.generated.Bytes3;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple5;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class AssemblyReturns extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061016d806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806326121ff014610030575b600080fd5b6100386100c3565b60405180868152602001857dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001847cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200183151581526020018273ffffffffffffffffffffffffffffffffffffffff1681526020019550505050505060405180910390f35b6000806000806000600294507fabcd00000000000000000000000000000000000000000000000000000000000093507f61626300000000000000000000000000000000000000000000000000000000009250600191507372ad2b713faa14c2c4cd2d7affe5d8f538968f5a9050909192939456fea264697066735822122029f20aee3ae4efe3136551ee902feb7a0b4f08aa33f0b9582fc15650c1eee53d64736f6c63430007010033";

    public static final String FUNC_F = "f";

    protected AssemblyReturns(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AssemblyReturns(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple5<BigInteger, byte[], byte[], Boolean, String>> f() {
        final Function function = new Function(FUNC_F, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Bytes2>() {}, new TypeReference<Bytes3>() {}, new TypeReference<Bool>() {}, new TypeReference<Address>() {}));
        return new RemoteCall<Tuple5<BigInteger, byte[], byte[], Boolean, String>>(
                new Callable<Tuple5<BigInteger, byte[], byte[], Boolean, String>>() {
                    @Override
                    public Tuple5<BigInteger, byte[], byte[], Boolean, String> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple5<BigInteger, byte[], byte[], Boolean, String>(
                                (BigInteger) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue(), 
                                (byte[]) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue(), 
                                (String) results.get(4).getValue());
                    }
                });
    }

    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssemblyReturns.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssemblyReturns.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AssemblyReturns load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AssemblyReturns(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AssemblyReturns load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AssemblyReturns(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
