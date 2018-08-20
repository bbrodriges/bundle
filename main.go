package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type generator struct {
	packageName   string
	removeSources bool
	filesArgs     []string
}

func run() error {
	flagPackage := flag.String("p", "", "Package name to be written into result file. Required")
	flagDeleteSource := flag.Bool("d", false, "Delete source files after bundling")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Usage: bundle -p <package_name> [-d] pattern... \n")
		fmt.Fprint(flag.CommandLine.Output(), "Available options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *flagPackage == "" {
		flag.Usage()
		return errors.New("Error: package name required")
	}

	argsFiles := flag.Args()
	if len(argsFiles) == 0 {
		flag.Usage()
		return errors.New("Error: no files given")
	}

	g := generator{
		packageName:   *flagPackage,
		removeSources: *flagDeleteSource,
		filesArgs:     argsFiles,
	}

	return g.makeBundle(os.Stdout)
}

func (g generator) makeBundle(w io.Writer) error {
	files, err := g.collectFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("Error: no files found in given path")
	}

	if len(files) == 1 {
		return errors.New("Covardly refusing to bundle single file")
	}

	if err := g.writeOutput(w, files); err != nil {
		return err
	}

	g.deleteSources(files)

	return nil
}

func (g generator) collectFiles() (files []string, err error) {
	files = make([]string, 0, len(g.filesArgs))

	for _, argFile := range g.filesArgs {
		paths, err := filepath.Glob(argFile)
		if err != nil {
			return nil, err
		}

		for _, path := range paths {
			pstat, err := os.Stat(path)
			if err != nil {
				return nil, err
			}

			if !pstat.Mode().IsRegular() {
				return nil, fmt.Errorf("Error: non-file \"%s\" found in args", path)
			}

			files = append(files, path)
		}
	}

	return
}

func (g generator) writeOutput(w io.Writer, files []string) error {
	// collect imports
	var imports []*ast.ImportSpec
	for _, path := range files {
		fset := token.NewFileSet()

		fileAST, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		imports = append(imports, fileAST.Imports...)
	}

	bundleBuf := bytes.NewBuffer(nil)

	// write gen notification
	bundleBuf.WriteString("// Code generated by bundle generation tool; DO NOT EDIT.\n\n")
	// write package name
	bundleBuf.WriteString("package " + g.packageName + "\n\n")

	// write imports
	if len(imports) > 0 {
		bundleBuf.WriteString("import (\n")
		for _, im := range imports {
			if im.Name != nil {
				bundleBuf.WriteString("\t" + im.Name.String() + " " + im.Path.Value + "\n")
			} else {
				bundleBuf.WriteString("\t" + im.Path.Value + "\n")
			}
		}
		bundleBuf.WriteString(")\n")
	}

	// write well-formed code
	for _, path := range files {
		fset := token.NewFileSet()

		fileAST, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// find and remove import declarations
		fileAST.Imports = nil
		for i := len(fileAST.Decls) - 1; i >= 0; i-- {
			decl, ok := fileAST.Decls[i].(*ast.GenDecl)
			if ok && decl.Tok.String() == "import" {
				fileAST.Decls = append(fileAST.Decls[:i], fileAST.Decls[i+1:]...)
			}
		}

		// write well-formed code into byte buffer
		buf := bytes.NewBuffer(nil)
		printer.Fprint(buf, fset, fileAST)

		bundleBuf.WriteString("\n// The code below has been bundled from \"" + path + "\" source file.\n")

		// remove old package declaration and write original code
		packagePos := fileAST.Package - 1
		if fileAST.Package > 1 {
			packagePos = fileAST.Package - 2
		}
		bundleBuf.Write(buf.Bytes()[:packagePos])
		bundleBuf.Write(buf.Bytes()[fileAST.Name.End():])

		buf.Reset()
	}

	// sort bundle imports
	fset := token.NewFileSet()
	bundleAST, err := parser.ParseFile(fset, "", bundleBuf, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.SortImports(fset, bundleAST)
	printer.Fprint(w, fset, bundleAST)

	w.Write([]byte("\n"))

	return nil
}

func (g generator) deleteSources(files []string) {
	if g.removeSources {
		for _, path := range files {
			os.Remove(path)
		}
	}
}
