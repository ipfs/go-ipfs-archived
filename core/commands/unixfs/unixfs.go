package unixfs

import (
	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
	cmds "github.com/ipfs/go-ipfs/commands"
)

var UnixFSCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Interact with IPFS objects representing Unix filesystems.",
		ShortDescription: `
'ipfs file' provides a familiar interface to file systems represented
by IPFS objects, which hides ipfs implementation details like layout
objects (e.g. fanout and chunking).
`,
		LongDescription: `
'ipfs file' provides a familiar interface to file systems represented
by IPFS objects, which hides ipfs implementation details like layout
objects (e.g. fanout and chunking).
`,
	},

	Subcommands: map[string]*cmds.Command{
		"ls": LsCmd,
	},
}

// copy+pasted from ../commands.go
func unwrapOutput(i interface{}) interface{} {
	var (
		ch <-chan interface{}
		ok bool
	)

	if ch, ok = i.(<-chan interface{}); !ok {
		if ch, ok = i.(chan interface{}); !ok {
			return nil
		}
	}

	return <-ch
}
