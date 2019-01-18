package ui

import (
	"fmt"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"k8s.io/api/core/v1"
)

type Table struct {
	Title        string
	Header       string
	c1Width      int
	c2Width      int
	c3Width      int
	c4Width      int
	formatString string
	spinner      *spinner.Spinner
	verbose      bool
}

// this is a UI thing as well - struct Table Row Header - add column widtch

func NewTable(c1Title, c2Title, c3Title, c4Title string, verbose bool) Table {
	table := Table{
		c1Width: 4,
		c2Width: 45,
		c3Width: 45,
		c4Width: 50,
		spinner: spinner.New(spinner.CharSets[9], 100*time.Millisecond),
		verbose: verbose,
	}
	table.formatString = fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds\n", table.c1Width, table.c2Width, table.c3Width, table.c4Width)
	if table.verbose {
		fmt.Printf(table.printRow(c1Title, c2Title, c3Title, c4Title))
	}
	return table
}

func (t Table) printRow(c1, c2, c3, c4 string) string {
	return fmt.Sprintf(t.formatString, c1, c2, c3, c4)
}

func (t Table) PrepareRow() {
	if t.verbose {
		t.spinner.Start()
	}
}

func (t Table) DiscardRow() {
	if t.verbose {
		t.spinner.Stop()
	}
}

func (t Table) CommitRow(c1, c2, c3, c4 string) {
	if t.verbose {
		t.spinner.FinalMSG = t.printRow(c1, c2, c3, c4)
		t.spinner.Stop()
	}
}

func PrintTitle(title string, verbose bool) {
	if verbose {
		color.Yellow(title)
	}
}

func Print(title string, verbose bool) {
	if verbose {
		fmt.Println(title)
	}
}

func PrintBanner(title string) {
	bannerString := banner.PrintS(title)
	color.Yellow(bannerString)
	fmt.Println("")
}

func PrintPodList(pods []v1.Pod, title string, namespace, verbose bool) {
	if verbose {
		color.Yellow(title + "\n")
		for _, pod := range pods {
			if namespace {
				fmt.Printf("> %s (%s)\n", pod.Name, pod.Namespace)
			} else {
				fmt.Printf("> %s\n", pod.Name)
			}
		}
		fmt.Println("")
	}
}

func AskForConfirmation() (bool, error) {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true, nil
	} else if containsString(nokayResponses, response) {
		return false, nil
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
