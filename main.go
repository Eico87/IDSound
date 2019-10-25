package main

import (
	"fmt"
	"time"
)

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
			if true { //DEBUG
				//read the last line
				tailLog(allMonitoredFiles[i]) //NOT WORKING: IT READS "RANDOM" LINE FROM FILE

				//perform all the tests relative to the log in question
				//this could be better if I could pass the file path inside the attack struct
				x := allMonitoredFiles[i].name

				switch x {

				case "auth.log":
					detectAttack(allAttacks[0], allMonitoredFiles[0])
					detectAttack(allAttacks[1], allMonitoredFiles[0])

				case "apache2 access.log":
					detectAttack(allAttacks[2], allMonitoredFiles[1])
					detectAttack(allAttacks[3], allMonitoredFiles[1])

				case "apache2 error.log":
					detectAttack(allAttacks[4], allMonitoredFiles[2])

				case "xplico_access.log":
					detectAttack(allAttacks[5], allMonitoredFiles[3])

				case "syslog":
					detectAttack(allAttacks[6], allMonitoredFiles[4])
				}
			}
		}
		//for any attack
		for j := 0; j < len(allAttacks); j++ {
			//if it has been proved
			if allAttacks[j].check {

				printEvidence(allAttacks[j]) //print evidence on log file (prints on the terminal now)
				playAlert(allAttacks[j])     //play audio alert

				//bruteforce alarm here

				//reset attack variables
				allAttacks[j].resetAttack()
			}
		}
		loopcounter++
		if loopcounter > 60 {
			loopcounter = 0
		}
		p("=========== Loop: ", loopcounter, "===========") //DEBUG
	}
}

//TODO

//https://github.com/hpcloud/tail
// read just the last line

// print evidence on terminal and in log files
// save evidence in /var/logs/idsound.log

//STRUCTURE
// import attacks parameters from json
// default values and constuctors

//UPGRADES
// find a better control fon nmap
// set clock:how many logs are there in a second? should I use millliseconds?
// add new log to monitor with relative attacks

//NEW FUNCTIONS
// wordlist spotter
// cpu monitor
// network traffic monitor
// cool interface
