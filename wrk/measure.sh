#! /usr/bin/env bash
THREADS=4
CONNECTIONS=400
DURATION_SECONDS=10

SUBJECT=$1


TSK_SRV="taskset -c 0,1,2,3"
TSK_LOAD="taskset -c 4,5,6,7"

if [ "$SUBJECT" = "" ] ; then
    echo "usage: $0 subject    # subject: zig or go"
    exit 1
fi

if [ "$SUBJECT" = "zig-zap" ] ; then
    zig build -Doptimize=ReleaseFast wrk > /dev/null
    $TSK_SRV ./zig-out/bin/wrk &
    PID=$!
    URL=http://127.0.0.1:3000
fi

if [ "$SUBJECT" = "zigstd" ] ; then
    zig build -Doptimize=ReleaseFast wrk_zigstd > /dev/null
    $TSK_SRV ./zig-out/bin/wrk_zigstd &
    PID=$!
    URL=http://127.0.0.1:3000
fi

if [ "$SUBJECT" = "go" ] ; then
    cd wrk/go && go build main.go 
    $TSK_SRV ./main &
    PID=$!
    URL=http://127.0.0.1:8090/hello
fi

if [ "$SUBJECT" = "go-fasthttp" ] ; then
    cd wrk/go-fasthttp && go mod tidy && go build main.go
    $TSK_SRV ./main &
    PID=$!
    URL=http://127.0.0.1:8090/hello
fi

if [ "$SUBJECT" = "go-uringnet" ] ; then
    cd wrk/go-uringnet && go mod tidy && go build main.go
    $TSK_SRV ./main 127.0.0.1:8090 &
    PID=$!
    URL=http://127.0.0.1:8090
fi

if [ "$SUBJECT" = "go-gnet" ] ; then
    cd wrk/go-gnet && go mod tidy && go build main.go
    $TSK_SRV ./main &
    PID=$!
    URL=http://127.0.0.1:8090/hello
fi

if [ "$SUBJECT" = "go-nbio" ] ; then
    cd wrk/go-nbio && go mod tidy && go build main.go
    $TSK_SRV ./main &
    PID=$!
    URL=http://127.0.0.1:8090/echo
fi

if [ "$SUBJECT" = "python" ] ; then
    $TSK_SRV python wrk/python/main.py &
    PID=$!
    URL=http://127.0.0.1:8080
fi

if [ "$SUBJECT" = "python-sanic" ] ; then
    $TSK_SRV python wrk/sanic/sanic-app.py &
    PID=$!
    URL=http://127.0.0.1:8000
fi

if [ "$SUBJECT" = "rust-bythebook" ] ; then
    cd wrk/rust/bythebook && cargo build --release
    $TSK_SRV ./target/release/hello &
    PID=$!
    URL=http://127.0.0.1:7878
fi

if [ "$SUBJECT" = "rust-bythebook-improved" ] ; then
    cd wrk/rust/bythebook-improved && cargo build --release
    $TSK_SRV ./target/release/hello &
    PID=$!
    URL=http://127.0.0.1:7878
fi


if [ "$SUBJECT" = "rust-clean" ] ; then
    cd wrk/rust/clean && cargo build --release
    $TSK_SRV ./target/release/hello &
    PID=$!
    URL=http://127.0.0.1:7878
fi

if [ "$SUBJECT" = "rust-axum" ] ; then
    cd wrk/axum/hello-axum && cargo build --release
    $TSK_SRV ./target/release/hello-axum &
    PID=$!
    URL=http://127.0.0.1:3000
fi

if [ "$SUBJECT" = "csharp" ] ; then
    cd wrk/csharp && dotnet publish csharp.csproj -o ./out
    $TSK_SRV ./out/csharp --urls "http://127.0.0.1:5026" &
    PID=$!
    URL=http://127.0.0.1:5026
fi

if [ "$SUBJECT" = "cpp-beast" ] ; then
    cd wrk/cpp && zig build -Doptimize=ReleaseFast
    $TSK_SRV ./zig-out/bin/cpp-beast 127.0.0.1 8070 . &
    PID=$!
    URL=http://127.0.0.1:8070
fi

sleep 1
echo "========================================================================"
echo "                          $SUBJECT"
echo "========================================================================"
$TSK_LOAD wrk -c $CONNECTIONS -t $THREADS -d $DURATION_SECONDS --latency $URL 

kill $PID

