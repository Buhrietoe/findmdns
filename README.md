# findmdns
Just request all mdns services and allow 5 seconds for them to respond.
Output in JSON.

## Install it

`go install -v github.com/Buhrietoe/findmdns@latest`

## Why?
Its simpler.

## Caveat
A correctly setup system for mdns.

Example:
Archlinux using systemd-resolved with /etc/resolv.conf symlinked to /run/systemd/resolve/stub-resolv.conf
