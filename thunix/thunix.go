// Skeleton plugin adds "plugin-test" public command to ircb
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

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
	name := "thunix" // no global variables in plugin
	c.Log.Println("plugin loading:", name)
	c.AddCommand("shells", commandShells)
	c.AddCommand("uptime", commandHostUptime)
	c.Log.Println("plugin loaded:", name)
	return nil
}

func commandShells(c *ircb.Connection, irc *ircb.IRC) {
	var shells int
	letters, err := ioutil.ReadDir("/home")
	if err != nil {
		c.Log.Println(err)
		return
	}

	for _, letter := range letters {
		if letter.IsDir() {
			userhomes, err := ioutil.ReadDir("/home/" + letter.Name())
			if err != nil {
				c.Log.Println(err)
				continue
			}
			shells += len(userhomes)
		}

	}

	irc.Reply(c,
		fmt.Sprintf("Current number of shell accounts: %v", shells))

}

func commandHostUptime(c *ircb.Connection, irc *ircb.IRC) {
	// output of uptime command
	uptime := exec.Command("/usr/bin/uptime")

	out, err := uptime.CombinedOutput()
	if err != nil {
		c.Log.Println(irc, err)
		c.SendMaster("%s", err)
	}

	output := strings.Split(string(out), "\n")[0]
	if strings.TrimSpace(output) != "" {
		irc.Reply(c, output)
	}
}
