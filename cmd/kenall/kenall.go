package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oirik/gokenall"
	"github.com/oirik/gosubcommand"
	"github.com/pkg/errors"
)

var (
	version  string
	revision string
)

func main() {
	gosubcommand.AppName = "kenall"
	gosubcommand.Version = fmt.Sprintf("version: %s\nrevision: %s", version, revision)
	gosubcommand.Summary = "kenall is a tool for managing ken_all.csv"

	download := &downloadCommand{}
	gosubcommand.Register("download", download)

	updated := &updatedCommand{}
	gosubcommand.Register("updated", updated)

	normalize := &normalizeCommand{}
	gosubcommand.Register("normalize", normalize)

	os.Exit(int(gosubcommand.Execute()))
}

type downloadCommand struct {
	extract bool
	output  string
}

func (download *downloadCommand) Summary() string {
	return "Download ken_all.zip from japanpost website"
}

func (download *downloadCommand) SetFlag(fs *flag.FlagSet) {
	fs.BoolVar(&download.extract, "x", false, "Extract file from an archive.")
	fs.StringVar(&download.output, "o", "", "Save file to <string> path instead of standard output.")
}

func (download *downloadCommand) Execute(fs *flag.FlagSet) gosubcommand.ExitCode {
	var w io.Writer

	if download.output == "" || download.output == "-" {
		w = os.Stdout
	} else {
		f, err := os.Create(download.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrapf(err, "failed to create file: %s", download.output))
			return gosubcommand.ExitCodeError
		}
		defer f.Close()
		w = f
	}

	if err := gokenall.Download(w, download.extract); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return gosubcommand.ExitCodeError
	}

	return gosubcommand.ExitCodeSuccess
}

type updatedCommand struct {
	print bool
}

func (updated *updatedCommand) Summary() string {
	return "Read updated date of data from japanpost website. Exit status 0 if later than [argument](yyyyMMdd) or exit status 1."
}

func (updated *updatedCommand) SetFlag(fs *flag.FlagSet) {
	fs.BoolVar(&updated.print, "p", false, "Print updated date.")
}

func (updated *updatedCommand) Execute(fs *flag.FlagSet) gosubcommand.ExitCode {
	t, err := time.Parse("20060102", fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "date string is wrong format. To compare, `yyyyMMdd` string is required: %s\n", fs.Arg(0))
		return gosubcommand.ExitCodeError
	}

	result, updatedTime, err := gokenall.Updated(t)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return gosubcommand.ExitCodeError
	}
	if updated.print {
		fmt.Fprintln(os.Stdout, updatedTime.Format("20060102"))
	}
	if !result {
		return gosubcommand.ExitCodeError
	}
	return gosubcommand.ExitCodeSuccess
}

type normalizeCommand struct {
	output string
	width  bool
	utf8   bool
	trim   bool
}

func (normalize *normalizeCommand) Summary() string {
	return "Normalize -make easy to use- input (file or standard input if no argument)"
}

func (normalize *normalizeCommand) SetFlag(fs *flag.FlagSet) {
	fs.StringVar(&normalize.output, "o", "", "Save file to <string> path instead of standard output.")
	fs.BoolVar(&normalize.width, "width", (gokenall.DefaultNormalizeOption&gokenall.NormalizeWidth) != 0, "Convert hankaku kana into zenkaku, ascii letters into hankaku")
	fs.BoolVar(&normalize.utf8, "utf8", (gokenall.DefaultNormalizeOption&gokenall.NormalizeUTF8) != 0, "Convert ShiftJIS into UTF8")
	fs.BoolVar(&normalize.trim, "trim", (gokenall.DefaultNormalizeOption&gokenall.NormalizeTrim) != 0, "Trim spaces from each text")
}

func (normalize *normalizeCommand) Execute(fs *flag.FlagSet) gosubcommand.ExitCode {
	input := fs.Arg(0)

	var r io.ReadCloser
	var w io.WriteCloser

	if input == "" || input == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrapf(err, "failed to open file: %s", input))
			return gosubcommand.ExitCodeError
		}
		r = f
	}
	defer r.Close()

	if normalize.output == "" || normalize.output == "-" {
		w = os.Stdout
	} else {
		f, err := os.Create(normalize.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrapf(err, "failed to create file: %s", normalize.output))
			return gosubcommand.ExitCodeError
		}
		w = f
	}
	defer w.Close()

	option := gokenall.DefaultNormalizeOption
	if normalize.width {
		option |= gokenall.NormalizeWidth
	} else {
		option &^= gokenall.NormalizeWidth
	}
	if normalize.utf8 {
		option |= gokenall.NormalizeUTF8
	} else {
		option &^= gokenall.NormalizeUTF8
	}
	if normalize.trim {
		option |= gokenall.NormalizeTrim
	} else {
		option &^= gokenall.NormalizeTrim
	}
	bufR := bufio.NewReader(r)
	bufW := bufio.NewWriter(w)
	if err := gokenall.Normalize(bufR, bufW, option); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return gosubcommand.ExitCodeError
	}
	err := bufW.Flush()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to flush: %v", err))
		return gosubcommand.ExitCodeError
	}

	return gosubcommand.ExitCodeSuccess
}
