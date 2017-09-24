# GoAdaline

GoAdaline is a Go implementation of an [Adaline neuron](https://en.wikipedia.org/wiki/ADALINE) using the Delta Rule. 

This project is helpful to demonstrate how Neural Networks work to those starting in Artificial Intelligence. 


## Getting Started

To try it out you just need to have Go installed. I used Go 1.9 so I can't assure that will work with lower versions
but it should work fine. 

To install it just run:

```
go get https://github.com/Plorenzo/goAdaline
``` 

## Prerequisites
The dataset has to be in csv format and the last column of it must be the expected output.

Also you need to have the data set in a folder divided into 3 files:
   * train.csv
   * validate.csv
   * test.csv
   
For best results train.csv should have about 70% of the rows and 15% the other files. 
   
   

   


## Flags
Required flags for the program to run: 

```
-path       "/path/to/folder" 
-cycles     100               //nº of training cycles
-lr         0.1               //Learning rate
```

Example:
```
./goAdaline -path=/Users/plorenzo/dev/data/ -cycles=100 -lr=0.1
```



# TODO's
Features to add, bugs to fix, improvements to make:

* ~~Add ability to pass file paths from terminal~~
* ~~Also cycles and learning rate~~
* Refactor readCSV() to convert to float64 whiles reads
* Change createCSV() to append to file instead of creating a new one each time