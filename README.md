# homedir

Forked from [mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)
after the original repository was archived.

## Overview
This is a Go library for detecting the user details (username and  home directory) without
the use of cgo, so the library can be used in cross-compilation environments.

Usage is incredibly simple:
* `user.Username()` returns the username of the current user
* `user.HomeDir()` returns the home directory for the current user
* `user.ExpandPath("~/.config")` expands the `~` in a path to the user's home directory

## Why not just use `os/user`?
The built-in `os/user` package requires cgo on Darwin systems. This means
that any Go code that uses that package cannot cross compile. I needed a
cross-compilable method to get the username and home directory for a few
of my packages, so I forked the original and expanded upon it.
