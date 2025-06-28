#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

printf "\n\n"
figlet -f chunky Broke
figlet -f chunky da
figlet -f chunky LINT

# This will fail if go mod changes
echo "Ensuring no go mod updates..."
[ -n "$(go mod tidy)" ] && exit 1

echo "Linting Go files..."
golangci-lint run --fix --config ./.golangci.yml --allow-serial-runners

printf "\n\n"

echo "Linting Protobuf files..."
protoc -I . --include_source_info "$(find . -name '*.proto')" -o /dev/stdout | buf lint -

printf "\n\n"

echo "Linting Licenses..."
reuse lint

printf "\n\n"
figlet -f chunky CLOC
tokei .

printf "\n\n"
figlet -f cricket allPau! | dotacat
