<!--
SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>

SPDX-License-Identifier: Apache-2.0
-->

# daKINE

daKINE is the monorepo for web infrastructure and operations at Broke da EAR!

## Tools and Dependencies

- Nix package manager

Nix is used to setup the development environment and build packages.

For your convenience, there are also some command line tools available in the
Nix environment, courtesy of [daFLAKE](https://github.com/brokedaear/daFLAKE):

- git
- fish
- ripgrep-all
- helix
- neovim
- jq
- yq

To format all files, use the command

```shell
nix fmt .
```
