{
  description = "Zhangyu Mao";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: rec {
        devenv.shells = {
          default = {
            languages = {
              go.enable = true;
              go.package = pkgs.go_1_24;
            };

            pre-commit.hooks = {
              actionlint.enable = true;
              nixpkgs-fmt.enable = true;
              yamllint.enable = true;
            };

            packages = with pkgs; [
              gosmee

              actionlint
              golangci-lint
              yamllint
            ];

            scripts = {
              versions.exec = ''
                go version
                golangci-lint version
                actionlint version
              '';
              smee-proxy.exec = ''
                smee -u "$WEBHOOK_PROXY_URL"
              '';
            };

            enterShell = ''
              versions
            '';

            # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
            containers = pkgs.lib.mkForce { };
          };

          ci = devenv.shells.default;
        };
      };
    };
}
