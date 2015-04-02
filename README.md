Ninja Sphere Development CLI
============================

Instructions for OSX and Linux

* install Go, git and mercurial, docker and boot2docker (OSX only)
* ensure you have defined GOPATH and that $GOPATH/bin is in your PATH
* run:

<code>
go get github.com/ninjasphere/ninja-dev-cli/ninja
</code>

* check out a Ninja Sphere driver or application

<code>
go get -d github.com/ninjasphere/driver-samsung-tv
</code>
  
* build it

<pre>
cd $GOPATH/src/github.com/ninjasphere/driver-samsung-tv
ninja build .
</pre>



