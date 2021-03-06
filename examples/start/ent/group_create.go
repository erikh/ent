// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/examples/start/ent/group"
)

// GroupCreate is the builder for creating a Group entity.
type GroupCreate struct {
	config
	name  *string
	users map[int]struct{}
}

// SetName sets the name field.
func (gc *GroupCreate) SetName(s string) *GroupCreate {
	gc.name = &s
	return gc
}

// AddUserIDs adds the users edge to User by ids.
func (gc *GroupCreate) AddUserIDs(ids ...int) *GroupCreate {
	if gc.users == nil {
		gc.users = make(map[int]struct{})
	}
	for i := range ids {
		gc.users[ids[i]] = struct{}{}
	}
	return gc
}

// AddUsers adds the users edges to User.
func (gc *GroupCreate) AddUsers(u ...*User) *GroupCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gc.AddUserIDs(ids...)
}

// Save creates the Group in the database.
func (gc *GroupCreate) Save(ctx context.Context) (*Group, error) {
	if gc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := group.NameValidator(*gc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	return gc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (gc *GroupCreate) SaveX(ctx context.Context) *Group {
	v, err := gc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (gc *GroupCreate) sqlSave(ctx context.Context) (*Group, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(gc.driver.Dialect())
		gr      = &Group{config: gc.config}
	)
	tx, err := gc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(group.Table).Default()
	if value := gc.name; value != nil {
		insert.Set(group.FieldName, *value)
		gr.Name = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(group.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	gr.ID = int(id)
	if len(gc.users) > 0 {
		for eid := range gc.users {

			query, args := builder.Insert(group.UsersTable).
				Columns(group.UsersPrimaryKey[0], group.UsersPrimaryKey[1]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return gr, nil
}
