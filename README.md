# PCS

## How it works
the ransomware will automatically check if pri key existed in the current folder. If not, ransomware will encrypt all the files in the root dir with extension txt, doc, pdf, csv, exe. 

If file "keys" founded in the current folder, the ransomware will automatically decrypt all the files.

The retrieve.go will use the sha256 checksum, which is indicated in the Readme, to get the corresponding encrypted key and then decrypt the files and output to "keys" file.
## How to run:
	$ go run TestEnv/ransomware.go
    $ go run Server/retrieve.go <sha256 checksum>
	
### notes:
I recommend to build the .go file first by using 
	
	$ go build ransomware.go
	$ ./ransomware
	