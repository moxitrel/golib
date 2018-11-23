#!/bin/sh
#
# go tool cover
#   -func <profile-out> # output coverage profile information for each function
#   -html <profile-out> # generate HTML representation of coverage profile
#
#   -o <outfile>        # save as, default stdout
#
#
# go tool pprof <format> [options] [binary] <source> ...
#   -text            # Outputs top entries in text form
#   -web             # Visualize graph through web browser
#
#   -nodecount       # Max number of nodes to show
#

# repository root
cd $(dirname $0)

# When invoking the generated test binary (the result of 'go test -c') directly, the 'test.' prefix is mandatory.
cmd=" go test               "   # 默认测试 pkg in current directory
cmd=" $cmd -run ^$          "   # ignore tests, run benchmarks only
#cmd=" $cmd -v               "   # print Log() msg

# By default, no benchmarks are run
# -bench . :  run all benchmarks
#cmd=" $cmd -bench                   ${2:-.}                 "   # run benchmark matching $BENCH_REGEX
#cmd=" $cmd -benchtime               $BENCH_TIME             "   # run each benchmark for $BENCH_TIME:1s seconds, e.g. -benchtime 1h30s
cmd=" $cmd -benchmem                                        "   # print memory allocation state for benchmarks

# Enable more than one -xxxprofile may skew each other
#cmd=" $cmd -cpuprofile              cpu.prof                "   # profile cpu
cmd=" $cmd -memprofile              mem.prof                 "   # profile memory
#cmd=" $cmd -blockprofile            ${OUT:-block.prof}       "   # profile goroutine blocking
#cmd=" $cmd -mutexprofile            $OUT                    "   # profile mutex
#cmd=" $cmd -trace                   $OUT                    "   # execution trace
cmd=" $cmd -memprofilerate          ${MEM_PROFILE_RATE:-1}      "   # collect allocation >= $MEM_PROFILE_RATE,        0: disable , 1: profile all memory allocations
cmd=" $cmd -blockprofilerate        ${BLOCK_PROFILE_RATE:-1}    "   # collect blocking   >= $BLOCK_PROFILE_RATE ns, <=0: turn off, 1: profile all blocking event
cmd=" $cmd -mutexprofilefraction    ${MUTEX_PROFILE_FRACTION:-1}"   # 1/n events are reported,                        0: turn off, 1: profile all blocking event
cmd=" $cmd -outputdir               /tmp                    "   # dir for saving profiling data, default current directory
cmd=" $cmd -o                       /tmp/benchmark          "   # where to save the compiled binary

# run each test and benchmark $COUNT:1 times
#cmd=" $cmd -count   $COUNT "
# run no longer than $TIMEOUT:10m
# 0: the timeout is disabled
# the default is 10 minutes (10m)
#cmd=" $cmd -timeout $TIMEOUT "

$cmd $*
