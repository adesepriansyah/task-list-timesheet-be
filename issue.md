# Implementasi Database Schema & Auth API (Login/Logout)

## Deskripsi
Implementasi skema database untuk tabel `user`, `user_token`, dan `task`, serta API untuk fitur **Login** dan **Logout** user.

Dokumen ini adalah panduan step-by-step. Ikuti setiap tahapan secara berurutan.

---

## 1. Database Schema (Migration)

Buat file migration di folder `migrations/`. Buat tiga tabel berikut:

### Tabel `users`
| Kolom        | Tipe        | Keterangan          |
|-------------|-------------|---------------------|
| id          | SERIAL (PK) | Auto increment      |
| email       | VARCHAR     | Unique, not null    |
| password    | VARCHAR     | Hashed, not null    |
| created_at  | TIMESTAMP   | Default NOW()       |
| updated_at  | TIMESTAMP   | Default NOW()       |

### Tabel `user_tokens`
| Kolom        | Tipe        | Keterangan                       |
|-------------|-------------|----------------------------------|
| id          | SERIAL (PK) | Auto increment                   |
| token       | VARCHAR     | Unique, not null                 |
| user_id     | INT (FK)    | References `users(id)`           |
| expired_at  | TIMESTAMP   | Waktu kadaluarsa token           |
| created_at  | TIMESTAMP   | Default NOW()                    |
| updated_at  | TIMESTAMP   | Default NOW()                    |

### Tabel `tasks`
| Kolom        | Tipe        | Keterangan                                      |
|-------------|-------------|--------------------------------------------------|
| id          | SERIAL (PK) | Auto increment                                   |
| user_id     | INT (FK)    | References `users(id)`                           |
| title       | VARCHAR     | Not null                                         |
| description | TEXT        | Nullable                                         |
| status      | VARCHAR     | Enum: `'pending'`, `'in_progress'`, `'completed'` |
| date        | DATE        | Tanggal task                                     |
| effort_time | INT         | Dalam menit                                      |
| created_at  | TIMESTAMP   | Default NOW()                                    |
| updated_at  | TIMESTAMP   | Default NOW()                                    |

---

## 2. Entity (Model)

Buat struct Go untuk setiap tabel di folder `internal/entity/`.

#### File yang harus dibuat:
- `internal/entity/user.go` → struct `User` (mapping ke tabel `users`)
- `internal/entity/user_token.go` → struct `UserToken` (mapping ke tabel `user_tokens`)
- `internal/entity/task.go` → struct `Task` (mapping ke tabel `tasks`)

---

## 3. Repository Layer

Buat file repository di masing-masing folder fitur. Repository berisi fungsi-fungsi query ke database.

#### File yang harus dibuat:
- `internal/auth/repository.go` → berisi query terkait auth:
  - `FindUserByEmail(email string)` → mencari user berdasarkan email
  - `CreateToken(userID int, token string, expiredAt time.Time)` → menyimpan token baru ke tabel `user_tokens`
  - `DeleteToken(token string)` → menghapus token dari tabel `user_tokens` (untuk logout)
  - `FindToken(token string)` → mencari token yang masih valid (belum expired)

---

## 4. Service Layer (Business Logic)

Buat file service di folder `internal/auth/`.

#### File yang harus dibuat:
- `internal/auth/service.go` → berisi logika bisnis:
  - **Login**: Terima `email` dan `password` → cari user di DB → bandingkan password (gunakan bcrypt) → jika cocok, generate token random → simpan ke tabel `user_tokens` → return token
  - **Logout**: Terima `token` dari header `Authorization` → hapus token dari tabel `user_tokens`

---

## 5. API / Router Layer

Buat file router/handler di folder `internal/auth/`.

#### File yang harus dibuat:
- `internal/auth/api.go` → berisi HTTP handler dan route registration

### Endpoint Login

```
POST /api/users/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "token": "generated-random-token"
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Unauthorized"
}
```

### Endpoint Logout

```
POST /api/users/logout
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "data": "success"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Unauthorized"
}
```

> **Catatan:** Saat logout berhasil, token yang digunakan harus **dihapus dari tabel `user_tokens`** sehingga tidak bisa digunakan lagi.

---

## 6. Wiring (Dependency Injection)

Hubungkan semua layer di `cmd/server/main.go`:

1. Buat koneksi database PostgreSQL menggunakan DSN dari config
2. Inisialisasi `auth.Repository` dengan koneksi DB
3. Inisialisasi `auth.Service` dengan repository
4. Register `auth.RegisterHandlers(router)` ke Chi router

---

## 7. Tahapan Implementasi (Urutan Kerja)

Ikuti urutan ini agar tidak ada dependency yang terlewat:

1. **Buat migration SQL** → jalankan di database untuk membuat tabel
2. **Buat entity** (`user.go`, `user_token.go`, `task.go`)
3. **Buat koneksi database** di `cmd/server/main.go` (gunakan `database/sql` + `lib/pq` atau `pgx`)
4. **Buat `internal/auth/repository.go`** → implementasi query ke DB
5. **Buat `internal/auth/service.go`** → implementasi logika login & logout
6. **Buat `internal/auth/api.go`** → implementasi handler dan daftarkan route di chi
7. **Update `cmd/server/main.go`** → wiring repository → service → handler
8. **Test manual** menggunakan `curl` atau Postman:
   - Login dengan email & password → pastikan dapat token
   - Logout dengan token → pastikan token terhapus

## Kriteria Penerimaan (Acceptance Criteria)
- [ ] Tabel `users`, `user_tokens`, `tasks` berhasil dibuat di database
- [ ] `POST /api/users/login` mengembalikan token jika email & password benar
- [ ] `POST /api/users/login` mengembalikan error `401` jika email/password salah
- [ ] `POST /api/users/logout` menghapus token dari database
- [ ] Password disimpan dalam bentuk hash (bcrypt), bukan plain text
