package network.platon.test.evm.v0_6_12.function.paramandreturns;

import com.alaya.protocol.core.methods.response.TransactionReceipt;
import network.platon.contracts.evm.v0_6_12.FunctionParam;
import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import org.junit.Before;
import org.junit.Test;
import java.math.BigInteger;


/**
 * @title 入参是函数的使用
 * @description:
 * @author: liweic
 * @create: 2020/01/11 20:20
 **/


public class FunctionParamTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.FunctionParamTest-参数是函数的类型测试", sourcePrefix = "evm/0.6.12")
    public void Functionparam() {
        try {

            FunctionParam functionparam = FunctionParam.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = functionparam.getContractAddress();
            TransactionReceipt tx = functionparam.getTransactionReceipt().get();
            collector.logStepPass("FunctionParam deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("FunctionParam deploy gasUsed:" + functionparam.getTransactionReceipt().get().getGasUsed());

            BigInteger t = functionparam.t().send();
            collector.logStepPass("FunctionParam函数返回值：" + t);
            collector.assertEqual("7",t.toString());

        } catch (Exception e) {
            collector.logStepFail("FunctionParamContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}



