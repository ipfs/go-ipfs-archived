package dagcmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	path "github.com/ipfs/go-ipfs/path"

	"github.com/ipfs/go-ipfs-cmds/cmdsutil"

	cid "gx/ipfs/QmV5gPoRsjN1Gid3LMdNZTyfCtP2DsvqEbMAmz82RmmiGk/go-cid"
	node "gx/ipfs/QmYDscK7dmdo2GZ9aumS8s5auUUAH5mR1jvj5pYhWusfK7/go-ipld-node"
	ipldcbor "gx/ipfs/QmdaC21UyoyN3t9QdapHZfsaUo3mqVf5p4CEuFaYVFqwap/go-ipld-cbor"
)

var DagCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Interact with ipld dag objects.",
		ShortDescription: `
'ipfs dag' is used for creating and manipulating dag objects.

This subcommand is currently an experimental feature, but it is intended
to deprecate and replace the existing 'ipfs object' command moving forward.
		`,
	},
	Subcommands: map[string]*cmds.Command{
		"put": DagPutCmd,
		"get": DagGetCmd,
	},
}

type OutputObject struct {
	Cid *cid.Cid
}

var DagPutCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Add a dag node to ipfs.",
		ShortDescription: `
'ipfs dag put' accepts input from a file or stdin and parses it
into an object of the specified format.
`,
	},
	Arguments: []cmdsutil.Argument{
		cmdsutil.FileArg("object data", true, false, "The object to put").EnableStdin(),
	},
	Options: []cmdsutil.Option{
		cmdsutil.StringOption("format", "f", "Format that the object will be added as.").Default("cbor"),
		cmdsutil.StringOption("input-enc", "Format that the input object will be.").Default("json"),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		fi, err := req.Files().NextFile()
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		ienc, _, _ := req.Option("input-enc").String()
		format, _, _ := req.Option("format").String()

		switch ienc {
		case "json":
			nd, err := convertJsonToType(fi, format)
			if err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}

			c, err := n.DAG.Add(nd)
			if err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}

			res.SetOutput(&OutputObject{Cid: c})
			return
		case "raw":
			nd, err := convertRawToType(fi, format)
			if err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}

			c, err := n.DAG.Add(nd)
			if err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}

			res.SetOutput(&OutputObject{Cid: c})
			return
		default:
			res.SetError(fmt.Errorf("unrecognized input encoding: %s", ienc), cmdsutil.ErrNormal)
			return
		}
	},
	Type: OutputObject{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())
			oobj, ok := v.(*OutputObject)
			if !ok {
				return nil, fmt.Errorf("expected a different object in marshaler")
			}

			return strings.NewReader(oobj.Cid.String()), nil
		},
	},
}

var DagGetCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Get a dag node from ipfs.",
		ShortDescription: `
'ipfs dag get' fetches a dag node from ipfs and prints it out in the specifed format.
`,
	},
	Arguments: []cmdsutil.Argument{
		cmdsutil.StringArg("ref", true, false, "The object to get").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		p, err := path.ParsePath(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		obj, rem, err := n.Resolver.ResolveToLastNode(req.Context(), p)
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		var out interface{} = obj
		if len(rem) > 0 {
			final, _, err := obj.Resolve(rem)
			if err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}
			out = final
		}

		res.SetOutput(out)
	},
}

func convertJsonToType(r io.Reader, format string) (node.Node, error) {
	switch format {
	case "cbor", "dag-cbor":
		return ipldcbor.FromJson(r)
	case "dag-pb", "protobuf":
		return nil, fmt.Errorf("protobuf handling in 'dag' command not yet implemented")
	default:
		return nil, fmt.Errorf("unknown target format: %s", format)
	}
}

func convertRawToType(r io.Reader, format string) (node.Node, error) {
	switch format {
	case "cbor", "dag-cbor":
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return ipldcbor.Decode(data)
	default:
		return nil, fmt.Errorf("unsupported target format for raw input: %s", format)
	}
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
