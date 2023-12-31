# pyreqs

`pyreqs` is a command line tool for auto-generating python dependencies file (requirements.txt) for your project.

![](https://i.imgur.com/FaYjTro.png)

| Create Command | Clean Command |
| --- | --- |
| ![](https://i.imgur.com/73bduN1.png) | ![](https://i.imgur.com/FV6kyLq.png) |


## Installation

### Linux and MacOS

- Download the latest compatible binary from the releases or from below:
 
| Platform | Release |
| ------------- | ------------- |
| Linux-x64  | [Link](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Linux_x86_64.tar.gz)  |
| Linux-arm64  | [Link](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Linux_arm64.tar.gz)  |
| Linux-i386  | [Link](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Linux_i386.tar.gz)  |
| MacOS-x64  | [Link](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Darwin_x86_64.tar.gz)  |
| MacOS-arm64  | [Link](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Darwin_arm64.tar.gz)  |

- Untar it;
- Run by `./pyreqs`.

> [!NOTE]
> Linux and MacOS users need to give read-write permission to `pyreqs`.
> 
> For that, do `chmod 755 pyreqs` before running.


### Windows

- Download the latest compatible window release: [windows-x64](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Windows_x86_64.zip), [windows-i386](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Windows_i386.zip), [windows-arm64](https://github.com/arnavrneo/pyreqs/releases/download/v1.2.2/pyreqs_Windows_arm64.zip);
- Unzip it;
- Open `cmd` in the directory where the `pyreqs.exe` is downloaded and run `pyreqs.exe` from the `cmd`.


> [!TIP]
> For running `pyreqs` anywhere from the `terminal` or `cmd`, add `pyreqs` path to the system's path variables.
> ### Linux
> - Copy `pyreqs` to `/usr/local/bin`: `cp ~/pyreqs /usr/local/bin/`.
> ### MacOS
> - `cp ~/pyreqs /etc/paths/`.
> ### Windows
> - Add `pyreqs.exe` to the set of environment variables.

## Contributing

If anything feels off, or if you feel that some functionality is missing, go and hit-up an issue! 🌟

