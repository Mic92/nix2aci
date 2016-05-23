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
  exec = ["/bin/dnsmasq"];

  mounts = {
    varlibmisc.path = "/var/lib/misc/";
    varrun.path = "/var/run/";
    pxe.path = "/pxe/";
  };

  isolators = {
    "os/linux/capabilities-retain-set" = {
      set = ["CAP_NET_BIND_SERVICE" "CAP_NET_ADMIN"];
    };
  };
}

