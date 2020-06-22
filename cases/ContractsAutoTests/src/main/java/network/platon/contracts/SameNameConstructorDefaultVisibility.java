package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
 * <p>Generated with web3j version 0.9.1.0-SNAPSHOT.
 */
public class SameNameConstructorDefaultVisibility extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b61026c8061001e6000396000f300606060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633b016c7e1461007d5780637f14d919146100b45780638d97752a146100eb578063ac84179514610122578063ba1ae46e14610159578063ba91daeb1461017c575b600080fd5b341561008857600080fd5b61009e60048080359060200190919050506101b3565b6040518082815260200191505060405180910390f35b34156100bf57600080fd5b6100d560048080359060200190919050506101c5565b6040518082815260200191505060405180910390f35b34156100f657600080fd5b61010c60048080359060200190919050506101d8565b6040518082815260200191505060405180910390f35b341561012d57600080fd5b61014360048080359060200190919050506101eb565b6040518082815260200191505060405180910390f35b341561016457600080fd5b61017a60048080359060200190919050506101fe565b005b341561018757600080fd5b61019d6004808035906020019091905050610208565b6040518082815260200191505060405180910390f35b60006101be8261021a565b9050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b8060008190555050565b60006102138261022d565b9050919050565b6000816000819055506000549050919050565b60008160008190555060005490509190505600a165627a7a72305820996aa1af116ffe8a11395e762e338440a36f20d0cc2deb8bf22d45bc7c987bd20029";

    public static final String FUNC_PRIVATEVISIBILITYCHECK = "privateVisibilityCheck";

    public static final String FUNC_DEFAULTVISIBILITY = "defaultVisibility";

    public static final String FUNC_PUBLICVISIBILITY = "publicVisibility";

    public static final String FUNC_EXTERNALVISIBILITY = "externalVisibility";

    public static final String FUNC_SAMENAMECONSTRUCTORVISIBILITY = "SameNameConstructorVisibility";

    public static final String FUNC_INTERNALVISIBILITYCHECK = "internalVisibilityCheck";

    @Deprecated
    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> privateVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_PRIVATEVISIBILITYCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> defaultVisibility(BigInteger param) {
        final Function function = new Function(FUNC_DEFAULTVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> publicVisibility(BigInteger param) {
        final Function function = new Function(FUNC_PUBLICVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> externalVisibility(BigInteger param) {
        final Function function = new Function(FUNC_EXTERNALVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> SameNameConstructorVisibility(BigInteger param) {
        final Function function = new Function(
                FUNC_SAMENAMECONSTRUCTORVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> internalVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_INTERNALVISIBILITYCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}