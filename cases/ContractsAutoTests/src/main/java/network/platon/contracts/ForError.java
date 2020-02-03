package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class ForError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063895e3ada1461003b578063ec56ae5d1461005d575b600080fd5b61004361007f565b604051808215151515815260200191505060405180910390f35b6100656100a4565b604051808215151515815260200191505060405180910390f35b60008061008a6100c9565b90508060000160009054906101000a900460ff1691505090565b6000806100af6100f5565b90508060000160009054906101000a900460ff1691505090565b60008090505b600115158160000160009054906101000a900460ff16151514156100f2576100cf565b90565b60005b600190508060000160009054906101000a900460ff1615610118576100f8565b9056fea265627a7a72315820755542c76c206552ba6bb62b010a725ba80cab7e275da129aacd5058b2978d4864736f6c634300050d0032";

    public static final String FUNC_GETFORCONTROLRES = "getForControlRes";

    public static final String FUNC_GETFORCONTROLRES1 = "getForControlRes1";

    @Deprecated
    protected ForError(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected ForError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected ForError(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected ForError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Boolean> getForControlRes() {
        final Function function = new Function(FUNC_GETFORCONTROLRES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<Boolean> getForControlRes1() {
        final Function function = new Function(FUNC_GETFORCONTROLRES1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(ForError.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ForError> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ForError.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(ForError.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ForError> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ForError.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static ForError load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new ForError(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static ForError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new ForError(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static ForError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new ForError(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static ForError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new ForError(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
