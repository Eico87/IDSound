package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	//"log"
	//"github.com/fsnotify/fsnotify"
)

func checkErrors(e error) {
	if e != nil {
		panic(e)
	}
}

func checkLastModTime(f monitoredFile) time.Time {
	var t time.Time
	fi, err2 := os.Stat(f.path)
	checkErrors(err2)
	t = fi.ModTime()
	return t
}

//check if the actual value of modification time is changed
func watchLog(f monitoredFile) bool {
	t := checkLastModTime(f)

	if f.lastMod != t {
		f.hasBeenModified = true
		f.lastMod = t
		p(string(f.name)+" Has been modified!", f.lastLog) //DEBUG
	} else {
		f.hasBeenModified = false
	}
	return f.hasBeenModified
}


func tailLog(f monitoredFile) {
	var newLog string
	t, err := tail.TailFile(f.path, tail.Config{Follow: false})
	checkErrors(err)
	for line := range t.Lines { // THE PROBLEM IS HERE
		newLog = line.Text
	}
	f.lastLog = newLog
}

func detectAttack(a attack, f monitoredFile) {
	p("detecting attack")
	p("searching : ", a.control)
	p("in", f.lastLog)

	if strings.Contains(f.lastLog, a.control) {
		a.check = true
		p(a.name, " Detected!") //DEBUG
	}
}

func pSeparator() {
	p("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-") //how do I multiply a string?
}
func printEvidence(a attack) {
	pSeparator()
	p(a.name, " detected!")
	pSeparator()
	p("Evidence: ", a.evidence)
	p("in ", a.refernceLog)
	pSeparator()
}

func playAlert(a attack) {
	cmd := exec.Command("espeak", "-p", "90", "-g", "3")
	alertAr := []string{"ATTENTION", a.message, "detected"}
	alert := strings.Join(alertAr, ", ")
	cmd.Stdin = strings.NewReader(alert)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	checkErrors(err)
	time.Sleep(4 * time.Second)
}
