# BE News-Feed

Tech stack yang digunakan:

- Go 1.24
- MySQL 8.4.3

Autentikasi menggunakan JWT `(token dan refresh_token)` disimpan pada `HTTP-Cookies`. Sehingga pada client akan otomatis menyertakan `cookies`.

## Menjalankan Program

1. Jalankan migrasi database

    Pastikan telah mengunduh golang migrate di OS <github.com/golang-migrate/migrate>

    ```bash
    migrate -database "mysql://root:@tcp(localhost:3306)/gofiber_restapi" -path db/migrations up
    ```

2. Instal modul `go mod` projek

3. Jalankan file `main.go`

4. Import postman (Opsional)

    Import collection postman dari file `postman.json`

## Update Repository-Service Pattern

Dalam update ini menggunakan `google-wire` untuk melakukan dependency injection `service-repository-pattern`, sehingga:

- bisa akses `repo` yang telah didefinisikan `service` dalam service.
- `repo` diinject oleh `DB` sehingga tidak perlu chaining dari `controller-service-repo`
- ada pemisah layer:

  - `I/O (controller)`
  - `validasi db (unique, exists, dll)`, `dto` dan `proses diluar db` di `service`
  - transaksi database `repository`

### Transaksi

Transaksi dibutuhkan jika ada lebih dari satu proses insert/delete/update db.
Beberapa gaya yang bisa disesuaikan jika dibutuhkan:

- jadikan config dependant dari `router/repo/service/controller`
- jika ada endpoint yang butuh transaksi, tinggal gunakan di layer repository, tapi bisa redundant (misal sudah bikin query di fungsi repo lain).
