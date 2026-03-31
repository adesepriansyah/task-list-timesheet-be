# Issue: Implementasi CRUD API untuk Tasks

## Deskripsi

Membuat modul Tasks lengkap dengan endpoint CRUD (Create, Read, Update, Delete) untuk mengelola data tugas (task). Modul ini mengikuti pola Clean Architecture yang sudah ada di proyek (lihat modul `auth` sebagai referensi).

---

## Tabel `tasks` (Sudah Ada)

Tabel ini sudah dibuat di migration `0001_initial.up.sql`:

| Kolom       | Tipe         | Keterangan                                |
|-------------|--------------|-------------------------------------------|
| id          | SERIAL       | Primary Key, auto increment               |
| user_id     | INT          | Foreign key ke tabel `users`              |
| title       | VARCHAR(255) | Judul task                                |
| description | TEXT         | Deskripsi task                            |
| status      | VARCHAR(50)  | `pending`, `in_progress`, `completed`     |
| date        | DATE         | Tanggal task                              |
| effort_time | INT          | Estimasi waktu (dalam menit)              |
| created_at  | TIMESTAMP    | Waktu dibuat                              |
| updated_at  | TIMESTAMP    | Waktu terakhir diperbarui                 |

> **Note**: Entity `Task` sudah ada di `internal/entity/task.go`.

---

## Endpoint

### 1. `POST /api/tasks` — Create Task

**Request Body:**
```json
{
  "title": "Belajar Go",
  "description": "Belajar dasar-dasar Go",
  "status": "pending",
  "user_id": 1,
  "date": "2026-04-01",
  "effort_time": 60
}
```

**Response Success (201):**
```json
{
  "data": "Ok"
}
```

**Response Error (400):**
```json
{
  "error": "Bad Request"
}
```

---

### 2. `GET /api/tasks` — Get Task List (dengan Filter)

**Query Parameters (filter):**
| Parameter   | Tipe   | Wajib? | Keterangan                              |
|-------------|--------|--------|-----------------------------------------|
| user_id     | int    | ✅ Ya  | Filter berdasarkan user                 |
| search      | string | Tidak  | Pencarian di kolom `title` & `description` |
| date_from   | string | Tidak  | Batas awal range tanggal (format: `YYYY-MM-DD`) |
| date_to     | string | Tidak  | Batas akhir range tanggal (format: `YYYY-MM-DD`) |

**Contoh Request:**
```
GET /api/tasks?user_id=1&search=belajar&date_from=2026-04-01&date_to=2026-04-30
```

**Response Success (200):**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "Belajar Go",
      "description": "Belajar dasar-dasar Go",
      "status": "pending",
      "date": "2026-04-01",
      "effort_time": 60,
      "created_at": "2026-04-01T00:00:00Z",
      "updated_at": "2026-04-01T00:00:00Z"
    }
  ]
}
```

---

### 3. `PUT /api/tasks/{id}` — Update Task

**Request Body (sama seperti create):**
```json
{
  "title": "Belajar Go Lanjutan",
  "description": "Belajar Go dengan project",
  "status": "in_progress",
  "user_id": 1,
  "date": "2026-04-02",
  "effort_time": 120
}
```

**Response Success (200):**
```json
{
  "data": "Ok"
}
```

---

### 4. `DELETE /api/tasks/{id}` — Delete Task

**Response Success (200):**
```json
{
  "data": "Ok"
}
```

---

## Tahapan Implementasi

### Tahap 1: Buat Folder & File Baru untuk Modul Task

Buat struktur folder baru:
```
internal/
  task/
    api.go        ← handler HTTP (router + handler functions)
    service.go    ← interface Service + implementasi logic
    request.go    ← struct request body
