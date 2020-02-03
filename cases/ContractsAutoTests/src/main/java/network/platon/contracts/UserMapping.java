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
 * <p>Generated with web3j version 0.7.5.0.
 */
public class UserMapping extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101b7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806316fa21101461003b578063ff7ac36d14610073575b600080fd5b6100716004803603604081101561005157600080fd5b8101908080359060200190929190803590602001909291905050506100b5565b005b61009f6004803603602081101561008957600080fd5b8101908080359060200190929190505050610166565b6040518082815260200191505060405180910390f35b8160008083815260200190815260200160002081905550600073d982c01a2ef865f046f684a1710f5b0b7edef8836312c487069091836040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b15801561012657600080fd5b505af415801561013a573d6000803e3d6000fd5b505050506040513d602081101561015057600080fd5b8101908080519060200190929190505050505050565b600080600083815260200190815260200160002054905091905056fea265627a7a72315820b56f91c0784dbfc2630d8ed85a8264d5f6a2baff8b6766d86c17e5c7c06b68ee64736f6c634300050d0032";

    public static final String FUNC_GETOUTUSER = "getOutUser";

    public static final String FUNC_SETOUTUSER = "setOutUser";

    @Deprecated
    protected UserMapping(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected UserMapping(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected UserMapping(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected UserMapping(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getOutUser(BigInteger _id) {
        final Function function = new Function(FUNC_GETOUTUSER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_id)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> setOutUser(BigInteger _age, BigInteger _id) {
        final Function function = new Function(
                FUNC_SETOUTUSER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_age), 
                new org.web3j.abi.datatypes.generated.Uint256(_id)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<UserMapping> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(UserMapping.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<UserMapping> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(UserMapping.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<UserMapping> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(UserMapping.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<UserMapping> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(UserMapping.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static UserMapping load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new UserMapping(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static UserMapping load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new UserMapping(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static UserMapping load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new UserMapping(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static UserMapping load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new UserMapping(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
