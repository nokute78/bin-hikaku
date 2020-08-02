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
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

func generateBytes(b *testing.B, size int, diff int) []byte {
	b.Helper()

	buf := bytes.NewBuffer(make([]byte, size))
	w := bufio.NewWriter(buf)

	for i := 0; i < size; i++ {
		if rand.Intn(100) > diff {
			w.WriteByte(0x88)
		} else {
			w.WriteByte(0xee)
		}
	}
	w.Flush()
	return buf.Bytes()
}

func BenchmarkCompareSimple_10(b *testing.B) {
	var err error

	size := 8 * 1024 * 1024
	abuf := bytes.NewReader(generateBytes(b, size, 0))
	bbuf := bytes.NewReader(generateBytes(b, size, 10))

	out := bytes.NewBuffer([]byte{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out.Reset()
		abuf.Seek(0, 0)
		bbuf.Seek(0, 0)
		if err = CompareSimple(out, abuf, bbuf, 0, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func TestCompareSimple(t *testing.T) {
	type testcase struct {
		name   string
		input1 []byte
		input2 []byte
		expect []string
	}
	cases := []testcase{
		{"1byte",
			[]byte{0x00, 0x00, 0x00, 0x00},
			[]byte{0x00, 0xFF, 0x00, 0x00},
			[]string{"1", "00", "ff"},
		},
		{"size mismatch1",
			[]byte{0x00, 0x00, 0x00},
			[]byte{0x00, 0xFF, 0x00, 0x00},
			[]string{"1", "00", "ff"},
		},
		{"size mismatch2",
			[]byte{0x00, 0x00, 0x00, 0x00},
			[]byte{0x00, 0xFF, 0x00},
			[]string{"1", "00", "ff"},
		},
	}

	buf := bytes.NewBuffer([]byte{})

	for _, v := range cases {
		buf.Reset()
		a := bytes.NewReader(v.input1)
		b := bytes.NewReader(v.input2)
		if err := CompareSimple(buf, a, b, 0, 0); err != nil {
			t.Errorf("%s: %s", v.name, err)
		}
		strs := strings.Split(buf.String(), " ")
		if len(strs) != len(v.expect) {
			t.Errorf("%s: size mismatch", v.name)
		} else {
			for i, str := range strs {
				if !strings.Contains(str, v.expect[i]) {
					t.Errorf("%s: given %s expect %s output:%s", v.name, str, v.expect[i], buf.String())
				}
			}
		}
	}

	// same case
	buf.Reset()
	err := CompareSimple(buf,
		bytes.NewReader([]byte{0xaa, 0xbb, 0xcc, 0xdd}),
		bytes.NewReader([]byte{0xaa, 0xbb, 0xcc, 0xdd}),
		0, 0)
	if err != nil {
		t.Errorf("samecase:%s", err)
	} else if buf.Len() > 0 {
		t.Errorf("samecase: diff: %d", buf.Len())
	}

	// skip case
	buf.Reset()
	err = CompareSimple(buf,
		bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00}),
		bytes.NewReader([]byte{0xff, 0xff, 0xff, 0x00, 0x00}),
		3, 0)
	if err != nil {
		t.Errorf("skipcase:%s", err)
	} else if buf.Len() > 0 {
		t.Errorf("skipcase: diff: %d", buf.Len())
	}

	// limit case
	buf.Reset()
	err = CompareSimple(buf,
		bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00}),
		bytes.NewReader([]byte{0x00, 0x00, 0x00, 0xff, 0xff}),
		0, 3)
	if err != nil {
		t.Errorf("limitcase:%s", err)
	} else if buf.Len() > 0 {
		t.Errorf("limitcase: diff: %d", buf.Len())
	}

	// skip & limit case
	buf.Reset()
	err = CompareSimple(buf,
		bytes.NewReader([]byte{0xff, 0x00, 0x00, 0x00, 0x00}),
		bytes.NewReader([]byte{0x00, 0x00, 0x00, 0xff, 0xff}),
		1, 2)
	if err != nil {
		t.Errorf("skip_limit case:%s", err)
	} else if buf.Len() > 0 {
		t.Errorf("skip_limit case: diff: %d", buf.Len())
	}
}
