package storage

import (
	"testing"

	"ohurlshortener/utils"
)

func TestNewUser(t *testing.T) {
	init4Test(t)
	// NewUser("ohUrlShortener", "-2aDzm=0(ln_9^1")
	NewUser("ohUrlShortener1", "-2aDzm=0(ln_9^1")
	NewUser("ohUrlShortener2", "-2aDzm=0(ln_9^1")
}

func init4Test(t *testing.T) {
	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = InitDatabaseService()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPasswordBase58Hash(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TestPasswordBase58Hash", want: "EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t", wantErr: false, args: args{password: "-2aDzm=0(ln_9^1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PasswordBase58Hash(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordBase58Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PasswordBase58Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
