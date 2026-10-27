# Folder `cmd`

Folder ini sengaja dikosongkan.

Titik masuk utama (entry point) aplikasi ini berada di **root project**:

```
d:\go_26_10\main.go
```

Jalankan aplikasi dari folder root:

```powershell
go run .
```

atau

```powershell
go build
```

## Tujuan folder `cmd/`

Folder `cmd/` disediakan untuk menyimpan executable tambahan di masa depan, misalnya:

- `cmd/server/main.go` — HTTP server
- `cmd/worker/main.go` — background job
- `cmd/cli/main.go` — command-line tool

Semua tetap dapat menggunakan logika inti yang ada di `internal/`, `domain/`, dan `config/`.

## Aturan

- Jangan menaruh logika bisnis di dalam `cmd/`. Folder ini hanya untuk **wiring** atau pemanggilan komponen.
- Satu folder `cmd/<nama>/` sebaiknya hanya berisi satu `package main` sebagai titik masuk tunggal.
