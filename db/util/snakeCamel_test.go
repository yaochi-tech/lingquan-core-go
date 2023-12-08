package util

import "testing"

func TestToCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"测试转驼峰命名",
			args{"hello_world"},
			"helloWorld",
		},
		{
			"测试转驼峰命名",
			args{"hello_123"},
			"hello123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamel(tt.args.s); got != tt.want {
				t.Errorf("ToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSnake(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"测试转蛇形命名",
			args{"helloWorld"},
			"hello_world",
		},
		{
			"测试转蛇形命名",
			args{"hello123"},
			"hello_123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnake(tt.args.s); got != tt.want {
				t.Errorf("ToSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}
