goreplace
=========
A find and replace tool written in Go, because I got tired of messing with sed and/or perl. It's crude and untested, so use at your own risk :)

Usage: goreplace TOFIND TOREPLACE PATTERN...

PATTERN can either be a file or a regex pattern. Goreplace will find any file with the name specified, AND any file in the current directory matching the pattern as regex.

So, stuff goes down as follows.
~~~ sh
$ cat /home/thedude/myfiles/world.txt
hey world!
$ goreplace hey hello /home/thedude/myfiles/world.txt
Done! Changed 1 file(s).
$ cat /home/thedude/myfiles/world.txt
hello world!

$ cd /home/thedude/myfiles
$ cat someotherfile.txt
hello hello hello
$ goreplace hello hey *.txt
Done! Changed 2 file(s)
$ cat world.txt
hey world!
$ cat someotherfile.txt
hey hey hey
~~~

It's like magic.
