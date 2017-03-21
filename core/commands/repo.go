package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
	bstore "github.com/ipfs/go-ipfs/blocks/blockstore"
	cmds "github.com/ipfs/go-ipfs/commands"
	corerepo "github.com/ipfs/go-ipfs/core/corerepo"
	config "github.com/ipfs/go-ipfs/repo/config"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	lockfile "github.com/ipfs/go-ipfs/repo/fsrepo/lock"

	u "gx/ipfs/QmZuY8aV7zbNXVy6DyN9SmnuH3o9nG852F4aTiSBpts8d1/go-ipfs-util"
)

type RepoVersion struct {
	Version string
}

var RepoCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Manipulate the IPFS repo.",
		ShortDescription: `
'ipfs repo' is a plumbing command used to manipulate the repo.
`,
	},

	Subcommands: map[string]*cmds.Command{
		"gc":      repoGcCmd,
		"stat":    repoStatCmd,
		"fsck":    RepoFsckCmd,
		"version": repoVersionCmd,
		"verify":  repoVerifyCmd,
	},
}

var repoGcCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Perform a garbage collection sweep on the repo.",
		ShortDescription: `
'ipfs repo gc' is a plumbing command that will sweep the local
set of stored objects and remove ones that are not pinned in
order to reclaim hard disk space.
`,
	},
	Options: []cmdsutil.Option{
		cmdsutil.BoolOption("quiet", "q", "Write minimal output.").Default(false),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			log.Debug("Run.if: getNode failed with ", err)
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		gcOutChan, err := corerepo.GarbageCollectAsync(n, req.Context())
		if err != nil {
			log.Debug("Run.if: gcAsync failed with ", err)
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		outChan := make(chan interface{})
		res.SetOutput(outChan)

		go func() {
			defer func() {
				close(outChan)
				log.Debug("Run.go.defer: closing outChan")
			}()

			for k := range gcOutChan {
				log.Debugf("Run.go.for: rx %v from gcOutChan, sending to %v", k, outChan)
				outChan <- k
				log.Debugf("Run.go.for: forwarded to outChan.")
			}
		}()
	},
	Type: corerepo.KeyRemoved{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())

			quiet, _, err := res.Request().Option("quiet").Bool()
			if err != nil {
				return nil, err
			}
			obj, ok := v.(*corerepo.KeyRemoved)
			if !ok {
				return nil, u.ErrCast()
			}

			buf := new(bytes.Buffer)
			if quiet {
				buf = bytes.NewBufferString(obj.Key.String() + "\n")
			} else {
				buf = bytes.NewBufferString(fmt.Sprintf("removed %s\n", obj.Key))
			}

			return buf, nil
		},
	},
}

var repoStatCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Get stats for the currently used repo.",
		ShortDescription: `
'ipfs repo stat' is a plumbing command that will scan the local
set of stored objects and print repo statistics. It outputs to stdout:
NumObjects      int Number of objects in the local repo.
RepoPath        string The path to the repo being currently used.
RepoSize        int Size in bytes that the repo is currently taking.
Version         string The repo version.
`,
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		stat, err := corerepo.RepoStat(n, req.Context())
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		res.SetOutput(stat)
	},
	Options: []cmdsutil.Option{
		cmdsutil.BoolOption("human", "Output RepoSize in MiB.").Default(false),
	},
	Type: corerepo.Stat{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())
			stat, ok := v.(*corerepo.Stat)
			if !ok {
				return nil, u.ErrCast()
			}

			human, _, err := res.Request().Option("human").Bool()
			if err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "NumObjects \t %d\n", stat.NumObjects)
			sizeInMiB := stat.RepoSize / (1024 * 1024)
			if human && sizeInMiB > 0 {
				fmt.Fprintf(buf, "RepoSize (MiB) \t %d\n", sizeInMiB)
			} else {
				fmt.Fprintf(buf, "RepoSize \t %d\n", stat.RepoSize)
			}
			fmt.Fprintf(buf, "RepoPath \t %s\n", stat.RepoPath)
			fmt.Fprintf(buf, "Version \t %s\n", stat.Version)

			return buf, nil
		},
	},
}

var RepoFsckCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Remove repo lockfiles.",
		ShortDescription: `
'ipfs repo fsck' is a plumbing command that will remove repo and level db
lockfiles, as well as the api file. This command can only run when no ipfs
daemons are running.
`,
	},
	Run: func(req cmds.Request, res cmds.Response) {
		configRoot := req.InvocContext().ConfigRoot

		dsPath, err := config.DataStorePath(configRoot)
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		dsLockFile := filepath.Join(dsPath, "LOCK") // TODO: get this lockfile programmatically
		repoLockFile := filepath.Join(configRoot, lockfile.LockFile)
		apiFile := filepath.Join(configRoot, "api") // TODO: get this programmatically

		log.Infof("Removing repo lockfile: %s", repoLockFile)
		log.Infof("Removing datastore lockfile: %s", dsLockFile)
		log.Infof("Removing api file: %s", apiFile)

		err = os.Remove(repoLockFile)
		if err != nil && !os.IsNotExist(err) {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}
		err = os.Remove(dsLockFile)
		if err != nil && !os.IsNotExist(err) {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}
		err = os.Remove(apiFile)
		if err != nil && !os.IsNotExist(err) {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		res.SetOutput(&MessageOutput{"Lockfiles have been removed.\n"})
	},
	Type: MessageOutput{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: MessageTextMarshaler,
	},
}

type VerifyProgress struct {
	Message  string
	Progress int
}

var repoVerifyCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Verify all blocks in repo are not corrupted.",
	},
	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdsutil.ErrNormal)
			return
		}

		out := make(chan interface{})
		go func() {
			defer close(out)
			bs := bstore.NewBlockstore(nd.Repo.Datastore())

			bs.HashOnRead(true)

			keys, err := bs.AllKeysChan(req.Context())
			if err != nil {
				log.Error(err)
				return
			}

			var fails int
			var i int
			for k := range keys {
				_, err := bs.Get(k)
				log.Debugf("Run.go.for: bs.Get(%s) returned err=%s", k, err)
				if err != nil {
					out <- &VerifyProgress{
						Message: fmt.Sprintf("block %s was corrupt (%s)", k, err),
					}
					fails++
				}
				i++
				out <- &VerifyProgress{Progress: i}
			}
			if fails == 0 {
				log.Debug("Run.go.else: all ok")
				out <- &VerifyProgress{Message: "verify complete, all blocks validated."}
			} else {
				//out <- &VerifyProgress{Message: "verify complete, some blocks were corrupt."}
				log.Debugf("Run.go.else: %T.SetError", res)
				res.SetError(fmt.Errorf("verify complete, some blocks were corrupt"), cmdsutil.ErrNormal)
			}
		}()

		res.SetOutput((<-chan interface{})(out))
	},
	Type: VerifyProgress{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())

			obj, ok := v.(*VerifyProgress)
			if !ok {
				return nil, u.ErrCast()
			}

			buf := new(bytes.Buffer)
			if obj.Message != "" {
				if strings.Contains(obj.Message, "blocks were corrupt") {
					return nil, fmt.Errorf(obj.Message)
				}
				if len(obj.Message) < 20 {
					obj.Message += "             "
				}
				fmt.Fprintln(buf, obj.Message)
				return buf, nil
			}

			fmt.Fprintf(buf, "%d blocks processed.\r", obj.Progress)
			return buf, nil
		},
	},
}

var repoVersionCmd = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline: "Show the repo version.",
		ShortDescription: `
'ipfs repo version' returns the current repo version.
`,
	},

	Options: []cmdsutil.Option{
		cmdsutil.BoolOption("quiet", "q", "Write minimal output."),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		res.SetOutput(&RepoVersion{
			Version: fmt.Sprint(fsrepo.RepoVersion),
		})
	},
	Type: RepoVersion{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v := unwrapOutput(res.Output())
			response := v.(*RepoVersion)

			quiet, _, err := res.Request().Option("quiet").Bool()
			if err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)
			if quiet {
				buf = bytes.NewBufferString(fmt.Sprintf("fs-repo@%s\n", response.Version))
			} else {
				buf = bytes.NewBufferString(fmt.Sprintf("ipfs repo version fs-repo@%s\n", response.Version))
			}
			return buf, nil

		},
	},
}
