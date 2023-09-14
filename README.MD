## Golang Temp Folder Clearer

A simple Go application for clearing files in the temp folder with a given prefix.

## Flags
```
prefix string optional
    The prefix of the files to be deleted
    If not provided, all files in the temp folder will be deleted
    Multiple prefixes split by comma can be provided
```

## Usage
```
go run . prefix=<prefix-value>
```
## Example
Clear cdk temp files
```
go run . prefix=cdk
```
## build
```
go build .
```
