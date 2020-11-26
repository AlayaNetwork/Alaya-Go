package network.platon.contracts.evm.v0_6_12;

import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class CalledContract extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b50607080601d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806335b09a6e14602d575b600080fd5b60336035565b005b600080fdfea264697066735822122030a208cc7ebfe1976bdde8bce23c1cd6824546d25d55938f12c2e3b921bac9c164736f6c634300060c0033";

    public static final String FUNC_SOMEFUNCTION = "someFunction";

    protected CalledContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CalledContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public void someFunction() {
        throw new RuntimeException("cannot call constant function with void return type");
    }

    public static RemoteCall<CalledContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CalledContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CalledContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CalledContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CalledContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CalledContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CalledContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CalledContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
