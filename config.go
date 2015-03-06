// Copyright (c) 2015 Bertrand Janin <b@janin.com>
// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jessevdk/go-flags"
)

// Cmd is a singleton storing all the command-line parameters.
type Cmd struct {
	ConfigFile string `short:"c" description:"Configuration file" default:"/etc/asoi.conf"`
}

// Cfg is a singleton storing all the config file parameters.
type Cfg struct {
	// IRCServer is the hostname and port of the IRC server.
	IRCServer string

	// IRCNickname is the nickname of the bot, passed upon connction.
	IRCNickname string

	// Channels is the list of channels to auto-matically join.
	Channels []string

	// Topics is the list of all the status topics on the AWS status page.
	// This is basically the name of the RSS feed, minus the path and minus
	// the .rss extention.
	Topics []string
}

var (
	cfg = Cfg{}
	cmd = Cmd{}
)

// Look in the current directory for an config.json file.
func parseConfigFile() error {
	file, err := os.Open(cmd.ConfigFile)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	if cfg.IRCNickname == "" {
		return errors.New("'IRCNickname' is not defined")
	}

	if cfg.IRCServer == "" {
		return errors.New("'IRCServer' is not defined")
	}

	if len(cfg.Channels) == 0 {
		return errors.New("'Channels' is not defined")
	}

	if len(cfg.Topics) == 0 {
		return errors.New("'Topics' is not defined")
	}

	return nil
}

// Parse the command line arguments and populate the global cmd struct.
func parseCommandLine() {
	flagParser := flags.NewParser(&cmd, flags.PassDoubleDash)
	_, err := flagParser.Parse()
	if err != nil {
		println("command line error: " + err.Error())
		flagParser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
}
