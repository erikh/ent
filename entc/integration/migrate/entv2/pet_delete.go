// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package entv2

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/entc/integration/migrate/entv2/pet"
	"github.com/facebookincubator/ent/entc/integration/migrate/entv2/predicate"
)

// PetDelete is the builder for deleting a Pet entity.
type PetDelete struct {
	config
	predicates []predicate.Pet
}

// Where adds a new predicate to the delete builder.
func (pd *PetDelete) Where(ps ...predicate.Pet) *PetDelete {
	pd.predicates = append(pd.predicates, ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *PetDelete) Exec(ctx context.Context) (int, error) {
	return pd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *PetDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *PetDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(pd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(pet.Table))
	for _, p := range pd.predicates {
		p(selector)
	}
	query, args := builder.Delete(pet.Table).FromSelect(selector).Query()
	if err := pd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// PetDeleteOne is the builder for deleting a single Pet entity.
type PetDeleteOne struct {
	pd *PetDelete
}

// Exec executes the deletion query.
func (pdo *PetDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{pet.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *PetDeleteOne) ExecX(ctx context.Context) {
	pdo.pd.ExecX(ctx)
}
