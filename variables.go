package main

import (
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

//SET
func (f *monitoredFile) sHasBeenModified(b bool) { f.hasBeenModified = b }
func (f *monitoredFile) sLastLog(s string)       { f.lastLog = s }
func (f *monitoredFile) sLastMod(t time.Time)    { f.lastMod = t }

//RESET
func (f *monitoredFile) resetMonitor() {
	f.hasBeenModified = false
	f.lastLog = ""
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

//SET
func (a *attack) sCheck(b bool)      { a.check = b }
func (a *attack) sEvidence(s string) { a.evidence = s }
func (a *attack) sRecursive(n int)   { a.recursive = n }

//RESET
func (a *attack) resetAttack() {
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
		refernceLog: "auth.log", //can I refer to a struct?? if I do how do I change variables?
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
