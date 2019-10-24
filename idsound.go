package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type monitoredFile struct {
	name string
	path string
	//default
	hasBeenModified bool
	lastLog         string
}

type attack struct {
	name string
	//refernceLog string
	refernceLog monitoredFile
	control     string
	message     string
	//default
	check     bool
	evidence  string
	recursive int
}

//GET monitoredFile
func (f monitoredFile) gName() string          { return f.name }
func (f monitoredFile) gPath() string          { return f.path }
func (f monitoredFile) gHasBeenModified() bool { return f.hasBeenModified }
func (f monitoredFile) gLastLog() string       { return f.lastLog }

//SET monitoredFile
func (f monitoredFile) sHasBeenModified(b bool) { f.hasBeenModified = b }
func (f monitoredFile) sLastLog(s string)       { f.lastLog = s }

//RESET
func (f monitoredFile) resetMonitor() {
	f.sHasBeenModified(false)
	f.sLastLog("")
}

//GET attack
func (a attack) gName() string { return a.name }

//func (a attack) gRefernceLog() string { return a.refernceLog }
func (a attack) gControl() string  { return a.control }
func (a attack) gMessage() string  { return a.message }
func (a attack) gCheck() bool      { return a.check }
func (a attack) gEvidence() string { return a.evidence }

//SET attack
func (a attack) sCheck(b bool)      { a.check = b }
func (a attack) sEvidence(s string) { a.evidence = s }
func (a attack) sRecursive(n int)   { a.recursive = n }

//RESET
func (a attack) resetAttack() {
	a.sCheck(false)
	a.sEvidence("")
	a.sRecursive(0)
}

var allMonitoredFiles = []monitoredFile{
	monitoredFile{
		name:            "auth.log",
		path:            "/var/log/auth.log",
		hasBeenModified: false,
		lastLog:         "",
	},
	monitoredFile{
		name:            "apache2 access.log",
		path:            "/var/log/apache2/access.log",
		hasBeenModified: false,
		lastLog:         "",
	},
	monitoredFile{
		name:            "apache2 error.log",
		path:            "/var/log/apache2/error.log",
		hasBeenModified: false,
		lastLog:         "",
	},
	monitoredFile{
		name:            "apache2 xplico_access.log",
		path:            "/var/log/apache2/xplico_access.log",
		hasBeenModified: false,
		lastLog:         "",
	},
	monitoredFile{
		name:            "syslog",
		path:            "/var/log/syslog",
		hasBeenModified: false,
		lastLog:         "",
	},
}

var allAttacks = []attack{
	attack{
		name:        "Failed Login",
		refernceLog: "auth.log", //can I refer to a struct??
		control:     "FAILED SU",
		message:     "Attempted login failed on the machine",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "Failed Login on SSH",
		refernceLog: "auth.log",
		control:     "Failed password",
		message:     "Attempted login failed on SSH",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "Subdomain Bruteforce",
		refernceLog: "apache2 access.log",
		control:     "gobuster",
		message:     "subdomain bruteforcing",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "NMAP Scan",
		refernceLog: "apache2 access.log",
		control:     "Nmap Scripting Engine",
		message:     "port scan",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "invalid HTTP Method",
		refernceLog: "apache2 error.log",
		control:     "Invalid method",
		message:     "invalid method request recived",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "NMAP Scan",
		refernceLog: "apache2 xplico_access.log",
		control:     "Nmap Scripting Engine",
		message:     "port scan",
		check:       false,
		evidence:    "",
	},
	attack{
		name:        "Segmentation Fault",
		refernceLog: "syslog",
		control:     "segfault at",
		message:     "Segmentation Fault",
		check:       false,
		evidence:    "",
	},
}

var p = fmt.Println

