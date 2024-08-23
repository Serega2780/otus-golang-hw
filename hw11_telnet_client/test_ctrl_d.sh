#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=2s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 2

rm -f go-telnet
echo "PASS"
