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
	"bufio"
	"fmt"
	"io"
)

func CompareSimple(out io.Writer, a io.ReaderAt, b io.ReaderAt, off int64, limit uint) error {
	var an, bn int
	var aerr, berr error
	var i int
	var imax int
	match := false

	bufout := bufio.NewWriter(out)
	defer bufout.Flush()

	bufsize := 4096

	offset := off
	abuf := make([]byte, bufsize)
	bbuf := make([]byte, bufsize)

	for {
		// read file
		an, aerr = a.ReadAt(abuf, offset)
		if aerr != nil && aerr != io.EOF {
			return aerr
		}
		bn, berr = b.ReadAt(bbuf, offset)
		if berr != nil && berr != io.EOF {
			return berr
		}

		// readsize
		if (limit > 0) && (int(offset-off)+bufsize > int(limit)) {
			imax = int(limit) - (int(offset - off))
		} else if an > bn {
			imax = bn
		} else {
			imax = an
		}

		for i = 0; i < imax; i++ {
			if abuf[i] != bbuf[i] {
				fmt.Fprintf(bufout, "0x%016x %02x %02x\n", offset+int64(i), abuf[i], bbuf[i])
				match = true
			} else {
				if match {
					fmt.Fprint(bufout, "\n")
				}
				match = false
			}
		}
		offset += int64(imax)

		if aerr == io.EOF || berr == io.EOF {
			break
		}
	}
	return bufout.Flush()
}
