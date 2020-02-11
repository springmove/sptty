package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestI18N(t *testing.T) {
	dir, _ := os.Getwd()
	i18n := path.Join(dir, "i18n.json")

	f, _ := os.Open(i18n)
	defer f.Close()

	content, _ := ioutil.ReadAll(f)

	lang := map[string]map[string]string{}
	json.Unmarshal(content, &lang)

	fmt.Println(lang)
}
