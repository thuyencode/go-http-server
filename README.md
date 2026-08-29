# Go HTTP Servers Tutorial

I was coding along the [Go HTTP Servers](https://www.youtube.com/playlist?list=PLlv5lnjOHQo7m2w-KjtTZq10ZNVfNHHdP) playlist to learn Go, obviously.

The current HTTP endpoints:

- `/`:  Return `"Welcome!"` on happy path
- `/goodbye`:  Return `"Goodbye!"` on happy path
- `/hello?name={name}`: Return `"Hello, {name}!"` on happy path
- `/param/{name}`: Return `"Hello, {name}!"` on happy path
- `/header`: Provide a key-value pair `"name: {name}"` will return `"Hello, {name}!"` on happy path
- `/json`: Provide a JSON with one key-value pair `"name: {name}"` will return `"Hello, {name}!"` on happy path
