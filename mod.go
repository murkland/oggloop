package oggloop

import (
	"io"
	"strconv"
	"strings"

	"github.com/jfreymuth/oggvorbis"
)

type Info struct {
	Channels   int
	LoopStart  int64
	LoopLength int64
	Length     int64
}

func ReadInfo(in io.Reader) (Info, error) {
	var info Info

	r, err := oggvorbis.NewReader(in)
	if err != nil {
		return info, err
	}

	info.Channels = r.Channels()
	info.Length = r.Length()

	for _, comment := range r.CommentHeader().Comments {
		sepIndex := strings.IndexRune(comment, '=')
		if sepIndex == -1 {
			continue
		}

		value := comment[sepIndex+1:]

		switch strings.ToLower(comment[:sepIndex]) {
		case "loopstart":
			loopStart, err := strconv.ParseUint(value, 10, 63)
			if err != nil {
				return info, err
			}
			info.LoopStart = int64(loopStart)

		case "looplength":
			loopLength, err := strconv.ParseUint(value, 10, 63)
			if err != nil {
				return info, err
			}
			info.LoopLength = int64(loopLength)
		}
	}

	return info, nil
}
