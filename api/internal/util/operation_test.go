package util

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// These guards run before RunOperation touches the database, so a nil db is
// enough - reaching it would panic.

func TestRunOperationRequiresRequestSource(t *testing.T) {
	err := RunOperation(context.Background(), nil, Operation{Kind: "sample.create", PerformedBy: uuid.New()},
		func(context.Context, *sqlx.Tx) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "no request source") {
		t.Errorf("err = %v, want a missing request source error", err)
	}
}

func TestRunOperationRequiresKind(t *testing.T) {
	ctx := WithRequestSource(context.Background(), RequestSource{Client: ClientWeb})
	err := RunOperation(ctx, nil, Operation{PerformedBy: uuid.New()},
		func(context.Context, *sqlx.Tx) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "kind is required") {
		t.Errorf("err = %v, want a missing kind error", err)
	}
}

func TestRunOperationRefusesNesting(t *testing.T) {
	ctx := WithRequestSource(context.Background(), RequestSource{Client: ClientWeb})
	ctx = context.WithValue(ctx, operationCtxKey{}, uuid.New())
	err := RunOperation(ctx, nil, Operation{Kind: "sample.create", PerformedBy: uuid.New()},
		func(context.Context, *sqlx.Tx) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "can't nest") {
		t.Errorf("err = %v, want a nesting error", err)
	}
}

func TestSystemContext(t *testing.T) {
	ctx, span := SystemContext(context.Background(), "cleanup")
	defer span.End()

	src, ok := requestSource(ctx)
	if !ok || src.Client != ClientSystem || src.Info["task"] != "cleanup" {
		t.Errorf("request source = %+v, %v", src, ok)
	}
}
