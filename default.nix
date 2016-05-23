{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/119b7c9863e.tar.gz") {}
, mkACI ? import lib/mkACI.nix
, mkPod ? import lib/mkPod.nix
}:

let
  dnsquirks = false;
in {
  acserver = import pkgs/linux/acserver.nix { inherit pkgs mkACI; static=false; };
  bash = import pkgs/linux/bash.nix { inherit pkgs mkACI; static=false; };
  busybox = import pkgs/linux/busybox.nix {
    inherit pkgs mkACI;
    static=false;
    authors="steveeJ";
    documentation="https://github.com/steveeJ/nix2aci";
    homepage="https://github.com/steveeJ/nix2aci";
  };
  busyboxPfwd = import pkgs/linux/busybox-pfwd.nix { inherit pkgs mkACI; static=false; };
  busyboxStatic = import pkgs/linux/busybox.nix { inherit pkgs mkACI; static=true; };
  busyboxThin = import pkgs/linux/busybox.nix { inherit pkgs mkACI; static=true; thin=true; };
  coreosIpxeServer = import pkgs/linux/coreos-ipxe-server.nix { inherit pkgs mkACI; };
  dnsmasq = import pkgs/linux/dnsmasq.nix { inherit pkgs mkACI; };
  dnsfail = import pkgs/linux/dnsfail.nix { inherit pkgs mkACI; };
  docker = import pkgs/linux/docker.nix { inherit pkgs mkACI; inherit dnsquirks; };
  dovecot = import pkgs/linux/dovecot.nix { inherit pkgs mkACI; inherit dnsquirks; };
  etcd2 = import pkgs/linux/etcd2.nix { inherit pkgs mkACI; };
  getmail = import pkgs/linux/getmail.nix { inherit pkgs mkACI; inherit dnsquirks; };
  iperf = import pkgs/linux/iperf.nix { inherit pkgs mkACI; static=false; };
  pixiecore = import pkgs/linux/pixiecore.nix { inherit pkgs mkACI; static=false; };
  qemu = import pkgs/linux/qemu.nix { inherit pkgs mkACI; };
  rkt = import pkgs/linux/rkt.nix { inherit pkgs mkACI; dnsquirks=false;};
  plex = import pkgs/linux/plex.nix { inherit pkgs mkACI; dnsquirks=false;};
}
