package main

import (
	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
)

var cli struct {
	Man  mangokong.ManFlag `help:"Print man page."`
	Dest string            `help:"Destination directory." placeholder:"DIR"`
	Log  string            `help:"Log level (${enum})." enum:"debug,info,warn,error" default:"debug"`

	Info   InfoCmd   `cmd:"" help:"Info command."`
	Delete DeleteCmd `cmd:"" help:"Delete command."`
}

type InfoCmd struct {
	System InfoSystemCmd `cmd:"" help:"Info on the system."`
	User   InfoUserCmd   `cmd:"" help:"Info on users."`
}

type InfoSystemCmd struct{}

type InfoUserCmd struct {
	User []string `arg:"" required:"" help:"User to retrieve info on."`
}

type DeleteCmd struct {
	User DeleteUserCmd `cmd:"" help:"Delete users."`
}

type DeleteUserCmd struct {
	Purge bool     `help:"Completely purge all history for the user."`
	User  []string `arg:"" required:"" help:"User to delete."`
}

func main() {
	kong.Parse(&cli, kong.HelpOptions{}, kong.Description(`
A test application for mango-kong.

And this is more detail about what this does.
`))
}
