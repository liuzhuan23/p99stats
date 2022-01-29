#### p99统计函数
```
部分代码摘自rpcx原作者smallnest的benchmark，附原文连接

https://github.com/rpcxio/rpcx-benchmark.git

p99stats这个代码用Go和c混编，对外提供c接口如下

P99初始化
extern void P99Init(int c_totalRequests, int c_jobCount);

P99获取平均事务数
extern int P99AvgTs();

P99在每一个事务启用前显示的调用该函数
extern long long int P99BeginTrans();

P99在每一个事务结束后显示的调用该函数
extern void P99EndTrans(int c_jobIndex, long long int c_transTimer, int c_transCode);

P99统计报告输出
extern void P99Stats();
```

#### 一个使用例子
```
void waitns(int ns)
{
	struct timeval delay;
	delay.tv_sec = 0;
	delay.tv_usec = ns;	//ns
	select(0, NULL, NULL, NULL, &delay);
}

int main()
{
	//初始化P99，总请求数，总计调度作业数量
	int reqCount = 15;
	int jobCount = 3;
	P99Init(reqCount, jobCount);

	//平均事务数
	int avg = P99AvgTs();

	//循环是job，在其内模拟每个job进行若干次平均事务
	long long int rTime = 0;
	for (int i = 0; i < jobCount; i++) {
		for (int j = 0; j < avg; j++) {
			rTime = P99BeginTrans();	//启用事务计时
			//printf("%llu\n", rTime);
			waitns(j+1);
			P99EndTrans(i, rTime, 1);	//统计事务时间差
		}
	}

	//汇总报告进行计算输出
	P99Stats();
	return 0;
}
```

#### 编译环境
```
1 首先在p99stats目录内执行make，go环境可以是1.17，14和13并未做测试
2 SO目录内含有编译后的俩文件，libp99stats.a  libp99stats.h
3 执行编译链接libp99stats示例程序即可
```
