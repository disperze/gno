package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func isGnoFile(f fs.DirEntry) bool {
	name := f.Name()
	return !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go") && !f.IsDir()
}

func gnoFilesFromArgs(args []string) ([]string, error) {
	paths := []string{}
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid file or package path: %w", err)
		}
		if !info.IsDir() {
			curpath := arg
			paths = append(paths, curpath)
		} else {
			err = filepath.WalkDir(arg, func(curpath string, f fs.DirEntry, err error) error {
				if err != nil {
					return fmt.Errorf("%s: walk dir: %w", arg, err)
				}

				if !isGnoFile(f) {
					return nil // skip
				}
				paths = append(paths, curpath)
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return paths, nil
}

func gnoPackagesFromArgs(args []string) ([]string, error) {
	paths := []string{}
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid file or package path: %w", err)
		}
		if !info.IsDir() {
			paths = append(paths, arg)
		} else {
			// if the passed arg is a dir, then we'll recursively walk the dir
			// and look for directories containing at least one .gno file.

			visited := map[string]bool{} // used to run the builder only once per folder.
			err = filepath.WalkDir(arg, func(curpath string, f fs.DirEntry, err error) error {
				if err != nil {
					return fmt.Errorf("%s: walk dir: %w", arg, err)
				}
				if f.IsDir() {
					return nil // skip
				}
				if !isGnoFile(f) {
					return nil // skip
				}

				parentDir := filepath.Dir(curpath)
				if _, found := visited[parentDir]; found {
					return nil
				}
				visited[parentDir] = true

				// cannot use path.Join or filepath.Join, because we need
				// to ensure that ./ is the prefix to pass to go build.
				pkg := "." + string(filepath.Separator) + parentDir
				paths = append(paths, pkg)
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return paths, nil
}

func fmtDuration(d time.Duration) string {
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func guessRootDir() string {
	cmd := exec.Command("go", "list", "-m", "-mod=mod", "-f", "{{.Dir}}", "github.com/gnolang/gno")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("can't guess --root-dir, please fill it manually.")
	}
	rootDir := strings.TrimSpace(string(out))
	return rootDir
}

func CpFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CpDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CpDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CpFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