```

> **Referensi**: Lihat folder `internal/auth/` sebagai contoh pola yang harus diikuti.

---

### Tahap 2: Buat Repository

Edit file `internal/repository/task.go` (file baru):

1. **Buat interface `TaskRepository`** dengan method:
   ```go
   type TaskRepository interface {
       Create(ctx context.Context, task *entity.Task) error
       FindByID(ctx context.Context, id int) (*entity.Task, error)
       FindAll(ctx context.Context, filter TaskFilter) ([]entity.Task, error)
       Update(ctx context.Context, task *entity.Task) error
       Delete(ctx context.Context, id int) error
   }
   ```

2. **Buat struct `TaskFilter`** untuk filter pencarian:
   ```go
   type TaskFilter struct {
       UserID   int
       Search   string    // pencarian di title & description (ILIKE '%search%')
       DateFrom string    // format: YYYY-MM-DD
       DateTo   string    // format: YYYY-MM-DD
   }
   ```

3. **Implementasi setiap method:**
   - `Create`: `INSERT INTO tasks (user_id, title, description, status, date, effort_time) VALUES ($1, $2, $3, $4, $5, $6)`
   - `FindByID`: `SELECT * FROM tasks WHERE id = $1`
   - `FindAll`: `SELECT * FROM tasks WHERE user_id = $1` + tambahkan kondisi filter secara dinamis:
     - Jika `Search` tidak kosong: `AND (title ILIKE '%search%' OR description ILIKE '%search%')`
     - Jika `DateFrom` tidak kosong: `AND date >= $N`
     - Jika `DateTo` tidak kosong: `AND date <= $N`
   - `Update`: `UPDATE tasks SET title=$1, description=$2, status=$3, date=$4, effort_time=$5, updated_at=NOW() WHERE id=$6`
   - `Delete`: `DELETE FROM tasks WHERE id = $1`

---

### Tahap 3: Buat Service

Buat file `internal/task/service.go`:

1. **Buat interface `Service`:**
   ```go
   type Service interface {
       CreateTask(ctx context.Context, req CreateTaskRequest) error
       GetTasks(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, error)
       UpdateTask(ctx context.Context, id int, req UpdateTaskRequest) error
       DeleteTask(ctx context.Context, id int) error
   }
   ```

2. **Implementasi:**
   - `CreateTask`: Validasi input, lalu panggil `repo.Create()`.
   - `GetTasks`: Panggil `repo.FindAll()` dengan filter dari query parameter.
   - `UpdateTask`: Cek apakah task ada via `repo.FindByID()`, lalu panggil `repo.Update()`.
   - `DeleteTask`: Panggil `repo.Delete()`.

---

### Tahap 4: Buat Request Struct

Buat file `internal/task/request.go`:
```go
type CreateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    UserID      int    `json:"user_id"`
    Date        string `json:"date"`        // format: YYYY-MM-DD
    EffortTime  int    `json:"effort_time"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    UserID      int    `json:"user_id"`
    Date        string `json:"date"`
    EffortTime  int    `json:"effort_time"`
}
```

---

### Tahap 5: Buat API Handler

Buat file `internal/task/api.go`:

1. **Daftarkan route di `RegisterHandlers`:**
   ```go
   func RegisterHandlers(r chi.Router, service Service) {
       res := &resource{service}
       r.Route("/api/tasks", func(r chi.Router) {
           r.Post("/", res.create)
           r.Get("/", res.list)
           r.Put("/{id}", res.update)
           r.Delete("/{id}", res.delete)
       })
   }
   ```

2. **Implementasi handler:**
   - `create`: Decode JSON body → panggil `service.CreateTask()` → return 201.
   - `list`: Ambil query params (`user_id`, `search`, `date_from`, `date_to`) → panggil `service.GetTasks()` → return 200 dengan list.
   - `update`: Ambil `{id}` dari URL + JSON body → panggil `service.UpdateTask()` → return 200.
   - `delete`: Ambil `{id}` dari URL → panggil `service.DeleteTask()` → return 200.

   > Untuk mengambil `{id}` dari URL, gunakan: `chi.URLParam(r, "id")` lalu konversi ke `int` dengan `strconv.Atoi()`.

---

### Tahap 6: Wiring di `main.go`

Edit `cmd/server/main.go`:

1. Import paket `task`:
   ```go
   "github.com/adesepriansyah/task-list-timesheet-be/internal/task"
   ```

2. Inisialisasi layer:
   ```go
   taskRepo := repository.NewTaskRepository(db)
   taskService := task.NewService(taskRepo)
   ```

3. Daftarkan handler:
   ```go
   task.RegisterHandlers(r, taskService)
   ```

---

### Tahap 7: Testing

1. **Unit Test**: Buat `internal/task/service_test.go` dengan mock repository.
2. **Manual Test**: Gunakan `curl` untuk menguji semua endpoint:
   - `POST /api/tasks` → buat task baru.
   - `GET /api/tasks?user_id=1` → list task berdasarkan user.
   - `GET /api/tasks?user_id=1&search=belajar` → pencarian berdasarkan judul/deskripsi.
   - `PUT /api/tasks/1` → update task.
   - `DELETE /api/tasks/1` → hapus task.

---

## Catatan Penting

- **Semua endpoint menggunakan JSON** untuk request dan response body.
- **`user_id` wajib** pada endpoint `GET /api/tasks` (sebagai query parameter).
- **Status hanya boleh**: `pending`, `in_progress`, `completed` (validasi di service layer).
- **Filter pencarian** menggunakan `ILIKE` untuk case-insensitive search.
- **Ikuti pola Clean Architecture** yang sudah ada di modul `auth` (lihat file-file di `internal/auth/` dan `internal/repository/`).
