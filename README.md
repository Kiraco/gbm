# GBM Challenge

This project is a challenge from GBM, the program takes a txt file with a list of paths to json files that holds operations.

# Usage

The program receive a path to a `.txt` file that holds all the json paths that holds the orders to be run. It will execute all json in the txt file.
The output of the program will be store in the root folder as `gbm_file_[indexOfFile]_[timestamp].json` along with `gbm_errors_[timestamp]`
for any error thrown while reading the files in the batch.

The `.txt` file that holds the path to the `.json` files, must include the absolute path to the file.

# Commands

Build: 
```go build```

Run:
```./gbm '[path to the txt file]"```
 
Test:
```go test ./... -coverprofile gbm.ou```

View Coverage from test:
```go tool cover -html=gbm.out```