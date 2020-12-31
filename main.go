package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

const (
	AlgoMD5  = "md5"
	AlgoSHA1 = "sha1"
)

func main() {
	var algo string
	flag.StringVar(&algo, "algo", AlgoMD5, "The hash algorithm to use.")
	flag.Parse()

	in := flag.Arg(0)
	out := flag.Arg(1)
	info, err := os.Stat(out)
	if err != nil {
		fmt.Printf("Issue reading output directory (%s) -> %s\n", out, err.Error())
		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Println("Output must be a directory.")
		os.Exit(1)
	}
	if err := hashem(in, out, algo); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func hashem(dir, out, algo string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		var h hash.Hash
		switch algo {
		case AlgoSHA1:
			h = sha1.New()
		default:
			h = md5.New()
		}
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		_, err = io.Copy(h, f)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
		ext := filepath.Ext(path)
		filename := filepath.Join(out, hex.EncodeToString(h.Sum(nil)) + ext)
		if err := os.Rename(path, filename); err != nil {
			fmt.Printf("%s: %s\n", path, err.Error())
		} else {
			fmt.Printf("%s -> %s\n", path, filename)
		}
		return nil
	})
}
