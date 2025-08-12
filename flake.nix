{
  description = "A Nix Flake for a Go development shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go and essential tools
            go
            gopls
            go-tools
            delve

            # Linting and formatting
            golangci-lint
            gofumpt
          ];

          # Environment variables for the development shell
          GOPATH = "$(mktemp -d)";
          GOCACHE = "$(mktemp -d)";

          # Run this when the shell is entered
          shellHook = ''
            echo "Entering Go development environment..."
          '';
        };
      });
}
