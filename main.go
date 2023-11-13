package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "v", false, "-v 1 to print detail")
}

func main() {
	flag.Parse()
	fmt.Println("Simple script language!")
	fmt.Println("input exit(); to quit")
	f := bufio.NewReader(os.Stdin)
	prompt := ">>"
	scriptText := ""
	calculator := SimpleParser{}
	script := NewSimpleScript(verbose)
	for {
		fmt.Print(prompt)
		input, _ := f.ReadString('\n')
		if len(input) == 1 {
			continue
		}
		input = strings.TrimSpace(input)
		if input == "exit();" {
			fmt.Println("good bye!")
			break
		}
		scriptText += input + "\n"
		if strings.HasSuffix(scriptText, ";\n") {
			fmt.Println("your input is: " + scriptText)
			root := calculator.Parse(scriptText)
			if verbose {
				DumpAST(*root, "")
			}
			script.Evaluate(*root, "")
			scriptText = ""
		}
	}
}
