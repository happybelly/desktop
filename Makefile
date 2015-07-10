all:	swartznotes-desktop
swartznotes-desktop:	main.go
	go build
	cp ../annotator/pkg/annotator.js static/viewer/web/crowd-annotator.js
	cp ../crowd-annotator/crowd-loader.js static/viewer/web/
clean:
	rm swartznotes-desktop
