


# API Sekretaris Desa 📄

API ini digunakan untuk pengelolaan persuratan di sebuah desa. Dibangun menggunakan **Golang** dan **MySQL** sebagai database untuk menyimpan data surat masuk dan surat keluar. API ini dirancang untuk memudahkan pengelolaan surat-menyurat di tingkat desa.

---

## 🚀 Cara Menjalankan API

### 1. **Persiapan**
- Pastikan Anda telah menginstal:
  - **Golang** (versi `go1.24.1` atau lebih baru)
  - **MySQL** (database sudah dikonfigurasi)
- Clone repository ini ke lokal:
  ```bash
  git clone https://github.com/Ahmadfaisal04/Backend-sekertaris.git
  ```
- Install dependensi:
  ```bash
  go mod tidy
  ```

### 2. **Tools yang Dibutuhkan**
- **Postman** (untuk testing API)
- **MySQL Client** (untuk mengelola database)


## 🛠️ Struktur Proyek

Berikut adalah struktur proyek API Sekretaris Desa:

```
my-go-project/
├── config/
│   └── database.go           # Konfigurasi koneksi database
├── controller/
│   └── suratkeluar_controller.go  # Controller untuk surat keluar
│   └── suratmasuk_controller.go   # Controller untuk surat masuk
├── models/
│   └── suratkeluar.go        # Model data surat keluar
│   └── suratmasuk.go         # Model data surat masuk
├── repository/
│   └── suratkeluar_repository.go  # Repository untuk surat keluar
│   └── suratmasuk_repository.go   # Repository untuk surat masuk
├── static/
│   └── suratkeluar/          # Folder penyimpanan file surat keluar
│   └── suratmasuk/           # Folder penyimpanan file surat masuk
├── .env                      # File environment variable
├── .gitignore                # File ignore Git
├── go.mod                    # File dependensi Go
├── go.sum                    # File checksum dependensi Go
├── main.go                   # Entry point aplikasi
└── README.md                 # Dokumentasi proyek
```

---

## 🏃‍♂️ Menjalankan Server

1. Masuk ke direktori proyek:
   ```bash
   cd Backend-sekertaris
   ```
2. Jalankan server:
   ```bash
   go run main.go
   ```
3. Server akan berjalan di **http://localhost:8088**.

---

## 📚 Fitur API

### **Surat Masuk**
| **Endpoint**                          | **Method** | **Deskripsi**                               |
|---------------------------------------|------------|---------------------------------------------|
| `/api/suratmasuk`                     | POST       | Menambahkan data surat masuk                |
| `/api/suratmasuk/update/:id`          | PUT        | Memperbarui data surat masuk berdasarkan ID |
| `/api/suratmasuk/get`                 | GET        | Mengambil semua data surat masuk           |
| `/api/suratmasuk/get/:id`             | GET        | Mengambil data surat masuk berdasarkan ID  |
| `/api/suratmasuk/count`               | GET        | Menghitung jumlah surat masuk              |

### **Surat Keluar**
| **Endpoint**                          | **Method** | **Deskripsi**                               |
|---------------------------------------|------------|---------------------------------------------|
| `/api/suratkeluar`                    | POST       | Menambahkan data surat keluar               |
| `/api/suratkeluar/:id`                | PUT        | Memperbarui data surat keluar berdasarkan ID|
| `/api/suratkeluar`                    | GET        | Mengambil semua data surat keluar           |
| `/api/suratkeluar/get/:id`            | GET        | Mengambil data surat keluar berdasarkan ID  |
| `/api/suratkeluar/count`              | GET        | Menghitung jumlah surat keluar              |

---

## 🛠️ Contoh Penggunaan API

### **Menambahkan Surat Masuk**
- **URL**: `http://localhost:8088/api/suratmasuk`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "nomor_surat": "001/SM/2023",
    "tanggal_masuk": "2023-10-01",
    "pengirim": "Pemerintah Kabupaten",
    "perihal": "Undangan Rapat"
  }
  ```

### **Mengambil Semua Surat Keluar**
- **URL**: `http://localhost:8088/api/suratkeluar`
- **Method**: `GET`
- **Response** (JSON):
  ```json
  [
    {
      "id": 1,
      "nomor_surat": "001/SK/2023",
      "tanggal_keluar": "2023-10-02",
      "tujuan": "Pemerintah Kecamatan",
      "perihal": "Laporan Kegiatan"
    }
  ]
  ```

---

## 📝 Catatan
- Pastikan file `.env` sudah dikonfigurasi dengan benar untuk koneksi database.
- Gunakan **Postman** atau tools sejenis untuk testing API.
- Folder `static/` digunakan untuk menyimpan file surat (jika ada).

---

## 🤝 Kontribusi
Jika Anda ingin berkontribusi pada proyek ini, silakan buka **Pull Request** atau laporkan masalah melalui **Issues**.

---

Dokumentasi ini dirancang untuk memudahkan pengembang dalam memahami dan menggunakan API Sekretaris Desa. Selamat mencoba! 🚀

