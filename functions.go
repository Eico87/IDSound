package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"time"
	//"github.com/hpcloud/tail"
	//"log"
	//"github.com/fsnotify/fsnotify"
)

func checkErrors(e error) {
	if e != nil {
		panic(e)
	}
}

//---------------------------------------------------------Replace?
func checkLastModTime(f *monitoredFile) time.Time {
	var t time.Time
	fi, err2 := os.Stat(f.path)
	checkErrors(err2)
	t = fi.ModTime()
	return t
}

//check if the actual value of modification time is changed

//values gets modified inside the function but not outside
func watchLog(f *monitoredFile) bool {
	t := checkLastModTime(f)

	if f.lastMod != t {
		f.sHasBeenModified(true)
		f.sLastMod(t)
	} else {
		f.sHasBeenModified(false)
	}
	return f.hasBeenModified
}

func tailLog(f *monitoredFile) {

	file, err := os.Open(f.path)
	checkErrors(err)
	defer file.Close()

	buf := make([]byte, 100)
	stat, err := os.Stat(f.path) //checking file dimension
	start := stat.Size() - 100   //subtracting 80 bits to get the last line
	_, err = file.ReadAt(buf, start)
	if err == nil {
		f.sLastLog(string(buf))
	}
}

func detectAttack(a *attack, f *monitoredFile) {

	p("detecting attack")        //DEBUG
	p("searching : ", a.control) //DEBUG
	p("in", f.lastLog)           //DEBUG

	//since f.lastLog is void I can't detect attacks
	if strings.Contains(f.lastLog, a.control) {
		a.sEvidence(f.lastLog)
		a.sCheck(true)
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

//---------------------------------------------------------NEW WATCHER BELOW

/*
func watch(f monitoredFile) bool {
	watcher, err := fsnotify.NewWatcher()
	checkErrors(err)
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(f.path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
*/
//---------------------------------------------------------------------

/*
func tailLog(f monitoredFile) {
	var newLog string
	t, err := tail.TailFile(f.path, tail.Config{Follow: false, MaxLineSize: 150})
	checkErrors(err)

	// line.Text it's a channel...whatever that means
	// but I need just the last lines
	// so f.lastLog remains blank

	for line := range t.Lines {
		newLog = line.Text
		p(newLog)                           //DEBUG
	}
	pSeparator()                                //DEBUG
	p("LAST LOG SET FOR", f.name, ": ", newLog) //DEBUG
	pSeparator()                                //DEBUG
	f.lastLog = newLog
}
*/
