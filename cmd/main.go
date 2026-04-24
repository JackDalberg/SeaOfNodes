package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	simple "github.com/JackDalberg/SeaOfNodes"
	"github.com/JackDalberg/SeaOfNodes/ir"
)

func main() {
	useGoAST := flag.Bool("a", false, "")
	printString := flag.Bool("s", false, "")
	disablePeephole := flag.Bool("d", false, "")
	flag.Usage = func() {
		fmt.Println("Simple compiler written in Go. Prints graph representation of IR.")
		fmt.Printf("Usage: %s [-a] [-d] [-s] <file>\n", os.Args[0])
		fmt.Println("\t-a\tUse Go AST parser")
		fmt.Println("\t-d\tDisable peephole optimizations")
		fmt.Println("\t-s\tPrint string visualization")
		fmt.Println("\t-h\tPrint this help and exit")
	}

	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Missing file argument")
		flag.Usage()
		return
	}

	if *disablePeephole {
		ir.DisablePeephole = true
	}

	file := flag.Args()[0]
	fp, err := os.Open(file)
	if err != nil {
		fmt.Printf("error opening file %q", file)
		return
	}
	fileBytes, err := io.ReadAll(fp)
	if err != nil {
		fmt.Println("error reading in file contents")
		return
	}

	// Since reading from file, '\n' behaves weirdly with current parser... will fix
	filtered := bytes.ReplaceAll(fileBytes, []byte("\n"), []byte(""))
	code := string(filtered)

	var node ir.Node
	var generator *ir.Generator
	if *useGoAST {
		node, generator, err = simple.GoSimple(code)
	} else {
		node, generator, err = simple.Simple(code)
	}

	if err != nil {
		log.Fatalf("Compiler error: %v", err)
	}

	if *printString {
		fmt.Printf("String:\n\n%s\n", ir.ToString(node))
	} else {
		fmt.Printf("Graph:\n\n%s\n", ir.Visualize(generator))
	}
}
