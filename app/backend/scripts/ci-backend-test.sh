#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

printf "\n\n"

figlet -f chunky BrokeDa
figlet -f chunky goTESTS

echo "Running Go tests..."
# gotestsum --format testdox ./...
go test -mod=readonly ./...

printf "\n\n"

echo "    ⡴⠦⡄⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⠀⠀⠀⠀⠀⠀⠀⣀  "
echo "    ⢣⠀⠼⢦⠀⠀⣀⣀⡴⠲⡞⠁⠈⢣⡀⠀⠀⣠⠞⠁⢀⣽ "
echo "     ⢣⡀⠈⠳⡎⠁⢹⡀⠀⣇⠀⠠⡆⢣⢀⡼⠁⠀⣰⠋  "⠀
echo "    ⠀ ⠙⡄⠀⠱⡔⠘⢇⠀⢹⣀⠀⣇⣸⠏⠀⠀⢰⠁   "
echo "    ⠀⠀ ⢸⡀⠀⠱⣤⠚⣞⠛⣯⠷⠋⠁⠀⠀⢀⠆⠀   "
echo "    ⠀⠀⠀ ⣧⠀⠀⠉⠉⢉⡟⠁⠀⠀⠀⠀⢠⠏⠀⠀ Kden -- Shoots..."
echo "    ⠀⠀⠀ ⠸⡆⠀⠀⠀⢸⡀⠀⠀⠀⠀⣠⠏⠀⠀⠀   "
echo "    ⠀⠀⠀⠀ ⠙⣄⣀⣀⠀⠙⠒⠒⢀⡞⠁⠀⠀⠀⠀   "
echo "    ⠀⠀⠀⠀⠀ ⢸⡈⠀⠀⠀⠀⠀⢸⡀⠀⠀⠀⠀⠀   "
echo "    ⠀⠀⠀⠀⠀⠀ ⣷⠀⠀⠀⠀⠀⠈⠃⠀⠀⠀⠀⠀   "
