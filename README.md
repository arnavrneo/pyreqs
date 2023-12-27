# pyreqs

`pyreqs` is a command line tool for auto-generating python dependencies file (requirements.txt) for your project.

![](/assets/main-win-crop.png)

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

