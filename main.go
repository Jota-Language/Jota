package main

import (
	"bufio"
	"fmt"
	"jota/errors"
	"jota/interpreter"
	"jota/parser"
	"jota/scanner"
	"jota/utils"
	"log"
	"os"
	"path/filepath"
)

// Globals
var (
	errHandler = errors.ErrorHandler{
		Error:        false,
		RuntimeError: false,
		Log:          *log.New(os.Stderr, "", 0),
	}
	globalInterpreter = interpreter.NewInterpreter(errHandler)
)

func main() {
	args := os.Args[1:]
	length := len(args)

	if length > 1 {
		fmt.Println(utils.Yellow + "Usage -> " + utils.White + "jota [file.jota]" + utils.Reset)
		os.Exit(0)
	} else if length == 1 {
		if filepath.Ext(args[0]) != ".jota" {
			fmt.Println(utils.Yellow + "Usage ->" + utils.White + " You must enter an existing .jota file" + utils.Reset)
			fmt.Println(utils.Yellow + "Usage ->" + utils.White + " jota [file.jota]" + utils.Reset)
			os.Exit(0)
		}
		err := runFile(args[0])
		// If the error is not nil, it means that the error is not about non-existing files, so we should send a different error message to the user
		if err != nil {
			fmt.Println(utils.Red+"Error ->"+utils.White+" There was an error not related to a non-existing file:\n", err, "\n\n"+utils.Magenta+"Suggestion -> "+utils.White+"If you believe that this is an issue with the interpreter, please send an issue at "+utils.Blue+"https://github.com/mattishere/jota/issues"+utils.Reset)
		}
	} else {
		runREPL()
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(utils.Yellow + "Usage ->" + utils.White + " You must enter an existing .jota file" + utils.Reset)
			fmt.Println(utils.Yellow + "Usage ->" + utils.White + " jota [file.jota]" + utils.Reset)
			return nil
		}
		return err
	}

	run(string(bytes))

	// Since this is reading from a file, we need to stop execution if we encounter an error (in the REPL, we don't need to do this)
	if errHandler.Error || errHandler.RuntimeError {
		os.Exit(0)
	}

	return nil
}

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(utils.Green + "->" + utils.Reset + " ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		run(line)
		// We reset the errors
		errHandler.Error = false
	}
}


func run(source string) {
	scanner := scanner.CreateScanner(source, errHandler)
	tokens := scanner.ScanTokens()
	parser := parser.NewParser(tokens, errHandler)
	statements := parser.Parse()

	if statements == nil {
		return
	}

	if errHandler.Error {
		return
	}

	globalInterpreter.Interpret(statements)
}
