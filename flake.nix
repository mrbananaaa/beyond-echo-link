{
  description = "Beyond echo link development shell (nix shell)";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import (inputs.nixpkgs) {inherit system;};
      in {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # golang
            go
            air
            sqlc
            goose
            gotestfmt

            # node
            nodejs
            corepack
            eslint_d
            prettierd
            tailwindcss-language-server
            nodePackages.vscode-langservers-extracted

            # utils
            jq
            bruno
            lazygit
          ];

          shellHook = ''
            export GOBIN=$HOME/go/bin
            export PATH=$GOBIN:$PATH

            echo "Beyond echo link dev shell activated! Happy coding 🚀."
          '';
        };
      }
    );
}
