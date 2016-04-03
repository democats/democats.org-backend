#!/bin/bash


PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

while true; do
  ${PROJECT_DIR}"/stats"
  echo ''
  sleep 15
done
