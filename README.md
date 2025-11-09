# BE News-Feed

Tech stack yang digunakan:

- Go 1.24
- MySQL 8.4.3

Autentikasi menggunakan JWT `(token dan refresh_token)` disimpan pada `HTTP-Cookies`. Sehingga pada client akan otomatis menyertakan `cookies`.

## Menjalankan Program

1. Jalankan migrasi database

    Pastikan telah mengunduh golang migrate di OS <github.com/golang-migrate/migrate>

    ```bash
    migrate -database DB URL: "mysql://[DB_USER]:[DB_PASSWORD]@tcp([DB_HOST]:[DB_PORT])/[DB_NAME]" -path db/migrations up
    ```

2. Instal modul `go mod` projek

3. Jalankan file `main.go`

4. Import postman (Opsional)

    Import collection postman dari file `postman.json`
