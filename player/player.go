package player

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/amrelhewy09/drive_audio_streamer/client"
	"github.com/amrelhewy09/drive_audio_streamer/gui"
)

type AudioPlayer struct {
	client client.RemoteFileClient
}

func NewPlayer(client client.RemoteFileClient) Player {

	return &AudioPlayer{
		client: client,
	}
}

func (r *AudioPlayer) Play(item *gui.Item, guiPlayer *gui.Gui) {
	rPipe, wPipe, cmd := initFFPLAY()
	const bufferSize = 1024 * 1024 * 2 // STILL UNDER EXPERIMENTATION
	var start int64 = 0
	dataChan := make(chan []byte, 20) // STILL UNDER EXPERIMENTATION
	errChan := make(chan error)

	terminateDownloadChan := make(chan int)
	terminateExecutionChan := make(chan int)
	go func() {
		for {
			select {
			case <-terminateDownloadChan:
				// Received termination signal, exit the loop
				rPipe.Close()
				wPipe.Close()
				close(errChan)
				return
			default:
				end := start + bufferSize
				if end > item.Size {
					end = item.Size
				}

				stream, err := r.client.Stream(item.Id, start, end)
				if err != nil {
					if err == io.EOF {
						close(dataChan)
						return
					}
					errChan <- err
					return
				}

				data, err := io.ReadAll(stream)
				if err != nil {
					errChan <- err
					return
				}

				dataChan <- data
				start += bufferSize
				guiPlayer.UpdateProgress(int(float64(start) / float64(item.Size) * 100))

				if start >= item.Size {
					close(dataChan)
					return
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					fmt.Println("Data channel closed")
				} else {
					if _, err := io.Copy(wPipe, bytes.NewReader(data)); err != nil {
						return
					}
				}
			case <-terminateExecutionChan:
				// Received termination signal, exit the loop
				return
			}
		}
	}()

	go func() {
		for err := range errChan {
			log.Fatalf("Error occurred: %v", err)
		}
	}()
	guiPlayer.RenderMediaPlayerUI(terminateDownloadChan, terminateExecutionChan, cmd)
}

func (r *AudioPlayer) Pause()  {}
func (r *AudioPlayer) Resume() {}

func initFFPLAY() (*io.PipeReader, *io.PipeWriter, *exec.Cmd) {
	rPipe, wPipe := io.Pipe()
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-i", "pipe:0")
	cmd.Stdin = rPipe
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start ffplay: %v", err)
	}
	return rPipe, wPipe, cmd
}
