package utils_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/utils"
)

func TestReadIntList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		msg     string
		reader  *bufio.Reader
		want    []int
		wantErr bool
	}{
		{
			name:    "valid ints",
			reader:  bufio.NewReader(strings.NewReader("1 2 3\n")),
			msg:     "testing",
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "invalid ints",
			reader:  bufio.NewReader(strings.NewReader("1.1 Two 3\n")),
			msg:     "testing",
			want:    []int{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "no input",
			reader:  bufio.NewReader(strings.NewReader("\n")),
			msg:     "testing",
			want:    []int{},
			wantErr: true,
		},
		{
			name:    "no input 1",
			reader:  bufio.NewReader(strings.NewReader("")),
			msg:     "testing",
			want:    []int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.ReadIntList(tt.msg, tt.reader)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadIntList() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadIntList() succeeded unexpectedly")
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ReadIntList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAndSanitizeInput(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		msg     string
		reader  *bufio.Reader
		want    string
		wantErr bool
	}{
		{
			name:    "non empty input",
			msg:     "testing",
			reader:  bufio.NewReader(strings.NewReader("hello\n")),
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "empty input",
			msg:     "testing",
			reader:  bufio.NewReader(strings.NewReader("")),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.ReadAndSanitizeInput(tt.msg, tt.reader)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadAndSanitizeInput() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadAndSanitizeInput() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ReadAndSanitizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintListInRows(t *testing.T) {
	osWrite := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	utils.PrintListInRows([]string{"hello", "hi", "a", "b", "c"})

	w.Close()
	os.Stdout = osWrite

	out, _ := io.ReadAll(r)

	want := fmt.Sprintf(" 1.%-15s 2.%-15s 3.%-15s 4.%-15s 5.%-15s\n\n", "hello", "hi", "a", "b", "c")
	if string(out) != want {
		t.Errorf("got %s wanted %s", string(out), want)
	}
}
