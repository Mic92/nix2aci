package main

import (
	"archive/tar"
	"bufio"
	"crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Image struct {
	thin, static, dnsquirks bool
	manifest, customEnv     string
	closures                []string
}

type ChecksumWriter struct {
	to   io.Writer
	Hash hash.Hash
}

func (c *ChecksumWriter) Write(p []byte) (int, error) {
	c.Hash.Write(p)
	return c.to.Write(p)
}

func addFile(writer *tar.Writer, path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("cannot read '%s': %v", path, err)
	}
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("cannot write aci image entry '%s', %v", path, err)
	}
	return nil
}

func addPath(writer *tar.Writer, rootPath string, prefixLen int, dereferenceEtc bool) error {
	walkFn := func(source string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory '%s': %v", rootPath, err)
		}
		header := tar.Header{}
		relPath := source[(prefixLen):]
		header.Name = path.Join("rootfs", relPath)
		header.ModTime = info.ModTime()
		header.Mode = int64(info.Mode().Perm())

		if info.Mode().IsRegular() ||
			(dereferenceEtc && info.Mode()&os.ModeSymlink != 0 && strings.HasPrefix(relPath, "etc/")) {
			header.Typeflag = tar.TypeReg
			header.Size = info.Size()
			if err := writer.WriteHeader(&header); err != nil {
				return fmt.Errorf("cannot write aci image entry header '%s': %v", source, err)
			}
			return addFile(writer, source)
		} else {
			if info.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(source)
				if err != nil {
					return fmt.Errorf("cannot read resolve symlink '%s', %v", target, err)
				}
				header.Linkname = target
				header.Typeflag = tar.TypeSymlink
			} else if info.Mode().IsDir() {
				header.Typeflag = tar.TypeDir
			} else {
				return fmt.Errorf("not implemented")
			}
			if err := writer.WriteHeader(&header); err != nil {
				return fmt.Errorf("cannot write aci image entry '%s': %v", source, err)
			}
		}
		return nil
	}
	return filepath.Walk(rootPath, walkFn)
}

func readManifest(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest '%s': %v", path, err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest '%s': %v", path, err)
	}
	var manifest map[string]interface{}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("cannot decode manifest '%s': %v", path, err)
	}
	return manifest, nil
}

func writeManifest(tw *tar.Writer, manifest map[string]interface{}) error {
	buf, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("cannot serialize manifest: %v", err)
	}
	header := tar.Header{}
	header.Name = "manifest"
	header.Mode = 0644
	header.Size = int64(len(buf))
	header.Typeflag = tar.TypeReg
	tw.WriteHeader(&header)
	if _, err := tw.Write(buf); err != nil {
		return fmt.Errorf("cannot write manifest: %v", err)
	}
	return nil
}

func (i *Image) addManifest(tw *tar.Writer, mountFlags io.Writer, storePaths []string) error {
	manifest, err := readManifest(i.manifest)
	if err != nil {
		return err
	}
	if i.thin {
		if err := addNixstoreMounts(manifest, mountFlags, storePaths); err != nil {
			return fmt.Errorf("cannot add nixstore mounts: %v", err)
		}
	}

	return writeManifest(tw, manifest)
}

func writeHosts(tw *tar.Writer) error {
	etcheader := tar.Header{}
	etcheader.Name = "rootfs/etc/"
	etcheader.Mode = 0644
	etcheader.Typeflag = tar.TypeDir
	tw.WriteHeader(&etcheader)

	content := []byte("127.0.0.1 localhost\n::1 localhost\n")
	header := tar.Header{}
	header.Name = "rootfs/etc/hosts"
	header.Mode = 0644
	header.Typeflag = tar.TypeReg
	header.Size = int64(len(content))

	tw.WriteHeader(&header)
	if _, err := tw.Write(content); err != nil {
		return fmt.Errorf("cannot write manifest: %v", err)
	}
	return nil
}

func sanitizeName(s string) string {
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, ".", "", -1)
	return strings.ToLower(s)
}

type MountPoint struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readOnly"`
}

func addNixstoreMounts(manifest map[string]interface{}, mountFlags io.Writer, storePaths []string) error {
	app, ok := manifest["app"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("manifest does not contain App key")
	}
	mounts, ok := app["mountPoints"].([]interface{})
	if !ok {
		return fmt.Errorf("manifest does not contain mountPoints key")
	}
	for _, path := range storePaths {
		name := sanitizeName(path)
		mounts = append(mounts, MountPoint{name, path, true})
		arg := fmt.Sprintf("--volume=%s,kind=host,source=%s ", name, path)
		if _, err := mountFlags.Write([]byte(arg)); err != nil {
			return fmt.Errorf("cannot write to mountflags: %v", err)
		}
	}
	app["mountPoints"] = mounts
	return nil
}

func (i *Image) build(writer, mountFlags io.Writer) (hash.Hash, error) {
	var storePaths []string
	var err error
	if i.static {
		storePaths = i.closures
	} else {
		if storePaths, err = pathsFromGraphs(i.closures); err != nil {
			return nil, err
		}
	}

	checksumWriter := ChecksumWriter{writer, sha512.New()}
	tw := tar.NewWriter(&checksumWriter)
	if err := i.addManifest(tw, mountFlags, storePaths); err != nil {
		return nil, err
	}

	if err := addPath(tw, i.customEnv, len(i.customEnv), true); err != nil {
		return nil, err
	}

	if !i.thin {
		for _, path := range storePaths {
			if err := addPath(tw, path, 0, false); err != nil {
				return nil, err
			}
		}
	}
	if i.dnsquirks {
		if err := writeHosts(tw); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("cannot close aci: %s", err)
	}
	return checksumWriter.Hash, nil
}

func parseFlags() *Image {
	var i Image
	flag.BoolVar(&i.thin, "thin", false, "thin images")
	flag.BoolVar(&i.dnsquirks, "dnsquirks", false, "write /etc/hosts")
	flag.BoolVar(&i.static, "static", false, "static executeable")
	flag.Parse()
	args := flag.Args()
	i.manifest = args[0]
	i.customEnv = args[1]
	i.closures = args[2:]
	return &i
}

func pathsFromGraphs(graphs []string) ([]string, error) {
	paths := make(map[string]bool)
	for _, graph := range graphs {
		f, err := os.Open(graph)
		if err != nil {
			return nil, fmt.Errorf("cannot open graph '%s': %v", graph, err)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			paths[scanner.Text()] = true
			scanner.Scan()
			scanner.Scan()
			count, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid graph file, expected number: %v", err)
			}
			for i := 0; i < count; i++ {
				scanner.Scan()
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("cannot read graph '%s': %v", err)
		}
	}

	keys := make([]string, 0, len(paths))
	for k := range paths {
		keys = append(keys, k)
	}
	return keys, nil
}

var tarStream = os.NewFile(3, "tar")
var checksumStream = os.NewFile(4, "checksum")
var mountFlagsStream = os.NewFile(5, "mountflags")

func main() {
	image := parseFlags()
	hash, err := image.build(tarStream, mountFlagsStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(-1)
	}
	fmt.Fprintf(checksumStream, "%x\n", hash.Sum(nil))
}
