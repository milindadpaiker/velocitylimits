# Velocitylimits

Velocitylimits is a command-line application in Golang that accepts or declines attempts to load funds into customers' accounts in real-time based on pre-defined rules.

Features:-

* Support for multiple sources and sink
* sqlite mode for resilient processing
* Extensible architecture. New rules can be easily added.
* Support for different currencies is a matter of configuration

Limitations:-
* Not safe for concurrent processing of incoming funds. Need to support "serializable" isolation level on DataAccessLayer(dal) to achieve concurrency.
## Installation

Pre-requisites: Go version 1.13 and above

1. Go to GOPATH/src and clone github repository

```bash
root@4cbdb7382f66:/go/src# git clone https://github.com/milindadpaiker/velocitylimits.git
Cloning into 'velocitylimits'...
remote: Enumerating objects: 160, done.
remote: Counting objects: 100% (160/160), done.
remote: Compressing objects: 100% (91/91), done.
remote: Total 160 (delta 74), reused 139 (delta 53), pack-reused 0
Receiving objects: 100% (160/160), 48.63 KiB | 262.00 KiB/s, done.
Resolving deltas: 100% (74/74), done.
```
2. cd to velocitylimits and run make. 

```bash
root@4cbdb7382f66:/go/src# cd velocitylimits/
root@847da267b049:/go/src/velocitylimits# make
echo "Installing dependencies"
Installing dependencies
go mod vendor
go: downloading github.com/pkg/errors v0.9.1
go: downloading gorm.io/gorm v1.20.2
go: downloading gorm.io/driver/sqlite v1.1.3
go: downloading github.com/jinzhu/inflection v1.0.0
go: downloading github.com/jinzhu/now v1.1.1
go: downloading github.com/mattn/go-sqlite3 v1.14.3
echo "**Building linux binary**"
**Building linux binary**
GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/velocity-limit-app ./cmd/velocitylimit
cp ./cmd/velocitylimit/config.json ./bin/linux_amd64/config.json
:
```
3. If make is successful, application binary and corresponding files (config.json and sample input) will be placed in ./bin/{os}_{arch} folder

```bash
root@847da267b049:/go/src/velocitylimits/bin/linux_amd64# ls
config.json  input.txt  velocity-limit-app
```

## Usage

Without any options, velocity-limit-app expects config.json and input.txt to be present at the some location. With these files in place, just running the application will generate output.txt and log file

```python
root@75d3a8ec19bf:/go/src/velocitylimits/bin/linux_amd64# ./velocity-limit-app
root@75d3a8ec19bf:/go/src/velocitylimits/bin/linux_amd64# ls
config.json  input.txt  output.txt  velocity-limit-app  velocitylimit.log
```

By default velocity-limit-app is an in-memory application. Although highly performant it is not resilient to crashes/interrupts. For reliable processing of incoming funds run the application with backend as sqlite. "append=true" appends output to the file.

```python
root@75d3a8ec19bf:/go/src/velocitylimits/bin/linux_amd64#  ./velocity-limit-app -backend="sqlite" -append=true
root@75d3a8ec19bf:/go/src/velocitylimits/bin/linux_amd64# ls
config.json  input.txt  output.txt  velocity-limit-app  velocitylimit.log  velocitylimits.db
```

velocity-limit-app is highly flexible and pluggable application. It is designed to accept input and emit output to/from various sources. Currently it supports file/terminal as input and sink, but supporting more options should be a breeze. See  [design doc](design.md) for details.
Following example shows reading input from file "new_input.txt" and sending output to stdput

```python
root@75d3a8ec19bf:/go/src/velocitylimits/bin/linux_amd64# ./velocity-limit-app -infile=new_input.txt -stdout=true
{"id":"15887","customer_id":"528","accepted":true}
{"id":"30081","customer_id":"154","accepted":true}
{"id":"26540","customer_id":"426","accepted":true}
{"id":"10694","customer_id":"1","accepted":true}
{"id":"15089","customer_id":"205","accepted":true}
{"id":"3211","customer_id":"409","accepted":true}
:
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)