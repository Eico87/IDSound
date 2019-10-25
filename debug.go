package main


func testingAttackValues(a attack) {
	p("-------------------------------------")
	p("testing values in ", a.name)
	p("-------------------------------------")
	p("detected", a.check)
	p("recursive", a.recursive)
	p("evidence", a.evidence)
}

func testingMonitoredFilesValue(f monitoredFile) {
	p("-------------------------------------")
	p("testing values in ", f.name)
	p("-------------------------------------")
	p("hasBeenModified ", f.hasBeenModified)
	p("last log", f.lastLog)
	p("last mod", f.lastMod)
}

func test() {
	p("-------------------------------------")
	p("RUNNING TEST")
	p("-------------------------------------")
	for i := 0; i < len(allAttacks); i++ {
		testingAttackValues(allAttacks[i])
	}
	for j := 0; j < len(allMonitoredFiles); j++ {
		testingMonitoredFilesValue(allMonitoredFiles[j])
	}
}
