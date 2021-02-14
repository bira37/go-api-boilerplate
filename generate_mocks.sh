#!/bin/bash

TYPES=( service repository )

for type in "${TYPES[@]}"
do
  FILES=$(find contract/$type -type f)
  for file in $FILES
  do
    echo "Generating Mock for" $(basename $file) "in" test/mock/$type/mock_$(basename $file)
    mockgen -source $file -destination test/mock/$type/mock_$(basename $file)
  done
done