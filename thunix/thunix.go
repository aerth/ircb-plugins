// Skeleton plugin adds "plugin-test" public command to ircb
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aerth/ircb"
)

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
	name := "thunar" // no global variables in plugin
	c.Log.Println("plugin loading:", name)
	c.AddCommand("shells", commandShells)
	c.Log.Println("plugin loaded:", name)
	return nil
}

func commandShells(c *ircb.Connection, irc *ircb.IRC) {
	dir, err := ioutil.ReadDir("/home")
	if err != nil {
		c.Log.Println(err)
		return
	}

	irc.Reply(c, fmt.Sprintf("Current number of shell accounts: %v", len(dir)))

}
