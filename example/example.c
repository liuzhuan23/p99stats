#include "libp99stats.h"

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>
#include <unistd.h>

int main()
{
	int reqCount = 10;
	int jobCount = 2;
	P99Init(reqCount, jobCount);

	long long int rTime = 0;

	rTime = P99BeginTrans();
	//printf("%llu\n", rTime);
    usleep(100);
	P99EndTrans(0, rTime, 1);

	rTime = P99BeginTrans();
    usleep(100);
	P99EndTrans(0, rTime, 1);

	rTime = P99BeginTrans();
    usleep(100);
	P99EndTrans(1, rTime, 1);

	rTime = P99BeginTrans();
    usleep(100);
	P99EndTrans(1, rTime, 1);

	P99Stats();
	return 0;
}
