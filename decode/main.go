package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/spectral"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Recording.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	nSamples := 0

	data := []float64{}

	chk(stream.Start())
	for {
		chk(stream.Read())
		for _, i := range in {
			data = append(data, float64(i))
		}
		nSamples += len(in)
		select {
		case <-sig:
			return
		default:
		}
		if nSamples > 2048 {
			p, freq := spectral.Pwelch(data, 44100.0, &spectral.PwelchOptions{})
			fmt.Printf("%+v %+v", p, freq)
			nSamples = 0
			data = []float64{}
			fmt.Printf("nsamples reset \n")
		}
	}
	chk(stream.Stop())
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
