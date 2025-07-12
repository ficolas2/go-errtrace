# Configure
You can create a custom tracer, to specify configuration or a custom formatter.
```go
tracer := errtrace.NewTracerBuilder()
    .TrimStackTrace("api/api.go")
    .MustBuild()
```
After that, you can set it as the default tracer
```go
errtrace.SetDefaultTracer(tracer)
```
Or  use it on its own
```go
tracer.Wrap(err)
```

## TrimStackTrace
Often, large parts of a stack trace — such as Go runtime internals or third-party library calls — aren’t useful when debugging your application.

`TrimStackTrace` allows you to specify a substring that marks the **last frame to keep** in the stack trace. It performs a simple `contains` check, so using a long, specific path helps ensure an accurate match.

For example, if you build and run your project with `--trimpath` (see [trimpath](#trimpath) section for details), you might write:
```go
traceBuilder.
    TrimStackTrace("github.com/ficolas2/example/api/api.go")
```

This will turn a stack trace like:
```
Stack trace:
  github.com/ficolas2/example/repository.(*Repository).GetUser
    github.com/ficolas2/example/repository/repository.go:13
  github.com/ficolas2/example/service.(*Service).GetUser
    github.com/ficolas2/example/service/service.go:13
  github.com/ficolas2/example/api.GetUser
    github.com/ficolas2/example/api/api.go:26
  main.main
    github.com/ficolas2/example/main.go:6
  runtime.main
    runtime/proc.go:283
  runtime.goexit
    runtime/asm_amd64.s:1700
```

Into simply:
```
Stack trace:
  github.com/ficolas2/example/repository.(*Repository).GetUser
    github.com/ficolas2/example/repository/repository.go:13
  github.com/ficolas2/example/service.(*Service).GetUser
    github.com/ficolas2/example/service/service.go:13
  github.com/ficolas2/example/api.GetUser
    github.com/ficolas2/example/api/api.go:26
```

## MaxStackDepth
```go
traceBuilder.
    MaxVarStackDepth(64)
```
The default is 1. Determines the max stack depth error.

## MaxVarStackDepth
```go
traceBuilder.
    MaxVarStackDepth(1)
```
The default is 1. Determines the max stack depth for the tracked variables. This will have no change in the default formatter, and is only useful if you use a custom formatter that takes advantage of a bigger variable stack depth.

## SetFormatter
Set a custom formatter.
```go
traceBuilder.
    SetFormatter(func (stacktrace []StackFrame, err error, vars []VarPoint) string { ... })
```
You can look at formatter.go to use the default formatter as an example.


# trimpath
When building or running Go programs, it's recommended to use the `--trimpath` flag. 
This removes absolute file paths from compiled binaries and replaces them with module-relative paths.

Without `--trimpath`, stack traces may include full local paths like:

`/home/user/dev/projects/example/api/api.go`

Which makes it harder to write portable and reusable stack trimming logic. With `--trimpath`, the same stack frame becomes:

`github.com/user/example/api/api.go`

To use `--trimpath`, pass it when building, running or testing:

```sh
go build --trimpath
go test --trimpath
go run --trimpath .
```

