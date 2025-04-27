package repository

import (
	"Sekertaris/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PermohonanSuratRepository struct {
	db *sql.DB
}

func NewPermohonanSuratRepository(db *sql.DB) *PermohonanSuratRepository {
	return &PermohonanSuratRepository{db: db}
}

func (r *PermohonanSuratRepository) AddPermohonanSurat(permohonan model.PermohonanSurat) (*model.PermohonanSurat, error) {
	query := `
        INSERT INTO permohonansurat (
            nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin, 
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan, 
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url, 
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah, 
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.NamaUsaha, permohonan.JenisUsaha, permohonan.AlamatUsaha, permohonan.AlamatTujuan,
		permohonan.AlasanPindah, permohonan.NamaAyah, permohonan.NamaIbu, permohonan.TglKematian,
		permohonan.PenyebabKematian, permohonan.Status,
	)
	if err != nil {
		log.Println("Error adding permohonan surat:", err)
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert ID:", err)
		return nil, err
	}

	var newPermohonan model.PermohonanSurat
	query = `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, status,
            created_at, updated_at
        FROM permohonansurat WHERE id = ?
    `
	err = r.db.QueryRow(query, lastInsertID).Scan(
		&newPermohonan.ID, &newPermohonan.NIK, &newPermohonan.NamaLengkap, &newPermohonan.TempatLahir,
		&newPermohonan.TanggalLahir, &newPermohonan.JenisKelamin, &newPermohonan.Pendidikan,
		&newPermohonan.Pekerjaan, &newPermohonan.Agama, &newPermohonan.StatusPernikahan,
		&newPermohonan.Kewarganegaraan, &newPermohonan.AlamatLengkap, &newPermohonan.JenisSurat,
		&newPermohonan.Keterangan, &newPermohonan.NomorHP, &newPermohonan.DokumenURL,
		&newPermohonan.NamaUsaha, &newPermohonan.JenisUsaha, &newPermohonan.AlamatUsaha,
		&newPermohonan.AlamatTujuan, &newPermohonan.AlasanPindah, &newPermohonan.NamaAyah,
		&newPermohonan.NamaIbu, &newPermohonan.TglKematian, &newPermohonan.PenyebabKematian,
		&newPermohonan.Status, &newPermohonan.CreatedAt, &newPermohonan.UpdatedAt,
	)
	if err != nil {
		log.Println("Error retrieving new permohonan surat:", err)
		return nil, err
	}

	return &newPermohonan, nil
}

func (r *PermohonanSuratRepository) GetPermohonanSurat() ([]model.PermohonanSurat, error) {
	query := `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, status,
            created_at, updated_at
        FROM permohonansurat
        ORDER BY created_at ASC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error retrieving permohonan surat:", err)
		return nil, err
	}
	defer rows.Close()

	var permohonanSuratList []model.PermohonanSurat
	for rows.Next() {
		var permohonan model.PermohonanSurat
		err := rows.Scan(
			&permohonan.ID, &permohonan.NIK, &permohonan.NamaLengkap, &permohonan.TempatLahir,
			&permohonan.TanggalLahir, &permohonan.JenisKelamin, &permohonan.Pendidikan,
			&permohonan.Pekerjaan, &permohonan.Agama, &permohonan.StatusPernikahan,
			&permohonan.Kewarganegaraan, &permohonan.AlamatLengkap, &permohonan.JenisSurat,
			&permohonan.Keterangan, &permohonan.NomorHP, &permohonan.DokumenURL,
			&permohonan.NamaUsaha, &permohonan.JenisUsaha, &permohonan.AlamatUsaha,
			&permohonan.AlamatTujuan, &permohonan.AlasanPindah, &permohonan.NamaAyah,
			&permohonan.NamaIbu, &permohonan.TglKematian, &permohonan.PenyebabKematian,
			&permohonan.Status, &permohonan.CreatedAt, &permohonan.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning permohonan surat row:", err)
			return nil, err
		}
		permohonanSuratList = append(permohonanSuratList, permohonan)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after retrieving permohonan surat:", err)
		return nil, err
	}

	return permohonanSuratList, nil
}

