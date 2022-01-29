#include "libp99stats.h"

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>
#include <unistd.h>

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
