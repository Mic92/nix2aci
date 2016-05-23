{ mkACI, pkgs, thin ? false, ... } @ args:
let 
  pkg = pkgs.docker;
in

mkACI rec {
  inherit pkgs;
  inherit thin;
  dnsquirks = args.dnsquirks;

  packages = [ pkg pkgs.busybox pkgs.cacert ];

  mounts = {
    libdocker.path = "/var/lib/docker";
    rundocker.path = "/var/run/docker";
  };
}
