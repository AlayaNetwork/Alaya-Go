package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes32;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
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
public class Blockhash extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061017c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630f7536281461005c57806366b3eb341461007a578063696d67c3146100985780639e1f194e146100b6578063e1b99d74146100d4575b600080fd5b6100646100f2565b6040518082815260200191505060405180910390f35b610082610103565b6040518082815260200191505060405180910390f35b6100a0610115565b6040518082815260200191505060405180910390f35b6100be610127565b6040518082815260200191505060405180910390f35b6100dc610138565b6040518082815260200191505060405180910390f35b60008060ff43034090508091505090565b60008061010043034090508091505090565b60008061010143034090508091505090565b600080601e43034090508091505090565b60008043409050809150509056fea2646970667358221220541cc19616209e6259ff556339d7df23f3fd72b0dfe9072b5d07fb832a9fed9e64736f6c634300060c0033";

    public static final String FUNC_GETBLOCKHASHBEFORE0 = "getBlockhashbefore0";

    public static final String FUNC_GETBLOCKHASHBEFORE255 = "getBlockhashbefore255";

    public static final String FUNC_GETBLOCKHASHBEFORE256 = "getBlockhashbefore256";

    public static final String FUNC_GETBLOCKHASHBEFORE257 = "getBlockhashbefore257";

    public static final String FUNC_GETBLOCKHASHBEFORE30 = "getBlockhashbefore30";

    protected Blockhash(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Blockhash(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> getBlockhashbefore0() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE0, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore255() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE255, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore256() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE256, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore257() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE257, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore30() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE30, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<Blockhash> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Blockhash.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Blockhash> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Blockhash.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Blockhash load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Blockhash(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Blockhash load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Blockhash(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
