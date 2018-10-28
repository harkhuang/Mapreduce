# 实现mapreduce模型框架

### 基础功能

1. mapf     处理单个kv
2. reducef  合并处理的结果
//3.schdule  任务调度  // 在论文中并没有描述  我们理解此处是客户端的调用 mapreduce
3. worker //可以详细的执行任务的的计划,如果想提升效率容错性 细化工作的粒度和容错机制
3.1 worker  call mapfunction // 处理单个kv
3.2.worker call reducefunction // 合并kv
4. master  主要的协调管理  和 work是同时fork出来的生产节点调度
master 的级别和worker级别是一样的,都是有userapi触发fork生成的rpc调度
4.1 切分文件   
4.2 根据切分的文件分发工作任务(rpc)  n->map     r->reduce
4.3 调度不同的节点处理任务  call map worker(词频统计)  ->  call reduce worker (合并)
  此处会涉及并行计算、容错、数据分布、负载均衡等复杂的细节，这些问题都被封装在 了一个库里面,暂时不用关心
4.4 完成所有的task 和 reduce  输出统计结果


### 高级功能  //TODO
1. work节点容错   
  1).保持work和master之间的心跳
  2)保留工作现场,恢复工作现场
2. master节点容错
  1).方法1 kill 所有jobs相关的work和产生的中间文件
  2).方法2保留task 任务调度器job状态 根据中间文件逐个回复work的工作现场.方式2工程上比较麻烦不推荐实现
3. 在处理失效的处理
  1) worker map 活着 worker reduce 是多个相同任务的调度 调度结果 我默认认为第一个正确// 高效率 ,准确性下降
  2) worker map 或者worker reduce执行结果进行比较 如果正确我们 我们取最优节点结果,如果不正确我们分配节点执行任务
4. 默认调度3个worker拷贝执行同一个map 或者reduce 的k/v
5. 分发粒度  文件大小16M到64M 
假设有2000个机器  20000个map 和 5000reduce
6. 对一些超时任务的处理,木桶原理,不要因为最慢的执行影响我们的效率,所以我们使用多个节点来派发任务解决这个问题

### 技巧:
1. 分区函数
hash 对 任务的平衡分区  //一致性hash
这样来自同一个urls的文件讲被写入同一个输出文件中
2. 顺序保证
按照hash值大小顺序执行work会为我们后续存取key随机存取有较大意义
3. 中间key重读值会有很多,< the,1 >单词统计中出现n多个the我们把这些常用的key分配给combiner函数处理输出到某个文件中,这里考虑对work任务做一个详细的分配,目的是更高效的调度执行计划
4. 通过reader接口  整合输出的中间文件  适应不同类型的输入输出
这里使用点是适配器模式
5. 利用worker生成辅助临时文件,更高效的工作
6. 跳过损坏的文件,检车哪些文件出的错.在出错后利用系统信号,发送最后一条出错信息到master进程
// 涉及到分布式处理错误查询很麻烦,我们如果有个本地环境调试起来可能会方便很多
7. 本地的gdb方式模拟调试
8. 主机状态监控  利用嵌入式http服务器 如jetty显示调度信息,出错信息
执行进度等数据
9. 监控中的计数器  map执行次数  reduce执行次数 等...





### rpc 思考
rpc调度模拟分布式a
问题有几个a
注册中心谁来做
服务由谁来启动
函数调用如何找到对应的rpc服务
我怎么知道哪些对应的节点有对应的rpc服务
虽然这些在dubbo中实现是zookeeper
但是这是我的rpc服务中心,注册中初步实现可以考虑但节点

那我们分配下rpc模块吧

//初阶要求
为了多进程工作  我们不得不引入io复用 ->复用io比较难以理解,go我们就理解为多个协程公用一个
io复用概念比较抽象  
是为了我们能在最短时间内消费内核已经缓冲好的io缓冲区,内核缓冲好的io在linux中我们理解为可读可写的一种状态变化
在golang中的select有所不同 是他仅仅去读取ch的数据  读取可以 我们可以决定是否消费这个数据  是否是边缘水平触发
分布式服务注册中心有一个 (暂时不用考虑分布式节点) 



///
分布式rpc调度过程中是否能够很好的利用资源是值得我们深度思考的问题
保持注册中心分布式一致性?
服务器注册服务采用hash一致性来分配不同节点去负载不同的函数调用
rpc和webservice最大区别就是不需要中间文件 ,完全通过参数传递来实现调用
potobuf这么看 跟 webservice真的有点像   有中间文件 有方法 有输入 有输出





























####### 下面是胡说八道

调度的work是并行的
多节点模拟我们使用rpc
schedule  可以认为是main函数  这里是客户调用的处理进程  
这里的任务我们理解为基于不同节点的wc
//4.1 ->task  将schedule 分配成若干节点(rpc)
//4.2 ->worker  单独执行的一个task的节点(rpc)
map:抽象多个kv的集合  这些kv分布在不同的机器上的不同的存储结构,可大可小,原本kv处理我们按照块大小处理
这里实现忽略此模块大小主要是通过实现mapreduce理解分布式计算上核心思想
// 读论文前基于我现有知识的理解  
mapf:  处理kv数据   ??  处理单个kv
reducef:  传入若干文件名称   根据文件名称分配不同的节点做计算   
master:  任务总的调度?  那么schdule干啥的??
task 干什么的?? 
task和worker关系
接口抽象
reduce
mapF(k1,v1)
reduceF(k1,vector< v1 >)
master(ks,vs)



### 理解master 调度中心
mapreduce的分布式计算框架模型很简单
应该说基础模型很简单, 但是附加上一个条件 `每个节点都是不可靠的`
这样一来,我们每次计算的节点就要附加保存很多的额外信息,比如每次计算分配的map任务我们可能要一个切分的map单元我们要分配给3个计算节点去处理,如果计算节点出来任何故障,我们可以利用心跳来排查问题,也可以用其他方式来排查问题,在有些问题处理的粒度上可能c语言更贴近操作系统,能让我们分布式计算系统更强悍更可靠,那我们不得不重构这个简单的模型,考虑这些问题 从map开始重构,后面的reduce处理也要重构,同样输出也重构



### 关于分布式系统效率问题的考虑
这里文件分发　　考虑到　用户态到内核态的拷贝　浪费系统io　我们做了ｓｅｎｄｆｉｌｅ的操作，让内核io缓冲直接和ｔｃｐ协议栈交互　　这样我们可以更高效的做文件传输





