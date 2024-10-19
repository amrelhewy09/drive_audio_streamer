package player

import (
	"github.com/amrelhewy09/drive_audio_streamer/gui"
)

type Player interface {
	Play(*gui.Item, *gui.Gui)
	Pause()
	Resume()
}
