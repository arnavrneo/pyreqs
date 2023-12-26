# pyreqs

`pyreqs` is a command line tool for auto-generating python dependencies file (requirements.txt) for your project.

```
Seamlessly generate python dependencies (requirements.txt) file. 
Call pyreqs with the available commands to see the magic done!

Usage:
  pyreqs [flags]
  pyreqs [command]

Available Commands:
  clean       Cleans up your requirements.txt
  create      Generates a requirements.txt file
  help        Help about any command

Flags:
  -d, --dirPath string    directory to .py files (default "./")
  -h, --help              help for pyreqs
  -i, --ignore string     ignore specific directories (each seperated by comma) (default " ")
  -p, --print             print requirements.txt to terminal
  -v, --venvPath string   directory to venv (virtual env) (default " ")

Use "pyreqs [command] --help" for more information about a command.

```

## Documentation

#### Creating the requirements.txt

- Searches for all python files in the current working directory
```bash
pyreqs create
```

- Searches for python files in specified directory
```bash
pyreqs create -d <path-to-your-project>
```

- Searches for python files and gets packages info from `venv` directory
```bash
pyreqs create -d <path-to-your-project> -v <path-to-venv-folder>
```

- Search for python files excluding specific directories (multiple dirs are seperated by comma)
```bash
pyreqs create -d <path-to-your-project> -i <dir1-to-ignore>,<dir2-to-ignore>
```

#### Cleaning the requirements.txt

```bash
pyreqs clean -r <path-to-requirements.txt> -d <path-to-project>
```

## Installation

### Linux & MacOS

- Download the latest compatible binary from the releases and add its path to the environment variable.


## Contributing

If anything feels off, or if you feel that some functionality is missing, go and hit-up an issue!
