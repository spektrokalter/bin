package main

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("slide: ")

	flag.Parse()
	var (
		action = flag.Arg(0)
		slide  = flag.Arg(1)
	)

	// validate action
	//
	if action != "-" && action != "+" {
		log.Fatal("invalid action")
	}

	// validate slide
	//
	slideInfo, err := os.Stat(slide)
	if err != nil || !slideInfo.Mode().IsRegular() {
		log.Fatal("invalid slide")
	}

	entries, err := os.ReadDir(filepath.Dir(slide))
	if err != nil {
		log.Fatal(err)
	}

	// keep only regular files
	//
	files := entries[:0]
	for _, e := range entries {
		if e.Type().IsRegular() {
			files = append(files, e)
		}
	}
	entries = files

	slices.SortFunc(
		entries,
		func(a, b os.DirEntry) int {
			return cmp.Compare(a.Name(), b.Name())
		},
	)

	i := sort.Search(
		len(entries),
		func(i int) bool {
			n := cmp.Compare(slideInfo.Name(), entries[i].Name())
			return n == -1 || n == 0
		},
	)

	var next int
	if action == "-" {
		next = i - 1
		if next < 0 {
			next = len(entries) - 1
		}
	} else {
		next = (i + 1) % len(entries)
	}

	fmt.Printf("e! %s\n", entries[next].Name())
}
