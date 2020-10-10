package network.platon.utils;

import lombok.extern.slf4j.Slf4j;
import network.platon.autotest.utils.FileUtil;

import java.io.File;
import java.nio.file.Paths;

/**
 * @title 将ABI文件和二进制文件生成包装类
 * @author: qcxiao
 * @create: 2019/12/27 13:53
 **/
@Slf4j
public class GeneratorUtil {

    /**
     * @description:
     * @author: qcxiao
     * @create: 2019/12/14 16:34
     **/
    public static void generator(File contractFile) throws Exception {
        try {
            String binPath = contractFile.getPath();
            String abiPath = contractFile.getPath().replace(".bin", ".abi");
            String compilerVersion = contractFile.getPath().replaceAll("(.*)(0\\..\\d*\\.\\d*)(.*$)", "$2");

            String outputPath = FileUtil.pathOptimization(System.getProperty("user.dir") + "/src/main/java");
            String packagePath = "network.platon.contracts.evm.v" + compilerVersion.replaceAll("\\.", "_");
            log.info("packagePath: " + packagePath);
            String os = System.getProperty("os.name");
            String[] args = null;
            String rootPath = FileUtil.pathOptimization(Paths.get("scripts", "platon-web3j", "bin").toUri().getPath());

            if (!os.startsWith("Linux") && !os.startsWith("Mac OS")) {
                if (os.startsWith("Windows")) {
                    args = new String[]{"cmd", "/C", rootPath + "platon-web3j.bat" + " " + "solidity" + " " + "generate" + " " + binPath + " " + abiPath + " " + "-o" + " " + outputPath + " " + "-p" + " " + packagePath};
                } else {
                    log.error("Not supported operate system platform");
                }
            } else {
                args = new String[]{"/bin/bash", "-c", rootPath + "platon-web3j" + " " + "solidity" + " " + "generate" + " " + binPath + " " + abiPath + " " + "-o" + " " + outputPath + " " + "-p" + " " + packagePath};
            }

            CompileUtil.execGenerate(args);
        } catch (Exception e) {
            log.error("contract generator error:{}", e.getMessage());
            throw new Exception(e);
        }
    }

    public void generatorWasm(String contractName) throws Exception {
        try {
            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "wasm", "build").toUri().getPath());
            String binPath = filePath + contractName + ".wasm";
            String abiPath = filePath + contractName + ".abi.json";

            String outputPath = FileUtil.pathOptimization(System.getProperty("user.dir") + "/src/main/java");
            String packagePath = "network.platon.contracts.wasm";
            String rootPath = System.getProperty("user.dir");
            String[] args = new String[]{"/bin/bash", "-c", rootPath + "/scripts/platon-web3j/bin/platon-web3j" + " " + "wasm" + " " + "generate" + " " + binPath + " " + abiPath + " " + "-o" + " " + outputPath + " " + "-p" + " " + packagePath};

            CompileUtil.execGenerate(args);
        } catch (Exception e) {
            log.error("contract generator error:{}", e.getMessage());
            throw new Exception(e);
        }
    }


}