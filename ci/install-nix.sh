#!/usr/bin/env bash
set -eux

NIX_LATEST=$(wget -qO- https://nixos.org/releases/nix/latest/| grep -- -x86_64-linux.tar.bz2\" | sed 's/.*href="\([^"]*\)".*/\1/')

wget https://nixos.org/releases/nix/latest/$NIX_LATEST -qO- | tar xj -C ~/

wget http://static.proot.me/proot-x86_64 -O ~/proot-x86_64
chmod u+x ~/proot-x86_64
~/proot-x86_64 -b ~/nix-${NIX_VERSION}-x86_64-linux/:/nix bash -c "/nix/install"
