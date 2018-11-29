#!/bin/sh
#
# go test       		// test pkg in current directory
# go test <pkg> 		// test pkg
#   -outputdir <.>      // dir for saving profiling data
#
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
PKG="./..."


# repository root
cd $(dirname $0)

# When invoking the generated test binary (the result of 'go test -c') directly, the 'test.' prefix is mandatory.
cmd=" go test               "   # 默认测试 pkg in current directory
cmd="$cmd -race             "   # enable data race detection
cmd=" $cmd -vet -all        "   # run "go vet $VET_FLAGS" during test
cmd=" $cmd -failfast        "   # exit if one of tests failed

cmd=" $cmd -cover                                           "   # coverage analysis, line number may be changed
# The default is "set" unless -race is enabled, in which case it is "atomic"
#   count : how many times does this statement run?
#   atomic: count, but correct in multithreaded tests; significantly more expensive
#   set   : does this statement run?
#cmd=" $cmd -covermode               $COVER_MODE             "
#cmd=" $cmd -coverpkg                $COVER_PKG              "   # apply to pkgs match "pattern1,pattern2,pattern3", default for the package being tested
#cmd=" $cmd -coverprofile            $OUT                    "   # profile coverage

# run no longer than $TIMEOUT:10m
# 0: the timeout is disabled
# the default is 10 minutes (10m)
#cmd=" $cmd -timeout $TIMEOUT "

#cmd=" $cmd -run $TEST_REGEX "   # run test matching $TEST_REGEX
#cmd=" $cmd -count    $COUNT "   # run each test and benchmark $COUNT:1 times
#cmd=" $cmd -v               "   # print Log() msg
$cmd $PKG $*
