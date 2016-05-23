{ mkACI, pkgs, thin ? false, ... } @ args:
let
  pkg = pkgs.dovecot;
in

mkACI rec {
  inherit pkgs;
  inherit thin;
  dnsquirks = args.dnsquirks;

  packages = [ pkg pkgs.dovecot_pigeonhole ];

  ports = {
    imaps.port = 993;
    sieve.port = 4190;
  };

  mountPoints = {
    mail.path = "/var/vmail";
    etc-dovecot.path = "/etc/dovecot";
  };

  environment = {
    LC_ALL = "en_US.UTF-8";
    LANG = "en_US.UTF-8";
  };
}
