package mangokong

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/muesli/mango"
	"github.com/muesli/roff"
)

// ManFlag can be used as a Kong flag that will write a default man page to stdout.
type ManFlag bool

func (m ManFlag) BeforeApply(app *kong.Kong) error {
	man := NewManPage(1, app.Model)
	fmt.Fprint(app.Stdout, man.Build(roff.NewDocument()))
	app.Exit(0)
	return nil
}

// NewManPage from a Kong Application.
func NewManPage(section uint, app *kong.Application) *mango.ManPage {
	help := strings.Split(strings.TrimSpace(app.Help), "\n")
	man := mango.NewManPage(section, app.Name, help[0])
	if app.Detail != "" {
		man = man.WithLongDescription(app.Detail)
	} else {
		man = man.WithLongDescription(strings.TrimSpace(strings.Join(help[1:], "\n")))
	}
	addCommand(man, app.Node, nil)
	return man
}

func addCommand(man *mango.ManPage, node *kong.Node, parent *mango.Command) {
	var item *mango.Command
	name := node.Summary()
	if parent == nil {
		name = node.Name
	}
	item = mango.NewCommand(name, node.Help, node.Detail)
	if parent == nil {
		man.Root = *item
		item = &man.Root
	} else {
		err := parent.AddCommand(item)
		if err != nil { // Kong already validates structure, so this should never fail.
			panic(err)
		}
	}

	for _, child := range node.Children {
		if child.Hidden {
			continue
		}
		addCommand(man, child, item)
	}

	for _, flag := range node.Flags {
		short := ""
		if flag.Short != 0 {
			short = fmt.Sprintf("%c", flag.Short)
		}
		item.AddFlag(mango.Flag{
			Name:  strings.TrimPrefix(flag.Summary(), "--"),
			Short: short,
			Usage: flag.Help,
			PFlag: true,
		})
	}
}
