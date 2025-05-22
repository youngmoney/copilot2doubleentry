#!/usr/bin/env bash

function run() {
	go run . --config tests/transactions.config.json $@ tests/transactions.csv
}

function missing() {
	go run . --config tests/empty.config.json tests/missing_"$1"._override.csv
}

diff <(run) <(cat tests/transactions.converted.csv)
diff <(run --firstDay=2025-01-03 --lastDay=2025-01-04) <(cat tests/transactions.converted.03-04.csv)
diff <(run --firstDay=2025-01-03 --lastDay=2025-01-03) <(cat tests/transactions.converted.03-03.csv)
diff <(run --firstDay=2025-01-03) <(cat tests/transactions.converted.03-.csv)
diff <(run --lastDay=2025-01-03) <(cat tests/transactions.converted.-03.csv)

missing expense 2>/dev/null && echo "missing expense should error"
missing income 2>/dev/null && echo "missing income should error"
missing transfer_in 2>/dev/null && echo "missing transfer in should error"
missing transfer_out 2>/dev/null && echo "missing transfer out should error"
