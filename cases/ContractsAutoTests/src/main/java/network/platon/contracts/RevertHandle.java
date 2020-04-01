package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class RevertHandle extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610157806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063852da1631461003b578063f76051e714610069575b600080fd5b6100676004803603602081101561005157600080fd5b8101908080359060200190929190505050610097565b005b6100956004803603602081101561007f57600080fd5b8101908080359060200190929190505050610111565b005b600a81111561010e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f636865636b20636174636820657863657074696f6e000000000000000000000081525060200191505060405180910390fd5b50565b600a81111561011f57600080fd5b5056fea265627a7a7231582005dbe82e36421be973147a0bff594da8228673b4a74792c54aedd9a26d94048e64736f6c634300050d0032";

    public static final String FUNC_REVERTCHECK = "revertCheck";

    public static final String FUNC_REVERTREASONCHECK = "revertReasonCheck";

    @Deprecated
    protected RevertHandle(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected RevertHandle(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected RevertHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected RevertHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> revertCheck(BigInteger param) {
        final Function function = new Function(
                FUNC_REVERTCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> revertReasonCheck(BigInteger param) {
        final Function function = new Function(
                FUNC_REVERTREASONCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<RevertHandle> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(RevertHandle.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RevertHandle> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RevertHandle.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<RevertHandle> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(RevertHandle.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RevertHandle> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RevertHandle.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RevertHandle load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new RevertHandle(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static RevertHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new RevertHandle(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static RevertHandle load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new RevertHandle(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static RevertHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new RevertHandle(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
