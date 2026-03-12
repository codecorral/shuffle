{
  description = "Shuffle – declarative configuration for agent-deck";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: {
      packages.default = nixpkgs.legacyPackages.${system}.buildGoModule {
        pname = "shuffle";
        version = "0.1.0";
        src = ./.;
        vendorHash = "sha256-2SXAu1fxiRbuMOKOoB8OVzTmtR3Os423j80En+SHnzU=";
      };
    });
}
