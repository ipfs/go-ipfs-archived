package commands

import (
	"fmt"
	"io"
	"os"

	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
	core "github.com/ipfs/go-ipfs/core"
	coreunix "github.com/ipfs/go-ipfs/core/coreunix"

	context "context"
)

const progressBarMinSize = 1024 * 1024 * 8 // show progress bar for outputs > 8MiB

var CatCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline:          "Show IPFS object data.",
		ShortDescription: "Displays the data contained by an IPFS or IPNS object(s) at the given path.",
	},

	Arguments: []cmdsutil.Argument{
		cmdsutil.StringArg("ipfs-path", true, true, "The path to the IPFS object(s) to be outputted.").EnableStdin(),
	},
	Run: func(req cmds.Request, re cmds.ResponseEmitter) {
		log.Debugf("cat: RespEm type is %T", re)
		node, err := req.InvocContext().GetNode()
		if err != nil {
			re.SetError(err, cmdsutil.ErrNormal)
			return
		}

		if !node.OnlineMode() {
			if err := node.SetupOfflineRouting(); err != nil {
				re.SetError(err, cmdsutil.ErrNormal)
				return
			}
		}

		readers, length, err := cat(req.Context(), node, req.Arguments())
		log.Debug("cat returned ", err)

		if err != nil {
			re.SetError(err, cmdsutil.ErrNormal)
			return
		}

		/*
			if err := corerepo.ConditionalGC(req.Context(), node, length); err != nil {
				res.SetError(err, cmdsutil.ErrNormal)
				return
			}
		*/

		re.SetLength(length)

		reader := io.MultiReader(readers...)
		re.Emit(reader)
		re.Close()
	},
	PostRun: map[cmds.EncodingType]func(cmds.Request, cmds.ResponseEmitter) cmds.ResponseEmitter{
		cmds.CLI: func(req cmds.Request, re cmds.ResponseEmitter) cmds.ResponseEmitter {
			re_, res := cmds.NewChanResponsePair(req)

			go func() {
				if res.Length() > 0 && res.Length() < progressBarMinSize {
					log.Debugf("cat/res.Length() == %v < progressBarMinSize", res.Length())
					cmds.Copy(re, res)
					return
				}

				// Copy closes by itself, so we must not do this before
				defer re.Close()

				v, err := res.Next()
				log.Debugf("cat/res.Next() returned (%v, %v)", v, err)
				if err != nil {
					if err == cmds.ErrRcvdError {
						re.SetError(res.Error().Message, res.Error().Code)
					} else {
						re.SetError(res.Error(), cmdsutil.ErrNormal)
					}

					return
				}

				r, ok := v.(io.Reader)
				if !ok {
					re.SetError(fmt.Sprintf("expected io.Reader, not %T", v), cmdsutil.ErrNormal)
					return
				}

				bar, reader := progressBarForReader(os.Stderr, r, int64(res.Length()))
				bar.Start()

				re.Emit(reader)
			}()

			return re_
		},
	},
}

func cat(ctx context.Context, node *core.IpfsNode, paths []string) ([]io.Reader, uint64, error) {
	readers := make([]io.Reader, 0, len(paths))
	length := uint64(0)
	for _, fpath := range paths {
		read, err := coreunix.Cat(ctx, node, fpath)
		if err != nil {
			return nil, 0, err
		}
		readers = append(readers, read)
		length += uint64(read.Size())
	}
	return readers, length, nil
}
