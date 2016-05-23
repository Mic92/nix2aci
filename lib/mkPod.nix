args @ { apps
, version
, annotations ? []
, volumes ? []
, ports ? []
, isolators ? []
}:

let
  readMounts = (aci: (builtins.readJSON (builtins.readFile ${aci}/mounts.json)));
in
  pkgs.stdenv.mkdDerivation rec {
}
