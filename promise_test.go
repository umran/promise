package promise

import (
	"errors"
	"testing"
)

func TestAwait(t *testing.T) {
	res, err := New(func(p *Promise) {
		p.Resolve("ok")
	}).Await()

	if err != nil {
		t.Error(err)
	}

	msg, ok := res.(string)
	if !ok {
		t.Error("unexpected type")
	}

	if msg != "ok" {
		t.Error("unexpected message")
	}
}

func TestAwaitAll(t *testing.T) {
	plist := make([]*Promise, 20000)
	for i := 0; i < 20000; i++ {
		plist[i] = New(func(p *Promise) {
			p.Resolve("ok")
		})
	}

	resList, err := AwaitAll(plist)
	if err != nil {
		t.Error(err)
	}

	if len(resList) != len(plist) {
		t.Error("unexpected result length")
	}

	for _, res := range resList {
		msg, ok := res.(string)
		if !ok {
			t.Error("unexpected type")
		}

		if msg != "ok" {
			t.Error("unexpected message")
		}
	}
}

func TestReject(t *testing.T) {
	_, err := New(func(p *Promise) {
		p.Reject(errors.New("ok"))
	}).Await()

	if err == nil {
		t.Error("expected an error, but received none")
	} else if err.Error() != "ok" {
		t.Error("unexpected error")
	}
}

func TestAwaitAllReject(t *testing.T) {
	plist := []*Promise{New(func(p *Promise) { p.Reject(errors.New("ok")) })}
	_, err := AwaitAll(plist)

	if err == nil {
		t.Error("expected an error, but received none")
	} else if err.Error() != "ok" {
		t.Error("unexpected error")
	}
}
