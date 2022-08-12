package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "read files error")
	}

	for _, f := range files {
		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			return nil, errors.Wrap(err, "open files error")
		}

		scanner := bufio.NewScanner(file)
		ok := scanner.Scan()
		if ok {
			// терминальные нули (0x00) заменяются на перевод строки (\n);
			s := strings.Join(strings.Split(scanner.Text(), "\x00"), "\n")
			// пробелы и табуляция в конце T удаляются
			env[f.Name()] = strings.TrimRight(s, " \t")
		} else {
			env[f.Name()] = ""
		}

		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "scanning error")
		}

		file.Close()
	}

	return env, nil
}
