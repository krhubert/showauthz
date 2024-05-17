package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/nsf/jsondiff"
	"github.com/r3labs/diff/v3"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Equal[T any](t *testing.T, got T, want T) {
	t.Helper()
	Equalf(t, got, want, "")
}

func Equalf[T any](t *testing.T, got T, want T, format string, a ...any) {
	t.Helper()
	if isNil(got) && isNil(want) {
		return
	}
	if isNil(got) || isNil(want) {
		t.Fatalf("not equal %s \n diff: %s\n", fmt.Sprintf(format, a...), prettydiff(got, want))
	}
	gotd := deref(got)
	wantd := deref(want)

	egot, ok := gotd.(interface{ Equal(T) bool })
	if ok {
		if !egot.Equal(wantd.(T)) {
			t.Fatalf("not equal %s \n diff: %s\n", fmt.Sprintf(format, a...), prettydiff(got, want))
		}
		return
	}

	if !reflect.DeepEqual(gotd, wantd) {
		t.Fatalf("not equal %s \n diff: %s\n", fmt.Sprintf(format, a...), prettydiff(got, want))
	}
}

func ErrorContains(t *testing.T, err error, target any) {
	t.Helper()
	if err == nil {
		t.Fatalf("got no error, but expected %v", target)
	}

	// catch any errors.Is/As panics
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s", r)
		}
	}()

	switch e := target.(type) {
	case string:
		// if this is a valid regexp, compile it and use it
		// otherwise, just use it as a string
		if re, err1 := regexp.Compile(e); err1 == nil {
			if !re.MatchString(err.Error()) {
				t.Fatalf("error %q does not match %q", err, e)
			}
		} else {
			if !strings.Contains(err.Error(), e) {
				t.Fatalf("error %q does not contain %q", err, e)
			}
		}

	case error:
		if !errors.Is(err, e) {
			t.Fatalf("error %q is not %T", err, e)
		}
	default:
		if !errors.As(err, e) {
			t.Fatalf("error %q as not %T", err, e)
		}
	}
}

func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("got no error")
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func Nil(t *testing.T, got any) {
	t.Helper()
	if !isNil(got) {
		t.Fatalf("not nil %v", got)
	}
}

func Len[T any](t *testing.T, got []T, l int) {
	t.Helper()
	if len(got) != l {
		var zero T
		t.Fatalf("[]%T got length %d, want %d", zero, len(got), l)
	}
}

func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Fatalf("got false, want true")
	}
}

func prettydiff[T any](a, b T) string {
	// for strings use diffmatchpatch to get better output
	switch any(a).(type) {
	case string:
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(any(b).(string), any(a).(string), true)
		return dmp.DiffPrettyText(diffs)
	}

	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	opts := jsondiff.DefaultConsoleOptions()
	d, ret := jsondiff.Compare(aj, bj, &opts)

	// some un exported fields are different
	// use fallback less readable output
	if d == jsondiff.FullMatch {
		changelog, _ := diff.Diff(a, b)
		ret = "\n"
		for _, c := range changelog {
			ret += fmt.Sprintf("[%s]%T path %s: %q -> %q\n", c.Type, a, strings.Join(c.Path, "."), c.From, c.To)
		}
		return ret
	}
	return ret
}

func deref(a any) any {
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface()
}

func isNil(obj any) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Pointer, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		return v.IsNil()
	}
	return false
}
