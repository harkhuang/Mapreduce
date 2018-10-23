

package main

import (
	"strings"
//	"io"
//	"hash/fnv"
//	"io/ioutil"
	"os"
//	"encoding/json"
	"fmt"
	"bufio"
	"strconv"
)
func mapF(file string, contens string) (map[string]int){
	inputFile ,inputError := os.Open(file)
	if inputError != nil{
		fmt.Println("open file error")
		return nil
	}
	defer inputFile.Close()


	retmap := make(map[string]int)

	scanner := bufio.NewScanner(inputFile)  

	// TODO: 1.换行单词统计   
	//       2.contens 整理工作环境
	for scanner.Scan(){
		var words []string = strings.Fields(scanner.Text())
		//fmt.Println(len(words))
		for _,word :=range words{
			retmap[word]+=1
		}
	}
	return retmap
}

func traveral(mapV map[string]int){
	for key,vale := range mapV{
		fmt.Println("key",key,":value",vale)
	}
}


func main(){

	var mapret map[string]int = mapF("main.go","xxxxx")
	//  S3： 将mapF返回的数据根据key分类,跟文件名对应(reduceName获取文件名)
	//var reducename string = reduceName(jobName,mapTaskNumber,nReduce)
	writeFile ,_ := os.OpenFile("out.file",os.O_WRONLY|os.O_CREATE, 0600)
	defer writeFile.Close()
	var buf string 
	for key,value := range mapret {
		buf += key
		buf += ":" 
		buf += strconv.Itoa(value) 
		buf += "\n"		
	}
	//writeFile.Write([]byte(buf))
	writeFile.WriteString(buf)
	// 　S4: 　将分类好的数据分别写入不同文件

}