# Go HTTP Servers Tutorial

I was coding along the [Go HTTP Servers](https://www.youtube.com/playlist?list=PLlv5lnjOHQo7m2w-KjtTZq10ZNVfNHHdP) playlist to learn Go, obviously.

## Prerequisite

Some targets in [`Makefile`](./Makefile) rely one [`air`](https://github.com/air-verse/air) and [`golangci-lint`](https://golangci-lint.run/docs/welcome/install/local/). Make sure to install them if you intend to use `Makefile`.

If you want to use my git hooks, copy them into your local `.git/hooks`.

```sh
cp script/githooks/* .git/hooks
```

Note: You will need the cli tools I mentioned above.

## Wiki

The current HTTP endpoints:

- `/`:  Return `"Welcome!"` on happy path
- `/goodbye`:  Return `"Goodbye!"` on happy path
- `/hello?name={name}`: Return `"Hello, {name}!"` on happy path
- `/param/{name}`: Return `"Hello, {name}!"` on happy path
- `/header`: Provide a key-value pair `"name: {name}"` will return `"Hello, {name}!"` on happy path
- `/json`: Provide a JSON with one key-value pair `"name: {name}"` will return `"Hello, {name}!"` on happy path
