#PCS

## How it words
the ransomware will automatically check if pri key existed in the current folder. If not, ransomware will encrypt all the files in the root dir with extension txt, doc, pdf, csv, exe. 

If file "keys" founded in the current folder, the ransomware will automatically decrypt all the files.

## How to run:
	```
	$ go run ransomware.go
	```
### notes:
	I recommend to build the .go file first by using 
	```
	$ go build ransomware.go
	$ ./ransomware
	```