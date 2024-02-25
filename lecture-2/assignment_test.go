package assignment

import (
	"github.com/stretchr/testify/mock"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

func Reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Palindrome(s []string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func Anagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	charCount := make(map[rune]int)
	for _, char := range s1 {
		charCount[char]++
	}
	for _, char := range s2 {
		charCount[char]--
		if charCount[char] < 0 {
			return false
		}
	}
	return true
}

func RemoveDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, s)
}

func ReplaceDigits(s, replacement string) string {
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			builder.WriteString(replacement)
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

type DataSource interface {
	ReadStudent(studentID int) (Student, error)
	ReadCourse(courseID int) (Course, error)
}

type Student interface {
	ID() int
	Name() string
}

type Course interface {
	ID() int
	Name() string
	EnrollStudent(student Student) error
}

func EnrollStudentToCourse(dataSource DataSource, studentID, courseID int) error {
	student, err := dataSource.ReadStudent(studentID)
	if err != nil {
		return err
	}
	course, err := dataSource.ReadCourse(courseID)
	if err != nil {
		return err
	}
	return course.EnrollStudent(student)
}

func TestReverse(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Reverse string",
			args: args{
				s: []string{"a", "b", "c"},
			},
			want: []string{"c", "b", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPalindrome(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Palindrome string",
			args: args{
				s: []string{"a", "b", "c", "b", "a"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Palindrome(tt.args.s); got != tt.want {
				t.Errorf("Palindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnagram(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Anagram string",
			args: args{
				s1: "abeceda",
				s2: "aebcdea",
			},
			want: true,
		},
		{
			name: "Not an anagram string",
			args: args{
				s1: "abeceda",
				s2: "abecede",
			},
			want: false,
		},
		{
			name: "Not an anagram string",
			args: args{
				s1: "abcd",
				s2: "abcdabcd",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Anagram(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("Anagram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDigits(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Remove digits",
			args: args{
				s: "a1b2c3",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDigits(tt.args.s); got != tt.want {
				t.Errorf("RemoveDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceDigits(t *testing.T) {
	type args struct {
		s string
		r string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Replace digits",
			args: args{
				s: "a1b2c3",
				r: "x",
			},
			want: "axbxcx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceDigits(tt.args.s, tt.args.r); got != tt.want {
				t.Errorf("ReplaceDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockDataSource struct {
	mock.Mock
}

func (m *mockDataSource) ReadStudent(studentID int) (Student, error) {
	args := m.Called(studentID)
	return args.Get(0).(Student), args.Error(1)
}

func (m *mockDataSource) ReadCourse(courseID int) (Course, error) {
	args := m.Called(courseID)
	return args.Get(0).(Course), args.Error(1)
}

type fakeStudent struct {
	id int
}

func (f *fakeStudent) ID() int {
	return f.id
}

func (*fakeStudent) Name() string {
	return "John Doe"
}

type fakeCourse struct {
	id int
}

func (f *fakeCourse) ID() int {
	return f.id
}

func (*fakeCourse) Name() string {
	return "Microservices 101"
}

func (*fakeCourse) EnrollStudent(s Student) error {
	return nil
}

func TestEnrollStudentToCourse(t *testing.T) {
	type args struct {
		dataSource DataSource
		sID        int // studentID
		cID        int // courseID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Enroll student to course",
			args: args{
				dataSource: &mockDataSource{},
				sID:        1,
				cID:        1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ds := tt.args.dataSource.(*mockDataSource)
		ds.
			On("ReadStudent", tt.args.sID).
			Return(&fakeStudent{tt.args.sID}, nil)
		ds.
			On("ReadCourse", tt.args.cID).
			Return(&fakeCourse{tt.args.cID}, nil)
		t.Run(tt.name, func(t *testing.T) {
			if err := EnrollStudentToCourse(tt.args.dataSource, tt.args.sID, tt.args.cID); (err != nil) != tt.wantErr {
				t.Errorf("EnrollStudentToCourse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.args.dataSource.(*mockDataSource).AssertExpectations(t)
	}
}
