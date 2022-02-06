package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
*/
import "C"

import (
	"github.com/montanaflynn/stats"
	"github.com/smallnest/rpcx/log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
}

var (
	startTime     int64     //调用进程事务开始时间
	totalRequests int       //总事务数
	jobCount      int       //进行事务的job总数
	avgRequests   int       //每个job的平均事务数量
	tookTimes     [][]int64 //每个事务的调用时间
	trans         uint64    //事务完成一次（无论失败成功）
	transOK       uint64    //成功事务数
	muX           sync.Mutex
)

//export P99Init
func P99Init(c_totalRequests C.int, c_jobCount C.int) {
	startTime = time.Now().UnixNano()                 //初始化进程事务开始时间
	totalRequests = int(c_totalRequests)              //配置总事务数
	jobCount = int(c_jobCount)                        //作业总数
	avgRequests = totalRequests / jobCount            //平均事务数
	tookTimes = make([][]int64, jobCount-1, jobCount) //初始化这个切片，它用来记录每个trans的时间差
	dt := make([]int64, 0, avgRequests)
	tookTimes = append(tookTimes, dt) //会重新构造切片结构
}

//export P99AvgTs
func P99AvgTs() C.int {
	return C.int(avgRequests)
}

//export P99BeginTrans
func P99BeginTrans() C.longlong {
	return C.longlong(time.Now().UnixNano())
}

//export P99EndTrans
func P99EndTrans(c_jobIndex C.int, c_transTimer C.longlong, c_transCode C.int) {
	muX.Lock()
	defer muX.Unlock()
	t := time.Now().UnixNano() - int64(c_transTimer)
	idx := int(c_jobIndex)
	tookTimes[idx] = append(tookTimes[idx], t)

	if c_transCode > 0 {
		atomic.AddUint64(&transOK, 1)
	}
	atomic.AddUint64(&trans, 1)
}

//export P99Stats
func P99Stats(c_putDiff C.int) {
	if c_putDiff == 1 {
		for n, v := range tookTimes {
			log.Infof("trans unit diff %v : %v", n, v)
		}
	}
	Stats(startTime, totalRequests, tookTimes, trans, transOK)
}

// Stats用作统计结果出口
func Stats(startTime int64, totalRequests int, tookTimes [][]int64, trans, transOK uint64) {
	// 测试总耗时
	totalTInNano := time.Now().UnixNano() - startTime
	totalT := totalTInNano / 1000000
	log.Infof("took %d ms for %d requests", totalT, totalRequests)

	// 汇总每个请求的耗时
	totalD := make([]int64, 0, totalRequests)
	for _, k := range tookTimes {
		totalD = append(totalD, k...)
	}
	// 将int64数组转换成float64数组，以便分析
	totalD2 := make([]float64, 0, totalRequests)
	for _, k := range totalD {
		totalD2 = append(totalD2, float64(k))
	}

	// 计算各个指标
	mean, _ := stats.Mean(totalD2)
	median, _ := stats.Median(totalD2)
	max, _ := stats.Max(totalD2)
	min, _ := stats.Min(totalD2)
	p999, _ := stats.Percentile(totalD2, 99.9)

	// 输出结果
	log.Infof("sent     requests    : %d\n", totalRequests)
	log.Infof("received requests    : %d\n", trans)
	log.Infof("received requests_OK : %d\n", transOK)
	if totalT == 0 {
		log.Infof("throughput  (TPS)    : %d\n", int64(totalRequests)*1000*1000000/totalTInNano)
	} else {
		log.Infof("throughput  (TPS)    : %d\n\n", int64(totalRequests)*1000/totalT)
	}

	log.Infof("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99.9: %.f ns\n", mean, median, max, min, p999)
	log.Infof("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99.9: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p999/1000000))
}
