package ioutil

import "github.com/pterm/pterm"

var (
	PTermParagraph   = pterm.DefaultParagraph
	PTermHeader      = pterm.DefaultHeader.WithFullWidth()
	PTermBarChart    = pterm.DefaultBarChart.WithHorizontalBarCharacter("▆").WithHorizontal().WithShowValue()
	PTermTable       = pterm.DefaultTable.WithHasHeader()
	PTermProgressbar = pterm.DefaultProgressbar.WithRemoveWhenDone()
)
