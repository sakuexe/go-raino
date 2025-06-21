{
  # using zsh, enter the shell using: nix develop --command zsh
  description = "Development flake for Go discord bot";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          # config.allowUnfree = true; # enable the use of proprietary packages
        };
      in
      {
        devShell =
          with pkgs;
          mkShell {
            buildInputs = [
              go
              sqlite
              air # hot reload of the go server

              libwebp

              # go static checking tools
              go-tools # staticcheck
              govulncheck # vulnerability checking
              gosec # security checker
            ];
            # non-secret env variables for development
            DEVELOPMENT = true;
          };
      }
    );
}
