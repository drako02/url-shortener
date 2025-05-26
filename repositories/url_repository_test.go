package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/drako02/url-shortener/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	dialector := postgres.New(postgres.Config{Conn: db})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open GORM db: %v", err)
	}

	return gormDB, mock, err
}

// func Test_URLRepository_Delete(t *testing.T) {
// 	type fields struct {
// 		DB *gorm.DB
// 	}

// 	type mockBehaviour func(mock sqlmock.Sqlmock, shortCode string)

// 	db, mock, _ := newMockDB(t)

// 	type args struct {
// 		ctx       context.Context
// 		shortCode string
// 	}
// 	tests := []struct {
// 		name         string
// 		fields       fields
// 		args         args
// 		setUpMock    mockBehaviour
// 		wantErr      bool
// 		wantNotFound bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name:   "delete valid field",
// 			fields: fields{db},
// 			args:   args{context.Background(), "exists123"},
// 			setUpMock: func(mock sqlmock.Sqlmock, shortCode string) {
// 				mock.ExpectBegin()
// 				// mock.ExpectExec(`DELETE FROM "urls" WHERE "short_code" = \$1`).WithArgs(shortCode).WillReturnResult(sqlmock.NewResult(0, 1))
// 				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE short_code = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), shortCode).WillReturnResult(sqlmock.NewResult(0, 1))

// 				mock.ExpectCommit()
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "delete non existing field",
// 			fields: fields{db},
// 			args:   args{context.Background(), "nonexistent123"},
// 			setUpMock: func(mock sqlmock.Sqlmock, shortCode string) {
// 				mock.ExpectBegin()
// 				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE short_code = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), shortCode).WillReturnResult(sqlmock.NewResult(0, 0))

// 				mock.ExpectCommit()
// 			},
// 			wantErr:      true,
// 			wantNotFound: true,
// 		},
// 		{
// 			name:   "error occurred when deleting",
// 			fields: fields{db},
// 			args:   args{context.Background(), "errorprone123"},
// 			setUpMock: func(mock sqlmock.Sqlmock, shortCode string) {
// 				mock.ExpectBegin()
// 				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE short_code = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), shortCode).WillReturnError(fmt.Errorf("some exec error"))
// 				// mock.ExpectCommit()
// 				mock.ExpectRollback()

// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "empty string",
// 			fields: fields{db},
// 			args:   args{context.Background(), ""},
// 			setUpMock: func(mock sqlmock.Sqlmock, shortCode string) {
// 				mock.ExpectBegin()
// 				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE short_code = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), shortCode).WillReturnResult(sqlmock.NewResult(0, 0))
// 				mock.ExpectCommit()
// 			},
// 			wantErr:      true,
// 			wantNotFound: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			tt.setUpMock(mock, tt.args.shortCode)
// 			r := &URLRepository{
// 				DB: tt.fields.DB,
// 			}

// 			_, err := r.Delete(tt.args.ctx, tt.args.shortCode)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("URLRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			if tt.wantNotFound {
// 				if !errors.Is(err, ErrNotFound) {
// 					t.Errorf("expected ErrNotFound, got %v", err)
// 				}
// 			}

// 			// verify all expectations
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("there were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

func TestDelete(t *testing.T) {
	type fields struct{ DB *gorm.DB }

	type args struct {
		ctx context.Context
		id  uint
	}

	gormDB, mock, err := newMockDB(t)

	if err != nil {
		t.Errorf("failed to create mockDB: %v", err)
	}

	now := time.Now()

	tests := []struct {
		name          string
		field         fields
		arg           args
		mockBehaviour func(mock sqlmock.Sqlmock, id uint)
		want          models.URL
		wantedErr     bool
	}{
		{
			"valid shortcode",
			fields{gormDB},
			args{context.Background(), 1},
			func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).WithArgs(id, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "short_code", "long_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, "code123", "https://example.com", now, now, nil))
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE id = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			models.URL{
				ID:        1, // Or expectedID
				ShortCode: "code123",
				LongUrl:   "https://example.com",
				CreatedAt: now,
				UpdatedAt: now,
				// CreatedAt and UpdatedAt should match 'now' if Delete returns the fetched object
				// Model definition will determine if these are time.Time or *time.Time
			},
			false,
		},
		{
			"non existent id",
			fields{gormDB},
			args{context.Background(), 1},
			func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).WithArgs(id, 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			models.URL{},
			true,
		},
		{
			"already deleted id",
			fields{gormDB},
			args{context.Background(), 3},
			func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			models.URL{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &URLRepository{DB: tt.field.DB}

			tt.mockBehaviour(mock, tt.arg.id)

			url, err := repo.Delete(tt.arg.ctx, tt.arg.id)
			t.Logf("url : %v", url)

			if (err != nil) != tt.wantedErr {
				t.Errorf("URLRepository.Delete() error  = %v, wantedError = %v", err, tt.wantedErr)
			}

			if url != tt.want {
				t.Errorf("expected: %v, got: %v", tt.want, url)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_IsValidUpdateField(t *testing.T) {
	type args struct {
		field string
	}

	tests := []struct {
		name string
		arg  args
		want bool
	}{
		{
			name: "valid field",
			arg:  args{"long_url"},
			want: true,
		},
		{
			name: "invalid field",
			arg:  args{"id"},
			want: false,
		},
	}

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IsValidUpdateField(tt.arg.field)

			if res != tt.want{
				t.Errorf("wanted=%v, got=%v", tt.want, res)
			}
		} )

	}

}
func Test_UpdateById(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    uint
		field string
	}

	type fields struct {
		DB *gorm.DB
	}

	gormDB, _, err := newMockDB(t)

	if err != nil {
		t.Errorf("Failed to mock database %v", err)
	}
	tests := []struct {
		name    string
		field   fields
		arg     args
		want    models.URL
		wantErr bool
	}{
		{
			"valid field",
			fields{gormDB},
			args{context.Background(), 123, "long_url"},
			models.URL{ID: 123},
			false,
		},
		{
			"invalid field",
			fields{gormDB},
			args{context.Background(), 123, "uid"},
			models.URL{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &URLRepository{DB: tt.field.DB}
			url, err := r.UpdateById(tt.arg.ctx, tt.arg.id, tt.arg.field)

			if (err != nil) != tt.wantErr {
				t.Errorf("update by id err = %v, wantedErr = %v", err, tt.wantErr)
			}

			if tt.want != url {
				t.Errorf("wanted = %v, got = %v", tt.want, url)
			}

		})

	}

}
