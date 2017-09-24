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

### Prerequisites
The dataset has to be in csv format and the last column of it must be the expected output.

Also you need to have the data set in a folder divided into 3 files:
   * train.csv
   * validate.csv
   * test.csv
   
For best results train.csv should have about 70% of the rows and 15% the other files. 
   
You can use this sample data to test the neuron. It has been normalized and randomized before 
splited into files. 

**Sample data**
* Train data:       [https://pastebin.com/H3YgFF0a](https://pastebin.com/H3YgFF0a)
* Validate data:    [https://pastebin.com/aeK6krxD](https://pastebin.com/aeK6krxD)
* Test data:        [https://pastebin.com/mt5P8AZS](https://pastebin.com/mt5P8AZS)

You can find out more info about the dataset here [http://sci2s.ugr.es/keel/dataset.php?cod=44](http://sci2s.ugr.es/keel/dataset.php?cod=44)
### Flags
Required flags for the program to run: 

```
-path       "/path/to/folder" 
-cycles     100                     //nº of training cycles
-lr         0.1                     //Learning rate
-out        "path/to/destination"
```

Example:
```
./goAdaline -path=/Users/plorenzo/dev/data/ -cycles=100 -lr=0.1
```
## Output
The neuron will output the squared error obtained in the test data and the weights used.
```
Test error:
0.017166281479163065
Weights:
[0.6283026308947595 0.5098346273068519 0.2182228283924059 -0.2505429228543444 0.11335767250598956 0.03879773117815263 0.09743364005207034 0.5027576658656207 -0.011776841084895283]
```
Also a **csv file** with the train and validate error for each cycle, the neuron output for the test data and the weights used. So you can make cool charts like this:

![Adaline neuron error chart](https://i.imgur.com/dM1xWom.png)


## TODO's
Features to add, bugs to fix, improvements to make:

* ~~Add ability to pass file paths from terminal~~
* ~~Also cycles and learning rate~~
* Refactor readCSV() to convert to float64 whiles reads
* ~~Add optional flag for path to error output csv~~
* ~~Print final weights~~
* Make readCSV() calls concurrent