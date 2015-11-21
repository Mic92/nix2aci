{ mkACI
, pkgs
, thin ? false
, ... }
@ args:

let
  pkg = pkgs.dnsmasq;
in

mkACI rec {
  inherit pkgs;
  inherit thin;
  dnsquirks = false;
  packages = [ pkg ];
  versionAddon = "";
  exec = ''/bin/dnsmasq'';

  labels = {
    "os"="linux";
    "arch"="amd64";
  };

  mounts = {
    "varlibmisc" = "/var/lib/misc/";
    "varrun" = "/var/run/";
    "pxe" = "/pxe/";
  };
}

