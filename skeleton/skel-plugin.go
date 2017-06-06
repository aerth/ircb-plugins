// +build cgo

// Skeleton plugin adds "plugin-test" public command to ircb
package main

import "github.com/aerth/ircb"

// empty main
func main() {}

// Init must be named Init taking an ircb.Connection and returning an error
// Useful methods of ircb.Connection include:
//  c.AddCommand(name, fn)
//  c.AddMasterCommand(name, fn)
//  c.RemoveCommand(name)
//  c.RemoveMasterCommand(name)
//  c.SendMaster(format, ...interface{})
func Init(c *ircb.Connection) error {
	name := "skeleton" // no global variables in plugin
	c.Log.Println("plugin loading:", name)
	c.AddCommand("plugin-test", bar)
	c.Log.Println("plugin loaded:", name)
	return nil
}

// commands must have the following signature
func bar(c *ircb.Connection, irc *ircb.IRC) {
	irc.Reply(c, "plugins work!")
}
