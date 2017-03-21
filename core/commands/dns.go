package commands

import (
	"io"
	"strings"

	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
	cmds "github.com/ipfs/go-ipfs/commands"
	namesys "github.com/ipfs/go-ipfs/namesys"
	util "gx/ipfs/QmZuY8aV7zbNXVy6DyN9SmnuH3o9nG852F4aTiSBpts8d1/go-ipfs-util"
)

var DNSCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Resolve DNS links.",
		ShortDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, IPFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.
`,
		LongDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, IPFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.

For example, with this DNS TXT record:

	> dig +short TXT _dnslink.ipfs.io
	dnslink=/ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy

The resolver will give:

	> ipfs dns ipfs.io
	/ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy

The resolver can recursively resolve:

	> dig +short TXT recursive.ipfs.io
	dnslink=/ipns/ipfs.io
	> ipfs dns -r recursive.ipfs.io
	/ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy
`,
	},

	Arguments: []cmdsutil.Argument{
		cmdsutil.StringArg("domain-name", true, false, "The domain-name name to resolve.").EnableStdin(),
	},
	Options: []cmdsutil.Option{
		cmdsutil.BoolOption("recursive", "r", "Resolve until the result is not a DNS link.").Default(false),
	},
	Run: func(req cmds.Request, res cmds.Response) {

		recursive, _, _ := req.Option("recursive").Bool()
		name := req.Arguments()[0]
		resolver := namesys.NewDNSResolver()

		depth := 1
		if recursive {
			depth = namesys.DefaultDepthLimit
		}
		output, err := resolver.ResolveN(req.Context(), name, depth)
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}
		res.SetOutput(&ResolvedPath{output})
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())
			output, ok := v.(*ResolvedPath)
			if !ok {
				return nil, util.ErrCast()
			}
			return strings.NewReader(output.Path.String() + "\n"), nil
		},
	},
	Type: ResolvedPath{},
}
