// SwartzNotes standalone package for clients
// ------------------------------------------
// Licensed under an Apache license

// The package provides minimal support for the browser extension to work
// Also serves as a bridge for the fact extractor java package to work


Architecture
---

This program listens on port 3333 on the local machine to intercept
PDF from the extension. Once it sees the PDF blob, it saves to a file
and then calls factExtractor.jar to run on it to produce a fact file.

Compiling
---

This program is written in Go so it can compile to a number of platforms
and architectures.

Download the Go package from golang.org, then to compile an executable 
native to your OS, just do
	go build

Cross-compiling: Follow the cross-compile 
instructions at `https://github.com/davecheney/golang-crosscompile`.