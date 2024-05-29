package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func ReadJson(path string, val any) (err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, val)
	return
}

func WriteJson(dir_path string, name string, val any) (err error) {
	_text, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		err = fmt.Errorf("json.MarshalIndent %s", err.Error())
		return
	}

	err = os.MkdirAll(dir_path, fs.FileMode(S_IRWXU|S_IRWXO))
	if err != nil && !errors.Is(err, os.ErrExist) {
		err = fmt.Errorf("os.MkdirAll %s", err.Error())
		return
	}
	err = WriteFile(dir_path+"\\"+name+".json", _text)
	if err != nil {
		err = fmt.Errorf("base.WriteFile %s", err.Error())
		return
	}

	return
}
