package keystorage

import (
	"context"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	"reflect"
	"testing"
)

func getDB() *gorm.DB {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := gorm.Open(sqlite.Open(path.Join(home, ".link/link.db")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func TestNewServiceImpl(t *testing.T) {
	db := getDB()
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{"test1",
			args{
				ctx: context.Background(),
				db:  db,
			},
			&ServiceImpl{
				db:    db,
				Error: nil,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServiceImpl(tt.args.ctx, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServiceImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImpl_Get(t *testing.T) {
	db := getDB()
	type fields struct {
		db    *gorm.DB
		Error error
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "key not found 1",
			fields: fields{
				db:    db,
				Error: nil,
			},
			args: args{
				key: "123",
			},
			want: "",
		},
		{
			name: "key not found 2",
			fields: fields{
				db:    db,
				Error: nil,
			},
			args: args{
				key: "--=-=",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ServiceImpl{
				db:    tt.fields.db,
				Error: tt.fields.Error,
			}
			if got := self.Get(tt.args.key); got != tt.want || self.Err() != gorm.ErrRecordNotFound {
				t.Errorf("Get() = %v, want %v, err %v", got, tt.want, self.Err())
			}
		})
	}
}

func TestServiceImpl_Set(t *testing.T) {
	db := getDB()
	type fields struct {
		tableName string
		db        *gorm.DB
		Error     error
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set 1",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{
				key:   "3721",
				value: "99999999999",
			},
		},
		{
			name: "set 2",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{
				key:   "键",
				value: "值",
			},
		},
		{
			name: "set 3",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{
				key:   "😀",
				value: "😀",
			},
		},
		{
			name: "set 4",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{
				key:   "😀",
				value: "update",
			},
		},
		{
			name: "set 5",
			fields: fields{
				tableName: "keys",
				db:        db,
				Error:     nil,
			},
			args: args{
				key:   "😀",
				value: "table name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ServiceImpl{
				tableName: tt.fields.tableName,
				db:        tt.fields.db,
				Error:     tt.fields.Error,
			}
			self.Set(tt.args.key, tt.args.value)
			if self.Err() != nil {
				t.Error(self.Err())
			}
		})
	}
}

func TestServiceImpl_Get1(t *testing.T) {
	db := getDB()
	type fields struct {
		tableName string
		db        *gorm.DB
		Error     error
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "get data",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{"3721"},
			want: "99999999999",
		},
		{
			name: "get data",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{"键"},
			want: "值",
		},
		{
			name: "get data",
			fields: fields{
				tableName: "",
				db:        db,
				Error:     nil,
			},
			args: args{"😀"},
			want: "update",
		},
		{
			name: "get data",
			fields: fields{
				tableName: "keys",
				db:        db,
				Error:     nil,
			},
			args: args{"😀"},
			want: "table name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ServiceImpl{
				tableName: tt.fields.tableName,
				db:        tt.fields.db,
				Error:     tt.fields.Error,
			}
			if got := self.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImpl_Delete(t *testing.T) {
	type fields struct {
		tableName string
		db        *gorm.DB
		Error     error
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "delete 1", fields: fields{tableName: "", db: getDB(), Error: nil}, args: args{"3721"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &ServiceImpl{
				tableName: tt.fields.tableName,
				db:        tt.fields.db,
				Error:     tt.fields.Error,
			}
			self.Delete(tt.args.key)
		})
	}
}
