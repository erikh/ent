// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sql

import (
	"context"
	"database/sql/driver"
	"regexp"
	"strings"
	"testing"

	"github.com/facebookincubator/ent/schema/field"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestNeighbors(t *testing.T) {
	tests := []struct {
		name      string
		input     *Step
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "O2O/1type",
			// Since the relation is on the same table,
			// V used as a reference value.
			input: NewStep(
				From("users", "id", 1),
				To("users", "id"),
				Edge(O2O, false, "users", "spouse_id"),
			),
			wantQuery: "SELECT * FROM `users` WHERE `spouse_id` = ?",
			wantArgs:  []interface{}{1},
		},
		{
			name: "O2O/1type/inverse",
			input: NewStep(
				From("nodes", "id", 1),
				To("nodes", "id"),
				Edge(O2O, true, "nodes", "prev_id"),
			),
			wantQuery: "SELECT * FROM `nodes` JOIN (SELECT `prev_id` FROM `nodes` WHERE `id` = ?) AS `t1` ON `nodes`.`id` = `t1`.`prev_id`",
			wantArgs:  []interface{}{1},
		},
		{
			name: "O2M/1type",
			input: NewStep(
				From("users", "id", 1),
				To("users", "id"),
				Edge(O2M, false, "users", "parent_id"),
			),
			wantQuery: "SELECT * FROM `users` WHERE `parent_id` = ?",
			wantArgs:  []interface{}{1},
		},
		{
			name: "O2O/2types",
			input: NewStep(
				From("users", "id", 2),
				To("card", "id"),
				Edge(O2O, false, "cards", "owner_id"),
			),
			wantQuery: "SELECT * FROM `card` WHERE `owner_id` = ?",
			wantArgs:  []interface{}{2},
		},
		{
			name: "O2O/2types/inverse",
			input: NewStep(
				From("cards", "id", 2),
				To("users", "id"),
				Edge(O2O, true, "cards", "owner_id"),
			),
			wantQuery: "SELECT * FROM `users` JOIN (SELECT `owner_id` FROM `cards` WHERE `id` = ?) AS `t1` ON `users`.`id` = `t1`.`owner_id`",
			wantArgs:  []interface{}{2},
		},
		{
			name: "O2M/2types",
			input: NewStep(
				From("users", "id", 1),
				To("pets", "id"),
				Edge(O2M, false, "pets", "owner_id"),
			),
			wantQuery: "SELECT * FROM `pets` WHERE `owner_id` = ?",
			wantArgs:  []interface{}{1},
		},
		{
			name: "M2O/2types/inverse",
			input: NewStep(
				From("pets", "id", 2),
				To("users", "id"),
				Edge(M2O, true, "pets", "owner_id"),
			),
			wantQuery: "SELECT * FROM `users` JOIN (SELECT `owner_id` FROM `pets` WHERE `id` = ?) AS `t1` ON `users`.`id` = `t1`.`owner_id`",
			wantArgs:  []interface{}{2},
		},
		{
			name: "M2O/1type/inverse",
			input: NewStep(
				From("users", "id", 2),
				To("users", "id"),
				Edge(M2O, true, "users", "parent_id"),
			),
			wantQuery: "SELECT * FROM `users` JOIN (SELECT `parent_id` FROM `users` WHERE `id` = ?) AS `t1` ON `users`.`id` = `t1`.`parent_id`",
			wantArgs:  []interface{}{2},
		},
		{
			name: "M2M/2type",
			input: NewStep(
				From("groups", "id", 2),
				To("users", "id"),
				Edge(M2M, false, "user_groups", "group_id", "user_id"),
			),
			wantQuery: "SELECT * FROM `users` JOIN (SELECT `user_groups`.`user_id` FROM `user_groups` WHERE `user_groups`.`group_id` = ?) AS `t1` ON `users`.`id` = `t1`.`user_id`",
			wantArgs:  []interface{}{2},
		},
		{
			name: "M2M/2type/inverse",
			input: NewStep(
				From("users", "id", 2),
				To("groups", "id"),
				Edge(M2M, true, "user_groups", "group_id", "user_id"),
			),
			wantQuery: "SELECT * FROM `groups` JOIN (SELECT `user_groups`.`group_id` FROM `user_groups` WHERE `user_groups`.`user_id` = ?) AS `t1` ON `groups`.`id` = `t1`.`group_id`",
			wantArgs:  []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := Neighbors("", tt.input)
			query, args := selector.Query()
			require.Equal(t, tt.wantQuery, query)
			require.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestSetNeighbors(t *testing.T) {
	tests := []struct {
		name      string
		input     *Step
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "O2M/2types",
			input: NewStep(
				From("users", "id", Select().From(Table("users")).Where(EQ("name", "a8m"))),
				To("pets", "id"),
				Edge(O2M, false, "users", "owner_id"),
			),
			wantQuery: `SELECT * FROM "pets" JOIN (SELECT "users"."id" FROM "users" WHERE "name" = $1) AS "t1" ON "pets"."owner_id" = "t1"."id"`,
			wantArgs:  []interface{}{"a8m"},
		},
		{
			name: "M2O/2types",
			input: NewStep(
				From("pets", "id", Select().From(Table("pets")).Where(EQ("name", "pedro"))),
				To("users", "id"),
				Edge(M2O, true, "pets", "owner_id"),
			),
			wantQuery: `SELECT * FROM "users" JOIN (SELECT "pets"."owner_id" FROM "pets" WHERE "name" = $1) AS "t1" ON "users"."id" = "t1"."owner_id"`,
			wantArgs:  []interface{}{"pedro"},
		},
		{
			name: "M2M/2types",
			input: NewStep(
				From("users", "id", Select().From(Table("users")).Where(EQ("name", "a8m"))),
				To("groups", "id"),
				Edge(M2M, false, "user_groups", "user_id", "group_id"),
			),
			wantQuery: `
SELECT *
FROM "groups"
JOIN
  (SELECT "user_groups"."group_id"
   FROM "user_groups"
   JOIN
     (SELECT "users"."id"
      FROM "users"
      WHERE "name" = $1) AS "t1" ON "user_groups"."user_id" = "t1"."id") AS "t1" ON "groups"."id" = "t1"."group_id"`,
			wantArgs: []interface{}{"a8m"},
		},
		{
			name: "M2M/2types/inverse",
			input: NewStep(
				From("groups", "id", Select().From(Table("groups")).Where(EQ("name", "GitHub"))),
				To("users", "id"),
				Edge(M2M, true, "user_groups", "user_id", "group_id"),
			),
			wantQuery: `
SELECT *
FROM "users"
JOIN
  (SELECT "user_groups"."user_id"
   FROM "user_groups"
   JOIN
     (SELECT "groups"."id"
      FROM "groups"
      WHERE "name" = $1) AS "t1" ON "user_groups"."group_id" = "t1"."id") AS "t1" ON "users"."id" = "t1"."user_id"`,
			wantArgs: []interface{}{"GitHub"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := SetNeighbors("postgres", tt.input)
			query, args := selector.Query()
			tt.wantQuery = strings.Join(strings.Fields(tt.wantQuery), " ")
			require.Equal(t, tt.wantQuery, query)
			require.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestHasNeighbors(t *testing.T) {
	tests := []struct {
		name      string
		step      *Step
		selector  *Selector
		wantQuery string
	}{
		{
			name: "O2O/1type",
			// A nodes table; linked-list (next->prev). The "prev"
			// node holds association pointer. The neighbors query
			// here checks if a node "has-next".
			step: NewStep(
				From("nodes", "id"),
				To("nodes", "id"),
				Edge(O2O, false, "nodes", "prev_id"),
			),
			selector:  Select("*").From(Table("nodes")),
			wantQuery: "SELECT * FROM `nodes` WHERE `nodes`.`id` IN (SELECT `nodes`.`prev_id` FROM `nodes` WHERE `nodes`.`prev_id` IS NOT NULL)",
		},
		{
			name: "O2O/1type/inverse",
			// Same example as above, but the neighbors
			// query checks if a node "has-previous".
			step: NewStep(
				From("nodes", "id"),
				To("nodes", "id"),
				Edge(O2O, true, "nodes", "prev_id"),
			),
			selector:  Select("*").From(Table("nodes")),
			wantQuery: "SELECT * FROM `nodes` WHERE `nodes`.`prev_id` IS NOT NULL",
		},
		{
			name: "O2M/2type2",
			step: NewStep(
				From("users", "id"),
				To("pets", "id"),
				Edge(O2M, false, "pets", "owner_id"),
			),
			selector:  Select("*").From(Table("users")),
			wantQuery: "SELECT * FROM `users` WHERE `users`.`id` IN (SELECT `pets`.`owner_id` FROM `pets` WHERE `pets`.`owner_id` IS NOT NULL)",
		},
		{
			name: "M2O/2type2",
			step: NewStep(
				From("pets", "id"),
				To("users", "id"),
				Edge(M2O, true, "pets", "owner_id"),
			),
			selector:  Select("*").From(Table("pets")),
			wantQuery: "SELECT * FROM `pets` WHERE `pets`.`owner_id` IS NOT NULL",
		},
		{
			name: "M2M/2types",
			step: NewStep(
				From("users", "id"),
				To("groups", "id"),
				Edge(M2M, false, "user_groups", "user_id", "group_id"),
			),
			selector:  Select("*").From(Table("users")),
			wantQuery: "SELECT * FROM `users` WHERE `users`.`id` IN (SELECT `user_groups`.`user_id` FROM `user_groups`)",
		},
		{
			name: "M2M/2types/inverse",
			step: NewStep(
				From("users", "id"),
				To("groups", "id"),
				Edge(M2M, true, "group_users", "group_id", "user_id"),
			),
			selector:  Select("*").From(Table("users")),
			wantQuery: "SELECT * FROM `users` WHERE `users`.`id` IN (SELECT `group_users`.`user_id` FROM `group_users`)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HasNeighbors(tt.selector, tt.step)
			query, args := tt.selector.Query()
			require.Equal(t, tt.wantQuery, query)
			require.Empty(t, args)
		})
	}
}

func TestHasNeighborsWith(t *testing.T) {
	tests := []struct {
		name      string
		step      *Step
		selector  *Selector
		predicate func(*Selector)
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "O2O",
			step: NewStep(
				From("users", "id"),
				To("cards", "id"),
				Edge(O2O, false, "cards", "owner_id"),
			),
			selector: Dialect("postgres").Select("*").From(Table("users")),
			predicate: func(s *Selector) {
				s.Where(EQ("expired", false))
			},
			wantQuery: `SELECT * FROM "users" WHERE "users"."id" IN (SELECT "cards"."owner_id" FROM "cards" WHERE "expired" = $1)`,
			wantArgs:  []interface{}{false},
		},
		{
			name: "O2O/inverse",
			step: NewStep(
				From("cards", "id"),
				To("users", "id"),
				Edge(O2O, true, "cards", "owner_id"),
			),
			selector: Dialect("postgres").Select("*").From(Table("cards")),
			predicate: func(s *Selector) {
				s.Where(EQ("name", "a8m"))
			},
			wantQuery: `SELECT * FROM "cards" WHERE "cards"."owner_id" IN (SELECT "users"."id" FROM "users" WHERE "name" = $1)`,
			wantArgs:  []interface{}{"a8m"},
		},
		{
			name: "O2M",
			step: NewStep(
				From("users", "id"),
				To("pets", "id"),
				Edge(O2M, false, "pets", "owner_id"),
			),
			selector: Dialect("postgres").Select("*").
				From(Table("users")).
				Where(EQ("last_name", "mashraki")),
			predicate: func(s *Selector) {
				s.Where(EQ("name", "pedro"))
			},
			wantQuery: `SELECT * FROM "users" WHERE "last_name" = $1 AND "users"."id" IN (SELECT "pets"."owner_id" FROM "pets" WHERE "name" = $2)`,
			wantArgs:  []interface{}{"mashraki", "pedro"},
		},
		{
			name: "M2O",
			step: NewStep(
				From("pets", "id"),
				To("users", "id"),
				Edge(M2O, true, "pets", "owner_id"),
			),
			selector: Dialect("postgres").Select("*").
				From(Table("pets")).
				Where(EQ("name", "pedro")),
			predicate: func(s *Selector) {
				s.Where(EQ("last_name", "mashraki"))
			},
			wantQuery: `SELECT * FROM "pets" WHERE "name" = $1 AND "pets"."owner_id" IN (SELECT "users"."id" FROM "users" WHERE "last_name" = $2)`,
			wantArgs:  []interface{}{"pedro", "mashraki"},
		},
		{
			name: "M2M",
			step: NewStep(
				From("users", "id"),
				To("groups", "id"),
				Edge(M2M, false, "user_groups", "user_id", "group_id"),
			),
			selector: Dialect("postgres").Select("*").From(Table("users")),
			predicate: func(s *Selector) {
				s.Where(EQ("name", "GitHub"))
			},
			wantQuery: `
SELECT *
FROM "users"
WHERE "users"."id" IN
  (SELECT "user_groups"."user_id"
  FROM "user_groups"
  JOIN "groups" AS "t0" ON "user_groups"."group_id" = "t0"."id" WHERE "name" = $1)`,
			wantArgs: []interface{}{"GitHub"},
		},
		{
			name: "M2M/inverse",
			step: NewStep(
				From("groups", "id"),
				To("users", "id"),
				Edge(M2M, true, "user_groups", "user_id", "group_id"),
			),
			selector: Dialect("postgres").Select("*").From(Table("groups")),
			predicate: func(s *Selector) {
				s.Where(EQ("name", "a8m"))
			},
			wantQuery: `
SELECT *
FROM "groups"
WHERE "groups"."id" IN
  (SELECT "user_groups"."group_id"
  FROM "user_groups"
  JOIN "users" AS "t0" ON "user_groups"."user_id" = "t0"."id" WHERE "name" = $1)`,
			wantArgs: []interface{}{"a8m"},
		},
		{
			name: "M2M/inverse",
			step: NewStep(
				From("groups", "id"),
				To("users", "id"),
				Edge(M2M, true, "user_groups", "user_id", "group_id"),
			),
			selector: Dialect("postgres").Select("*").From(Table("groups")),
			predicate: func(s *Selector) {
				s.Where(And(NotNull("name"), EQ("name", "a8m")))
			},
			wantQuery: `
SELECT *
FROM "groups"
WHERE "groups"."id" IN
  (SELECT "user_groups"."group_id"
  FROM "user_groups"
  JOIN "users" AS "t0" ON "user_groups"."user_id" = "t0"."id" WHERE ("name" IS NOT NULL) AND ("name" = $1))`,
			wantArgs: []interface{}{"a8m"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HasNeighborsWith(tt.selector, tt.step, tt.predicate)
			query, args := tt.selector.Query()
			tt.wantQuery = strings.Join(strings.Fields(tt.wantQuery), " ")
			require.Equal(t, tt.wantQuery, query)
			require.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestCreateNode(t *testing.T) {
	tests := []struct {
		name    string
		spec    *CreateSpec
		expect  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "fields",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "age", Type: field.TypeInt, Value: 30},
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`age`, `name`) VALUES (?, ?)")).
					WithArgs(30, "a8m").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "fields/user-defined-id",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id", Value: 1},
				Fields: []*FieldSpec{
					{Column: "age", Type: field.TypeInt, Value: 30},
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`age`, `name`, `id`) VALUES (?, ?, ?)")).
					WithArgs(30, "a8m", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "fields/json",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "json", Type: field.TypeJSON, Value: struct{}{}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`json`) VALUES (?)")).
					WithArgs([]byte("{}")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/m2o",
			spec: &CreateSpec{
				Table: "pets",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "pedro"},
				},
				Edges: []*EdgeSpec{
					{Rel: M2O, Columns: []string{"owner_id"}, Inverse: true, Target: &EdgeTarget{Nodes: []driver.Value{2}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `pets` (`name`, `owner_id`) VALUES (?, ?)")).
					WithArgs("pedro", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/o2o/inverse",
			spec: &CreateSpec{
				Table: "cards",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "number", Type: field.TypeString, Value: "0001"},
				},
				Edges: []*EdgeSpec{
					{Rel: O2O, Columns: []string{"owner_id"}, Inverse: true, Target: &EdgeTarget{Nodes: []driver.Value{2}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `cards` (`number`, `owner_id`) VALUES (?, ?)")).
					WithArgs("0001", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/o2m",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
				Edges: []*EdgeSpec{
					{Rel: O2M, Table: "pets", Columns: []string{"owner_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("a8m").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("UPDATE `pets` SET `owner_id` = ? WHERE (`id` = ?) AND (`owner_id` IS NULL)")).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/o2m",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
				Edges: []*EdgeSpec{
					{Rel: O2M, Table: "pets", Columns: []string{"owner_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2, 3, 4}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("a8m").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("UPDATE `pets` SET `owner_id` = ? WHERE (`id` IN (?, ?, ?)) AND (`owner_id` IS NULL)")).
					WithArgs(1, 2, 3, 4).
					WillReturnResult(sqlmock.NewResult(1, 3))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/o2o",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
				Edges: []*EdgeSpec{
					{Rel: O2O, Table: "cards", Columns: []string{"owner_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("a8m").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("UPDATE `cards` SET `owner_id` = ? WHERE (`id` = ?) AND (`owner_id` IS NULL)")).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/o2o/bidi",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "a8m"},
				},
				Edges: []*EdgeSpec{
					{Rel: O2O, Bidi: true, Table: "users", Columns: []string{"spouse_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`, `spouse_id`) VALUES (?, ?)")).
					WithArgs("a8m", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("UPDATE `users` SET `spouse_id` = ? WHERE (`id` = ?) AND (`spouse_id` IS NULL)")).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/m2m",
			spec: &CreateSpec{
				Table: "groups",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "GitHub"},
				},
				Edges: []*EdgeSpec{
					{Rel: M2M, Table: "group_users", Columns: []string{"group_id", "user_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `groups` (`name`) VALUES (?)")).
					WithArgs("GitHub").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("INSERT INTO `group_users` (`group_id`, `user_id`) VALUES (?, ?)")).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/m2m/inverse",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "mashraki"},
				},
				Edges: []*EdgeSpec{
					{Rel: M2M, Inverse: true, Table: "group_users", Columns: []string{"group_id", "user_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("mashraki").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("INSERT INTO `group_users` (`group_id`, `user_id`) VALUES (?, ?)")).
					WithArgs(2, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/m2m/bidi",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "mashraki"},
				},
				Edges: []*EdgeSpec{
					{Rel: M2M, Bidi: true, Table: "user_friends", Columns: []string{"user_id", "friend_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("mashraki").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("INSERT INTO `user_friends` (`user_id`, `friend_id`) VALUES (?, ?), (?, ?)")).
					WithArgs(1, 2, 2, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "edges/m2m/bidi/batch",
			spec: &CreateSpec{
				Table: "users",
				ID:    &FieldSpec{Column: "id"},
				Fields: []*FieldSpec{
					{Column: "name", Type: field.TypeString, Value: "mashraki"},
				},
				Edges: []*EdgeSpec{
					{Rel: M2M, Bidi: true, Table: "user_friends", Columns: []string{"user_id", "friend_id"}, Target: &EdgeTarget{Nodes: []driver.Value{2}, IDSpec: &FieldSpec{Column: "id"}}},
					{Rel: M2M, Bidi: true, Table: "user_friends", Columns: []string{"user_id", "friend_id"}, Target: &EdgeTarget{Nodes: []driver.Value{3}, IDSpec: &FieldSpec{Column: "id"}}},
					{Rel: M2M, Inverse: true, Table: "group_users", Columns: []string{"group_id", "user_id"}, Target: &EdgeTarget{Nodes: []driver.Value{4}, IDSpec: &FieldSpec{Column: "id"}}},
					{Rel: M2M, Inverse: true, Table: "group_users", Columns: []string{"group_id", "user_id"}, Target: &EdgeTarget{Nodes: []driver.Value{5}, IDSpec: &FieldSpec{Column: "id"}}},
				},
			},
			expect: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(escape("INSERT INTO `users` (`name`) VALUES (?)")).
					WithArgs("mashraki").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("INSERT INTO `user_friends` (`user_id`, `friend_id`) VALUES (?, ?), (?, ?), (?, ?), (?, ?)")).
					WithArgs(1, 2, 2, 1, 1, 3, 3, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(escape("INSERT INTO `group_users` (`group_id`, `user_id`) VALUES (?, ?), (?, ?)")).
					WithArgs(4, 1, 5, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			tt.expect(mock)
			err = CreateNode(context.Background(), OpenDB("", db), tt.spec)
			require.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}

func escape(query string) string {
	rows := strings.Split(query, "\n")
	for i := range rows {
		rows[i] = strings.TrimPrefix(rows[i], " ")
	}
	query = strings.Join(rows, " ")
	return regexp.QuoteMeta(query)
}
