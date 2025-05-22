{
  description = "testr";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      forAllSystems = nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            name = "testr";
            shellHook = ''
              git rev-parse --is-inside-work-tree >/dev/null 2>&1 || git init
              git config pull.rebase true
              ${pkgs.neo-cowsay}/bin/cowsay -f sage "testr"
            '';
            buildInputs = with pkgs; [
              editorconfig-checker
              go
              gopls
              delve
            ];
          };
        }
      );
    };
}
