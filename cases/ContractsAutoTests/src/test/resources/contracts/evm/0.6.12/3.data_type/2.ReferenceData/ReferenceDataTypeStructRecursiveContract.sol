pragma solidity ^0.6.12;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 *测试结构体递归
 */

contract ReferenceDataTypeStructRecursiveContract {

      //定义结构体，嵌套递归结构体数组
     struct Person {
          Person[] children;
     } 
     Person person;
     //构造函数赋值
     //length只读，无法调整数组大小
     //constructor() public {
     //    person.children.length = 2;
     //    person.children[0].children.length = 10;
     //    person.children[1].children.length = 20;
     //}
     //获取结构数组长度
     function getStructPersonLength() public view returns (uint256, uint256, uint256) {
        Person memory memoryPerson;
        memoryPerson = person;
        return(memoryPerson.children.length,
               memoryPerson.children[0].children.length,
               memoryPerson.children[1].children.length);
     }
}



