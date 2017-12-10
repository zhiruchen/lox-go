package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/zhiruchen/lox-go/interpreter"
	"github.com/zhiruchen/lox-go/parser"
	"github.com/zhiruchen/lox-go/scanner"
)

func runPrompt() {

	reader := bufio.NewReader(os.Stdin)
	itp := interpreter.NewInterpreter()

	for {
		fmt.Print("code > ")
		source, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		run(itp, source)
	}
}

func run(itp *interpreter.Interpreter, source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	itp.Interpret(p.Parse())
}

func main() {
	runPrompt()
}