func main() {
	loopcounter := 0
	//for each second in infinite loop
	for {
		time.Sleep(700 * time.Millisecond) //wait 1 sec
		//for each logfile watched
		for i := 0; i < len(allMonitoredFiles); i++ {

			//check if it has been modified
			if watchLog(allMonitoredFiles[i]) { //NOT WORKING: ALWAYS RETURN TRUE

				//read the last line
				tailLog(allMonitoredFiles[i]) //NOT WORKING: IT JUST READ THE WHOLE FILE

				//perform all the tests relative to the log in question

				//this could be better if I could pass the file path inside the attack struct

				switch allMonitoredFiles[i].gName {

				case "auth.log":
					p("searching attack evidence in  auth.log")
					detectAttack(allAttacks[0], allMonitoredFiles[0])
					detectAttack(allAttacks[1], allMonitoredFiles[0])

				case "apache2 access.log":
					p("searching attack evidence in apache2 access.log")
					detectAttack(allAttacks[2], allMonitoredFiles[1])
					detectAttack(allAttacks[3], allMonitoredFiles[1])

				case "apache2 error.log":
					p("searching attack evidence in  apache2 error.log")
					detectAttack(allAttacks[4], allMonitoredFiles[2])

				case "xplico_access.log":
					p("searching attack evidence in  apache2 xplico_access.log")
					detectAttack(allAttacks[5], allMonitoredFiles[3])

				case "syslog":
					p("searching attack evidence in  syslog")
					detectAttack(allAttacks[6], allMonitoredFiles[4])

					//if an attack is detected
					if allAttacks[i].gCheck() { //NOW WORKING
						printEvidence(allAttacks[i]) //print evidence on log file and NOW WORKING it prints on the terminal
						playAlert(allAttacks[i])     //play audio alert
						//add bruteforce allarm here

						//reset attack variables
						allAttacks[i].resetAttack()
					}
				}
			}
		}
		loopcounter++
		if loopcounter > 60 {
			loopcounter = 0
		}
		p("Loop: ", loopcounter)
	}
}

func watchLog(f monitoredFile) bool {

	p("Watching", f.gName()) //DEBUG
	/*
		var t time.Time
		if fi, err2 := os.Stat(f.gPath()); err2 == nil {
			t := fi.ModTime()
		}

		//----pass time
		//add condition here
	*/
	if 1 > 0 {
		f.sHasBeenModified(true)
		p(string(f.gName()) + " Has been modified.") //DEBUG
	}
	return f.hasBeenModified
}

func checErrors(e error) {
	if e != nil {
		panic(e)
	}
}
func tailLog(f monitoredFile) string {

	p("getting the last line from", f.gName()) //DEBUG

	logFile, err := ioutil.ReadFile(f.gPath())
	checErrors(err)

	f.sLastLog(string(logFile))

	return string(logFile)
}

func printEvidence(a attack) {
	p("------------------------------------")
	p(a.gName())
	p("------------------------------------")
	p("Evidence:") //DEBUG
	p(a.gEvidence())
	p("in")
	p(a.gRefernceLog())
	p("------------------------------------")
}
func detectAttack(a attack, f monitoredFile) bool {
	vulnerable := false
	/*
		if strings.Contains(string(f.gLastLog), string(a.gControl)) {
			vulnerable = true
			a.sCheck(true)
			p(string(a.gName()) + " detected!") //DEBUG
		}
	*/
	return vulnerable
}
func playAlert(a attack) {
	cmd := exec.Command("espeak", "-p", "90", "-g", "3")
	alertAr := []string{"ATTENTION", a.message, "detected"}
	alert := strings.Join(alertAr, ", ")
	cmd.Stdin = strings.NewReader(alert)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(4 * time.Second)
}

//TODO

// FUNCTIONS
// read just the last line
// print evidence on terminal
// save evidence in /var/logs/idsound.log
// check if the file has been modified and then do the rest

//STRUCTURE
// import attacks parameters from json
// default values and constuctors
// subdivide in files

//UPGRADES
// find a better control fon nmap
// set clock:how many logs are there in a second? should I use millliseconds?

//NEW FUNCTIONS
// wordlist spotter
// cpu monitor
// network traffic monitor
// cool interface

