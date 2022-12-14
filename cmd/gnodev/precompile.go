package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnolang/gno/pkgs/command"
	"github.com/gnolang/gno/pkgs/errors"
	gno "github.com/gnolang/gno/pkgs/gnolang"
)

type precompileOptions struct {
	Verbose     bool   `flag:"verbose" help:"verbose"`
	SkipFmt     bool   `flag:"skip-fmt" help:"do not check syntax of generated .go files"`
	GoBinary    string `flag:"go-binary" help:"go binary to use for building"`
	GofmtBinary string `flag:"go-binary" help:"gofmt binary to use for syntax checking"`
	Output      string `flag:"output" help:"output directory"`
	Test        bool   `flag:"test" help:"include test files"`
}

var defaultPrecompileOptions = precompileOptions{
	Verbose:     false,
	SkipFmt:     false,
	GoBinary:    "go",
	GofmtBinary: "gofmt",
	Output:      ".",
	Test:        false,
}

func precompileApp(cmd *command.Command, args []string, iopts interface{}) error {
	opts := iopts.(precompileOptions)
	if len(args) < 1 {
		cmd.ErrPrintfln("Usage: precompile [precompile flags] [packages]")
		return errors.New("invalid args")
	}

	// precompile .gno files.
	paths, err := gnoFilesFromArgs(args)
	if err != nil {
		return fmt.Errorf("list paths: %w", err)
	}

	errCount := 0
	for _, filepath := range paths {
		err = precompileFile(filepath, opts)
		if err != nil {
			err = fmt.Errorf("%s: precompile: %w", filepath, err)
			cmd.ErrPrintfln("%s", err.Error())
			errCount++
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d precompile errors", errCount)
	}

	return nil
}

func precompilePkg(pkgPath string, opts precompileOptions) error {
	files, err := filepath.Glob(filepath.Join(pkgPath, "*.go"))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if err = precompileFile(file, opts); err != nil {
			return fmt.Errorf("%s: %v", file, err)
		}
	}

	return nil
}

func precompileFile(srcPath string, opts precompileOptions) error {
	// shouldCheckFmt := !opts.SkipFmt
	verbose := opts.Verbose
	gofmt := opts.GofmtBinary
	if gofmt == "" {
		gofmt = defaultPrecompileOptions.GofmtBinary
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "%s\n", srcPath)
	}

	// parse .gno.
	source, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// compute attributes based on filename.
	var targetFilename string
	var tags string
	nameNoExtension := strings.TrimSuffix(filepath.Base(srcPath), ".go")
	targetFilename = nameNoExtension + ".gno"
	switch {
	case strings.HasSuffix(srcPath, "_filetest.go"):
		if !opts.Test {
			return nil
		}
		tags = "gno,filetest"
	case strings.HasSuffix(srcPath, "_test.go"):
		if !opts.Test {
			return nil
		}
		tags = "gno,test"
	default:
		tags = "gno"
	}

	// preprocess.
	transformed, err := gno.Precompile(string(source), tags, srcPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// write .go file.
	dir := filepath.Dir(srcPath)
	buildDir := filepath.Join(dir, "build")
	if err = os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("%w", err)
	}
	targetPath := filepath.Join(buildDir, targetFilename)
	err = os.WriteFile(targetPath, []byte(transformed), 0o644)
	if err != nil {
		return fmt.Errorf("write .go file: %w", err)
	}

	// check .go fmt.
	// if shouldCheckFmt {
	// 	err = gno.PrecompileVerifyFile(targetPath, gofmt)
	// 	if err != nil {
	// 		return fmt.Errorf("check .go file: %w", err)
	// 	}
	// }

	return nil
}
