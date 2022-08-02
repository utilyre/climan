# CLIMan

Parse and run `.http` files from command line blazingly fast! 🦔 (idk I
couldn't find a gopher emoji)

## 📦 Installation

Install via go

```bash
go install github.com/utilyre/climan
```

If you're using neovim you might need to add http file detection

```lua
-- ~/.config/nvim/ftdetect/http.lua

vim.api.nvim_create_augroup("fthttp", {})
vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  group = "fthttp",
  pattern = "*.http",
  callback = function(a) vim.bo[a.buf].filetype = "http" end,
})
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

Whoa whoa whoa, we got a lot to talk about here. `###` is a special type of
comment that separates requests. Here we have a simple `GET` request to
`jsonplaceholder` API which is pretty similar to the previous example.

But we also got a `PUT` request after the `GET` request (notice they are
separated by `###`). And right after the URL line you can put as much headers
as you want in a form like `name: value`.

In this case, since we set `application/json` as `Content-Type` we need to
specify a request body. As you've probably noticed the request body should (at
least) have a single empty line before its beginning.

See [examples](/examples)

## 🔖 Todos

- [ ] Support for other request body types (like `text/xml` and `text/plain`).
- [ ] Parse the response body base on `Content-Type` header.

## 📢 Credits

- [REST Client](https://github.com/Huachao/vscode-restclient)
