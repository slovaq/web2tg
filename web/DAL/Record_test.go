package DAL

import (
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestCreateOrGetCity(t *testing.T) {
	type args struct {
		name   string
		chatID int
	}
	tests := []struct {
		name string
		args args
		want *City
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateOrGetCity(tt.args.name, tt.args.chatID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateOrGetCity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	type args struct {
		login    string
		name     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUser(tt.args.login, tt.args.name, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_CreateRecord(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Login    string
		Password string
	}
	type args struct {
		date time.Time
		dest []City
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Record
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Login:    tt.fields.Login,
				Password: tt.fields.Password,
			}
			got, err := u.CreateRecord(tt.args.date, tt.args.dest, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}
