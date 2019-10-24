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
	name        string
	refernceLog string
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
func (a attack) gName() string        { return a.name }
func (a attack) gRefernceLog() string { return a.refernceLog }
func (a attack) gControl() string     { return a.control }
func (a attack) gMessage() string     { return a.message }
func (a attack) gCheck() bool         { return a.check }
func (a attack) gEvidence() string    { return a.evidence }

//SET attack
func (a attack) sCheck(b bool)      { a.check = b }
func (a attack) sEvidence(s string) { a.evidence = s }
func (a attack) sRecursive(n int)   { a.recursive = n }

//RESET
func (a attack) resetAttacks() {
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
		refernceLog: "auth.log",
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

func main() {
	loopcounter := 0
	//for each second in infinite loop
	for {
		time.Sleep(1 * time.Second) //wait 1 sec
		//for each logfile watched
		for i := 0; i < len(allMonitoredFiles); i++ {

			//check if it has been modified
			if watchLog(allMonitoredFiles[i]) { //NOT WORKING: ALWAYS RETURN TRUE
				//read the last line
				tailLog(allMonitoredFiles[i]) //NOT WORKING: IT JUST READ THE WHOLE FILE

				//perform all the tests relative to the log in question
				switch i {
				case 0:
					fmt.Println("auth.log")
				case 1:
					fmt.Println("apache2 access.log")
				case 2:
					fmt.Println("apache2 error.log")
				case 3:
					fmt.Println("apache2 xplico_access.log")
				case 4:
					fmt.Println("syslog")
				}

				//if attack is detected print evidence on log file and play audio alert
				/*
					if detectAttack() { //NOW WORKING
						printEvidence()	//NOW WORKING it prints on the terminal
						playAlert()
						//resetValues()
				*/
			}
		}
	}
	loopcounter++
	if loopcounter > 60 {
		loopcounter = 0
	}
}

func watchLog(f monitoredFile) bool {
	fmt.Println("watching log %s", f.gName()) //DEBUG

	/*
		var t time.Time
		if fi, err2 := os.Stat(f.gPath()); err2 == nil {
			t := fi.ModTime()
		}

		//----pass time
		//add condition here
	*/

	if 1 > 0 {
		f.hasBeenModified = true
		fmt.Println("%s has been modified", f.gName()) //DEBUG
	}
	return f.hasBeenModified
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func tailLog(f monitoredFile) string {

	fmt.Println("tail log") //DEBUG

	logFile, err := ioutil.ReadFile(f.gPath())
	check(err)

	//now it just reads the full file
	//fmt.Printf("Log: \n %s", string(logFile)) //DEBUG

	f.sLastLog(string(logFile))

	return string(logFile)
}

func printEvidence() {
	fmt.Println("print Evidence") //DEBUG
	//fmt.Println(a.logFilePath) //DEBUG
	//fmt.Println(a.message)     //DEBUG
	//fmt.Println(a.evidence)    //DEBUG
	fmt.Println("------------------------------------")
}

func detectAttack(a attack) bool {
	vulnerable := false
	/*
		if strings.Contains( , a.gControl) {
			vulnerable = true
		}*/
	fmt.Println("detecting %s ... %b", a.name, vulnerable) //DEBUG
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

//read just the last line
// default values and constuctors
// import attacks parameters from json
// check if the file has been modified and then do the rest
// print evidence on terminal
// save evidence in /var/logs/idsound.log
// reset values
// find a better control fon nmap
// wordlist spotter
// set clock:how many logs are there in a second? should I use millliseconds?
// cpu monitor
// network traffic monitor
// cool interface
