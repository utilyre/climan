<h1 align="center">CLI Man</h1>

<p align="center">
  Make HTTP requests from command line <strong>BLAZINGLY FASTER</strong>!
</p>

## 📦 Installation

- [Latest Release](https://github.com/utilyre/climan/releases/latest).

- Manual

  Clone the latest version of climan

  ```bash
  git clone --depth=1 --branch=v0.3.0 https://github.com/utilyre/climan.git
  ```

  Hop into the cloned repo and build

  ```bash
  cd climan
  go build
  ```

## 🌟 Integration

### Neovim

In order to get syntax highlighting with nvim-treesitter plugin create a file
at `$XDG_CONFIG_HOME/nvim/ftdetect/http.lua` with the following content

```lua
-- ~/.config/nvim/ftdetect/http.lua

vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  group = vim.api.nvim_create_augroup("fthttp", {}),
  pattern = { "*.http", "*.rest" },
  callback = function(a)
    vim.bo[a.buf].filetype = "http"
    vim.bo[a.buf].commentstring = "#%s"
  end,
})
```

## 🚀 Usage

HTTP file example (see [examples](/examples) for more)

```http
# example.http

# Method URL
POST http://localhost:8080
# Header: Value
Content-Type: application/json
User-Agent: Mozilla/5.0 (Windows NT 6.0; en-US; rv:1.9.1.20) Gecko/20140827 Firefox/35.0

# Body
{
  "title": "climan",
  "description": "A file based HTTP client"
}

# The next line is an special type of comment that separates requests
###

GET http://localhost:8080
User-Agent: Mozilla/5.0 (Windows NT 6.0; en-US; rv:1.9.1.20) Gecko/20140827 Firefox/35.0
```

**NOTE**: Trailing comment is _NOT_ available.

---

Make the first (and maybe the only) request of `example.http`

```bash
climan example.http
```

Make the second request of `example.http`

```bash
climan -i 2 example.http
```

Make the last request of `example.http`

```bash
climan -i -1 example.http
```

Make the second request of `example.http` and output verbosely

```bash
climan -vi 2 example.http
```

Learn more

```bash
climan -h
man climan
```

## 🔖 Todos

- [x] Support for other request body types (like `text/xml` and `text/plain`).
- [x] Parse the response body base on `Content-Type` header.
- [x] Show response details.
- [x] Colored output.
- [ ] Colored response body.
- [x] Environment variable support.
- [ ] Add more `http` examples.

## 📢 Credits

- [REST Client](https://github.com/Huachao/vscode-restclient)
