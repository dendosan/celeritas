package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the app name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}

	log.Println("App name is", appName)

	// git clone the skeleton app
	color.Green("\tCloning repository...")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/dendosan/celeritas-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to go .env file
	color.Yellow("\tCreating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", cel.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a makefile
	var source *os.File
	if runtime.GOOS == "windows" {
		source, err = os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
	} else {
		source, err = os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
	}
	if err != nil {
		exitGracefully(err)
	}
	defer source.Close()

	dest, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
	if err != nil {
		exitGracefully(err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		exitGracefully(err)
	}

	_ = os.Remove("./" + appName + "/Makefile.windows")
	_ = os.Remove("./" + appName + "/Makefile.mac")

	// update the go.mod file
	color.Yellow("\tCreating go.mod file...")
	_ = os.Remove("./" + appName + "/go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), "./" + appName + "/go.mod")
	if err != nil {
		exitGracefully(err)
	}

	// update existing .go files with correct name/imports

	// run go mod tidy in the project directory
}