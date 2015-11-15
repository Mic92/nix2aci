# nix2aci
Let's use Nix' super powers to build [App Container Images](http://github.com/appc/spec)!

This project should be understood as a proof of concept until stated otherwise.
You can expect this README to be minimal but it should always contain working examples.

## Build Requirements
* local copy of this repository
* [nix](http://www.nixos.org/nix) plus the skills to query package names

## Runtime Requirements
* [coreos/rkt](https://github.com/coreos/rkt/)

## Signing Requirements
Including the signing process into the nix workflow seems quite tedious and is not fully satisfactory at this point. I chose to setup a key for the nix build environment. The downside is that every build can read and use the key which is bad if the build system is compromised in any way.

Please take a look at [the included script](scripts/setup-gpg.sh), which can be used to do the preparation.


# Building ACIs
There's more than one way to build and use ACIs with Nix, because the filesystem structure allows for side-by-side installations of almost any package. Every package (version) is stored at $NIX\_STORE identified by hash, and can be pulled into different profiles independently. These profiles could be copied, but it should also be possible to bind-mount the host versions into the containers.

## Thin ACIs
* Status: working and under development

Thin ACIs don't contain any binary files, but for the most part just the manifest file and a directory skeleton.
The manifest file specifies one host type mount per package, representing the effectively available packages for the ACI.
These mountpoints can add up to a few dozen depending on the target package, and they all have to be passed to the container runtime that consumes the ACI, supplying the correct path from the host to the package's mount.

In order to make this usable, a file that can be `cat`ed into the rkt cmdline will be generated alongside the ACI when using the build script.

### Usage
This section will surely ***CHANGES TO FAST TO UPDATE***.


## Fat ACIs
* Status: working and under development

Fat ACIs contain all files that are needed to run the contained application. This is the choice if you want to move the ACI onto a system where for whichever reason the nix store outputs are not available.


### Usage Example
This section will surely ***CHANGES TO FAST TO UPDATE***.

```
$ nix-build . -A busybox
```
