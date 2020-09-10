robots
=========

robots  — is parser for [robots.txt](www.robotstxt.org) files for Go 

Installation
------------

```bash
$ go get github.com/gvenkat/robots 
```

Usage
-----

Here's an example of using robots

```go
import "github.com/gvenkat/robots"

...

// Instantiate from any io.Reader instance
parser := robots.FromReader(...)

// Or from a URL
parser := robots.FromURL(...)

// Or from a file 
parser := robots.FromFile(...)



```

Default crawler user-agent is:

    Mozilla/5.0 (X11; Linux i686; rv:5.0) Gecko/20100101 Firefox/5.0

License
-------

See [LICENSE](https://github.com/gvenkat/robots/blob/master/LICENSE)
file.


Resources
=========

  * [Robots.txt Specifications by Google](http://code.google.com/web/controlcrawlindex/docs/robots_txt.html)
  * [Robots.txt parser for JavaScript](https://github.com/ekalinin/robots.js)
  * [A Standard for Robot Exclusion](http://www.robotstxt.org/orig.html)
