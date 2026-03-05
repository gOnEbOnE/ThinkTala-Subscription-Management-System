package login

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User merepresentasikan data dari tabel account.za_users
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
}

// Struct untuk payload input dari Controller -> Service
type LoginPayload struct {
	LoginID   string `json:"login_id"`
	Password  string `json:"password"`
	IP        string `json:"ip"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Struct output Service -> Controller
type LoginResult struct {
	Token    string         `json:"token"`
	UserData map[string]any `json:"user_data"`
}
