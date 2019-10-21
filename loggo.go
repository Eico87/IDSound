package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const authLog = "/root/custom/loggo/auth.log.test" //logfile auth.log
//const apacheLog = "var/log/apache2/access.log"
//const apacheLog = "var/log/syslog"

//check
const checkLogin = "FAILED SU"
const checkLoginSSH = "Failed password " //for root from [IP]

//const checkNmap = "Nmap Scripting Engine" //apache
//const checkSegFault = "segfault at" //in syslog
//const checkDictionary =
//const checkCPU =
//const checkServer =
//const checktraffic =

//alerts
const failedLogin = "Attempted login failed on the machine"
const sshFailedLogin = "Attempted login failed on SSH"
const machineBruteforce = "Brute force attack on the machine"
const sshBruteforce = "Brute force attack on ssh"

//const sshConnection = "SSH connection enstablished"
//const nmapscan = "NMAP scan detected"
//const traffic = "Unusual packets traffic detected"
//const dictionary = "Dictionary attack detected"
//const cpuActivity = "unusual CPU activity detected"
//const serverStarted = "Server started"

func main() {
	loopcounter := 0
	machineBruteforceCounter := 0
	sshBruteforceCounter := 0

	for {
		time.Sleep(1 * time.Second) //wait 1 sec

		file, err := os.Open("/var/log/auth.log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			authLogFile := scanner.Text()

			//fmt.Printf("Log: \n %s", authLogFile) //DEBUG

			//select only last string
			//something equivalent to tail -n 1 -f auth.log

			//perform tests

			if strings.Contains(authLogFile, checkLogin) {
				machineBruteforceCounter++
				if machineBruteforceCounter > 6 {
					playAlert(machineBruteforce)
					machineBruteforceCounter = 0
				} else {
					playAlert(failedLogin)
				}
			}

			if strings.Contains(authLogFile, checkLoginSSH) {
				sshBruteforceCounter++
				if sshBruteforceCounter > 6 {
					playAlert(sshBruteforce)
					sshBruteforceCounter = 0
				} else {
					playAlert(sshFailedLogin)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		loopcounter++
		if loopcounter > 60 {
			loopcounter = 0
		}
	}
}

func playAlert(x string) {
	cmd := exec.Command("espeak", "-p", "90", "-g", "3")
	alertAr := []string{"ATTENTION", x, "detected"}
	alert := strings.Join(alertAr, ", ")
	cmd.Stdin = strings.NewReader(alert)
	//cmd.Stdin = strings.NewReader(nmapscan)
	//cmd.Stdin = strings.NewReader(dictionary)
	//cmd.Stdin = strings.NewReader(sshConnection)
	//cmd.Stdin = strings.NewReader(cpuActivity)
	//cmd.Stdin = strings.NewReader(serverStarted)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(4 * time.Second)
}

//TODO
// check if the file has been modified and then do the rest
// read just the last line
// 6 failed login in one minute will trigger bruteforce alert.
// save constant in a file

//if the file auth.log has been modified
//lastModTime := file.Atime_ns
//fmt.Printf("New LOG!") //DEBUG

//Investigate failed login attempts
//Investigate brute-force attacks and other vulnerabilities related to user authorization mechanism.

//syslog
//This log file contains generic system activity logs.
