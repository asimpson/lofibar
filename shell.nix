{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.gopls
    pkgs.mpv
    pkgs.pkg-config
    pkgs.gtk3
    pkgs.libappindicator-gtk3
  ];

  shellHook = ''
    mkdir -p .go
    export GOPATH=$PWD/.go
  '';
}
