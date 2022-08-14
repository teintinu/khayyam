package internal

import (
	"errors"
	"fmt"
	"os"
)

type InitOptions struct {
	CleanArchitecture bool
}

func Init(opts InitOptions) error {
	if _, err := os.Stat("khayyam.yml"); !errors.Is(err, os.ErrNotExist) {
		return errors.New("khayyam.yml already exists")
	}
	var content string
	if opts.CleanArchitecture {
		content = "a"
	} else {
		content = `workspace:
  name: "example"
  version: "1.0.0"

packages:
  "@example/a":
    folder: "a"
    dependencies:
      "@example/b": "*"

  "@example/b":
    folder: "b"`
	}
	f, err := os.Create("khayyam.yml")

	if err != nil {
		return err
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		return err
	}

	fmt.Println("khayyam.yml created")
	return nil
}
