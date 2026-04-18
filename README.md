# WUPE: Windows Update Pause Extender

Windows Update unapologetically restarts the operating system when an update is
due, causing unsaved work to be lost, long-running processes to be stopped.
This causes many users lots of grief.

Windows does not have a built-in way of permanently stopping the update
process. While there are ways to stop the system services completely, it
cripples some of the Windows's functionality (especially with 3rd party
utilities that come with hardware and expect WU to be up and running).

This is a system service written in Go that extends the Windows update pause
period every day so it is always 7 days in the future. It will prevent updates
from executing perpetually until the service is stopped and updates are
resumed.

## Installation

Get the wupe.exe binary from the [latest
release](https://github.com/foxbunny/wupe/releases/latest).

Copy wupe.exe to any location. Run:

```cmd
powershell -Command "New-EventLog -LogName Application -Source WUPEService"
sc create WUPEService binPath="C:\Path\To\wupe.exe"
sc description WUPEService "Extend Windows Update Pause indefinitely"
```

## Usage

To start:

```cmd
sc start WUPEService
```

To stop:

```cmd
sc stop WUPEService
```

To query state:

```cmd
sc query WIPEService
```

## Uninstallation

To uninstall:

```cmd
sc delete WUPEServices
powershell -Command "Remove-EventLog -Source WUPEService"
```

And delte the .exe file.


## Issues

No known issues.
