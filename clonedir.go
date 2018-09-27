package clonedir

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var debugging = os.Getenv("CLONEDIR_DEBUG") == "1"

func Clone(from, to string) error {
	dirs := map[string][]string{}
	err := filepath.Walk(from, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		mode := fi.Mode()
		perm := mode.Perm()
		if mode.IsDir() {
			p = path.Join(to, strings.TrimPrefix(p, from))
			os.MkdirAll(p, perm)
			dirs[p] = []string{}
		} else if mode.IsRegular() {
			dir := path.Dir(path.Join(to, strings.TrimPrefix(p, from)))
			dirs[dir] = append(dirs[dir], p)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for dir, files := range dirs {
		if len(files) == 0 {
			continue
		}
		args := append(files, dir)
		if runtime.GOOS == "darwin" {
			args = append([]string{"-c"}, args...)
		} else {
			args = append([]string{"--reflink=auto"}, args...)
		}
		if debugging {
			log.Printf("cp %q\n", args)
		}
		cmd := exec.Command("cp", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
