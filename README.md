# WUPE: Windows Update Pause Extender

Windows Update unapologetically restarts the operating system when an update is due, causing unsaved work to be lost, long-running processes to be stopped. This causes many users lots of grief.

Windows does no have a built-in way of permanently stopping the update process. While there are ways to stop the system services completely, it cripples some of the Windows's functionality (especially with 3rd party utilities that come with hardware and expect WU to be up and running).

This is a system service written in Go that extends the Windows update pause period by 7 days every day. It will prevent updates from executing perpetually until the service is stopped and updates are resumed.

## Installation


## Usage


## Issues