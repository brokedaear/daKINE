# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: CC0-1.0

{ pkgs, ... }:
{
  # See https://github.com/numtide/treefmt-nix#supported-programs
  projectRootFile = ".git/config";
  settings.global.includes = [
    "*.go"
    ".toml"
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
  programs.gofumpt = {
    enable = true;
    package = pkgs.gofumpt;
  };
  programs.goimports.enable = true;
  programs.golines.enable = true;
  programs.buf = {
    enable = true;
    package = pkgs.buf;
  };
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
  programs.nixfmt = {
    enable = true;
    package = pkgs.nixfmt-rfc-style;
  };
  # Shell
  programs.shfmt.enable = true;
}
