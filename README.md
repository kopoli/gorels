# Gorels

A git-tag generator with semantic versioning

## Installation

```
$ go get github.com/kopoli/gorels
```

## Description

This automates the version bumping and tag creation when dealing with git
repositories.

This is quite experimental as of yet and the abstractions are leaky.

For example creating the `v0.2.1` tag for this repository can be done with the
following:

```
$ gorels -verbose bump-patch tag

>> Found 0.2.0 as previous version
>> Bumping patch level
>> Creating the git tag: v0.2.1
>> Running: git tag --annotate -m gorels v0.2.1

Kalle Kankare (4):
      LICENSE: Add MIT license
      gorels: Remove currently unimplemented amend operation
      gorels: Rename Commands to Operations
      README: Add

 v0.2.1 HEAD
```

## Usage

```
$ gorels -h

Command line options:
  -D	Don't actually run any operations. Implies -verbose. (shorthand for -dryrun)
  -V	Enable verbose output. (shorthand for -verbose)
  -dryrun
    	Don't actually run any operations. Implies -verbose.
  -l	List operations. (shorthand for -list)
  -list
    	List operations.
  -v	Display version. (shorthand for -version)
  -verbose
    	Enable verbose output.
  -version
    	Display version.
```

After the command line options have been given, a space separated list of
operations are to be given. They are executed sequentially.

## Supported operations

The following operations are supported:

- **bump-major**: Bump the major version number.
- **bump-minor**: Bump the minor version number.
- **bump-patch**: Bump the patch level version number.
- **commit=**: Commit to operate on. Default: `HEAD`
- **git=**: Git program to use.
- **message=**: Message to inject into the tag
- **set-build=**: Set version build field.
- **set-prerelease=**: Set version pre-release field.
- **set-tag-prefix=**: Set tag prefix. Default `v`.
- **set-version=**: Set explicit version.
- **tag**: Create a tag using `git-tag`.

## License

MIT license
