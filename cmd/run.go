/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/amrelhewy09/drive_audio_streamer/client"
	"github.com/amrelhewy09/drive_audio_streamer/gui"
	"github.com/amrelhewy09/drive_audio_streamer/player"

	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "Run attempts to initialize oauth with your remote client and browse audio files to stream afterwards",
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewGoogleDriveClient()
		err := client.ConfigureRemoteFileDestination()
		if err != nil {
			log.Fatalf("Unable to configure remote file destination: %v", err)
		}

		files, err := client.List()

		if err != nil {
			log.Fatalf("Unable to list files: %v", err)
		}

		audioPlayer := player.NewPlayer(client)

		for {
			selectedFile := gui.RenderInteractiveList(items(files))
			audioPlayer.Play(&gui.Item{Id: selectedFile.Id, Name: selectedFile.Name, Size: selectedFile.Size}, gui.CreatePlayerGui(selectedFile.Name))
		}
	},
}

func init() {
	rootCmd.AddCommand(run)
}

func items(files []*client.File) []gui.Item {
	items := make([]gui.Item, len(files))

	for i, file := range files {
		items[i] = gui.Item{
			Name: file.Name,
			Id:   file.Id,
			Size: file.Size,
		}
	}

	return items

}
