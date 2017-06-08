// Skeleton plugin adds "plugin-test" public command to ircb
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	name := "go play" // no global variables in plugin
	c.Log.Println("plugin loading:", name)
	c.AddCommand("go", commandGo)
	c.Log.Println("plugin loaded:", name)
	return nil
}

// commands must have the following signature
func commandGo(c *ircb.Connection, irc *ircb.IRC) {
	if len(irc.Arguments) < 1 {
		return
	}
	resp, err := sendToCompiler(c.HTTPClient, lineToMainFunc([]byte(strings.Join(irc.Arguments, " "))))
	if err != nil {
		c.Log.Println(err)
		return
	}
	var output string
	for _, ev := range resp.Events {
		if ev.Message != "" {
			output += "  " + ev.Message
		}
	}
	if resp.Errors != "" {
		if strings.Contains(resp.Errors, "/") {
			i := strings.Index(resp.Errors, "/")
			resp.Errors = resp.Errors[i:]
			i = strings.Index(resp.Errors, "/")
			resp.Errors = resp.Errors[i:]
			output += "  " + resp.Errors
		}

	}
	irc.Reply(c, output)
}

type Request struct {
	Code string `json:'code'`
}
type Event struct {
	Message string
}

type Response struct {
	Output string `json:'output'`
	Error  string `json:'compiler_errors'`
	Errors string
	Events []Event
}

func sendToCompiler(client *http.Client, gocode []byte) (Response, error) {
	playground := "http://golang.org/compile"
	if gocode == nil || len(gocode) < len(`println("æ")`) {
		return Response{}, fmt.Errorf("not enough bytes")
	}
	//gocode = []byte(strings.Replace(string(gocode), "\"", "\\\"", -1))
	gocode = []byte("version=2&body=" + url.QueryEscape(string(gocode)))
	req, err := http.NewRequest(http.MethodPost, playground, bytes.NewReader(gocode))
	if err != nil {
		return Response{}, fmt.Errorf("%v", err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("%v", err)
	}

	var response Response
	err = json.Unmarshal(j, &response)
	if err != nil {
		return Response{}, fmt.Errorf("%v", err)
	}
	return response, nil
}

func lineToMainFunc(code []byte) []byte {

	if !strings.Contains(string(code), "func main()") {
		tmpl := `
																																																func main(){
																																																		%s      
																																																	}
																																																	`
		code = []byte(fmt.Sprintf(tmpl, string(code)))
	}

	if !strings.Contains(string(code), "package main") {
		code = []byte("package main\n" + string(code))
	}
	return code
}