# ttee【tiːtiː】

tee command with time

# Usage

```bash
$ ttee -h
Usage:
  ttee [OPTION]... [FILE]...

Application Options:
  -a, --append             Append to the given FILEs, do not overwrite
  -i, --ignore-interrupts  Ignore interrupt signals
  -c, --clock-time         Display clock time

Help Options:
  -h, --help               Show this help message

$ seq 3 | xargs -i sh -c "echo sleep {} sec ; sleep {} ; echo done" | ttee foo.log
[00:00:00.000] sleep 1 sec
[00:00:00.993] done
[00:00:00.996] sleep 2 sec
[00:00:02.998] done
[00:00:03.001] sleep 3 sec
[00:00:06.003] done
```
