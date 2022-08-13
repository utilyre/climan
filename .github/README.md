# ✨ CLI Man

Make http requests from command line **BLAZINGLY FASTER**!

## 📦 Installation

Install via go

```bash
go install github.com/utilyre/climan@latest
```

If you're using neovim you might need to add http file detection as well

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

Http file example (see [examples](/examples) for more)

```http
# example.http

# Method URL
POST http://localhost:8080
# Header: Value
Content-Type: application/json
User-Agent: Mozilla/5.0 (Windows NT 6.0; en-US; rv:1.9.1.20) Gecko/20140827 Firefox/35.0

# Body (only json is supported for now)
{
  "title": "climan",
  "description": "A file based HTTP client"
}

# The next line is an special type of comment that separates requests
###

GET http://localhost:8080
User-Agent: Mozilla/5.0 (Windows NT 6.0; en-US; rv:1.9.1.20) Gecko/20140827 Firefox/35.0
```

**NOTE**: Trailing comments are _NOT_ available.

---

Make the first request

```bash
climan example.http
```

Make the second request

```bash
climan -r 2 example.http
```

Show response headers of the first request
```bash
climan -v example.http
```

## 🔖 Todos

- [ ] Support for other request body types (like `text/xml` and `text/plain`).
- [x] Parse the response body base on `Content-Type` header.
- [x] Show response details.
- [ ] Colored output.
- [ ] Add more `http` examples.

## 📢 Credits

- [REST Client](https://github.com/Huachao/vscode-restclient)
