package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

// func TestDelete(t *testing.T) {
// 	type fields struct{ DB *gorm.DB }

// 	type args struct {
// 		ctx context.Context
// 		id  uint
// 	}

// 	gormDB, mock, err := newMockDB(t)

// 	if err != nil {
// 		t.Errorf("failed to create mockDB: %v", err)
// 	}

// 	now := time.Now()

// 	tests := []struct {
// 		name          string
// 		field         fields
// 		arg           args
// 		mockBehaviour func(mock sqlmock.Sqlmock, id uint)
// 		want          models.URL
// 		wantedErr     bool
// 	}{
// 		{
// 			"valid shortcode",
// 			fields{gormDB},
// 			args{context.Background(), 1},
// 			func(mock sqlmock.Sqlmock, id uint) {
// 				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).WithArgs(id, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "short_code", "long_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, "code123", "https://example.com", now, now, nil))
// 				mock.ExpectBegin()
// 				mock.ExpectExec(`UPDATE "urls" SET "deleted_at"=\$1 WHERE id = \$2 AND "urls"."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(0, 1))
// 				mock.ExpectCommit()
// 			},
// 			models.URL{
// 				ID:        1, // Or expectedID
// 				ShortCode: "code123",
// 				LongUrl:   "https://example.com",
// 				CreatedAt: now,
// 				UpdatedAt: now,
// 				// CreatedAt and UpdatedAt should match 'now' if Delete returns the fetched object
// 				// Model definition will determine if these are time.Time or *time.Time
// 			},
// 			false,
// 		},
// 		{
// 			"non existent id",
// 			fields{gormDB},
// 			args{context.Background(), 1},
// 			func(mock sqlmock.Sqlmock, id uint) {
// 				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).WithArgs(id, 1).WillReturnError(gorm.ErrRecordNotFound)
// 			},
// 			models.URL{},
// 			true,
// 		},
// 		{
// 			"already deleted id",
// 			fields{gormDB},
// 			args{context.Background(), 3},
// 			func(mock sqlmock.Sqlmock, id uint) {
// 				mock.ExpectQuery(`SELECT \* FROM "urls" WHERE id = \$1 AND "urls"."deleted_at" IS NULL ORDER BY "urls"."id" LIMIT \$2`).
// 					WithArgs(id, 1).
// 					WillReturnError(gorm.ErrRecordNotFound)
// 			},
// 			models.URL{},
// 			true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := &URLRepository{DB: tt.field.DB}

// 			tt.mockBehaviour(mock, tt.arg.id)

// 			url, err := repo.Delete(tt.arg.ctx, tt.arg.id)
// 			t.Logf("url : %v", url)

// 			if (err != nil) != tt.wantedErr {
// 				t.Errorf("URLRepository.Delete() error  = %v, wantedError = %v", err, tt.wantedErr)
// 			}

// 			if url != tt.want {
// 				t.Errorf("expected: %v, got: %v", tt.want, url)
// 			}

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("there were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

func TestIsValidUpdateField(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IsValidUpdateField(tt.arg.field)

			if res != tt.want {
				t.Errorf("wanted=%v, got=%v", tt.want, res)
			}
		})

	}

}
func TestURLRepository_UpdateById(t *testing.T) {

	type args struct {
		ctx context.Context
		id  uint
		// field string
		data Data
	}

	type fields struct {
		DB *gorm.DB
	}

	gormDB, mock, err := newMockDB(t)
	now := time.Now()
	t.Log(now)

	if err != nil {
		t.Errorf("Failed to mock database %v", err)
	}

	tests := []struct {
		name          string
		field         fields
		arg           args
		mockBehaviour func(mock sqlmock.Sqlmock, id uint, field string, value any)
		want          error
		// wantErr       bool
	}{
		{
			"should allow valid field",
			fields{gormDB},
			args{context.Background(), 123, Data{Field: "long_url"}},
			nil,
			nil,
			// false,
		},
		{
			"should prevent invalid field",
			fields{gormDB},
			args{context.Background(), 123, Data{Field: "uid"}},
			nil,
			nil,
			// true,
		},
		{
			"should update db with new field if exists",
			fields{gormDB},
			args{context.Background(), 123, Data{LongUrl, "www.example.com/1"}},
			func(mock sqlmock.Sqlmock, id uint, field string, value any) {
				mock.ExpectBegin()
				query := fmt.Sprintf(`UPDATE "urls" SET "%s"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "urls"."deleted_at" IS NULL`, field)
				mock.ExpectExec(query).WithArgs(value, sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(123, 1))
			},
			nil,
			// false,
		},
		{
			"should return error for failed update operation",
			fields{gormDB},
			args{context.Background(), 123, Data{LongUrl, "www.example.com"}},
			func(mock sqlmock.Sqlmock, id uint, field string, value any) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			nil,
			// false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &URLRepository{DB: tt.field.DB}

			// t.Logf("mockBehaviour: %v")

			if tt.mockBehaviour != nil {
				tt.mockBehaviour(mock, tt.arg.id, tt.arg.data.Field, tt.arg.data.Value)
				t.Log("mock func did run")
			}

			err := r.UpdateById(tt.arg.ctx, tt.arg.id, tt.arg.data)

			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("update by id err = %v, wantedErr = %v", err, tt.wantErr)
			// }

			if tt.want != nil {
				t.Errorf("update by id err = %v, wantedErr = %v", err, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}

		})

	}

}
