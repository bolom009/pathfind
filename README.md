# pathfind

Package pathfind finds the shortest path based on graph of squares.

The algorithm works as follows:

First setup:
- determine all squares wrapped polygon with holes
- build a visibility graph based on squares edges

Path search:
- add start and end points to visibility graph (include dynamic obstacles are optional)
- use the A* search algorithm (package [astar](https://github.com/fzipp/astar))
  on the visibility graph to find the shortest path

## Requirements for executing demo

##### Ubuntu

    apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev

##### Windows

###### cgo

On Windows you need C compiler, like [Mingw-w64](https://mingw-w64.org) or [TDM-GCC](http://tdm-gcc.tdragon.net/).
You can also build binary in [MSYS2](https://msys2.github.io/) shell.

To remove console window, build with `-ldflags "-H=windowsgui"`.

## Demo

https://github.com/user-attachments/assets/d25f2d41-68bd-49be-84b7-0ed3bef3814f
