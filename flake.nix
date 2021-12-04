{
  description = "Nix flake for scytale";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";

  outputs = { self, nixpkgs }: let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      defaultPackage.x86_64-linux =
        with pkgs; buildGoModule {
          pname = "scytale";
          version = "0.0.0";
          src = ./.;

          vendorSha256 = "0c01csr9rmp4yizqnps7q1pdadsnin3hmliqgp3z0f5z56xw9lky";

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
