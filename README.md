Temporary Workspaces Manager (ws)
=================================

Introduction:
-------------

`ws` is a lightweight bash program designed to create and manage temporary workspaces. It proves useful in command-line intensive workflows, allowing users to quickly create temporary directories for tasks like prototyping programs, transcoding audio, or downloading files. This program was developed as a means to sharpen my skills in the Go programming language. It incorporates a shell wrapper to access shell built-ins within the current session.

Installation:
-------------

1. Run the following commands
```bash
make
sudo make install
mkdir ~/ws
```

2. Add the following line to your .zshrc or .bashrc file:
```bash
eval "$(ws_internal activate)"
```

Usage:
-----

```text
Usage: ws [command]

ws is a program for managing temporary workspaces.

Commands:
  list|ls          List all available workspaces
  new              Create a new workspace and change into it
  print_current_workspace|pcw
                   Print the current workspace path
  recent|rec       Change to the most recently used workspace
  next|n           Change to the next workspace in the history
  prev|p           Change to the previous workspace in the history
  remove|rm        Remove the current workspace and return to the previous one
  return|ret       Return to the previous workspace
  usage            Show this usage information
```

Examples:
---------

Create a new workspace and change into it:

```text
[user@ubuntu ~] ws new
[user@ubuntu 2023-06-25_07:35:20.2685626823]
```
Return to the previous location with `ret`:

```text
[user@ubuntu 2023-06-25_07:35:20.2685626823] ws ret
[user@ubuntu ~]
```

List workspaces with `ls`:

```text
[user@ubuntu ~] ws ls
1 2023-06-25_07:29:02.3918850963 0 days 0 hours 6 min 58 secs
2 2023-06-25_07:35:20.2685626823 0 days 0 hours 0 min 41 secs
```

Change into the most recently accessed workspace with `ws`:

```text
[user@ubuntu ~] ws
[user@ubuntu 2023-06-25_07:35:20.2685626823]
```

Move back one workspace with `prev` or `p`:

```text
[user@ubuntu 2023-06-25_07:35:20.2685626823] ws p
[user@ubuntu 2023-06-25_07:29:02.3918850963]
```

Move forward one workspace with `next` or `n`:

```text
[user@ubuntu 2023-06-25_07:29:02.3918850963] ws n
[user@ubuntu 2023-06-25_07:35:20.2685626823]
```
