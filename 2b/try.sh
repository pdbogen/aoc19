#!/bin/sh

go build .

for i in $(seq 1 100); do
  for j in $(seq 1 100); do
    result="$(echo "$(cut -d, -f1 input),$i,$j,$(cut -d, -f4- input)" | ./2b | cut -d',' -f1)"
    echo "$i $j => $result"
    if [ "$result" = "19690720" ]; then
      exit 0
    fi
  done
done
