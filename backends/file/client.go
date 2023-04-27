package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
	"github.com/haad/confd/log"
	util "github.com/haad/confd/util"
	"github.com/nqd/flat"
)

// Client provides a shell for the yaml client
type Client struct {
	filepath []string
	filter   string
}

type ResultError struct {
	response uint64
	err      error
}

func NewFileClient(filepath []string, filter string) (*Client, error) {
	return &Client{filepath: filepath, filter: filter}, nil
}

func readFile(path string, vars map[string]string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	jsonDoc, err := yaml.YAMLToJSON(data)
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
		return err
	}
	var yamlMap map[string]interface{}
	err = json.Unmarshal(jsonDoc, &yamlMap)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err.Error())
		return err
	}
	out, err := flat.Flatten(yamlMap, &flat.Options{Delimiter: "/"})
	if err != nil {
		fmt.Printf("Error flatten yaml: %s\n", err.Error())
		return err
	}
	for key, v := range out {
		fixed_key := "/" + key
		switch v := v.(type) {
		case string:
			vars[fixed_key] = v
		case int:
			vars[fixed_key] = strconv.Itoa(v)
		case bool:
			vars[fixed_key] = strconv.FormatBool(v)
		case float64:
			vars[fixed_key] = strconv.FormatFloat(v, 'f', -1, 64)
		}
	}
	return nil
}

func (c *Client) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	var filePaths []string
	for _, path := range c.filepath {
		p, err := util.RecursiveFilesLookup(path, c.filter)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, p...)
	}

	for _, path := range filePaths {
		err := readFile(path, vars)
		if err != nil {
			return nil, err
		}
	}
VarsLoop:
	for k := range vars {
		for _, key := range keys {
			if strings.HasPrefix(k, key) {
				continue VarsLoop
			}
		}
		delete(vars, k)
	}
	log.Debug(fmt.Sprintf("Key Map: %#v", vars))
	return vars, nil
}

func (c *Client) watchChanges(watcher *fsnotify.Watcher, stopChan chan bool) ResultError {
	outputChannel := make(chan ResultError)
	go func() error {
		defer close(outputChannel)
		for {
			select {
			case event := <-watcher.Events:
				log.Debug(fmt.Sprintf("Event: %s", event))
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Create == fsnotify.Create {
					outputChannel <- ResultError{response: 1, err: nil}
				}
			case err := <-watcher.Errors:
				outputChannel <- ResultError{response: 0, err: err}
			case <-stopChan:
				outputChannel <- ResultError{response: 1, err: nil}
			}
		}
	}()
	return <-outputChannel
}

func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	if waitIndex == 0 {
		return 1, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return 0, err
	}
	defer watcher.Close()
	for _, path := range c.filepath {
		isDir, err := util.IsDirectory(path)
		if err != nil {
			return 0, err
		}
		if isDir {
			dirs, err := util.RecursiveDirsLookup(path, "*")
			if err != nil {
				return 0, err
			}
			for _, dir := range dirs {
				err = watcher.Add(dir)
				if err != nil {
					return 0, err
				}
			}
		} else {
			err = watcher.Add(path)
			if err != nil {
				return 0, err
			}
		}
	}
	output := c.watchChanges(watcher, stopChan)
	if output.response != 2 {
		return output.response, output.err
	}
	return waitIndex, nil
}
