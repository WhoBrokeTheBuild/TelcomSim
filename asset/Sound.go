package asset

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

func init() {
}

// Sound represents a beep Sound
type Sound struct {
	Stream beep.StreamSeekCloser
	Format beep.Format
}

// NewSoundFromFile returns a new Sound from the given file
func NewSoundFromFile(filename string) (*Sound, error) {
	s := &Sound{}
	err := s.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// LoadFromFile loads a Sound from a given file
func (s *Sound) LoadFromFile(filename string) error {
	log.Loadf("asset.Sound [%v]", filename)
	b, err := data.Asset(filename)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filename)
	if ext == ".mp3" {
		s.Stream, s.Format, err = mp3.Decode(ioutil.NopCloser(bytes.NewReader(b)))
		if err != nil {
			return err
		}
	} else if ext == ".wav" {
		s.Stream, s.Format, err = wav.Decode(ioutil.NopCloser(bytes.NewReader(b)))
		if err != nil {
			return err
		}
	} else if ext == ".ogg" {
		s.Stream, s.Format, err = vorbis.Decode(ioutil.NopCloser(bytes.NewReader(b)))
		if err != nil {
			return err
		}
	} else if ext == ".flac" {
		s.Stream, s.Format, err = flac.Decode(ioutil.NopCloser(bytes.NewReader(b)))
		if err != nil {
			return err
		}
	}

	return nil
}

// Play plays the sound on the default speaker
func (s *Sound) Play() {
	speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/10))
	speaker.Play(s.Stream)
}
