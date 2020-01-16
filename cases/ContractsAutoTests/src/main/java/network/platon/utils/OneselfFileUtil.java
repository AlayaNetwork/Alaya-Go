package network.platon.utils;

import network.platon.autotest.utils.FileUtil;

import java.io.File;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public class OneselfFileUtil {
    private static List<String> list = new ArrayList<>();

    /**
     * @title OneselfFileUtil
     * @description 获取所有sol源文件
     * @author qcxiao
     * @updateTime 2019/12/27 14:22
     */
    public static List<String> getResourcesFile(String path, int deep) {
        // 获得指定文件对象
        File file = new File(path);
        // 获得该文件夹内的所有文件
        File[] files = file.listFiles();
        for (int i = 0; i < files.length; i++) {
            if (files[i].isFile()) {
                if (files[i].getName().substring(files[i].getName().lastIndexOf(".") + 1).equals("sol")) {
                    list.add(files[i].getPath());
                }
            } else if (files[i].isDirectory()) {
                //文件夹需要调用递归 ，深度+1
                getResourcesFile(files[i].getPath(), deep + 1);
            }
        }
        return list;
    }

    /**
     * @title OneselfFileUtil
     * @description 获取所有二进制文件，并返回文件名称列表
     * @author qcxiao
     * @updateTime 2019/12/27 14:24
     */
    public static List<String> getBinFileName() {
        List<String> files = new ArrayList<>();
        String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "build").toUri().getPath());
        File file = new File(filePath);
        File[] tempList = file.listFiles();

        for (int i = 0; i < tempList.length; i++) {
            if (tempList[i].isFile()) {
                String fileName = tempList[i].getName();
                if (fileName.substring(fileName.lastIndexOf(".") + 1).equals("bin")) {
                    fileName = fileName.substring(0, fileName.lastIndexOf("."));
                    files.add(fileName);
                }
            }
        }
        return files;
    }
}
