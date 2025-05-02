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

func (r *PermohonanSuratRepository) AddPermohonanSurat(permohonan model.PermohonanSurat) (int64, error) {
	query := `
        INSERT INTO permohonansurat (
            nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
            pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
            alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
            nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
            nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, ditujukan, status,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.NamaUsaha, permohonan.JenisUsaha, permohonan.AlamatUsaha, permohonan.AlamatTujuan,
		permohonan.AlasanPindah, permohonan.NamaAyah, permohonan.NamaIbu, permohonan.TglKematian,
		permohonan.PenyebabKematian, permohonan.Ditujukan, permohonan.Status,
		time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Error inserting permohonan surat: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	return id, nil
}

func (r *PermohonanSuratRepository) GetPermohonanSurat() ([]model.PermohonanSurat, error) {
	query := `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
               pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
               alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
               nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
               nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, ditujukan, status,
               created_at, updated_at
        FROM permohonansurat
    `
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying permohonan surat: %v", err)
		return nil, err
	}
	defer rows.Close()

	var permohonans []model.PermohonanSurat
	for rows.Next() {
		var p model.PermohonanSurat
		err := rows.Scan(
			&p.ID, &p.NIK, &p.NamaLengkap, &p.TempatLahir, &p.TanggalLahir, &p.JenisKelamin,
			&p.Pendidikan, &p.Pekerjaan, &p.Agama, &p.StatusPernikahan, &p.Kewarganegaraan,
			&p.AlamatLengkap, &p.JenisSurat, &p.Keterangan, &p.NomorHP, &p.DokumenURL,
			&p.NamaUsaha, &p.JenisUsaha, &p.AlamatUsaha, &p.AlamatTujuan, &p.AlasanPindah,
			&p.NamaAyah, &p.NamaIbu, &p.TglKematian, &p.PenyebabKematian, &p.Ditujukan, &p.Status,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning permohonan surat: %v", err)
			return nil, err
		}
		permohonans = append(permohonans, p)
	}
	return permohonans, nil
}

func (r *PermohonanSuratRepository) GetPermohonanSuratByID(id int64) (model.PermohonanSurat, error) {
	query := `
        SELECT id, nik, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
               pendidikan, pekerjaan, agama, status_pernikahan, kewarganegaraan,
               alamat_lengkap, jenis_surat, keterangan, nomor_hp, dokumen_url,
               nama_usaha, jenis_usaha, alamat_usaha, alamat_tujuan, alasan_pindah,
               nama_ayah, nama_ibu, tgl_kematian, penyebab_kematian, ditujukan, status,
               created_at, updated_at
        FROM permohonansurat WHERE id = ?
    `
	var p model.PermohonanSurat
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.NIK, &p.NamaLengkap, &p.TempatLahir, &p.TanggalLahir, &p.JenisKelamin,
		&p.Pendidikan, &p.Pekerjaan, &p.Agama, &p.StatusPernikahan, &p.Kewarganegaraan,
		&p.AlamatLengkap, &p.JenisSurat, &p.Keterangan, &p.NomorHP, &p.DokumenURL,
		&p.NamaUsaha, &p.JenisUsaha, &p.AlamatUsaha, &p.AlamatTujuan, &p.AlasanPindah,
		&p.NamaAyah, &p.NamaIbu, &p.TglKematian, &p.PenyebabKematian, &p.Ditujukan, &p.Status,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return p, fmt.Errorf("permohonan surat with ID %d not found", id)
	}
	if err != nil {
		log.Printf("Error querying permohonan surat by ID: %v", err)
		return p, err
	}
	return p, nil
}

func (r *PermohonanSuratRepository) UpdatePermohonanSurat(permohonan model.PermohonanSurat) error {
	query := `
        UPDATE permohonansurat SET
            nik = ?, nama_lengkap = ?, tempat_lahir = ?, tanggal_lahir = ?, jenis_kelamin = ?,
            pendidikan = ?, pekerjaan = ?, agama = ?, status_pernikahan = ?, kewarganegaraan = ?,
            alamat_lengkap = ?, jenis_surat = ?, keterangan = ?, nomor_hp = ?, dokumen_url = ?,
            nama_usaha = ?, jenis_usaha = ?, alamat_usaha = ?, alamat_tujuan = ?, alasan_pindah = ?,
            nama_ayah = ?, nama_ibu = ?, tgl_kematian = ?, penyebab_kematian = ?, ditujukan = ?, status = ?,
            updated_at = ?
        WHERE id = ?
    `
	result, err := r.db.Exec(query,
		permohonan.NIK, permohonan.NamaLengkap, permohonan.TempatLahir, permohonan.TanggalLahir,
		permohonan.JenisKelamin, permohonan.Pendidikan, permohonan.Pekerjaan, permohonan.Agama,
		permohonan.StatusPernikahan, permohonan.Kewarganegaraan, permohonan.AlamatLengkap,
		permohonan.JenisSurat, permohonan.Keterangan, permohonan.NomorHP, permohonan.DokumenURL,
		permohonan.NamaUsaha, permohonan.JenisUsaha, permohonan.AlamatUsaha, permohonan.AlamatTujuan,
		permohonan.AlasanPindah, permohonan.NamaAyah, permohonan.NamaIbu, permohonan.TglKematian,
		permohonan.PenyebabKematian, permohonan.Ditujukan, permohonan.Status,
		time.Now(), permohonan.ID,
	)
	if err != nil {
		log.Printf("Error updating permohonan surat: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("permohonan surat with ID %d not found", permohonan.ID)
	}
	return nil
}

func (r *PermohonanSuratRepository) DeletePermohonanSurat(id int64) error {
	query := `DELETE FROM permohonansurat WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting permohonan surat: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("permohonan surat with ID %d not found", id)
	}
	return nil
}

func (r *PermohonanSuratRepository) UpdateStatus(id int64, status model.Status) error {
	query := `UPDATE permohonansurat SET status = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		log.Printf("Error updating status: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("permohonan surat with ID %d not found", id)
	}
	return nil
}
