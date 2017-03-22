// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package client

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
)

// CmdSimpleFSMove is the 'fs list' command.
type CmdSimpleFSMove struct {
	libkb.Contextified
	src         []keybase1.Path
	dest        keybase1.Path
	interactive bool
	force       bool
	opCanceler  *OpCanceler
}

var _ Canceler = (*CmdSimpleFSMove)(nil)

// NewCmdSimpleFSMove creates a new cli.Command.
func NewCmdSimpleFSMove(cl *libcmdline.CommandLine, g *libkb.GlobalContext) cli.Command {
	return cli.Command{
		Name:         "mv",
		ArgumentHelp: "<source> [source] <dest>",
		Usage:        "move one or more directory elements to dest",
		Action: func(c *cli.Context) {
			cl.ChooseCommand(&CmdSimpleFSMove{
				Contextified: libkb.NewContextified(g),
				opCanceler:   NewOpCanceler(g),
			}, "mv", c)
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "i, interactive",
				Usage: "Prompt before overwrite",
			},
			cli.BoolFlag{
				Name:  "f, force",
				Usage: "force overwrite",
			},
		},
	}
}

// Run runs the command in client/server mode.
func (c *CmdSimpleFSMove) Run() error {
	cli, err := GetSimpleFSClient(c.G())
	if err != nil {
		return err
	}

	ctx := context.TODO()

	isDestDir, destPathString, err := checkPathIsDir(ctx, cli, c.dest)
	if err != nil {
		return err
	}

	destPaths, err := doSimpleFSGlob(c.G(), ctx, cli, c.src)
	if err != nil {
		return err
	}

	for _, src := range destPaths {
		c.G().Log.Debug("SimpleFSMove %s -> %s, %v", pathToString(src), destPathString, isDestDir)

		dest, err := makeDestPath(c.G(), ctx, cli, src, c.dest, isDestDir, destPathString)

		if err == ErrTargetFileExists {
			if c.interactive == true {
				err = doOverwritePrompt(c.G(), pathToString(dest))
			} else if c.force == true {
				err = nil
			}
		}

		if err != nil {
			return err
		}
		c.G().Log.Debug("SimpleFSMove %s -> %s", pathToString(src), pathToString(dest))

		// Don't spawn new jobs if we've been cancelled.
		// TODO: This is still a race condition, if we get cancelled immediately after.
		if c.opCanceler.IsCancelled() {
			break
		}

		opid, err := cli.SimpleFSMakeOpid(ctx)
		if err != nil {
			return err
		}
		c.opCanceler.AddOp(opid)
		defer cli.SimpleFSClose(ctx, opid)

		err = cli.SimpleFSMove(ctx, keybase1.SimpleFSMoveArg{
			OpID: opid,
			Src:  src,
			Dest: dest,
		})
		if err != nil {
			break
		}
		err = cli.SimpleFSWait(ctx, opid)
		if err != nil {
			break
		}
	}
	return err
}

// ParseArgv does nothing for this command.
func (c *CmdSimpleFSMove) ParseArgv(ctx *cli.Context) error {
	var err error
	c.interactive = ctx.Bool("interactive")
	c.force = ctx.Bool("force")

	if c.force && c.interactive {
		return errors.New("force and interactive are incompatible")
	}

	c.src, c.dest, err = parseSrcDestArgs(c.G(), ctx, "mv")
	return err
}

// GetUsage says what this command needs to operate.
func (c *CmdSimpleFSMove) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:    true,
		KbKeyring: true,
		API:       true,
	}
}

func (c *CmdSimpleFSMove) Cancel() error {
	return c.opCanceler.Cancel()
}
