package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint16;
import com.alaya.abi.solidity.datatypes.generated.Uint64;
import com.alaya.abi.solidity.datatypes.generated.Uint8;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class TypeConversionBytesToUintContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061019b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634e9189bc146100465780636ab28181146100685780638acc06e014610090575b600080fd5b61004e6100b1565b604051808261ffff16815260200191505060405180910390f35b6100706100ed565b604051808267ffffffffffffffff16815260200191505060405180910390f35b61009861012f565b604051808260ff16815260200191505060405180910390f35b6000807f6162636400000000000000000000000000000000000000000000000000000000905060008160e01c9050600081905080935050505090565b6000807f6162636400000000000000000000000000000000000000000000000000000000905060008160e01c905060008163ffffffff16905080935050505090565b6000807f6100000000000000000000000000000000000000000000000000000000000000905060008160f81c905080925050509056fea2646970667358221220397e3b09c09d9dc5a75dfa3db34e4f240829d5bfc2b0d5dc2948710b88db950f64736f6c63430007010033";

    public static final String FUNC_BYTESTOBIGUINT = "bytesToBigUint";

    public static final String FUNC_BYTESTOSAMEUINT = "bytesToSameUint";

    public static final String FUNC_BYTESTOSMALLUINT = "bytesToSmallUint";

    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> bytesToBigUint() {
        final Function function = new Function(FUNC_BYTESTOBIGUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint64>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> bytesToSameUint() {
        final Function function = new Function(FUNC_BYTESTOSAMEUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> bytesToSmallUint() {
        final Function function = new Function(FUNC_BYTESTOSMALLUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint16>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
