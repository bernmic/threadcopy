package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func (c *Context) copyTarget(t *Target) error {
	if c.Verbose {
		fmt.Printf("start copy to target %s\n", t.Path)
	}
	for _, e := range t.Entries {
		fi, err := os.Stat(e)
		if err != nil {
			return fmt.Errorf("error stat file %s: %v", e, err)
		}
		if fi.IsDir() {
			_, err := c.copydir(e, t)
			if err != nil {
				fmt.Printf("error copying dir: %v\n", err)
			}
		} else {
			_, err := c.copyfile(e, filepath.Join(t.Path, fi.Name()))
			if err != nil {
				fmt.Printf("error copying file: %v\n", err)
			}
		}
	}
	if c.Verbose {
		fmt.Printf("finished copy to target %s\n", t.Path)
	}
	return nil
}

func (c *Context) copydir(src string, t *Target) (int, error) {
	parent := filepath.Dir(src)
	count := 0
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		subpath := path[len(parent):]
		targetpath := filepath.Join(t.Path, subpath)
		if info.IsDir() {
			fmt.Printf("create dir %s\n", targetpath)
			err := os.MkdirAll(targetpath, info.Mode())
			if err != nil {
				return fmt.Errorf("error creating dir %s: %v", targetpath, err)
			}
		} else {
			_, err := c.copyfile(path, targetpath)
			if err != nil {
				return fmt.Errorf("error copy file %s to %s: %v", path, targetpath, err)
			}
			count++
		}
		return nil
	})
	return count, err
}

// copyfile copies one file
// src absolute path of the file be be copied
// dst absolute path including filename of the copy
// returns number of bytes copied and an error if oone occured
func (c *Context) copyfile(src string, dst string) (int, error) {
	if c.Verbose {
		fmt.Printf("copy file %s to %s\n", src, dst)
	}
	bs, err := toByteSize(c.Buffersize)
	if err != nil {
		return 0, err
	}
	Buffer := make([]byte, bs)

	fin, ferr := os.Open(src)
	if ferr != nil {
		log.Fatal(ferr)
	}

	defer fin.Close()

	fout, ferr := os.Create(dst)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer fout.Close()
	copiedBytes := 0
	for {
		n, ferr := fin.Read(Buffer)
		if ferr != nil && ferr != io.EOF {
			log.Fatal(ferr)
		}
		if n == 0 {
			break
		}
		if _, ferr := fout.Write(Buffer[:n]); ferr != nil {
			log.Fatal(ferr)
		}
		copiedBytes += n
	}
	if c.Verbose {
		fmt.Printf("copied %d bytes\n", copiedBytes)
	}

	return copiedBytes, nil
}
