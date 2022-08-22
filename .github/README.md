# ✨ CLI Man

Make http requests from command line **BLAZINGLY FASTER**!

## 📦 Installation

### GitHub Releases

1. Go to [latest](https://github.com/utilyre/climan/releases/latest) page.

2. Download the appropriate `.tar.gz` archive that suits your operating system and architecture.

3. Extract the archive

```bash
tar -xzf climan-[OS]-[ARCH].tar.gz
```

4. Install the included files (\*nix only)

```bash
sudo install -Dm755 bin/climan /usr/local/bin/climan
```

### Go Modules

1. Have golang setup on your system

```bash
export GOPATH="$XDG_DATA_HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

2. Install climan globally

```bash
go install github.com/utilyre/climan@latest
```

## 🌠 Integration

### Neovim

Create a file at `$XDG_CONFIG_HOME/nvim/ftdetect/http.lua` with the following content

```lua
-- ~/.config/nvim/ftdetect/http.lua

vim.api.nvim_create_augroup("fthttp", {})
vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  group = "fthttp",
  pattern = "*.http",
  callback = function(a)
    vim.bo[a.buf].filetype = "http"
    vim.bo[a.buf].commentstring = "#%s"
  end,
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
climan -index 2 example.http
```

Make the last request of `example.http`

```bash
climan -index -1 example.http
```

Make the second request of `example.http` and output verbosely

```bash
climan -verbose -index 2 example.http
```

## 🔖 Todos

- [x] Support for other request body types (like `text/xml` and `text/plain`).
- [x] Parse the response body base on `Content-Type` header.
- [x] Show response details.
- [x] Colored output.
- [ ] Colored response body.
- [ ] Environment variable support.
- [ ] Add more `http` examples.

## 📢 Credits

- [REST Client](https://github.com/Huachao/vscode-restclient)
