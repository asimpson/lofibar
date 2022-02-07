{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.gopls
    pkgs.ffmpeg-full
    pkgs.pkg-config
    pkgs.gtk3
    pkgs.libappindicator-gtk3
  ];

  shellHook = ''
    mkdir -p .go
    export GOPATH=$PWD/.go
    export PATH=$PWD/.go/bin:$PATH
  '';
}
