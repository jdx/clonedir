package clonedir

import (
	"encoding/json"
	"os"
	"path"
	"testing"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func TestClonedir(t *testing.T) {
	Clone(path.Join("fixtures", "1-foo"), "tmp")
	f, err := os.Open(path.Join("fixtures/1-foo/node_modules/edon-test-c/package.json"))
	must(err)
	var pjson map[string]interface{}
	must(json.NewDecoder(f).Decode(&pjson))
	if pjson["name"] != "edon-test-c" {
		t.Fail()
	}
}
