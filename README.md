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
