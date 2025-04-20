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
            status, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.Status, permohonan.CreatedAt, permohonan.UpdatedAt,
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
            status, created_at, updated_at
        FROM permohonansurat WHERE id = ?
    `
	err = r.db.QueryRow(query, lastInsertID).Scan(
		&newPermohonan.ID, &newPermohonan.NIK, &newPermohonan.NamaLengkap, &newPermohonan.TempatLahir,
		&newPermohonan.TanggalLahir, &newPermohonan.JenisKelamin, &newPermohonan.Pendidikan,
		&newPermohonan.Pekerjaan, &newPermohonan.Agama, &newPermohonan.StatusPernikahan,
		&newPermohonan.Kewarganegaraan, &newPermohonan.AlamatLengkap, &newPermohonan.JenisSurat,
		&newPermohonan.Keterangan, &newPermohonan.NomorHP, &newPermohonan.DokumenURL,
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
            status, created_at, updated_at
        FROM permohonansurat
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
            status, created_at, updated_at
        FROM permohonansurat WHERE id = ?
    `
	err := r.db.QueryRow(query, id).Scan(
		&permohonan.ID, &permohonan.NIK, &permohonan.NamaLengkap, &permohonan.TempatLahir,
		&permohonan.TanggalLahir, &permohonan.JenisKelamin, &permohonan.Pendidikan,
		&permohonan.Pekerjaan, &permohonan.Agama, &permohonan.StatusPernikahan,
		&permohonan.Kewarganegaraan, &permohonan.AlamatLengkap, &permohonan.JenisSurat,
		&permohonan.Keterangan, &permohonan.NomorHP, &permohonan.DokumenURL,
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

func (r *PermohonanSuratRepository) UpdatePermohonanSuratByID(id int64, permohonan model.PermohonanSurat) error {
	query := `
        UPDATE permohonansurat 
        SET nik = ?, nama_lengkap = ?, tempat_lahir = ?, tanggal_lahir = ?, jenis_kelamin = ?,
            pendidikan = ?, pekerjaan = ?, agama = ?, status_pernikahan = ?, kewarganegaraan = ?,
            alamat_lengkap = ?, jenis_surat = ?, keterangan = ?, nomor_hp = ?, dokumen_url = ?,
            status = ?, updated_at = ?
        WHERE id = ?
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.Status, permohonan.UpdatedAt, id,
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
        SET status = "Selesai", updated_at = ?
        WHERE id = ?
    `
    // Tes dengan nilai hard-coded untuk memastikan query bekerja
    // result, err := r.db.Exec("UPDATE permohonansurat SET status = 'Selesai', updated_at = ? WHERE id = ?", updatedAt, id)
    result, err := r.db.Exec(query, string(status), updatedAt, id) // Kembalikan ke ini setelah tes berhasil
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
