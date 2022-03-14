//go:build !nobeep

package oggloop

import (
	"io"

	"github.com/faiface/beep"
	"github.com/faiface/beep/vorbis"
)

type streamerAndCloser struct {
	beep.Streamer
	io.Closer
}

type interval struct {
	s     beep.StreamSeeker
	start int
	end   int
}

func (i *interval) Err() error {
	return i.s.Err()
}

func (i *interval) Len() int {
	return i.end - i.start
}

func (i *interval) Position() int {
	return i.s.Position() - i.start
}

func (i *interval) Seek(p int) error {
	return i.s.Seek(i.start + p)
}

func (i *interval) Stream(samples [][2]float64) (int, bool) {
	if i.Position() < 0 {
		if err := i.Seek(0); err != nil {
			return 0, false
		}
	}

	n := len(samples)
	if m := i.Len() - i.Position(); n > m {
		n = m
	}

	if n == 0 {
		return 0, false
	}

	return i.s.Stream(samples[:n])
}

func Load(in io.ReadSeekCloser) (beep.StreamCloser, beep.Format, Info, error) {
	info, err := ReadInfo(in)
	if err != nil {
		return nil, beep.Format{}, info, err
	}

	in.Seek(0, io.SeekStart)
	stream, format, err := vorbis.Decode(in)
	if err != nil {
		return nil, format, info, err
	}

	if info.LoopLength == 0 {
		return stream, format, info, nil
	}

	return Wrap(stream, info), format, info, nil
}

func Wrap(stream beep.StreamSeekCloser, info Info) beep.StreamCloser {
	return &streamerAndCloser{beep.Seq(beep.Take(int(info.LoopStart), stream), beep.Loop(-1, &interval{stream, int(info.LoopStart), int(info.LoopStart + info.LoopLength)})), stream}
}
