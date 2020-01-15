package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple6;
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
public class SameNameConstructorPublicVisibility extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b60405160208061015b833981016040528080519060200190919050508060008190555050610119806100426000396000f300606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634c9b47a614604b575b3415604957600080fd5b005b3415605557600080fd5b605b6094565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b60008060008060008060008060008060006301e133809450680dd2d5fcf3bc9c0000935060ff925060ff9150680dd2d5fcf3bc9c0000905060005485858585859a509a509a509a509a509a5050505050509091929394955600a165627a7a72305820896255a4d81ff7788df46c5eef6355b6a8a7f9784348ac868c30ca29c3c0e26e0029";

    public static final String FUNC_DISCARDLITERALSANDSUFFIXES = "discardLiteralsAndSuffixes";

    @Deprecated
    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>> discardLiteralsAndSuffixes() {
        final Function function = new Function(FUNC_DISCARDLITERALSANDSUFFIXES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (BigInteger) results.get(3).getValue(), 
                                (BigInteger) results.get(4).getValue(), 
                                (BigInteger) results.get(5).getValue());
                    }
                });
    }

    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
