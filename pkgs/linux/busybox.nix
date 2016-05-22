{ mkACI
, pkgs
, thin ? false
, documentation ? null
, homepage ? null
, authors ? null
, static
, ... }
@ args:

let
  pkg = if static == true
    then
      (pkgs.busybox.override {
       # enableStatic + glibc is broken at the moment
       enableStatic = true;
       enableMinimal = true;
       useMusl = true;
      })
    else pkgs.busybox;
in

mkACI rec {
  inherit pkgs static thin documentation homepage authors;
  packages = [ pkg pkgs.eject ];
  versionAddon = if static == true then "-static" else "";

  exec = [
    "/bin/busybox"
    "sh" "-c" "busybox mkdir -p /sbin; /bin/busybox --install -s; sh"
  ];

  isolators = {
      "os/linux/capabilities-retain-set" = { "set" = [ "CAP_NET_ADMIN" ]; };
  };
}

