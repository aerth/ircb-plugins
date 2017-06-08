package main

import (
	"log"
	"net/http"
	"testing"
)

func TestFormatCompile(t *testing.T) {
	tc := []struct {
		input  string
		output string
	}{
		{`println("hi")`, "hi"},
		{`fmt.Println("hi")`, "hi"},
	}

	for _, test := range tc {
		result, err := sendToCompiler(http.DefaultClient, lineToMainFunc([]byte(test.input)))
		log.Println(result)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		if result.CombinedOutput() != test.output {

			t.Logf("wanted %q, got %q", test.output, result.CombinedOutput())
			t.FailNow()
		}
	}
}
