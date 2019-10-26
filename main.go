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

		loopcounter++
		
		if loopcounter > 60 {
			loopcounter = 0
		}

		p("=========== Loop: ", loopcounter, "===========") //DEBUG

		for jj := 0; jj < len(allMonitoredFiles); jj++ {
			allMonitoredFiles[jj].resetMonitor()
		}

		time.Sleep(50 * time.Millisecond) //wait
		//for each logfile watched
		for i := 0; i < len(allMonitoredFiles); i++ {
			//check if it has been modified

			if watchLog(allMonitoredFiles[i]) { //REAL LINE
				//if true { //DEBUG

				tailLog(allMonitoredFiles[i]) //read the last line

				//perform all the tests
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
	}
}
