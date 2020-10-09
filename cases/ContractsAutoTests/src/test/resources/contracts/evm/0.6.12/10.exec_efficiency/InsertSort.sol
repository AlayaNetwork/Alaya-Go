pragma solidity ^0.6.12;

/**
 * EVM 插入排序算法复杂度验证
 **/

contract InsertSort{

    int[] result_arr;

    function OuputArrays(int[] memory arr, uint n) public payable{
        uint i;
        uint k;
        uint j;
        for(i=1;i<n;i++)
        {
            uint k;
            int temp=arr[i];
            j=i;
            while(j>=1 && temp<arr[j-1])
            {
                arr[j]=arr[j-1];
                j--;
            }
            arr[j]=temp;
        }

        result_arr = arr;
    }

    function get_arr() public view returns(int[] memory){
        return result_arr;
    }

}