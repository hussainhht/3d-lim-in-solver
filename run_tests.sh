#!/bin/bash

export PATH=$PATH:/usr/local/go/bin  # adjust if needed

for i in {00..05}; do
    testfile="test/example$i.txt"
    resultfile="result/result$i.txt"

    echo -e "\n============================================================================="
    echo -e "▶️  Running: example$i.txt"
    echo -e "\n"

    if [[ -f "$testfile" ]]; then
        cat "$testfile"
    else
        echo "Result file not found: $testfile"
    fi

    echo -e "\n"

    # Run Go program and print output
    go run . "$testfile"

    # Print two new lines
    echo -e "\n"

    # Print "Right Ans:" and then contents of the result file
    echo -e "Right Ans:"
    if [[ -f "$resultfile" ]]; then
        cat "$resultfile"
    else
        echo "Result file not found: $resultfile"
    fi
    echo -e "\n"
done
