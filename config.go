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
	"errors"
	"flag"
	"io/ioutil"
)

// ConfigArgsMissing represents no Args error
var ConfigNoArgs error = errors.New("No Args")
var ConfigInputFileSize error = errors.New("Invalid Input")

type Config struct {
	showVersion   bool
	readSize      uint
	startOffset   uint
	histogramMode bool
	unitSize      uint
	inputFiles    []string
}

// Pass os.Args[1:]
// silent is to suppress help message for testing.
func Configure(args []string, silent bool) (*Config, error) {
	ret := &Config{}
	if len(args) < 1 {
		return nil, ConfigNoArgs
	}

	opt := flag.NewFlagSet("bin-hikaku", flag.ContinueOnError)
	opt.BoolVar(&ret.showVersion, "V", false, "show Version")
	opt.BoolVar(&ret.histogramMode, "H", false, "histogram mode")
	opt.UintVar(&ret.readSize, "r", 0, "read size to compare.")
	opt.UintVar(&ret.startOffset, "s", 0, "skip size")
	opt.UintVar(&ret.unitSize, "u", 4096, "unit size")

	if silent {
		opt.SetOutput(ioutil.Discard)
	}

	err := opt.Parse(args)
	ret.inputFiles = opt.Args()
	return ret, err
}
