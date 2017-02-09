# nkill

Kills all processes listening on the given TCP ports.

#### Install

```bash
   $ git clone https://github.com/ridha/nkill
   $ cd nkill; go build
   $ sudo ln -s /path/to/nkill /usr/bin/nkill
```

#### Usage

To kill any process listening to the port 8080:

```bash
    nkill 8080
```

Sometimes, process fork and will need to be killed many times, just run something like:

```bash
    watch -n 0 "nkill 8080"
```

##### Inspiration

http://voorloopnul.com/blog/a-python-netstat-in-less-than-100-lines-of-code/
