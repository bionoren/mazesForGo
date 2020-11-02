This repo implements the concepts in the book [Mazes for Programmers](https://www.amazon.com/Mazes-Programmers-Twisty-Little-Passages/dp/1680500554) in golang.

Requirements
------------
golang 1.15+
Dependencies for the [pixel game library](https://github.com/faiface/pixel) - mostly opengl support. On Mac OS X, this is simply the xcode command line tools

Running
-------
interactive gui
```go
go build ./ && ./mazes
```

print statistics and exit
```go
go build ./ && ./mazes --stats
```
