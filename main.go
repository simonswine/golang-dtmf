package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"math"
	"os"
	"time"
)

var x = []float64{1209, 1336, 1477, 1663}
var y = []float64{697, 770, 852, 941}

const sampleRate = 44100

func playTone(char rune) {
	var toneA float64
	var toneB float64

	if char == '0' {
		toneA = x[1]
		toneB = y[3]
	} else if char == '1' {
		toneA = x[0]
		toneB = y[0]
	} else if char == '2' {
		toneA = x[1]
		toneB = y[0]
	} else if char == '3' {
		toneA = x[2]
		toneB = y[0]
	} else if char == '4' {
		toneA = x[0]
		toneB = y[1]
	} else if char == '5' {
		toneA = x[1]
		toneB = y[1]
	} else if char == '6' {
		toneA = x[2]
		toneB = y[1]
	} else if char == '7' {
		toneA = x[0]
		toneB = y[2]
	} else if char == '8' {
		toneA = x[1]
		toneB = y[2]
	} else if char == '9' {
		toneA = x[2]
		toneB = y[2]
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newDtmf(toneA, toneB, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(200 * time.Millisecond)
	chk(s.Stop())
	time.Sleep(40 * time.Millisecond)
}

func dial(nr string) {
	for _, r := range nr {
		playTone(r)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [number]\n", os.Args[0])
		os.Exit(2)
	}
	dial(os.Args[1])
}

type dtmfSine struct {
	*portaudio.Stream
	step1, phase1 float64
	step2, phase2 float64
}

func newDtmf(freq1, freq2, sampleRate float64) *dtmfSine {
	s := &dtmfSine{nil, freq1 / sampleRate, 0, freq2 / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *dtmfSine) processAudio(out [][]float32) {
	for i := range out[0] {
		value := float32(math.Sin(2*math.Pi*g.phase1) + math.Sin(2*math.Pi*g.phase2))
		out[0][i] = value
		_, phase1 := math.Modf(g.phase1 + g.step1)
		_, phase2 := math.Modf(g.phase2 + g.step2)

		g.phase1 = phase1
		out[1][i] = value
		g.phase2 = phase2
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
