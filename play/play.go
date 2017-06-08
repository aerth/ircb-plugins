package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/aerth/ircb"
)

// empty main
func main() {}

// Init must be named Init taking an ircb.Connection and returning an error
func Init(c *ircb.Connection) error {
	name := "go play" // no global variables in plugin
	c.Log.Println("plugin loading:", name)
	c.AddCommand("go", commandGo)
	c.Log.Println("plugin loaded:", name)
	return nil
}

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
			resp.Errors = resp.Errors[i+1:]
			i = strings.Index(resp.Errors, "/")
			resp.Errors = resp.Errors[i+1:]
			output += "  " + resp.Errors
		}

	}
	if resp.Error != "" {
		output += "  " + resp.Error

	}
	output = strings.TrimSuffix(strings.TrimSpace(output), "\n")
	output = strings.Replace(output, "\n", " -- ", -1)
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

func (resp Response) CombinedOutput() string {

	var output string
	for _, ev := range resp.Events {
		if ev.Message != "" {
			output += "  " + ev.Message
		}
	}
	if resp.Errors != "" {
		if strings.Contains(resp.Errors, "/") {
			i := strings.Index(resp.Errors, "/")
			resp.Errors = resp.Errors[i+1:]
			i = strings.Index(resp.Errors, "/")
			resp.Errors = resp.Errors[i+1:]
			output += "  " + resp.Errors
		}

	}
	if resp.Error != "" {
		output += "  " + resp.Error

	}
	output = strings.TrimSuffix(strings.TrimSpace(output), "\n")
	output = strings.Replace(output, "\n", " -- ", -1)
	return output
}

func sendToCompiler(client *http.Client, code []byte) (Response, error) {
	compiler := "http://golang.org/compile"
	if code == nil || len(code) < len(`println("æ")`) {
		return Response{}, fmt.Errorf("not enough bytes")
	}

	// run goimports on it
	code, err := imports.Process("", code, nil)
	if err != nil {
		return Response{}, fmt.Errorf("imports: %v", err)

	}
	gocode := "version=2&body=" + url.QueryEscape(string(code))
	req, err := http.NewRequest(http.MethodPost,
		compiler,
		bytes.NewReader([]byte(gocode)))
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
