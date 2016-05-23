{ mkACI, pkgs, thin ? false, ... } @ args:
let
  pkg = pkgs.plex;
in mkACI rec {
  inherit pkgs;
  inherit thin;
  packages = [ pkg ];

  mountPoints = {
    config.path = "/var/lib/plexmediaserver/Library/Application Support";
    media.path = "/media";
  };

  ports = {
    https.port = 32400;
  };

  exec = ["/usr/lib/plexmediaserver/Plex Media Server"];

  environment = {
    PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR = "/var/lib/plexmediaserver/Library/Application Support";
    PLEX_MEDIA_SERVER_HOME = "/usr/lib/plexmediaserver";
    PLEX_MEDIA_SERVER_MAX_PLUGIN_PROCS = "6";
    PLEX_MEDIA_SERVER_TMPDIR = "/tmp";
    LD_LIBRARY_PATH = "/usr/lib/plexmediaserver";
    LC_ALL = "en_US.UTF-8";
    LANG = "en_US.UTF-8";
  };
}
