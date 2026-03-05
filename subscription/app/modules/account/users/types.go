package users

type UserListDTO struct {
	ID        string `json:"id"`
	Fullname  string `json:"fullname"` // Nama Lengkap
	NRK       string `json:"nrk"`
	NIK       string `json:"nik"`
	NIKKI     string `json:"nikki"`
	NPSN      string `json:"npsn"`
	GroupName string `json:"group_name"` // Kelas / Group
	RoleName  string `json:"role_name"`  // Peran
	LevelName string `json:"level_name"` // Level
	Status    string `json:"status"`     // Status (Aktif/Pending)
	IsActive  bool   `json:"is_active"`
}

type UserUpdateDTO struct {
	Fullname string
	NIK      string
	NRK      string
	IsActive bool
	// Tambahkan field lain jika perlu
}

type User struct {
	ID            string `json:"id"`
	Fullname      string `json:"fullname"`
	LevelID       int    `json:"level_id"`
	RoleID        string `json:"role_id"`
	GroupID       string `json:"group_id"`
	Status        string `json:"status"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Level         string `json:"level"` // Join result
	Role          string `json:"role"`  // Join result
	Group         string `json:"group"` // Join result
	Nrk           string `json:"nrk"`
	Nik           string `json:"nik"`
	Npsn          string `json:"npsn"`
	Nikki         string `json:"nikki"`
	Nisn          string `json:"nisn"`
	Nip           string `json:"nip"`
	Photo         string `json:"photo"`
	Jabatan       string `json:"jabatan"`
	MultipleLogin int    `json:"multiple_login"`
	JenisKelamin  string `json:"jenis_kelamin"`
	JenisAkun     string `json:"jenis_akun"`
	IsActive      bool   `json:"is_active"`
	DokumenPendukun string `json:"dokumen_pendukung"`
}
