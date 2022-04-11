{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.gh
    pkgs.gopls
    pkgs.ffmpeg-full
  ]
  ++ pkgs.lib.optionals pkgs.stdenv.isDarwin [
    pkgs.darwin.apple_sdk.frameworks.Cocoa
    pkgs.darwin.apple_sdk.frameworks.WebKit
  ]
  ++ pkgs.lib.optionals (!pkgs.stdenv.isDarwin) [
    pkgs.gtk3
    pkgs.pkg-config
    pkgs.libappindicator-gtk3
  ];

  shellHook = ''
    mkdir -p .go
    export GOPATH=$PWD/.go
    export PATH=$PWD/.go/bin:$PATH
  '';
}
