// Copyright 2016 The Ebiten Authors

package twenty48

import (
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/ebitenutil"

)

type sounds struct {
	music *audio.Player
	death *audio.Player
	score *audio.Player
}

func newAudio() (*sounds, error) {
	audioContext, err := audio.NewContext(sampleRate)
	if err != nil {
		return nil, err
	}

	f1, err := ebitenutil.OpenFile("chiptronical.ogg")
	if err != nil {
		return nil, err
	}
	f2, err := ebitenutil.OpenFile("death.wav")
	if err != nil {
		return nil, err
	}
	f3, err := ebitenutil.OpenFile("chaChing.wav")
	if err != nil {
		return nil, err
	}

	d1, err := vorbis.Decode(audioContext, f1)
	if err != nil {
		return nil, err
	}
	d2, err := wav.Decode(audioContext, f2)
	if err != nil {
		return nil, err
	}
	d3, err := wav.Decode(audioContext, f3)
	if err != nil {
		return nil, err
	}

	music, err := audio.NewPlayer(audioContext, d1)
	if err != nil {
		return nil, err
	}
	death, err := audio.NewPlayer(audioContext, d2)
	if err != nil {
		return nil, err
	}
	score, err := audio.NewPlayer(audioContext, d3)
	if err != nil {
		return nil, err
	}
	music.SetVolume(0.5)

	return &sounds{
		music: music, death: death, score: score,	
	}, nil
}