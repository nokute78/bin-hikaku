/*
   Copyright 2020 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const version string = "0.0.1"

// Exit status
const (
	ExitOK int = iota
	ExitArgError
	ExitCmdError
)

// CLI has In/Out/Err streams.
type CLI struct {
	OutStream io.Writer
	InStream  io.Reader
	ErrStream io.Writer
	quiet     bool // for testing to suppress output
}

// Run executes real main function.
func (cli *CLI) Run(args []string) (ret int) {
	cnf, err := Configure(args[1:], cli.quiet)
	if err != nil {
		if err == flag.ErrHelp {
			return ExitOK
		}
		fmt.Fprintf(cli.ErrStream, "%s\n", err)
		return ExitArgError
	}

	if cnf.showVersion {
		fmt.Fprintf(cli.OutStream, "Ver: %s\n", version)
		return ExitOK
	}

	if len(cnf.inputFiles) != 2 {
		fmt.Fprintf(cli.ErrStream, "invalid input len=%d\n", len(cnf.inputFiles))
		return ExitArgError
	}

	file1, err := os.Open(cnf.inputFiles[0])
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "%s", err)
		return ExitCmdError
	}
	defer file1.Close()
	file2, err := os.Open(cnf.inputFiles[1])
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "%s", err)
		return ExitCmdError
	}
	defer file2.Close()

	err = CompareSimple(cli.OutStream, file1, file2, int64(cnf.startOffset), cnf.readSize)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "%s", err)
		return ExitCmdError
	}

	return ExitOK
}

func main() {
	cli := &CLI{OutStream: os.Stdout, InStream: os.Stdin, ErrStream: os.Stderr}

	os.Exit(cli.Run(os.Args))
}
