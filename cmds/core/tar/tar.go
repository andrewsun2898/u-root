// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Create and extract tar archives.
//
// Synopsis:
//
//	tar [OPTION...] [FILE]...
//
// Description:
//
//	This command line can be used only in the following ways:
//	   tar -cvf x.tar directory/         # create
//	   tar -cvf x.tar file1 file2 ...    # create
//	   tar -tvf x.tar                    # list
//	   tar -xvf x.tar directory/         # extract
//
// Options:
//
//	-c: create a new tar archive from the given directory
//	-x: extract a tar archive to the given directory
//	-v: verbose, print each filename (optional)
//	-f: tar filename (required)
//	-t: list the contents of an archive
//
// TODO: The arguments deviates slightly from gnu tar.
package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/tarutil"
)

type cmd struct {
	p    params
	args []string
}

type params struct {
	file        string
	create      bool
	extract     bool
	list        bool
	noRecursion bool
	verbose     bool
}

var (
	errCreateAndExtract     = fmt.Errorf("cannot supply both -c and -x")
	errCreateAndList        = fmt.Errorf("cannot supply both -c and -t")
	errExtractAndList       = fmt.Errorf("cannot supply both -x and -t")
	errEmptyFile            = fmt.Errorf("file is required")
	errMissingMandatoryFlag = fmt.Errorf("must supply at least one of: -c, -x, -t")
	errExtractArgsLen       = fmt.Errorf("args length should be 1")
)

func command(p params, args []string) (*cmd, error) {
	if p.create && p.extract {
		return nil, errCreateAndExtract
	}
	if p.create && p.list {
		return nil, errCreateAndList
	}
	if p.extract && p.list {
		return nil, errExtractAndList
	}
	if p.extract && len(args) != 1 {
		return nil, errExtractArgsLen
	}
	if !p.extract && !p.create && !p.list {
		return nil, errMissingMandatoryFlag
	}
	if p.file == "" {
		return nil, errEmptyFile
	}

	return &cmd{
		p:    p,
		args: args,
	}, nil
}

func (c *cmd) run() error {
	opts := &tarutil.Opts{
		NoRecursion: c.p.noRecursion,
	}
	if c.p.verbose {
		opts.Filters = []tarutil.Filter{tarutil.VerboseFilter}
	}

	switch {
	case c.p.create:
		f, err := os.Create(c.p.file)
		if err != nil {
			return err
		}
		if err := tarutil.CreateTar(f, c.args, opts); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	case c.p.extract:
		f, err := os.Open(c.p.file)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := tarutil.ExtractDir(f, c.args[0], opts); err != nil {
			return err
		}
	case c.p.list:
		f, err := os.Open(c.p.file)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := tarutil.ListArchive(f); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	create := flag.BoolP("create", "c", false, "create a new tar archive from the given directory")
	extract := flag.BoolP("extract", "x", false, "extract a tar archive from the given directory")
	file := flag.StringP("file", "f", "", "tar file")
	list := flag.BoolP("list", "t", false, "list the contents of an archive")
	noRecursion := flag.Bool("no-recursion", false, "do not automatically recurse into directories")
	verbose := flag.BoolP("verbose", "v", false, "print each filename")

	flag.Parse()
	cmd, err := command(params{file: *file, create: *create, extract: *extract, list: *list, noRecursion: *noRecursion, verbose: *verbose}, flag.Args())
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	if err := cmd.run(); err != nil {
		log.Fatal(err)
	}
}
