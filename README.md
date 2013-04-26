# gorax

The Go ecosystem seems to lack a comprehensive cloud services API (at the time this README was first written).
As both Go and cloud services are trending in many businesses, and with Go used increasingly in infrastructure, it seems like an odd omission.
To fill this gap, gorax provides a Go binding to the Rackspace cloud APIs.
Rackspace offers many APIs that are compatible with OpenStack, and thus provides an ideal springboard for wider OpenStack technology adoption in the Go community.

The name, gorax, is derived from its sister (and, presently, far more complete) project Pyrax, which offers similar capabilities for Python.

**This library is still in the very early stages of development. Unless you want to contribute, it probably isn't what you want**

## Installation and Testing

To install:

```bash
go get github.com/racker/gorax
```

To run unit tests:

```bash
go test github.com/racker/gorax/v2.0/cloud/server
go test github.com/racker/gorax/v2.0/identity
```

**To run integration tests and examples, you'll need a Rackspace cloud user account and password.**
You may find the tests and examples in the `github.com/racker/gorax/exampels` directory.
Note that most, if not all, of the examples therein require a username and password to run.
For this reason, no scripting to invoke these tests exist.

## Contributing

The following guidelines are preliminary, as this project is just starting out.
However, this should serve as a working first-draft.

### Branching

The master branch must always be a valid build.
The `go get` command will not work otherwise.
Therefore, development must occur on a different branch.
As with Pyrax, we choose to stage all changes in a branch named "working."

When creating a feature branch, do so off the working branch:

```bash
git checkout working
git pull
git checkout -b featureBranch
git checkout -b featureBranch-wip
```

Perform all your editing and testing in the WIP-branch.
Feel free to make as many commits as you see fit.
You may even open "WIP" pull requests from your feature branch to seek feedback.
WIP pull requests will **never** be merged, however.

To get code merged, you'll need to "squash" your changes into a single commit.
These steps should be followed:

```bash
git checkout featureBranch
git merge --squash featureBranch-wip
git commit -a
git push origin featureBranch
```

You may now open a nice, clean, self-contained pull request from featureBranch to working.

The `git commit -a` command above will open a text editor so that
you may provide a comprehensive description of the changes.

In general, when submitting a pull request against working,
be sure to answer the following questions:

- What is the problem?
- Why is it a problem?
- What is your solution?
- How does your solution work?  (Recommended for non-trivial changes.)
- Why should we use your solution over someone elses?  (Recommended especially if multiple solutions being discussed.)

Remember that monster-sized pull requests are a bear to code-review,
so having helpful commit logs are an absolute must to review changes as quickly as possible.

Finally, (s)he who breaks working is ultimately responsible for fixing working.

### Source Representation

The Go community firmly believes in a consistent representation for all Go source code.
We do too.
Make sure all source code is passed through "go fmt" *before* you create your pull request.

Please note, however, that we fully acknowledge and recognize that we no longer rely upon punch-cards for representing source files.
Therefore, no 80-column limit exists.
However, if a line exceeds 132 columns, you may want to consider splitting the line.

### Unit and Integration Tests

Pull requests that include non-trivial code changes without accompanying unit tests will be flatly rejected.
While we have no way of enforcing this practice,
you can ensure your code is thoroughly tested by always [writing tests first by intention.](http://en.wikipedia.org/wiki/Test-driven_development)

When creating a pull request, if even one test fails, the PR will be rejected.
Make sure all unit tests pass.
Make sure all integration tests pass.

### Documentation

Private functions and methods which are obvious to anyone unfamiliar with gorax needn't be accompanied by documentation.
However, this is a code-smell; if submitting a PR, expect to justify your decision.

Public functions, regardless of how obvious, **must** have accompanying godoc-style documentation.
This is not to suggest you should provide a tome for each function, however.
Sometimes a link to more information is more appropriate, provided the link is stable, reliable, and pertinent.

Changing documentation often results in bizarre diffs in pull requests, due to text often spanning multiple lines.
To work around this, put [one logical thought or sentence on a single line.](http://rhodesmill.org/brandon/2012/one-sentence-per-line/)
While this looks weird in a plain-text editor,
remember that both godoc and HTML viewers will reflow text.
The source code and its comments should be easy to edit with minimal diff pollution.
Let software dedicated to presenting the documentation to human readers deal with its presentation.

## Examples

### Image List from a Region

```go
package main

import (
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("Usage: I need both username and API key on CLI, in that order.")
	}
	username := os.Args[1]
	password := os.Args[2]

	id := identity.NewIdentity(username, password, "")
	err := id.Authenticate()
	if err != nil {
		panic(err)
	}

	region, err := servers.RegionByName(id, "dfw")
	if err != nil {
		panic(err)
	}

	images, err := region.Images()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%-36s   Name\n", "UUID")
	for _, i := range images {
		fmt.Printf("%36s - %s\n", string(i.Id[0:36]), i.Name)
	}
}
```

### Cloud Monitoring

Preliminary.

```go
package main

import (
  "github.com/racker/gorax/monitoring"
  "github.com/racker/gorax/identity"
  "fmt"
)

func main() {
  cm := monitoring.MakePasswordMonitoringClient("https://monitoring.api.rackspacecloud.com/v1.0", identity.USIdentityService, "username", "password")
  cm.SetDebug(true)
  checks, err := cm.ListChecks("enuzk2tiph")
  fmt.Printf("%s\n", checks)
  fmt.Printf("%s\n", err)
}
```
