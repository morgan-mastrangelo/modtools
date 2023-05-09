modtools - Keep your go.mod dependencies up-to-date.
====

modtools is a command-line toolset designed to help keeping the go.mod dependencies up-to-date while still allowing exceptions.

Why?
---

Dependency pinning has its obvious advantages: for a project with dozens of dependencies you don't want your
build to break unpredictably due to a change in one of them. The flip side is that dependencies tend to become
stale. The changes within the same major version should maintain compatibility and in most cases they don't actually
break anything, but instead provide new features, bugfixes and security updates. That's why upgrading should be
done on a regular basis.

Modtools is designed to make the process of upgrading dependencies easier. Running `modtools check` shows the list
of direct and indirect dependencies that can be upgraded. It prints the list of commands that need to be run and
exits with a non-zero code if such dependencies exist:

```console
$ go install github.com/dop251/modtools@latest
$ modtools check
Some dependencies are out-of-date. Please upgrade by running 'modtools update' or the following commands:

go get github.com/kr/pretty@v0.2.1
go get github.com/kr/text@v0.2.0
go get gopkg.in/check.v1@v1.0.0-20201130134442-10cb98267c6c

Error: check has failed
```

This command could be added to a CI pipeline running on a schedule or to a commit hook.

In case a new version of a dependency causes a problem it can be added to the exception list by running
`modtools freeze modpath` so that it's ignored for up to the specified number of days:

```console
$ modtools freeze gopkg.in/check.v1 14
Don't forget to add modtools_frozen.yml to the repository.
```

During this time the necessary adjustments need to be made to accommodate for the change (if it was deliberate),
or a bug report should be raised if it wasn't. When the problem is fixed
`modtools thaw modpath` can be used to remove `modpath` from the list of exceptions.