func (r *PermohonanSuratRepository) GetPermohonanSuratByID(id int64) (*model.PermohonanSurat, error) {
	var permohonan model.PermohonanSurat
	query := `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, status,
            created_at, updated_at
        FROM permohonansurat WHERE id = ?
    `
	err := r.db.QueryRow(query, id).Scan(
		&permohonan.ID, &permohonan.NIK, &permohonan.NamaLengkap, &permohonan.TempatLahir,
		&permohonan.TanggalLahir, &permohonan.JenisKelamin, &permohonan.Pendidikan,
		&permohonan.Pekerjaan, &permohonan.Agama, &permohonan.StatusPernikahan,
		&permohonan.Kewarganegaraan, &permohonan.AlamatLengkap, &permohonan.JenisSurat,
		&permohonan.Keterangan, &permohonan.NomorHP, &permohonan.DokumenURL,
		&permohonan.NamaUsaha, &permohonan.JenisUsaha, &permohonan.AlamatUsaha,
		&permohonan.AlamatTujuan, &permohonan.AlasanPindah, &permohonan.NamaAyah,
		&permohonan.NamaIbu, &permohonan.TglKematian, &permohonan.PenyebabKematian,
		&permohonan.Status, &permohonan.CreatedAt, &permohonan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permohonan surat dengan ID %d tidak ditemukan", id)
		}
		log.Printf("Error retrieving permohonan surat by ID %d: %v", id, err)
		return nil, fmt.Errorf("gagal mengambil permohonan surat: %v", err)
	}
	return &permohonan, nil
}

func (r *PermohonanSuratRepository) GetOldestPendingPermohonan() (*model.PermohonanSurat, error) {
	var permohonan model.PermohonanSurat
	query := `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, status,
            created_at, updated_at
        FROM permohonansurat
        WHERE status = 'Diproses'
        ORDER BY created_at ASC
        LIMIT 1
    `
	err := r.db.QueryRow(query).Scan(
		&permohonan.ID, &permohonan.NIK, &permohonan.NamaLengkap, &permohonan.TempatLahir,
		&permohonan.TanggalLahir, &permohonan.JenisKelamin, &permohonan.Pendidikan,
		&permohonan.Pekerjaan, &permohonan.Agama, &permohonan.StatusPernikahan,
		&permohonan.Kewarganegaraan, &permohonan.AlamatLengkap, &permohonan.JenisSurat,
		&permohonan.Keterangan, &permohonan.NomorHP, &permohonan.DokumenURL,
		&permohonan.NamaUsaha, &permohonan.JenisUsaha, &permohonan.AlamatUsaha,
		&permohonan.AlamatTujuan, &permohonan.AlasanPindah, &permohonan.NamaAyah,
		&permohonan.NamaIbu, &permohonan.TglKematian, &permohonan.PenyebabKematian,
		&permohonan.Status, &permohonan.CreatedAt, &permohonan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Tidak ada permohonan dengan status Diproses
		}
		log.Println("Error retrieving oldest pending permohonan surat:", err)
		return nil, fmt.Errorf("gagal mengambil permohonan surat tertua: %v", err)
	}
	return &permohonan, nil
}

func (r *PermohonanSuratRepository) UpdatePermohonanSuratByID(id int64, permohonan model.PermohonanSurat) error {
	query := `
        UPDATE permohonansurat 
        SET nik = ?, nama_lengkap = ?, tempat_lahir = ?, tanggal_lahir = ?, jenis_kelamin = ?,
            pendidikan = ?, pekerjaan = ?, agama = ?, status_pernikahan = ?, kewarganegaraan = ?,
            alamat_lengkap = ?, jenis_surat = ?, keterangan = ?, nomor_hp = ?, dokumen_url = ?,
            nama_usaha = ?, jenis_usaha = ?, alamat_usaha = ?, alamat_tujuan = ?, alasan_pindah = ?,
            nama_ayah = ?, nama_ibu = ?, tgl_kematian = ?, penyebab_kematian = ?, ditujukan = ?, status = ?
        WHERE id = ?
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.NamaUsaha, permohonan.JenisUsaha, permohonan.AlamatUsaha, permohonan.AlamatTujuan,
		permohonan.AlasanPindah, permohonan.NamaAyah, permohonan.NamaIbu, permohonan.TglKematian,
		permohonan.PenyebabKematian,permohonan.Ditujukan, permohonan.Status, id,
	)
	if err != nil {
		log.Println("Error updating permohonan surat:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada permohonan surat dengan ID %d yang ditemukan", id)
	}

	return nil
}

func (r *PermohonanSuratRepository) UpdateStatusByID(id int64, status model.Status, updatedAt time.Time) error {
	log.Printf("Updating status for ID %d to %s at %s", id, status, updatedAt.Format(time.RFC3339))
	query := `
        UPDATE permohonansurat 
        SET status = ?, updated_at = ?
        WHERE id = ?
    `
	result, err := r.db.Exec(query, status, updatedAt, id)
	if err != nil {
		log.Println("Error updating permohonan surat status:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	log.Printf("Rows affected: %d", rowsAffected)
	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada permohonan surat dengan ID %d yang ditemukan", id)
	}

	return nil
}

func (r *PermohonanSuratRepository) DeletePermohonanSurat(id int64) error {
	query := "DELETE FROM permohonansurat WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting permohonan surat:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada permohonan surat dengan ID %d yang ditemukan", id)
	}

	return nil
}