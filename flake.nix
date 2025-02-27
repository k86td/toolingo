{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
  let
    pkgs = import nixpkgs {
      system = "x86_64-linux";
    };
  in
  {

    devShells.x86_64-linux.default = pkgs.mkShell {
      buildInputs = with pkgs; [
        beancount
        fava
        python312
        python312Packages.beancount
        go
        gopls
      ];
    };

  };
}
