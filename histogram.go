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
	"errors"
	"fmt"
	"io"
	"strings"
)

func PrintHistogram(out io.Writer, histogram []uint, start uint, unitsize int) error {
	bufout := bufio.NewWriter(out)
	defer bufout.Flush()

	for i := 0; i < len(histogram); i++ {
		fmt.Fprintf(bufout, "0x%016x-0x%016x %3d%% ", start+uint(i*unitsize), start+uint((i+1)*unitsize), histogram[i])
		count := histogram[i] / 10
		if histogram[i]%10 > 0 {
			count += 1
		}
		if count > 0 {
			fmt.Fprintf(bufout, "%s\n", strings.Repeat("#", int(count)))
		} else {
			fmt.Fprint(bufout, "\n")
		}
	}
	return nil
}

func Histogram(a io.ReaderAt, b io.ReaderAt, off int64, limit uint, unitsize uint) ([]uint, error) {
	arrsize := 0
	if limit == 0 {
		return []uint{}, errors.New("size is zero")
	}
	arrindex := 0
	arrsize = int(limit / unitsize)
	if limit%unitsize != 0 {
		arrsize += 1
	}
	ret := make([]uint, arrsize)
	abuf := make([]byte, unitsize)
	bbuf := make([]byte, unitsize)

	var imax int
	offset := off
	for {
		an, aerr := a.ReadAt(abuf, offset)
		if aerr != nil && aerr != io.EOF {
			return ret, aerr
		}
		bn, berr := b.ReadAt(bbuf, offset)
		if berr != nil && berr != io.EOF {
			return ret, berr
		}
		//		fmt.Printf("%d:an=%d bn=%d aerr=%s berr=%s\n", arrindex, an, bn, aerr, berr)
		if aerr == io.EOF && berr == io.EOF {
			break
		}

		// readsize
		if uint(offset-off)+unitsize > limit {
			imax = int(limit) - (int(offset - off))
		} else if an > bn {
			imax = bn
		} else {
			imax = an
		}
		offset += int64(imax)

		for i := 0; i < imax; i++ {
			if abuf[i] != bbuf[i] {
				ret[arrindex] += 1
			}
		}

		if ret[arrindex] > 0 && ret[arrindex]*100 < uint(imax) {
			ret[arrindex] = 1
		} else {
			ret[arrindex] = ret[arrindex] * 100 / uint(imax)
		}

		arrindex += 1
		if aerr == io.EOF || berr == io.EOF {
			break
		}
	}
	return ret, nil
}
