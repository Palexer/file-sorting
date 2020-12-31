# File-Sorting

File-Sorting looks for newly created files in a specified directory and automatically sorts them
after their year and month of creation. 

**Example:**

You put a file created on the 31st December 2020 in a directory. The programm now automatically creates
creates the subdirectory "2020" and inside the subdirectory a new directory named "12". Then it moves the file to
./2020/12/

*Note: Make sure to only add/create one file after another, otherwise the program might not correctly*

**Usage:**

``` bash

./filesorting -dir /path/to/my/directory

```

**Supported systems:**

- Linux amd64/arm64
- Windows amd64
- macOS amd64

**License:**

MIT
