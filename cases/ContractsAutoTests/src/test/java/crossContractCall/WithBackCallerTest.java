package crossContractCall;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DelegatecallCallee;
import network.platon.contracts.DelegatecallCaller;
import network.platon.contracts.WithBackCallee;
import network.platon.contracts.WithBackCaller;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigDecimal;
import java.math.BigInteger;


/**
 * @title 0.5.13跨合约调用者,接收返回值
 * @description:
 * @author: hudenian
 * @create: 2019/12/28
 */
public class WithBackCallerTest extends ContractPrepareTest {

    //需要进行double的值
    private String doubleValue = "10";

    //需要前缀hello的值
    private String helloValue = "hudenian";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "WithBackCallerTest-跨合约调用者对反回值进行编码与解码")
    public void crossContractCaller() {
        try {
            //调用者合约地址
            WithBackCaller withBackCaller = WithBackCaller.deploy(web3j, transactionManager, provider).send();
            String callerContractAddress = withBackCaller.getContractAddress();
            TransactionReceipt tx = withBackCaller.getTransactionReceipt().get();
            collector.logStepPass("WithBackCaller deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());


            //被调用者合约地址
            WithBackCallee withBackCallee = WithBackCallee.deploy(web3j, transactionManager, provider).send();
            String calleeContractAddress = withBackCallee.getContractAddress();
            TransactionReceipt tx1 = withBackCallee.getTransactionReceipt().get();
            collector.logStepPass("WithBackCallee deploy successfully.contractAddress:" + calleeContractAddress + ", hash:" + tx1.getTransactionHash());

            //数值类型跨合约调用
            TransactionReceipt tx2 = withBackCaller.callDoublelTest(calleeContractAddress,new BigInteger(doubleValue)).send();
            collector.logStepPass("WithBackCaller callDoublelTest successfully hash:" + tx2.getTransactionHash());
            //获取数值类型跨合约调用的结果
            String chainDoubleValue = withBackCaller.getuintResult().send().toString();
            collector.logStepPass("获取数值类型跨合约调用的结果值为:" + chainDoubleValue);
            collector.assertEqual(new BigDecimal(doubleValue).multiply(new BigDecimal("2")),new BigDecimal(chainDoubleValue));


            //字符串类型跨合约调用
            tx2 = withBackCaller.callgetNameTest(calleeContractAddress,helloValue).send();
            collector.logStepPass("WithBackCaller callDoublelTest successfully hash:" + tx2.getTransactionHash());
            //获取字符串类型跨合约调用
            String callerStringResult = withBackCaller.getStringResult().send().toString();
            collector.logStepPass("获取字符串类型跨合约调用的结果值为:" + callerStringResult);
            collector.assertEqual(callerStringResult,"hello"+helloValue);

        } catch (Exception e) {
            collector.logStepFail("WithBackCallerTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
