# nkill

Kills all processes listening on the given TCP ports.

#### Install

You need go installed and GOBIN in your PATH. Once that is done, run the command:

```bash
   $ go get -u github.com/ridha/nkill
```

#### Usage

To kill any process listening to the port 8080:

```bash
    nkill 8080
```

Sometimes process fork and will need to be killed many times:

```bash
    watch -n 0 "nkill 8080"
```

##### Inspiration

http://voorloopnul.com/blog/a-python-netstat-in-less-than-100-lines-of-code/
