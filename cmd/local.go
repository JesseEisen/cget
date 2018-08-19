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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	slibname string
	ilibname string
	fuzzy    bool
	sresult  []string
)

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
		if slibname != "" {
			searchLocal(slibname)
		}
		if ilibname != "" {
			installLocal(ilibname)
		}
	},
}

/*
 * search in CGET_PATH, the PATH is like the PATH
 */
func searchLocal(slibname string) {
	path := os.Getenv("CGET_PATH")
	if path == "" {
		path = "~/.cget/"
	}

	fmt.Printf("\nSearching...\n\n")
	paths := strings.Split(path, ":")
	wg.Add(len(paths))
	for _, target := range paths {
		go searchPath(target, slibname)
	}

	wg.Wait()
}

func searchPath(dir, libname string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == libname {
			Mutex.Lock()
			sresult = append(sresult, dir+file.Name())
			Mutex.Unlock()
			fmt.Printf("\t Find %s: %s\n", dir, file.Name())
		} else if fuzzy == true {
			if strings.Contains(file.Name(), libname) {
				Mutex.Lock()
				sresult = append(sresult, dir+file.Name())
				Mutex.Unlock()
				fmt.Printf("\t Find %s: %s\n", dir, file.Name())
			}
		}
	}

	defer wg.Done()
}

func copyToCurrent(source string) {
	fmt.Printf("dir is %s\n", source)
	current, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("current dir is: %s\n", current)
}

func installLocal(ilibname string) {
	// for install, we cannot use fuzzy, so if user define the -f
	// we need to toggle it to false
	fuzzy = false

	searchLocal(ilibname)
	if len(sresult) == 0 {
		fmt.Println("Can not install this lib, please check the libname")
		return
	}

	fmt.Println("Install...")
	for _, source := range sresult {
		copyToCurrent(source)
	}

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
	localCmd.Flags().StringVarP(&slibname, "search", "s", "", "search library")
	localCmd.Flags().StringVarP(&ilibname, "install", "i", "", "install library")
	localCmd.Flags().BoolVarP(&fuzzy, "fuzzy", "f", false, "fuzzy search and install")
}
