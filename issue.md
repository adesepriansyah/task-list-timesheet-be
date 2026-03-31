# Issue: Implementasi API Register & Get User Info

## Deskripsi

Menambahkan dua endpoint baru pada modul `auth`:
1. **Register** — Membuat user baru.
2. **Get User Info** — Mengambil informasi user yang sedang login berdasarkan token.

Fitur ini juga memerlukan penambahan kolom `name` pada tabel `users`.

---

## Endpoint

### 1. `POST /api/users/register`

**Request Body:**
```json
{
  "name": "name",
  "email": "email",
  "password": "password"
}
```

**Response Success (201):**
```json
{
  "data": "Ok"
}
```

**Response Error (400 - email sudah terdaftar):**
```json
{
  "error": "Email already registered"
}
```

---

### 2. `GET /api/users/info`

**Headers:**
```
Authorization: Bearer <token>
```

**Response Success (200):**
```json
{
  "data": {
    "id": 1,
    "name": "name",
    "email": "email",
    "expired_token": "2026-04-01T00:00:00Z"
  }
}
```

**Response Error (401 - token tidak valid / expired):**
```json
{
  "error": "Unauthorized"
}
```

---

## Tahapan Implementasi

### Tahap 1: Database Migration

Buat file migration baru `migrations/0002_add_name_to_users.up.sql`:
```sql
ALTER TABLE users ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';
```

Jalankan migration ke database:
```bash
cat migrations/0002_add_name_to_users.up.sql | docker exec -i task-list-timesheet-be-db-1 psql -U postgres -d taskdb
```

---

### Tahap 2: Update Entity

Edit file `internal/entity/user.go`, tambahkan field `Name`:
```go
type User struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`        // <-- TAMBAHKAN INI
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

---

### Tahap 3: Update Repository

Edit file `internal/repository/auth.go`:

1. **Tambahkan method baru ke interface `AuthRepository`:**
   ```go
   CreateUser(ctx context.Context, name, email, hashedPassword string) error
   FindUserByID(ctx context.Context, id int) (*entity.User, error)
   ```

2. **Implementasikan kedua method tersebut:**
   - `CreateUser`: Jalankan query `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`.
   - `FindUserByID`: Jalankan query `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`.

3. **Update query `FindUserByEmail`** agar juga mengambil kolom `name`:
   ```sql
   SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1
   ```
   dan tambahkan `&user.Name` pada `.Scan(...)`.

---

### Tahap 4: Update Service

Edit file `internal/auth/service.go`:

1. **Tambahkan method baru ke interface `Service`:**
   ```go
   Register(ctx context.Context, name, email, password string) error
   GetUserInfo(ctx context.Context, token string) (*entity.User, time.Time, error)
   ```

2. **Implementasi `Register`:**
   - Cek apakah email sudah terdaftar via `repo.FindUserByEmail()`.
   - Jika sudah ada, return error `"email already registered"`.
   - Hash password menggunakan `bcrypt.GenerateFromPassword()`.
   - Simpan user baru via `repo.CreateUser()`.

3. **Implementasi `GetUserInfo`:**
   - Cari token di database via `repo.FindToken()`.
   - Jika token tidak ditemukan atau sudah expired, return error `"unauthorized"`.
   - Ambil data user via `repo.FindUserByID()` menggunakan `userToken.UserID`.
   - Return data user dan `expiredAt` dari token.

---

### Tahap 5: Update Request Struct

Edit file `internal/auth/request.go`, tambahkan struct baru:
```go
type registerRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

---

### Tahap 6: Update API Handler

Edit file `internal/auth/api.go`:

1. **Tambahkan route baru di `RegisterHandlers`:**
   ```go
   r.Post("/register", res.register)
   r.Get("/info", res.getUserInfo)
   ```

2. **Buat handler `register`:**
   - Decode request body ke `registerRequest`.
   - Panggil `service.Register(ctx, req.Name, req.Email, req.Password)`.
   - Jika error `"email already registered"`, return 400.
   - Jika sukses, return `{"data": "Ok"}`.

3. **Buat handler `getUserInfo`:**
   - Ambil token dari header `Authorization: Bearer <token>`.
   - Panggil `service.GetUserInfo(ctx, token)`.
   - Jika error `"unauthorized"`, return 401.
   - Jika sukses, return:
     ```json
     {
       "data": {
         "id": 1,
         "name": "name",
         "email": "email",
         "expired_token": "2026-04-01T00:00:00Z"
       }
     }
     ```

---

### Tahap 7: Testing

1. **Unit Test**: Tambahkan test case untuk `Register` dan `GetUserInfo` di `internal/auth/service_test.go`.
2. **Manual Test**: Gunakan `curl` untuk menguji endpoint baru (lihat contoh di `api_testing.md`).

---

## Catatan Penting

- Password harus di-hash menggunakan `bcrypt` sebelum disimpan ke database.
- Token diambil dari header `Authorization: Bearer <token>`.
- Gunakan `time.Now()` untuk mengecek apakah token sudah expired.
- Pastikan semua query yang mengambil data `users` sudah menyertakan kolom `name`.
