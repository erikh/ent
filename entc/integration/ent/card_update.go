// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/gremlin"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/__"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/g"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/p"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/entc/integration/ent/card"
	"github.com/facebookincubator/ent/entc/integration/ent/predicate"
	"github.com/facebookincubator/ent/entc/integration/ent/user"
)

// CardUpdate is the builder for updating Card entities.
type CardUpdate struct {
	config

	update_time *time.Time

	name         *string
	clearname    bool
	owner        map[string]struct{}
	clearedOwner bool
	predicates   []predicate.Card
}

// Where adds a new predicate for the builder.
func (cu *CardUpdate) Where(ps ...predicate.Card) *CardUpdate {
	cu.predicates = append(cu.predicates, ps...)
	return cu
}

// SetName sets the name field.
func (cu *CardUpdate) SetName(s string) *CardUpdate {
	cu.name = &s
	return cu
}

// SetNillableName sets the name field if the given value is not nil.
func (cu *CardUpdate) SetNillableName(s *string) *CardUpdate {
	if s != nil {
		cu.SetName(*s)
	}
	return cu
}

// ClearName clears the value of name.
func (cu *CardUpdate) ClearName() *CardUpdate {
	cu.name = nil
	cu.clearname = true
	return cu
}

// SetOwnerID sets the owner edge to User by id.
func (cu *CardUpdate) SetOwnerID(id string) *CardUpdate {
	if cu.owner == nil {
		cu.owner = make(map[string]struct{})
	}
	cu.owner[id] = struct{}{}
	return cu
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (cu *CardUpdate) SetNillableOwnerID(id *string) *CardUpdate {
	if id != nil {
		cu = cu.SetOwnerID(*id)
	}
	return cu
}

// SetOwner sets the owner edge to User.
func (cu *CardUpdate) SetOwner(u *User) *CardUpdate {
	return cu.SetOwnerID(u.ID)
}

// ClearOwner clears the owner edge to User.
func (cu *CardUpdate) ClearOwner() *CardUpdate {
	cu.clearedOwner = true
	return cu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cu *CardUpdate) Save(ctx context.Context) (int, error) {
	if cu.update_time == nil {
		v := card.UpdateDefaultUpdateTime()
		cu.update_time = &v
	}
	if cu.name != nil {
		if err := card.NameValidator(*cu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(cu.owner) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	switch cu.driver.Dialect() {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		return cu.sqlSave(ctx)
	case dialect.Gremlin:
		return cu.gremlinSave(ctx)
	default:
		return 0, errors.New("ent: unsupported dialect")
	}
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CardUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CardUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CardUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CardUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(cu.driver.Dialect())
		selector = builder.Select(card.FieldID).From(builder.Table(card.Table))
	)
	for _, p := range cu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := cu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(card.Table)
	)
	updater = updater.Where(sql.InInts(card.FieldID, ids...))
	if value := cu.update_time; value != nil {
		updater.Set(card.FieldUpdateTime, *value)
	}
	if value := cu.name; value != nil {
		updater.Set(card.FieldName, *value)
	}
	if cu.clearname {
		updater.SetNull(card.FieldName)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if cu.clearedOwner {
		query, args := builder.Update(card.OwnerTable).
			SetNull(card.OwnerColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(cu.owner) > 0 {
		for _, id := range ids {
			eid, serr := strconv.Atoi(keys(cu.owner)[0])
			if serr != nil {
				return 0, rollback(tx, err)
			}
			query, args := builder.Update(card.OwnerTable).
				Set(card.OwnerColumn, eid).
				Where(sql.EQ(card.FieldID, id).And().IsNull(card.OwnerColumn)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(cu.owner) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"owner\" %v already connected to a different \"Card\"", keys(cu.owner))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

func (cu *CardUpdate) gremlinSave(ctx context.Context) (int, error) {
	res := &gremlin.Response{}
	query, bindings := cu.gremlin().Query()
	if err := cu.driver.Exec(ctx, query, bindings, res); err != nil {
		return 0, err
	}
	if err, ok := isConstantError(res); ok {
		return 0, err
	}
	return res.ReadInt()
}

func (cu *CardUpdate) gremlin() *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.V().HasLabel(card.Label)
	for _, p := range cu.predicates {
		p(v)
	}
	var (
		rv = v.Clone()
		_  = rv

		trs []*dsl.Traversal
	)
	if value := cu.update_time; value != nil {
		v.Property(dsl.Single, card.FieldUpdateTime, *value)
	}
	if value := cu.name; value != nil {
		v.Property(dsl.Single, card.FieldName, *value)
	}
	var properties []interface{}
	if cu.clearname {
		properties = append(properties, card.FieldName)
	}
	if len(properties) > 0 {
		v.SideEffect(__.Properties(properties...).Drop())
	}
	if cu.clearedOwner {
		tr := rv.Clone().InE(user.CardLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for id := range cu.owner {
		v.AddE(user.CardLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(user.CardLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(card.Label, user.CardLabel, id)),
		})
	}
	v.Count()
	if len(constraints) > 0 {
		constraints = append(constraints, &constraint{
			pred: rv.Count(),
			test: __.Is(p.GT(1)).Constant(&ErrConstraintFailed{msg: "update traversal contains more than one vertex"}),
		})
		v = constraints[0].pred.Coalesce(constraints[0].test, v)
		for _, cr := range constraints[1:] {
			v = cr.pred.Coalesce(cr.test, v)
		}
	}
	trs = append(trs, v)
	return dsl.Join(trs...)
}

// CardUpdateOne is the builder for updating a single Card entity.
type CardUpdateOne struct {
	config
	id string

	update_time *time.Time

	name         *string
	clearname    bool
	owner        map[string]struct{}
	clearedOwner bool
}

// SetName sets the name field.
func (cuo *CardUpdateOne) SetName(s string) *CardUpdateOne {
	cuo.name = &s
	return cuo
}

// SetNillableName sets the name field if the given value is not nil.
func (cuo *CardUpdateOne) SetNillableName(s *string) *CardUpdateOne {
	if s != nil {
		cuo.SetName(*s)
	}
	return cuo
}

// ClearName clears the value of name.
func (cuo *CardUpdateOne) ClearName() *CardUpdateOne {
	cuo.name = nil
	cuo.clearname = true
	return cuo
}

// SetOwnerID sets the owner edge to User by id.
func (cuo *CardUpdateOne) SetOwnerID(id string) *CardUpdateOne {
	if cuo.owner == nil {
		cuo.owner = make(map[string]struct{})
	}
	cuo.owner[id] = struct{}{}
	return cuo
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (cuo *CardUpdateOne) SetNillableOwnerID(id *string) *CardUpdateOne {
	if id != nil {
		cuo = cuo.SetOwnerID(*id)
	}
	return cuo
}

// SetOwner sets the owner edge to User.
func (cuo *CardUpdateOne) SetOwner(u *User) *CardUpdateOne {
	return cuo.SetOwnerID(u.ID)
}

// ClearOwner clears the owner edge to User.
func (cuo *CardUpdateOne) ClearOwner() *CardUpdateOne {
	cuo.clearedOwner = true
	return cuo
}

// Save executes the query and returns the updated entity.
func (cuo *CardUpdateOne) Save(ctx context.Context) (*Card, error) {
	if cuo.update_time == nil {
		v := card.UpdateDefaultUpdateTime()
		cuo.update_time = &v
	}
	if cuo.name != nil {
		if err := card.NameValidator(*cuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(cuo.owner) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	switch cuo.driver.Dialect() {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		return cuo.sqlSave(ctx)
	case dialect.Gremlin:
		return cuo.gremlinSave(ctx)
	default:
		return nil, errors.New("ent: unsupported dialect")
	}
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CardUpdateOne) SaveX(ctx context.Context) *Card {
	c, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// Exec executes the query on the entity.
func (cuo *CardUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CardUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CardUpdateOne) sqlSave(ctx context.Context) (c *Card, err error) {
	var (
		builder  = sql.Dialect(cuo.driver.Dialect())
		selector = builder.Select(card.Columns...).From(builder.Table(card.Table))
	)
	card.ID(cuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		c = &Card{config: cuo.config}
		if err := c.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Card: %v", err)
		}
		id = c.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Card with id: %v", cuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Card with the same id: %v", cuo.id)
	}

	tx, err := cuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(card.Table)
	)
	updater = updater.Where(sql.InInts(card.FieldID, ids...))
	if value := cuo.update_time; value != nil {
		updater.Set(card.FieldUpdateTime, *value)
		c.UpdateTime = *value
	}
	if value := cuo.name; value != nil {
		updater.Set(card.FieldName, *value)
		c.Name = *value
	}
	if cuo.clearname {
		var value string
		c.Name = value
		updater.SetNull(card.FieldName)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if cuo.clearedOwner {
		query, args := builder.Update(card.OwnerTable).
			SetNull(card.OwnerColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(cuo.owner) > 0 {
		for _, id := range ids {
			eid, serr := strconv.Atoi(keys(cuo.owner)[0])
			if serr != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(card.OwnerTable).
				Set(card.OwnerColumn, eid).
				Where(sql.EQ(card.FieldID, id).And().IsNull(card.OwnerColumn)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(cuo.owner) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"owner\" %v already connected to a different \"Card\"", keys(cuo.owner))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return c, nil
}

func (cuo *CardUpdateOne) gremlinSave(ctx context.Context) (*Card, error) {
	res := &gremlin.Response{}
	query, bindings := cuo.gremlin(cuo.id).Query()
	if err := cuo.driver.Exec(ctx, query, bindings, res); err != nil {
		return nil, err
	}
	if err, ok := isConstantError(res); ok {
		return nil, err
	}
	c := &Card{config: cuo.config}
	if err := c.FromResponse(res); err != nil {
		return nil, err
	}
	return c, nil
}

func (cuo *CardUpdateOne) gremlin(id string) *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.V(id)
	var (
		rv = v.Clone()
		_  = rv

		trs []*dsl.Traversal
	)
	if value := cuo.update_time; value != nil {
		v.Property(dsl.Single, card.FieldUpdateTime, *value)
	}
	if value := cuo.name; value != nil {
		v.Property(dsl.Single, card.FieldName, *value)
	}
	var properties []interface{}
	if cuo.clearname {
		properties = append(properties, card.FieldName)
	}
	if len(properties) > 0 {
		v.SideEffect(__.Properties(properties...).Drop())
	}
	if cuo.clearedOwner {
		tr := rv.Clone().InE(user.CardLabel).Drop().Iterate()
		trs = append(trs, tr)
	}
	for id := range cuo.owner {
		v.AddE(user.CardLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(user.CardLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(card.Label, user.CardLabel, id)),
		})
	}
	v.ValueMap(true)
	if len(constraints) > 0 {
		v = constraints[0].pred.Coalesce(constraints[0].test, v)
		for _, cr := range constraints[1:] {
			v = cr.pred.Coalesce(cr.test, v)
		}
	}
	trs = append(trs, v)
	return dsl.Join(trs...)
}
