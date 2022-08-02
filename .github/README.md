# CLIMan

Parse and run `.http` files from command line blazingly fast! 🦔 (idk I
couldn't find a gopher emoji)

## 📦 Installation

Intall via go:

```bash
go install github.com/utilyre/climan
```

## 🚀 Usage

Create a file with the following content

```http
# test.http

GET https://jsonplaceholder.typicode.com/posts
```

And run this command on the file

```bash
climan -n 1 test.http
```

Now let's see what we just ran

- `-n 1`: Tells climan to run the first request
- `test.http`: The file we've just created

So you might be thinking

> Can I have multiple requests in a single file?

Yes, you have the option to do that. Let's move on to another example.

Paste this in another file

```http
# multi-request.http

GET https://jsonplaceholder.typicode.com/comments?postId=1

###

PUT https://jsonplaceholder.typicode.com/posts/1
Content-Type: application/json

{
  "title": "Edited"
}
```

Woah woah woah, we got a lot to talk about here. `###` is a special type of
comment that separates requests. Here we have a simple *GET* request to
`jsonplaceholder` api which is pretty similiar to the previous example.

But we also got a *PUT* request after the *GET* request (notice they are
separated by `###`). And right after the url line you can put as much headers
as you want in a form like `name: value`.

In this case, since we set application/json as `Content-Type` we need to
specify a request body. As you've probably noticed the request body should (at
least) have a single empty line before its beginning.

## 🔖 Todos

- [ ] Support for other request body types (like xml and plain text).
- [ ] Parse the response body base on `Content-Type` header.

## 📢 Credits

- [REST Client](https://github.com/Huachao/vscode-restclient)
