# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Unlicense

{
  # See https://github.com/numtide/treefmt-nix#supported-programs

  projectRootFile = "./flake.nix";

  settings.global.includes = [
    "*.go"
    "*.yaml"
    "*.yml"
    "*.md"
    "*.nix"
    "*.proto"
    "*.sql"
    "*.sh"
  ];

  # settings.global.fail-on-change = false;
  # settings.global.no-cache = true;

  programs.gofumpt.enable = true;
  programs.goimports.enable = true;
  programs.golines.enable = true;
  programs.buf.enable = true; # protobuf
  programs.sql-formatter = {
    enable = true;
    dialect = "postgresql";
  };

  # GitHub Actions
  programs.yamlfmt.enable = true;
  programs.actionlint.enable = true;

  programs.taplo.enable = true;

  # Markdown
  programs.mdformat.enable = true;

  # Nix
  programs.nixfmt.enable = true;

  # Shell
  programs.shfmt.enable = true;
}
