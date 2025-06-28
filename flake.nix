# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: CC0-1.0

{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    da-flake = {
      url = "github:brokedaear/daFLAKE";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    pre-commit-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
      flake-utils,
      pre-commit-hooks,
      da-flake,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;

        # This script checks if `go-tidy` outputs something. It fails if it does.
        go-tidy-check-name = "go-tidy-check";
        go-tidy-check-script = da-flake.lib.${system}.mkScriptFromContent {
          name = go-tidy-check-name;
          content = ''
            #!/bin/bash
            if [ -n "$(go mod tidy)" ]; then
                echo "Looks like go mod tidy needs to be tidy'd up..."
                exit 1
            fi
            exit 0
          '';
        };

        ci-backend-lint-cloc-name = "ci-backend-lint-cloc";
        ci-backend-lint-cloc-script = da-flake.lib.${system}.mkScript {
          name = ci-backend-lint-cloc-name;
          scriptPath = ./app/backend/scripts/ci-backend-lint-cloc.sh;
        };

        ci-backend-test-script-name = "ci-backend-test";
        ci-backend-test-script = da-flake.lib.${system}.mkScript {
          name = ci-backend-test-script-name;
          scriptPath = ./app/backend/scripts/ci-backend-test.sh;
        };

        ci-backend-benchmark-name = "ci-backend-benchmark";
        ci-backend-benchmark-script = da-flake.lib.${system}.mkScript {
          name = ci-backend-benchmark-name;
          scriptPath = ./app/backend/scripts/ci-backend-benchmark.sh;
        };

        ciPackagesBackend =
          with pkgs;
          [
            go # Need that obviously
            gofumpt # Go formatter
            golangci-lint # Local/CI linter
            gotestsum # Pretty tester
            goperf # Go performance suite
            buf # protobuf linter/formatter
            protobuf # protobuf utils
            protoc-gen-go
          ]
          ++ da-flake.lib.${system}.ciPackages;

        ciPackagesFrontend =
          with pkgs;
          [
            nodejs_22
            yarn-berry
            typescript
            prettierd
          ]
          ++ da-flake.lib.${system}.ciPackages;

        devPackages =
          with pkgs;
          [
            gopls
            gotools
            stripe-cli # Stripe integration
            protoc-gen-go
          ]
          ++ da-flake.lib.${system}.devPackages;
      in
      {
        formatter = treefmtEval.config.build.wrapper;

        checks = {
          # Throws an error if any of the source files are not correctly formatted
          # when you run `nix flake check --print-build-logs`. Useful for CI
          treefmt = treefmtEval.config.build.check self;
          pre-commit-check = pre-commit-hooks.lib.${system}.run {
            src = ./.;
            hooks = {
              format = {
                enable = true;
                name = "Format with treefmt";
                pass_filenames = false;
                entry = "${treefmtEval.config.build.wrapper}/bin/treefmt";
                stages = [
                  "pre-commit"
                  "pre-push"
                ];
              };
              go-tidy-check = {
                enable = true;
                name = "`go mod tidy` check";
                entry = "${go-tidy-check-name}";
                stages = [ "pre-push" ];
              };
              lint-go = {
                enable = true;
                name = "Lint Go files";
                entry = "golangci-lint run --config ./.golangci.yml";
                pass_filenames = false;
                types = [ "go" ];
                stages = [ "pre-push" ];
              };
              lint-reuse = {
                enable = true;
                name = "Lint licenses using reuse";
                pass_filenames = false;
                entry = "reuse lint";
                stages = [ "pre-push" ];
              };
              unit-tests-go = {
                enable = true;
                name = "Run unit tests";
                entry = "gotestsum --format testdox ./...";
                pass_filenames = false;
                stages = [ "pre-push" ];
              };
            };
          };

          # ci =
          #   pkgs.runCommand "ci-runner"
          #     {
          #       nativeBuildInputs = [ ci-script ] ++ ciPackages; # Include ci-script and ciPackages
          #       inherit (self.checks.${system}.pre-commit-check) shellHook; # Reuse pre-commit shellHook if needed
          #     }
          #     ''
          #       echo "Running CI script with essential packages..."
          #       # Ensure the ci-script is executable and in PATH for this check
          #       export PATH=$out/bin:$PATH
          #       run-ci
          #       touch $out # Ensure the output path is created
          #     '';
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs =
              [
                go-tidy-check-script
                ci-backend-lint-cloc-script
                ci-backend-test-script
                ci-backend-benchmark-script
              ]
              ++ ciPackagesBackend
              ++ ciPackagesFrontend
              ++ devPackages
              ++ self.checks.${system}.pre-commit-check.enabledPackages;

            inherit (da-flake.lib.${system}.envVars) REUSE_COPYRIGHT REUSE_LICENSE;

            shellHook = ''
              ${self.checks.${system}.pre-commit-check.shellHook}
              # eval "$(starship init bash)"
              export PS1='$(printf "\033[01;34m(nix) \033[00m\033[01;32m[%s] \033[01;33m\033[00m$\033[00m " "\W")'
            '';
          };

          ci-frontend = pkgs.mkShell {
            buildInputs = [
              # TODO: Add frontend CI scripts here
            ] ++ ciPackagesFrontend;
            CI = true;
            shellHook = ''
              echo "Entering CI frontend shell. Only essential CI tools available."
            '';
          };
          ci-backend = pkgs.mkShell {
            buildInputs = [
              go-tidy-check-script
              ci-backend-lint-cloc-script
              ci-backend-test-script
              ci-backend-benchmark-script
            ] ++ ciPackagesBackend;
            CI = true;
            shellHook = ''
              echo "Entering CI backend shell. Only essential CI tools available."
            '';
          };
        };
      }
    );
}
