package main

import (
	"fmt"
	"io"
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
			fmt.Printf("copy dir %s\n", e)
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

func (c *Context) copyfile(src string, dst string) (int, error) {
	fmt.Printf("copy file %s to %s\n", src, dst)
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
	fmt.Printf("copied %d bytes\n", copiedBytes)
	return copiedBytes, nil
}
