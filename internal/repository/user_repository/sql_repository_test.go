package userrepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewSQLUserRepository(t *testing.T) {

	db, _, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *SQLUserRepository
	}{
		{
			name: "valid db",
			args: args{
				db: gormDb,
			},
			want: &SQLUserRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLUserRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLUserRepository_CreateAdminOffice(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("ADMIN_OFFICE", 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "buildings" \("building_name"\) VALUES \(\$1\) RETURNING "building_id"`).
					WithArgs("ADMIN_BUILDING").
					WillReturnRows(sqlmock.NewRows([]string{"building_id"}).AddRow("123e4567-e89b-12d3-a456-426614174000"))
				mock.ExpectExec(`INSERT INTO "floors" \("building_id","floor_number"\) VALUES \(\$1,\$2\) ON CONFLICT \("building_id","floor_number"\) DO UPDATE SET "building_id"="excluded"\."building_id"`).
					WithArgs("123e4567-e89b-12d3-a456-426614174000", 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "offices" \("office_name","building_id","floor_number"\) VALUES \(\$1,\$2,\$3\) ON CONFLICT \("office_id"\) DO UPDATE SET "building_id"="excluded"\."building_id","floor_number"="excluded"\."floor_number" RETURNING "office_id"`).
					WithArgs("ADMIN_OFFICE", "123e4567-e89b-12d3-a456-426614174000", 0).
					WillReturnRows(sqlmock.NewRows([]string{"office_id"}).AddRow("123e4567-e89b-12d3-a456-426614174002"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "office exists",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("ADMIN_OFFICE", 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id"}).AddRow("123e4567-e89b-12d3-a456-426614174002"))
			},
			wantErr: false,
		},
		{
			name: "database error when checking office",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("ADMIN_OFFICE", 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "building creation error",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("ADMIN_OFFICE", 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "buildings" \("building_name"\) VALUES \(\$1\) RETURNING "building_id"`).
					WithArgs("ADMIN_BUILDING").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			if err := sqlur.CreateAdminOffice(); (err != nil) != tt.wantErr {
				t.Errorf("CreateAdminOffice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLUserRepository_CreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		name     string
		email    string
		password string
		office   string
		role     roles.Role
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				name:     "kaushik",
				email:    "kaushik@a.com",
				password: "password",
				office:   "wg",
				role:     roles.Customer,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("wg", 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id"}).AddRow("123e4567-e89b-12d3-a456-426614174000"))
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" \("name","email","password","role","is_active","office_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "user_id"`).
					WithArgs("kaushik", "kaushik@a.com", "password", roles.Customer, true, "123e4567-e89b-12d3-a456-426614174000").
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("123e4567-e89b-12d3-a456-426614174003"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "office not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				name:     "kaushik",
				email:    "kaushik@a.com",
				password: "password",
				office:   "nonexistent",
				role:     roles.Customer,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("nonexistent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "user creation error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				name:     "kaushik",
				email:    "kaushik@a.com",
				password: "password",
				office:   "wg",
				role:     roles.Customer,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "offices" WHERE office_name = \$1 ORDER BY "offices"\."office_id" LIMIT \$2`).
					WithArgs("wg", 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id"}).AddRow("123e4567-e89b-12d3-a456-426614174000"))
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" \("name","email","password","role","is_active","office_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "user_id"`).
					WithArgs("kaushik", "kaushik@a.com", "password", roles.Customer, true, "123e4567-e89b-12d3-a456-426614174000").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			if err := sqlur.CreateUser(tt.args.name, tt.args.email, tt.args.password, tt.args.office, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLUserRepository_GetAllUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		want      []models.User
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			want: []models.User{
				{Name: "kaushik", Email: "kaushik@a.com"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"name", "email"}).
					AddRow("kaushik", "kaushik@a.com")
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE is_active = true`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			want: nil,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE is_active = true`).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			got, err := sqlur.GetAllUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLUserRepository_GetUserByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.User
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				email: "kaushik@a.com",
			},
			want: models.User{Name: "kaushik", Email: "kaushik@a.com"},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"name", "email"}).
					AddRow("kaushik", "kaushik@a.com")
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("kaushik@a.com", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "admin success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				email: "admin@a.com",
			},
			want: models.User{Name: "admin", Email: "admin@a.com"},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"name", "email"}).
					AddRow("admin", "admin@a.com")
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("admin@a.com", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				email: "notfound@example.com",
			},
			want: models.User{},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("notfound@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				email: "kaushik@a.com",
			},
			want: models.User{},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("kaushik@a.com", 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			got, err := sqlur.GetUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLUserRepository_GetUserById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.User
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				id: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: models.User{Name: "kaushik", Email: "kaushik@a.com"},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"name", "email"}).
					AddRow("kaushik", "kaushik@a.com")
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "invalid uuid",
			fields: fields{
				db: gormDb,
			},
			args: args{
				id: "invalid-uuid",
			},
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "user not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				id: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: models.User{},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 AND is_active = true ORDER BY "users"\."user_id" LIMIT \$2`).
					WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			got, err := sqlur.GetUserById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLUserRepository_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user models.User
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				user: models.User{Name: "kaushik", Email: "kaushik@a.com"},
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" \("name","email","password","role","is_active","office_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "user_id"`).
					WithArgs("kaushik", "kaushik@a.com", "", roles.Customer, true, "00000000-0000-0000-0000-000000000000").
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("123e4567-e89b-12d3-a456-426614174003"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				user: models.User{Name: "kaushik", Email: "kaushik@a.com"},
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" \("name","email","password","role","is_active","office_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "user_id"`).
					WithArgs("kaushik", "kaushik@a.com", "", roles.Customer, true, "00000000-0000-0000-0000-000000000000").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			sqlur := &SQLUserRepository{
				db: tt.fields.db,
			}
			if err := sqlur.Save(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
