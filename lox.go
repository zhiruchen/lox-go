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

func runPromt() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("code > ")
		source, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		run(source)
	}
}

func run(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	// for _, tk := range tokens {
	// 	fmt.Println(tk.ToString())
	// }
	p := parser.NewParser(tokens)

	itp := &interpreter.Interpreter{}
	itp.Interprete(p.Parse())
}

func main() {
	runPromt()
}
