#### p99统计函数
```
go和c混编，主要的Stats函数用go写的，然后对外提供c接口
P99Stats(long long int c_startTime, int c_totalRequests, long long int c_tookTimes[][], int c_tr, int c_td, long long unsigned int c_trans, long long unsigned int c_transOK);
```

#### p99入参解释
```
long long int c_startTime, 

int c_totalRequests, 

long long int [][]c_tookTimes, 

int c_tr, 

int c_td, 

unsigned long long int c_trans, 

unsigned long long int c_transOK,
```

#### p99测试代码
```
```
