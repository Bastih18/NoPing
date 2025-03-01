# 🛠️ Local Development Guide for NoPing

This guide explains how to set up NoPing for **local development**.

---

## 📥 Clone the Repository
First, clone the NoPing repository to your local machine:

```sh
git clone https://github.com/Bastih18/NoPing.git
cd NoPing
```

---

## 🔧 Enable Local Module Import

Inside `go.mod`, you will see this line:

```go
// replace github.com/Bastih18/NoPing => ./
```

Uncomment it by removing `//` so it becomes:

```go
replace github.com/Bastih18/NoPing => ./
```

This ensures Go uses the local NoPing modules instead of fetching it from GitHub.

After modifying `go.mod`, run:
```sh
go mod tidy
```

---

## 🚀 Building and Running NoPing

### **Build the Executable**
```sh
go build -o noping .
```

### **Run NoPing**
```sh
./noping
```

---

## 🔄 Updating Dependencies

If you modify dependencies in `go.mod`, update them using:
```sh
go mod tidy
```

If dependencies behave unexpectedly, clear the Go module cache:
```sh
go clean -modcache
```

---

## ✅ Summary

| Task | Command |
|------|---------|
| Clone NoPing | `git clone https://github.com/Bastih18/NoPing.git` |
| Enable Local Imports | Uncomment `replace github.com/Bastih18/NoPing => ./` in `go.mod` |
| Update Dependencies | `go mod tidy` |
| Build NoPing | `go build -o noping .` |
| Run NoPing | `./noping` |
| Clear Module Cache | `go clean -modcache` |

---

Happy coding! 🚀
