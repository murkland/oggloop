package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/nbarena/formats/oggloop"
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	m, format, info, err := oggloop.Load(f)
	if err != nil {
		log.Fatalf("failed to load music: %v", err)
	}

	log.Printf("playing %s, sr: %d Hz, length: %s, loop: [%s, %s)", f.Name(), format.SampleRate, format.SampleRate.D(int(info.Length)), format.SampleRate.D(int(info.LoopStart)), format.SampleRate.D(int(info.LoopStart+info.LoopLength)))

	resampler := beep.ResampleRatio(4, 1.0, m)
	speaker.Init(format.SampleRate, 128)
	speaker.Play(resampler)

	// Also need to downpitch by 1 semitone, but this is actually kinda hard.
	spedUp := false
	for {
		fmt.Scanln()
		if !spedUp {
			resampler.SetRatio(1.1)
		} else {
			resampler.SetRatio(1.)
		}
		spedUp = !spedUp
		log.Printf("sped up: %v", spedUp)
	}
}
