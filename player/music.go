package player

import (
	"embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io/fs"
	"log"
	"time"
)

//go:embed sound
var Sound embed.FS

var done = make(chan bool)
var initialized = false

func init() {

}

func play(file fs.File) {
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		log.Println(err)
	}

	if !initialized {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			panic(err)
		}
		initialized = true
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {

		done <- true
	})))

	<-done

}

func PlaySound(str string) {
	switch str {
	case "J":
		fileJ, _ := Sound.Open("sound/j.mp3")
		play(fileJ)
	case "N":
		fileN, _ := Sound.Open("sound/n.mp3")
		play(fileN)
	case "T":
		fileT, _ := Sound.Open("sound/t.mp3")
		play(fileT)
	case "M":
		fileM, _ := Sound.Open("sound/m.mp3")
		play(fileM)
	case "JNTM":
		fileNGM, _ := Sound.Open("sound/ngm.mp3")
		play(fileNGM)
	}

}
