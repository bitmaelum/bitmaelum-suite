package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strings"
)

func SelectAndRun(cmd *cobra.Command, args []string) {
	items := make([]*cobra.Command, len(cmd.Commands()))
	copy(items, cmd.Commands())

	if cmd == cmd.Root() {
		items = append(items, &cobra.Command{
			Use:   "quit",
			Short: "Quit program",
		})
	} else {
		items = append(items, &cobra.Command{
			Use:   "back",
			Short: "Back to previous menu",
		})
	}


	for {
		// Generate path
		pc := cmd
		var path []string
		//path = append(path, pc.Use)
		for pc != nil {
			path = append(path, pc.Use)
			pc = pc.Parent()
		}

		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}



		fmt.Println("\033[2J")
		prompt := promptui.Select{
			Label: `Welcome to MailV2. Please select your commands below`,
			Items: items,
			Templates: &promptui.SelectTemplates{
				Label:    fmt.Sprintf("Path: %s", strings.Join(path[1:
					], " / ")),
				Active:   `❯ {{ printf "%-20s" .Use | cyan | red }}  {{ printf "%-30s" .Short | yellow | red }}`,
				Inactive: `  {{ printf "%-20s" .Use | cyan  }}  {{ printf "%-30s" .Short | yellow }}`,
				Selected: ``,
			},
		}
		idx, _, err := prompt.Run()
		if err != nil {
			panic(err)
		}

		if (idx >= len(cmd.Commands())) {
			return
		}

		if items[idx].Use == "quit" {
			// Quit program
			return
		} else if items[idx].Use == "back" {
			// return
			return
		} else {
			cmd.Commands()[idx].Run(cmd.Commands()[idx], args)
		}
	}
}