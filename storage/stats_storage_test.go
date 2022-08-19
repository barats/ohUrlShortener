package storage

import (
	"testing"
)

func TestCallProcedureStatsIPSum(t *testing.T) {
	init4Test(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "TestCase CallProcedureStatsIPSum", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CallProcedureStatsIPSum(); (err != nil) != tt.wantErr {
				t.Errorf("CallProcedureStatsIPSum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCallProcedureStatsTop25(t *testing.T) {
	init4Test(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "TestCase CallProcedureStatsTop25", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CallProcedureStatsTop25(); (err != nil) != tt.wantErr {
				t.Errorf("CallProcedureStatsTop25() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCallProcedureStatsSum(t *testing.T) {
	init4Test(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "TestCase  CallProcedureStatsSum", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CallProcedureStatsSum(); (err != nil) != tt.wantErr {
				t.Errorf("CallProcedureStatsSum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
