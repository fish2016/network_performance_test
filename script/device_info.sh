#!/bin/bash

echo "========== CPU Info =========="
lscpu | grep "Model name" | awk -F ': ' '{print $2}'
lscpu | grep -E "Socket|Core|Thread"

echo -e "\n========== Memory Info =========="
awk '/MemTotal/ {printf "%.0fGB RAM\n", $2 / 1024 / 1024}' /proc/meminfo