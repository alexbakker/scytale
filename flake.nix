{
  description = "Nix flake for scytale";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }: let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      defaultPackage.x86_64-linux =
        with pkgs; buildGoModule {
          name = "scytale";
          src = ./.;

          vendorHash = "sha256-ftLEuym/OPDHfTjSCoeNVjfVbsBHX4t/9OTWnLJmATA=";

          subPackages = [
            "./cmd/scycli"
            "./cmd/scyserver"
          ];
        };
      devShell.x86_64-linux = with pkgs; mkShell {
        buildInputs = [
          go
        ];
      };
    };
}
