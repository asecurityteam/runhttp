package runhttp

// Reporter manages a background task which reports metrics outside of the main process
type Reporter interface {
	Report()
	Close()
}

// MultiReporter manages multiple reporters
type MultiReporter []Reporter

// Report executes the Report function for all Reporters managed by the MultiReporter
func (mr MultiReporter) Report() {
	for _, r := range mr {
		go r.Report()
	}
}

// Close executes the Close function for all Reporters managed by the MultiReporter
func (mr MultiReporter) Close() {
	for _, r := range mr {
		r.Close()
	}
}
