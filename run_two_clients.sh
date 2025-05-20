#!/bin/bash

set -e

# Store PIDs of launched processes
PIDS=()

cleanup() {
  echo "Caught interrupt. Killing all processes..."
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  exit 0
}

trap cleanup SIGINT

echo "=== Building server ==="
cd server
make
cd ..

echo "=== Building client ==="
cd client
make
cd ..

echo "=== Launching server ==="
./server/xrossover-server &
PIDS+=($!)

sleep 1  # give server time to start

echo "=== Launching clients ==="
./client/xrossover-client Alice &
PIDS+=($!)
./client/xrossover-client Bob &
PIDS+=($!)

echo "=== All processes launched === (Press Ctrl+C to stop)"
wait
