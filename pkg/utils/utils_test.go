package utils

import "testing"

func Test_base62Encode(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{
			name: "1",
			args: args{
				input: "123",
			},
			wantOutput: "DWjr",
		},
		{
			name: "2",
			args: args{
				input: "abc",
			},
			wantOutput: "QmIN",
		},
		{
			name: "3",
			args: args{
				input: `{"slideRefId:"s0000000","RefId":"sp000064","mediaGuid":"EVqmLth8aOimOW5l","mediaPlay":true}`,
			},
			wantOutput: "1Rei7k9ZzZeRbOvyDfi4sZ0fdNdKHUP8EVBRVofFLl5dfriJZ7rhsGwB08JxcvgvyG5x7Sgm1Ujx6O1r68UAxg73KeMzqm7w8AhaQOfluD3ar6c75FRO4PdmfOv",
		},
		{
			name: "4",
			args: args{
				input: "",
			},
			wantOutput: "",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotOutput := Base62Encode(tt.args.input)
				if gotOutput != tt.wantOutput {
					t.Errorf("base62Encode() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
				}
			},
		)
	}
}

func TestDecodeStr(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{
			name: "1",
			args: args{
				input: "DWjr",
			},
			wantOutput: "123",
		},
		{
			name: "2",
			args: args{
				input: "QmIN",
			},
			wantOutput: "abc",
		},
		{
			name: "3",
			args: args{
				input: "1Rei7k9ZzZeRbOvyDfi4sZ0fdNdKHUP8EVBRVofFLl5dfriJZ7rhsGwB08JxcvgvyG5x7Sgm1Ujx6O1r68UAxg73KeMzqm7w8AhaQOfluD3ar6c75FRO4PdmfOv",
			},
			wantOutput: `{"slideRefId:"s0000000","RefId":"sp000064","mediaGuid":"EVqmLth8aOimOW5l","mediaPlay":true}`,
		},
		{
			name: "4",
			args: args{
				input: "",
			},
			wantOutput: "",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotOutput := Base62Decode(tt.args.input)
				if gotOutput != tt.wantOutput {
					t.Errorf("DecodeStr() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
				}
			},
		)
	}
}
