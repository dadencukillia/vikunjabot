{
  description = "Isolated development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go_1_26
          docker
          docker-buildx
        ];

        shellHook = ''
          echo "❄️ Welcome back, stakeholders developer ❄️"

          export ISOLATED_DIR="$(pwd)/.envgarbage"
          mkdir -p "$ISOLATED_DIR"

          # Go isolation
          export GO_GARBAGE="$ISOLATED_DIR/go"
          export GOPATH="$GO_GARBAGE/path"
          export GOBIN="$GO_GARBAGE/bin"
          export GOCACHE="$GO_GARBAGE/cache"

          mkdir -p "$GOPATH"
          mkdir -p "$GOBIN"
          mkdir -p "$GOCACHE"

          export PATH="$GOBIN:$PATH"
        '';
      };
    };
}
