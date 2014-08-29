# syntaxhighlight

Package syntaxhighlight provides syntax highlighting for code. It currently uses
a language-independent lexer and performs decently on JavaScript, Java, Ruby,
Python, Go, and C.

The main [`AsHTML(src []byte) ([]byte,
error)`](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight/symbols/go/github.com/sourcegraph/syntaxhighlight/AsHTML) function outputs HTML that uses the same CSS classes as [google-code-prettify](https://code.google.com/p/google-code-prettify/), so any stylesheets for that should also work with this package.

**[Documentation on Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)**

[![Build Status](https://travis-ci.org/sourcegraph/syntaxhighlight.png?branch=master)](https://travis-ci.org/sourcegraph/syntaxhighlight)
[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/syntaxhighlight/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)
[![authors](https://sourcegraph.com/api/repos/github.com/sourcegraph/syntaxhighlight/badges/authors.png)](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)
[![top func](https://sourcegraph.com/api/repos/github.com/sourcegraph/syntaxhighlight/badges/top-func.png)](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)
[![funcs](https://sourcegraph.com/api/repos/github.com/sourcegraph/syntaxhighlight/badges/funcs.png)](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)
[![Total views](https://sourcegraph.com/api/repos/github.com/sourcegraph/syntaxhighlight/counters/views.png)](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight)


## Installation

```
go get github.com/sourcegraph/syntaxhighlight
```


## Example usage

The function [`AsHTML(src []byte) ([]byte,
error)`](https://sourcegraph.com/github.com/sourcegraph/syntaxhighlight/symbols/go/github.com/sourcegraph/syntaxhighlight/AsHTML)
returns an HTML-highlighted version of `src`. The input source code can be in
any language; the lexer is language independent.

```go
package syntaxhighlight_test

import (
	"fmt"
	"github.com/sourcegraph/syntaxhighlight"
	"os"
)

func Example() {
	src := []byte(`
/* hello, world! */
var a = 3;

// b is a cool function
function b() {
  return 7;
}`)

	highlighted, err := syntaxhighlight.AsHTML(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(highlighted))

	// output:
	// <span class="com">/* hello, world! */</span>
	// <span class="kwd">var</span> <span class="pln">a</span> <span class="pun">=</span> <span class="dec">3</span><span class="pun">;</span>
	//
	// <span class="com">// b is a cool function
	// </span><span class="kwd">function</span> <span class="pln">b</span><span class="pun">()</span> <span class="pun">{</span>
	//   <span class="kwd">return</span> <span class="dec">7</span><span class="pun">;</span>
	// <span class="pun">}</span>

}
```


## Contributors

* [Quinn Slack](https://sourcegraph.com/sqs)

Contributions are welcome! Submit a pull request on GitHub.