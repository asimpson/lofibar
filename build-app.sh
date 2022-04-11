#!/bin/sh

mkdir -p lofibar.app/Contents/MacOS/
mkdir -p lofibar.app/Contents/Resources/

cp lofibar lofibar.app/Contents/MacOS/lofibar
cp lofibar.icns lofibar.app/Contents/Resources/lofibar.icns
cp Info.plist lofibar.app/Contents/Info.plist
