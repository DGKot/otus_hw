#!/usr/bin/env bash
set -euo pipefail

go build -o go-telnet

(echo -e "SERVER\nDATA\n" && sleep 1) | nc -l localhost 4243 >/tmp/nc_ctrl_d.out &
NC_PID=$!

sleep 1
(echo -e "CLIENT\nDATA\n") | ./go-telnet localhost 4243 >/tmp/telnet_ctrl_d.out

wait ${NC_PID} 2>/dev/null || true

expected_nc_out='CLIENT
DATA'
expected_telnet_out='SERVER
DATA'

if [ "$(cat /tmp/nc_ctrl_d.out)" != "${expected_nc_out}" ]; then
  echo "unexpected nc output:"
  cat /tmp/nc_ctrl_d.out
  exit 1
fi

if [ "$(cat /tmp/telnet_ctrl_d.out)" != "${expected_telnet_out}" ]; then
  echo "unexpected telnet output:"
  cat /tmp/telnet_ctrl_d.out
  exit 1
fi

rm -f go-telnet
echo "PASS ctrl+d"