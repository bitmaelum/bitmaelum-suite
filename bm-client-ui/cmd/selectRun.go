package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strconv"
	"strings"
)

type commandSorterByPosition []*cobra.Command

func (c commandSorterByPosition) Len() int           { return len(c) }
func (c commandSorterByPosition) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c commandSorterByPosition) Less(i, j int) bool {
	si, _ := strconv.Atoi(c[i].Annotations["position"])
	sj, _ := strconv.Atoi(c[j].Annotations["position"])
	return si < sj
}

func SelectAndRun(cmd *cobra.Command, args []string) {
	// Add all non-hidden commands (gets rid of "help")
	var items []*cobra.Command
	for i := range cmd.Commands() {
		if cmd.Commands()[i].Hidden {
			continue
		}
		items = append(items, cmd.Commands()[i])
	}

	// sort commands
	sort.Sort(commandSorterByPosition(items))

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

		prompt := promptui.Select{
			Label: `Welcome to BitMaelum. Please select your commands below`,
			Items: items,
			Templates: &promptui.SelectTemplates{
				Label:    fmt.Sprintf("Path: %s", strings.Join(path[1:], " / ")),
				Active:   `❯ {{ printf "%-20s" .Use | cyan | red }}  {{ printf "%-30s" .Short | yellow | red }}`,
				Inactive: `  {{ printf "%-20s" .Use | cyan  }}  {{ printf "%-30s" .Short | yellow }}`,
				Selected: ``,
			},
		}
		idx, _, err := prompt.Run()
		if err != nil {
			if err != promptui.ErrInterrupt {
				log.Fatal(err)
			}
			continue
		}

		if idx >= len(cmd.Commands()) {
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
