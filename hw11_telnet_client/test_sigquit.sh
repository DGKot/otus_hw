#!/usr/bin/env bash
set -euo pipefail

go build -o go-telnet

(echo -e "HELLO\n" && sleep 10) | nc -l localhost 4244 >/tmp/nc_sigint.out &
NC_PID=$!

sleep 1
(echo "PING" && sleep 10) | ./go-telnet localhost 4244 >/tmp/telnet_sigint.out &
TL_PID=$!

sleep 1
kill -quit ${TL_PID}

wait ${TL_PID} 2>/dev/null || true
wait ${NC_PID} 2>/dev/null || true

if ps -p ${TL_PID} >/dev/null 2>&1; then
  echo "telnet client did not terminate on SIGQUIET"
  exit 1
fi

rm -f go-telnet
echo "PASS sigquiet"