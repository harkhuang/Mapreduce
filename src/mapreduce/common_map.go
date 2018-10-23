package mapreduce

import (
	"strings"
	"io"
	"hash/fnv"
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"bufio"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// s1; 读文件
	var mapret map[string]int = mapF(inFile,"xxxxx")
	//  S3： 将mapF返回的数据根据key分类,跟文件名对应(reduceName获取文件名)
	var reducename string = reduceName(jobName,mapTaskNumber,nReduce)
	filename = reducename + “.txt”
	writeFile ,writeError := open(reducename,os.O_WRONLY|os.O_CREATE, 0666)
	defer writeFile.close
	for key,value range mapret {
		writeFile.Write(key,value)
	}
	// 　S4: 　将分类好的数据分别写入不同文件
	defer writeFile.Close()

	var buf string 
	for key,value := range mapret {
		buf += key
		buf += ":" 
		buf += strconv.Itoa(value) 
		buf += "\n"		
	}
	writeFile.WriteString(buf)

}

 
// word counts from file
func mapF(file string, contens string) (map[string]int){
	inputFile ,inputError := os.Open(file)
	if inputError != nil{
		fmt.Println("open file error")
		return nil
	}
	defer inputFile.Close()

	retmap := make(map[string]int)

	scanner := bufio.NewScanner(inputFile)  

	for scanner.Scan(){
		fmt.Println(scanner.Text())
		words := strings.Fields(scanner.Text())
		for _,word :=range words{
			retmap[word]+=1
		}
	}
	return retmap
}



// 此处为何用到hash??
func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
