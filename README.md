# go-ipfs-archive

In order to keep the [ipfs/go-ipfs](https://github.com/ipfs/go-ipfs) project relatively easy to navigate, manage, and quick to clone, we periodically remove branches that nobody is actively working on. We think keeping that work available is important, though, so we archive those branches here.

Branches archived here get renamed to: `<branchname>/<archival date>`. For example, `feat/fix-offline-mount` → `feat/fix-offline-mount/2017-08-10` if it was archived on August 10, 2017.

If you need to use any of this archived work, add this repo as a remote to your local clone of `go-ipfs` and and then check out the branch you need:

```sh
# Add this repo as a remote
$ git remote add archived git@github.com:ipfs/go-ipfs-archived.git
# Check out the branch you need, for example:
$ git checkout feat/fix-offline-mount/2017-08-10
```

If you’re curious about or want to propose changes to the archival process, take a look at the [`bin/archive-branches.sh` script](https://github.com/ipfs/go-ipfs/blob/master/bin/archive-branches.sh) in `go-ipfs`.
