Runtil
======================================================================

Execute command until sometime in a day.

```
NAME:
   runtil - execute command until sometime in a day

USAGE:
   runtil [global options] [arguments...]

GLOBAL OPTIONS:
   --until HH:MM:SS, -u HH:MM:SS       process running until HH:MM:SS [$RUNTIL_UNTIL]
   --signal signal, -s signal          signal to send at HH:MM:SS (default: SIGINT) [$RUNTIL_SIGNAL]
   --kill-after duration, -k duration  waiting duration before sending KILL after the signal sent (default: 10s) [$RUNTIL_KILL_AFTER]
   --log-level level, -l level         level [debug|info|warn|error|dpanic|panic|fatal] (default: info) [$RUNTIL_LOG_LEVEL]
   --debug, -d                         use debug config of zap logger (default: false) [$RUNTIL_DEBUG]
   --help, -h                          show help (default: false)
   --version, -v                       print the version (default: false)
```
