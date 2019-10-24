package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	lastMod         time.Time
}

type attack struct {
	name string
	//refernceLog string
	refernceLog string
	control     string
	message     string
	//default
	check     bool
	evidence  string
	recursive int
}

//RESET
func (f monitoredFile) resetMonitor() {
	f.hasBeenModified = false
	f.lastLog = ""
}

//RESET
func (a attack) resetAttack() {
	a.check = false
	a.evidence = ""
	a.recursive = 0
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

	//save when was the last time that the file has been modified
	for ii := 0; ii < len(allMonitoredFiles); ii++ {
		t := checkLastModTime(allMonitoredFiles[ii])
		allMonitoredFiles[ii].lastMod = t
	}

	//in infinite loop
	for {

		time.Sleep(700 * time.Millisecond) //wait
		//for each logfile watched
		for i := 0; i < len(allMonitoredFiles); i++ {
			//check if it has been modified

			//if watchLog(allMonitoredFiles[i]) { //REAL LINE
			if true {
				//read the last line
				tailLog(allMonitoredFiles[i]) //NOT WORKING: IT JUST READ THE WHOLE FILE

				//perform all the tests relative to the log in question
				//this could be better if I could pass the file path inside the attack struct
				x := allMonitoredFiles[i].name
				switch x {

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
				}
			}
		}
		//for any attack
		for j := 0; j < len(allAttacks); j++ {

			//if it has been proved
			if allAttacks[j].check { //NOW WORKING

				printEvidence(allAttacks[j]) //print evidence on log file and NOW WORKING it prints on the terminal
				playAlert(allAttacks[j])     //play audio alert

				//add bruteforce alarm here

				//reset attack variables
				allAttacks[j].resetAttack()
			}
		}

		loopcounter++
		if loopcounter > 60 {
			loopcounter = 0
		}

		p("Loop: ", loopcounter)
	}
}

func checkLastModTime(f monitoredFile) time.Time {
	var t time.Time

	if fi, err2 := os.Stat(f.path); err2 == nil {
		t = fi.ModTime()
	}
	return t
}

//check if the actual value of modification time is changed
func watchLog(f monitoredFile) bool {
	p("Watching", f.name) //DEBUG
	t := checkLastModTime(f)

	if f.lastMod != t {
		f.hasBeenModified = true
		f.lastMod = t
		p(string(f.name) + " Has been modified") //DEBUG
	}

	return f.hasBeenModified
}

func checErrors(e error) {
	if e != nil {
		panic(e)
	}
}

func tailLog(f monitoredFile) string {

	p("getting the last line from", f.name) //DEBUG

	logFile, err := ioutil.ReadFile(f.path)
	checErrors(err)
	f.lastLog = (string(logFile))

	return string(logFile)
}

func printEvidence(a attack) {
	p("------------------------------------")
	p(a.name)
	p("------------------------------------")
	p("Evidence:") //DEBUG
	p(a.evidence)
	p("in")
	p(a.refernceLog)
	p("------------------------------------")
}
func detectAttack(a attack, f monitoredFile) bool {
	vulnerable := true
	/*
		if strings.Contains(string(f.LastLog), string(a.Control)) {
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

// read just the last line
// print evidence on terminal and in log files
// save evidence in /var/logs/idsound.log

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
