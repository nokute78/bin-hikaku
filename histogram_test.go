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
	"bytes"
	"strings"
	"testing"
)

func TestPrintHistogram(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	input := []uint{0, 1, 30, 32, 100}
	expect := []int{0, 1, 3, 4, 10}

	if err := PrintHistogram(buf, input, 0, 1); err != nil {
		t.Fatalf("%s", err)
	}

	strs := strings.Split(buf.String(), "\n")
	t.Logf("\n%s\n", buf.String())
	if len(strs) < len(expect) {
		t.Fatalf("size mismatch given:%d expect:%d", len(strs), len(expect))
	}

	for i := 0; i < len(expect); i++ {
		given := strings.Count(strs[i], "#")
		if given != expect[i] {
			t.Errorf("%d: given %d expect %d. str=%s\n", i, given, expect[i], strs[i])
		}
	}
}

func TestHistogram(t *testing.T) {
	unitsize := 1000
	unit := bytes.NewBuffer(make([]byte, unitsize))

	abuf := bytes.NewBuffer([]byte{})
	bbuf := bytes.NewBuffer([]byte{})

	expect := []uint{0, 5, 64, 100}

	t.Log("generate data")
	for i := 0; i < len(expect); i++ {
		t.Logf("%d start", i)
		unit.Reset()
		unit.Write(bytes.Repeat([]byte{0x00}, unitsize))
		_, err := abuf.Write(unit.Bytes())
		if err != nil {
			t.Fatalf("abuf.Write:%s", err)
		}

		t.Logf("%d mid", i)
		count := int(expect[i]) * unitsize / 100
		unit.Reset()
		unit.Write(bytes.Repeat([]byte{0xee}, count))
		unit.Write(bytes.Repeat([]byte{0x00}, unitsize-count))
		_, err = bbuf.Write(unit.Bytes())
		if err != nil {
			t.Fatalf("bbuf.Write:%s", err)
		}
	}
	t.Logf("test Histogram abuf.Len()=%d bbuf.Len()=%d", abuf.Len(), bbuf.Len())
	u, err := Histogram(bytes.NewReader(abuf.Bytes()), bytes.NewReader(bbuf.Bytes()), 0, uint(abuf.Len()), uint(unitsize))
	if err != nil {
		t.Fatalf("Histogram:%s", err)
	}

	if len(u) != len(expect) {
		t.Fatalf("array size mismatch")
	}

	for i := 0; i < len(u); i++ {
		if u[i] != expect[i] {
			t.Errorf("%d:given %d expect %d", i, u[i], expect[i])
		}
	}

}
