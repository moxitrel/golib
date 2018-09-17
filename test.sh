#!/bin/sh
#
# go test       		//test pkg in current directory
# go test <pkg> 		//test pkg
#	-test.v      	    //print Log() msg
#	-test.benchmem      //print memory allocation state
#
#	-run <regex> 	    //test functions match <regex>
#	-bench <regexp>     //run benchmark matching <regexp>
#	-bench .            //run all benchmarks
#
#	-cpuprofile   <out> //profile cpu
#	-memprofile   <out> //profile memory
#	-blockprofile <out> //profile goroutine blocking
#	-mutexprofile <out> //profile mutex
#	-coverprofile <out> //profile coverage
#	-trace        <out> //execution trace
#

PKG="./..."
TEST_REGEX=""
BENCH_REGEX=""


# repository root
cd $(dirname $0)

cmd=" go test "                                     # 默认测试 pkg in current directory
#cmd=" $cmd -test.v "                                # print Log() msg
#cmd=" $cmd -test.benchmem "                         # print memory allocation state
cmd=" $cmd ${TEST_REGEX:+  -run   $TEST_REGEX} "    # run test      matching $TEST_REGEX
cmd=" $cmd ${BENCH_REGEX:+ -bench $BENCH_REGEX} "   # run benchmark matching $BENCH_REGEX
cmd=" $cmd -cover "                                 # coverage analysis, line number may be changed
$cmd $PKG
