package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkLastModTime(f monitoredFile) time.Time {
	var t time.Time
	if fi, err2 := os.Stat(f.path); err2 == nil {
		t = fi.ModTime()
	}
	return t
}

//check if the actual value of modification time is changed
func watchLog(f monitoredFile) bool {
	t := checkLastModTime(f)

	if f.lastMod != t {
		f.hasBeenModified = true
		f.lastMod = t
		p(string(f.name) + " Has been modified!") //DEBUG
	}
	return f.hasBeenModified
}

func checkErrors(e error) {
	if e != nil {
		panic(e)
	}
}

func tailLog(f monitoredFile) {

	logFile, err := os.Open(f.path)
	checkErrors(err)
	o2, err := logFile.Seek(600, 0)
	checkErrors(err)
	b2 := make([]byte, 200)
	n2, err := logFile.Read(b2)
	checkErrors(err)
	p(n2, "bytes @", o2)
	p("the last line from ", f.name, " is:") //DEBUG
	p(string(b2[:n2]))
	//f.lastLog = logFile
}

func detectAttack(a attack, f monitoredFile) {
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
