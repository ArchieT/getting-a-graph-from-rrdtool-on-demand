# getting-a-graph-from-rrdtool-on-demand

Generating a few graphs every 30sec on a Raspberry Pi is not necessarily a good idea.
The graphs were CPU temperature graphs, and there were peaks clearly visible from it.

I wanted to generate those graphs on demand.

Install it like

`go get -u -v github.com/ArchieT/getting-a-graph-from-rrdtool-on-demand`

How to use this daemon? Point it to use your rrd files and leave it running

`getting-a-graph-from-rrdtool-on-demand <name> <filepath> <name> <filepath> ...`

where `name` is both an alias and the RRA name.

I use ArchLinuxARM, and I put it into a systemd service like this

`/etc/systemd/system/graphgetting.service`

```
[Unit]
Description=getting-a-graph-from-rrdtool-on-demand
After=syslog.target network.target

[Service]
ExecStart=/home/mf/gopath/bin/getting-a-graph-from-rrdtool-on-demand temp /home/mf/temperatury.rrd

[Install]
WantedBy=multi-user.target
```

The code is really bad, so I licensed it with BSD2. But if you make any improvements to it, please share them. Even if you feel that they are even worse. And don't forget to give it a star if you somehow found it useful (because I don't expect anyone to).

This program uses github.com/ziutek/rrd library (https://godoc.org/github.com/ziutek/rrd#Grapher)
