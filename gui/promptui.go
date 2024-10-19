package gui

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

func RenderInteractiveList(items []Item) Item {
	prompt := promptui.Select{
		Label: "🎧🎧🎧🎧🎧🎧 Select Audio 🎧🎧🎧🎧🎧🎧",
		Items: items,
		Size:  15,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Active:   "🎵🎵🎵 {{ .Name | yellow }}",
			Inactive: "  {{ .Name | green }}",
			Selected: "🎵🎵🎵 {{ .Name | red | cyan }}",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Sad to see u go 😢")
			os.Exit(0)
		}
		log.Fatalf("Prompt failed %v\n", err)
	}

	return items[i]
}
