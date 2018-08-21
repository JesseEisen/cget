// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"cget/copy"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type LocalInfo struct {
	slibname string
	ilibname string
	fuzzy    bool
	sresult  []string
}

var Li *LocalInfo

var (
	wg    sync.WaitGroup
	Mutex sync.Mutex
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "search and install from local",
	Long:  `use this command to install local library`,
	Run: func(cmd *cobra.Command, args []string) {
		if Li.slibname != "" {
			Li.SearchLocal()
		}
		if Li.ilibname != "" {
			Li.InstallLocal()
		}
	},
}

func (li *LocalInfo) SearchLocal() {
	path := os.Getenv("CGET_PATH")
	if path == "" {
		path = "~/.cget/"
	}

	fmt.Printf("\nSearching...\n\n")
	paths := strings.Split(path, ":")
	wg.Add(len(paths))
	for _, target := range paths {
		go li.searchPath(target)
	}
	wg.Wait()
}

func (li *LocalInfo) searchPath(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		ddir := filepath.Join(dir, file.Name())
		if file.Name() == li.slibname {
			Mutex.Lock()
			li.sresult = append(li.sresult, ddir)
			Mutex.Unlock()
		} else if li.fuzzy == true {
			if strings.Contains(file.Name(), li.slibname) {
				Mutex.Lock()
				li.sresult = append(li.sresult, ddir)
				Mutex.Unlock()
			}
		}
	}

	defer wg.Done()
}

func (li *LocalInfo) InstallLocal() {
	// for install, we cannot use fuzzy, so if user define the -f
	// we need to toggle it to false
	li.fuzzy = false

	li.SearchLocal()
	if len(li.sresult) == 0 {
		fmt.Println("Can not install this lib, please check the libname")
		return
	}

	fmt.Println("Install...")
	for _, source := range li.sresult {
		copyToCurrent(source)
	}
}

func copyToCurrent(source string) {
	fmt.Printf("dir is %s\n", source)
	//current, err := filepath.Abs(filepath.Dir(os.Args[0]))
	current, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	copy.Copy(source, current)
	fmt.Printf("current dir is: %s\n", current)
}

func init() {
	rootCmd.AddCommand(localCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//localCmd.Flags().BoolVarP(&tg, "toggle", "t", false, "Help message for toggle")
	localCmd.Flags().StringVarP(&Li.slibname, "search", "s", "", "search library")
	localCmd.Flags().StringVarP(&Li.ilibname, "install", "i", "", "install library")
	localCmd.Flags().BoolVarP(&Li.fuzzy, "fuzzy", "f", false, "fuzzy search and install")
}
