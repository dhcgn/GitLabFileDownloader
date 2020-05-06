package updater

import "testing"

func TestEquinoxUpdate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "Trivial test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EquinoxUpdate(); (err != nil) != tt.wantErr {
				t.Errorf("EquinoxUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}