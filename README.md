# Terraria Server Wrapper

This is a simple go program designed to provide regular autosaving and
a save on SIGTERM for the vanilla Terraria server. The SIGTERM save
enables the server to be run headless in a docker container; `docker
stop` sends a SIGTERM to the PID 1 on stop, but will kill the
container if it continues to run. The vanilla server does not respond
to SIGTERM, but this wrapper will cause the server to save and exit
gracefully. Docker is the intended use-case for this wrapper, and
this, combined with go's excellent concurrency, is why the wrapper is
written in go. Go is statically-linked by default, so this wrapper can
be used in a docker container without requiring any additional
libraries.

The wrapper passes all command line arguments to the server
executable, and can effectively be called as though it were the server
executable. The server can still be controlled via the command line
UI; all stdin to the wrapper is passed through to the server.

## Configuration ##

All configuration is done via environment variables.

| Env Variable                 | Default                   | Description                                                                                                              |
|:----------------------------:|:-------------------------:|--------------------------------------------------------------------------------------------------------------------------|
| AUTOSAVE\_TIME               | 5m                        | How often to autosave. Autosaves are only performed if nothing has been typed in stdin before the autosave is performed. |
| SERVER\_EXECUTABLE | ./TerrariaServer.bin.x86_64 | The server executable to run.                                                                                            |
