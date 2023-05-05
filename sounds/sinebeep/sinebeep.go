// Copyright 2019 The Oto Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file is a copy of:
// https://github.com/hajimehoshi/oto/blob/main/example/main.go
// ... and then hacked to its current form.

// Yeah, Alle meine Entchen!
// https://www.spiellieder.de/kinderlied-standards/alle_meine_entchen.php
// sinebeep -duration 240 cdefgg,aaaag,aaaag,ffffee,ggggc

package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate   = 44100
	channelCount = 2
)

var (
	octaveShift int = 0
	duration        = flag.Int("duration", 1000, "duration in ms")
	keyMap          = map[rune]int{
		// https://de.wikipedia.org/wiki/Frequenzen_der_gleichstufigen_Stimmung
		'c': 40, // C4
		'C': 41,
		'd': 42,
		'D': 43,
		'e': 44,
		'f': 45,
		'F': 46,
		'g': 47,
		'G': 48,
		'a': 49,
		'A': 50,
		'b': 51, // B4 🇩🇪 h^1
	}
	opMap = map[rune]func() string{
		'^': func() string { octaveShift += 12; return "Octave raised" },
		'V': func() string { octaveShift -= 12; return "Octave lowered" },
		'X': func() string { *duration *= 2; return "Duration * 2" },
		':': func() string { *duration /= 2; return "Duration / 2" },
	}
)

type SineWave struct {
	freq   float64
	length int64
	pos    int64

	channelCount int
	format       int

	remaining []byte
}

func NewSineWave(freq float64, duration time.Duration, channelCount int, format int) *SineWave {
	l := int64(channelCount) * int64(4) * int64(sampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		freq:         freq,
		length:       l,
		channelCount: channelCount,
		format:       format,
	}
}

func (s *SineWave) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		copy(s.remaining, s.remaining[n:])
		s.remaining = s.remaining[:len(s.remaining)-n]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(sampleRate) / float64(s.freq)

	num := 4 * s.channelCount
	p := s.pos / int64(num)

	for i := 0; i < len(buf)/num; i++ {
		bs := math.Float32bits(float32(math.Sin(2*math.Pi*float64(p)/length) * 0.3))
		for ch := 0; ch < channelCount; ch++ {
			buf[num*i+4*ch] = byte(bs)
			buf[num*i+1+4*ch] = byte(bs >> 8)
			buf[num*i+2+4*ch] = byte(bs >> 16)
			buf[num*i+3+4*ch] = byte(bs >> 24)
		}
		p++
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func play(context *oto.Context, freq float64, duration time.Duration, channelCount int, format int) oto.Player {
	p := context.NewPlayer(NewSineWave(freq, duration, channelCount, format))
	p.Play()
	return p
}

func run(c *oto.Context, frequency float64) error {
	p := play(c, frequency, time.Duration(*duration)*time.Millisecond, channelCount, oto.FormatFloat32LE)
	time.Sleep(time.Duration(*duration) * time.Millisecond)
	p.Close()
	return nil
}

func main() {
	frequency := flag.Float64("frequency", 523.3, "frequency in Hz (float)")
	flag.Parse()
	ctx, ready, err := oto.NewContext(sampleRate, channelCount, oto.FormatFloat32LE)
	if err != nil {
		panic(err)
	}
	<-ready
	if len(flag.Args()) > 0 {
		for _, args := range flag.Args() {
			for _, c := range args {
				if keyNum, ok := keyMap[c]; ok {
					f := math.Pow(2, ((float64(keyNum+octaveShift)-49)/12)) * 440
					fmt.Printf("%c -> %d => f=%f Hz. %d ms\n", c, keyNum+octaveShift, f, *duration)
					run(ctx, f)
				} else if op, ok := opMap[c]; ok {
					fmt.Printf("%s\n", op())
				} else {
					time.Sleep(time.Duration(*duration) * time.Millisecond)
				}
			}
		}
		return
	}
	if err := run(ctx, *frequency); err != nil {
		panic(err)
	}
}
