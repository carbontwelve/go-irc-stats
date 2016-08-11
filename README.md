# Go IRC Stats

This is an attempted port of the Ruby based IRC stats generator [0x263b/Stats](https://github.com/0x263b/Stats) by [Abram Kash](https://kash.im/). At time of writing I knew neither Ruby nor Golang and this has been an excuse in learning more about both.

## Build instructions

To build on linux run `make clean-run` this will build the executable `logstats`. You will need to place `logstats`, `template.html` and `config.yaml` into your choice of destination folder. Amend `config.yaml` and give either a relative or absolute path to the logfile location before executing `logstats`.

To build on windows run `make build` (I dont have make installed on my windows machine so I wrote the make.bat with similar functionality). This will create a folder locally called `bin` where `logstats.exe`, `template.html` and `config.yaml` can be found.

To generate a logfile to test this with you can run `php createtestlog.php > irctest.log` to generate a logfile containing 720 days worth of randomised data.

## To do

Before undertaking this project I hadn't touched golang, ruby, windows bat files or linux Makefiles; this has been a successful learning experience for me. Below are the items that need finishing for this to become a full port of Abram's code.

* Implement most active times output
* Implement active users output
* Implement top users output
* Implement times of day output
* Implement share of activity output
* Implement activity over time output

Additional functionality I hope to add:

* Configurable input format, either via user defined regex or predefined parsers
* To be able to be used with a folder as input, and it will loop over all files in the folder and successfully parse their lines
* To have last parsed and time taken to parse output available to template
* Have the output be a single page html app with css and javascript embedded into `stats.html` and an additional `stats.json` output containing the parsed data.
* Write tests (as a way of learning how to test golang) and attempt 100% coverage (for no other reason than for the act of trying)
* Have createtestlog.php re-written in golang

## Contributing

Please don't allow the fact that this is a personal learning project deter you from forking and having a go at either implementing the missing functionality or adding the additional functionality. I really like to see how other people program.