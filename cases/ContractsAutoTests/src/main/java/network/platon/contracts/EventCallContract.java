package network.platon.contracts;

import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.*;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;
import rx.Observable;
import rx.functions.Func1;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class EventCallContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610316806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806336aacafa146100515780637b0cb8391461005b5780639d6d4cde14610079578063b7301a1f14610097575b600080fd5b6100596100b5565b005b6100636101f8565b6040518082815260200191505060405180910390f35b610081610266565b6040518082815260200191505060405180910390f35b61009f610287565b6040518082815260200191505060405180910390f35b7f9f252e5d94c6346e0073dfdaa81c6bba97bc07b05f8378efc62d77d157e1b0116000604051808215151515815260200191505060405180910390a17f9f252e5d94c6346e0073dfdaa81c6bba97bc07b05f8378efc62d77d157e1b0116001604051808215151515815260200191505060405180910390a17ffc3a67c9f0b5967ae4041ed898b05ec1fa49d2a3c22336247201d71be6f9712033604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1600c6040518082815260200191505060405180910390a03373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c600c6040518082815260200191505060405180910390a2565b60007ffc3a67c9f0b5967ae4041ed898b05ec1fa49d2a3c22336247201d71be6f9712033604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a160018101905090565b6000600181019050806040518082815260200191505060405180910390a090565b60003373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c600c6040518082815260200191505060405180910390a26001810190509056fea265627a7a72315820b4783d18b79d8a37b9ff236102cc8b05cb46bea1a7d84b8fb9dfbdf72026a95464736f6c634300050d0032";

    public static final String FUNC_ANONYMOUSEVENT = "anonymousEvent";

    public static final String FUNC_EMITEVENT = "emitEvent";

    public static final String FUNC_INDEXEDEVENT = "indexedEvent";

    public static final String FUNC_TESTBOOL = "testBool";

    public static final Event ANONYMOUS_EVENT = new Event("Anonymous", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
    ;

    public static final Event BOOLEVENT_EVENT = new Event("BoolEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
    ;

    public static final Event DEPOSIT_EVENT = new Event("Deposit", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event INCREMENT_EVENT = new Event("Increment", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
    ;

    @Deprecated
    protected EventCallContract(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected EventCallContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected EventCallContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected EventCallContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public List<AnonymousEventResponse> getAnonymousEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ANONYMOUS_EVENT, transactionReceipt);
        ArrayList<AnonymousEventResponse> responses = new ArrayList<AnonymousEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            AnonymousEventResponse typedResponse = new AnonymousEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._id = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<AnonymousEventResponse> anonymousEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, AnonymousEventResponse>() {
            @Override
            public AnonymousEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ANONYMOUS_EVENT, log);
                AnonymousEventResponse typedResponse = new AnonymousEventResponse();
                typedResponse.log = log;
                typedResponse._id = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<AnonymousEventResponse> anonymousEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ANONYMOUS_EVENT));
        return anonymousEventObservable(filter);
    }

    public List<BoolEventEventResponse> getBoolEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(BOOLEVENT_EVENT, transactionReceipt);
        ArrayList<BoolEventEventResponse> responses = new ArrayList<BoolEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            BoolEventEventResponse typedResponse = new BoolEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.result = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<BoolEventEventResponse> boolEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, BoolEventEventResponse>() {
            @Override
            public BoolEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(BOOLEVENT_EVENT, log);
                BoolEventEventResponse typedResponse = new BoolEventEventResponse();
                typedResponse.log = log;
                typedResponse.result = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<BoolEventEventResponse> boolEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(BOOLEVENT_EVENT));
        return boolEventEventObservable(filter);
    }

    public List<DepositEventResponse> getDepositEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(DEPOSIT_EVENT, transactionReceipt);
        ArrayList<DepositEventResponse> responses = new ArrayList<DepositEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            DepositEventResponse typedResponse = new DepositEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._from = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<DepositEventResponse> depositEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, DepositEventResponse>() {
            @Override
            public DepositEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(DEPOSIT_EVENT, log);
                DepositEventResponse typedResponse = new DepositEventResponse();
                typedResponse.log = log;
                typedResponse._from = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<DepositEventResponse> depositEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(DEPOSIT_EVENT));
        return depositEventObservable(filter);
    }

    public List<IncrementEventResponse> getIncrementEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(INCREMENT_EVENT, transactionReceipt);
        ArrayList<IncrementEventResponse> responses = new ArrayList<IncrementEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            IncrementEventResponse typedResponse = new IncrementEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.who = (String) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<IncrementEventResponse> incrementEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, IncrementEventResponse>() {
            @Override
            public IncrementEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(INCREMENT_EVENT, log);
                IncrementEventResponse typedResponse = new IncrementEventResponse();
                typedResponse.log = log;
                typedResponse.who = (String) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<IncrementEventResponse> incrementEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(INCREMENT_EVENT));
        return incrementEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> anonymousEvent() {
        final Function function = new Function(
                FUNC_ANONYMOUSEVENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> emitEvent() {
        final Function function = new Function(
                FUNC_EMITEVENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> indexedEvent() {
        final Function function = new Function(
                FUNC_INDEXEDEVENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testBool() {
        final Function function = new Function(
                FUNC_TESTBOOL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<EventCallContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(EventCallContract.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<EventCallContract> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(EventCallContract.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<EventCallContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(EventCallContract.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<EventCallContract> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(EventCallContract.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static EventCallContract load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new EventCallContract(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static EventCallContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new EventCallContract(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static EventCallContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new EventCallContract(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static EventCallContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new EventCallContract(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static class AnonymousEventResponse {
        public Log log;

        public BigInteger _id;
    }

    public static class BoolEventEventResponse {
        public Log log;

        public Boolean result;
    }

    public static class DepositEventResponse {
        public Log log;

        public String _from;

        public BigInteger _value;
    }

    public static class IncrementEventResponse {
        public Log log;

        public String who;
    }
}
