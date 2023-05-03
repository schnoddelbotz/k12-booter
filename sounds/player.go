package sounds

import (
	"bytes"
	"embed"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

//go:embed *.mp3
var content embed.FS

const (
	Maelstrom_LaserZap        = "snd_100.wav.mp3"
	Maelstrom_LaserMetalPling = "snd_101.wav.mp3"
	Maelstrom_Explosion       = "snd_102.wav.mp3"
	Maelstrom_SpaceScream     = "snd_103.wav.mp3"
	Maelstrom_Pong            = "snd_104.wav.mp3"
	Maelstrom_Ping            = "snd_105.wav.mp3"
	Maelstrom_MaleOop         = "snd_106.wav.mp3"
	Maelstrom_FemaleMoaning   = "snd_107.wav.mp3"
	Maelstrom_IronBall        = "snd_108.wav.mp3"
	Maelstrom_Blip            = "snd_109.wav.mp3"
	Maelstrom_HardRockIntro   = "snd_110.wav.mp3"
	Maelstrom_Beamer          = "snd_111.wav.mp3"
	Maelstrom_AllRight        = "snd_112.wav.mp3"
	Maelstrom_Laughing        = "snd_113.wav.mp3"
	Maelstrom_Yo              = "snd_114.wav.mp3"
	Maelstrom_Warp            = "snd_115.wav.mp3"
	Maelstrom_Yahoo           = "snd_116.wav.mp3"
	Maelstrom_AreYouAsleep    = "snd_117.wav.mp3"
	Maelstrom_AlertSound      = "snd_118.wav.mp3"
	Maelstrom_Beamer2         = "snd_119.wav.mp3"
	Maelstrom_SpringBack      = "snd_120.wav.mp3"
	Maelstrom_Crushing        = "snd_121.wav.mp3"
	Maelstrom_BrokenGlass     = "snd_122.wav.mp3"
	Maelstrom_ExplosionLong   = "snd_123.wav.mp3"
	Maelstrom_Sweet           = "snd_124.wav.mp3"
	Maelstrom_HelpMe          = "snd_125.wav.mp3"
	Maelstrom_ThankGoodness   = "snd_126.wav.mp3"
	Maelstrom_LaserSplat      = "snd_127.wav.mp3"
	Maelstrom_Siren           = "snd_128.wav.mp3"
	Maelstrom_HotNail         = "snd_131.wav.mp3"
	Maelstrom_Thrust          = "snd_132.wav.mp3"
	Maelstrom_LaserZapShort   = "snd_133.wav.mp3"
	Maelstrom_Chord           = "snd_134.wav.mp3"
	Maelstrom_YouIdiot        = "snd_135.wav.mp3"
	Maelstrom_RisingBeeps     = "snd_136.wav.mp3"
)

// InitAudio() and PlayIt() came here by splitting
// https://github.com/hajimehoshi/oto/blob/main/example/main.go
// Thank you!

func InitAudio(disableAudio bool) *oto.Context {
	if disableAudio {
		return nil
	}
	samplingRate := 44100
	numOfChannels := 2
	audioBitDepth := 2
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan
	return otoCtx
}

func PlayIt(mp3Name string, otoCtx *oto.Context) {
	if otoCtx == nil {
		return
	}
	fileBytes, err := content.ReadFile(mp3Name)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}
	fileBytesReader := bytes.NewReader(fileBytes)
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	player := otoCtx.NewPlayer(decodedMp3)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
