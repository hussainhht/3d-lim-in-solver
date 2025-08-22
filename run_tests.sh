#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🧪 Running all test cases from test/ directory..."

for file in test/*.txt; do
    echo -e "\n=============================="
    echo -e "▶️  Running: $file"
    output=$(go run . "$file" 2>&1)

    # Check for error message
    if echo "$output" | grep -q "ERROR:"; then
        echo -e "${RED}❌ Error output:${NC}"
        echo "$output"
    else
        echo -e "${GREEN}✅ Success. Output:${NC}"
        echo "$output"
    fi
done

echo -e "\n🎉 All tests finished."
