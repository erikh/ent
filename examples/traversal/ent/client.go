// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/facebookincubator/ent/examples/traversal/ent/migrate"

	"github.com/facebookincubator/ent/examples/traversal/ent/group"
	"github.com/facebookincubator/ent/examples/traversal/ent/pet"
	"github.com/facebookincubator/ent/examples/traversal/ent/user"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Group is the client for interacting with the Group builders.
	Group *GroupClient
	// Pet is the client for interacting with the Pet builders.
	Pet *PetClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	c := config{log: log.Println}
	c.options(opts...)
	return &Client{
		config: c,
		Schema: migrate.NewSchema(c.driver),
		Group:  NewGroupClient(c),
		Pet:    NewPetClient(c),
		User:   NewUserClient(c),
	}
}

// Open opens a connection to the database specified by the driver name and a
// driver-specific data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil

	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug}
	return &Tx{
		config: cfg,
		Group:  NewGroupClient(cfg),
		Pet:    NewPetClient(cfg),
		User:   NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Group.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true}
	return &Client{
		config: cfg,
		Schema: migrate.NewSchema(cfg.driver),
		Group:  NewGroupClient(cfg),
		Pet:    NewPetClient(cfg),
		User:   NewUserClient(cfg),
	}
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// GroupClient is a client for the Group schema.
type GroupClient struct {
	config
}

// NewGroupClient returns a client for the Group from the given config.
func NewGroupClient(c config) *GroupClient {
	return &GroupClient{config: c}
}

// Create returns a create builder for Group.
func (c *GroupClient) Create() *GroupCreate {
	return &GroupCreate{config: c.config}
}

// Update returns an update builder for Group.
func (c *GroupClient) Update() *GroupUpdate {
	return &GroupUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *GroupClient) UpdateOne(gr *Group) *GroupUpdateOne {
	return c.UpdateOneID(gr.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *GroupClient) UpdateOneID(id int) *GroupUpdateOne {
	return &GroupUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Group.
func (c *GroupClient) Delete() *GroupDelete {
	return &GroupDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *GroupClient) DeleteOne(gr *Group) *GroupDeleteOne {
	return c.DeleteOneID(gr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *GroupClient) DeleteOneID(id int) *GroupDeleteOne {
	return &GroupDeleteOne{c.Delete().Where(group.ID(id))}
}

// Create returns a query builder for Group.
func (c *GroupClient) Query() *GroupQuery {
	return &GroupQuery{config: c.config}
}

// Get returns a Group entity by its id.
func (c *GroupClient) Get(ctx context.Context, id int) (*Group, error) {
	return c.Query().Where(group.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *GroupClient) GetX(ctx context.Context, id int) *Group {
	gr, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return gr
}

// QueryUsers queries the users edge of a Group.
func (c *GroupClient) QueryUsers(gr *Group) *UserQuery {
	query := &UserQuery{config: c.config}
	id := gr.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = group.Table
	step.From.Column = group.FieldID
	step.To.Table = user.Table
	step.To.Column = user.FieldID
	step.Edge.Rel = sql.M2M
	step.Edge.Inverse = false
	step.Edge.Table = group.UsersTable
	step.Edge.Columns = append(step.Edge.Columns, group.UsersPrimaryKey...)
	query.sql = sql.Neighbors(gr.driver.Dialect(), step)

	return query
}

// QueryAdmin queries the admin edge of a Group.
func (c *GroupClient) QueryAdmin(gr *Group) *UserQuery {
	query := &UserQuery{config: c.config}
	id := gr.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = group.Table
	step.From.Column = group.FieldID
	step.To.Table = user.Table
	step.To.Column = user.FieldID
	step.Edge.Rel = sql.M2O
	step.Edge.Inverse = false
	step.Edge.Table = group.AdminTable
	step.Edge.Columns = append(step.Edge.Columns, group.AdminColumn)
	query.sql = sql.Neighbors(gr.driver.Dialect(), step)

	return query
}

// PetClient is a client for the Pet schema.
type PetClient struct {
	config
}

// NewPetClient returns a client for the Pet from the given config.
func NewPetClient(c config) *PetClient {
	return &PetClient{config: c}
}

// Create returns a create builder for Pet.
func (c *PetClient) Create() *PetCreate {
	return &PetCreate{config: c.config}
}

// Update returns an update builder for Pet.
func (c *PetClient) Update() *PetUpdate {
	return &PetUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *PetClient) UpdateOne(pe *Pet) *PetUpdateOne {
	return c.UpdateOneID(pe.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PetClient) UpdateOneID(id int) *PetUpdateOne {
	return &PetUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Pet.
func (c *PetClient) Delete() *PetDelete {
	return &PetDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PetClient) DeleteOne(pe *Pet) *PetDeleteOne {
	return c.DeleteOneID(pe.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PetClient) DeleteOneID(id int) *PetDeleteOne {
	return &PetDeleteOne{c.Delete().Where(pet.ID(id))}
}

// Create returns a query builder for Pet.
func (c *PetClient) Query() *PetQuery {
	return &PetQuery{config: c.config}
}

// Get returns a Pet entity by its id.
func (c *PetClient) Get(ctx context.Context, id int) (*Pet, error) {
	return c.Query().Where(pet.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PetClient) GetX(ctx context.Context, id int) *Pet {
	pe, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pe
}

// QueryFriends queries the friends edge of a Pet.
func (c *PetClient) QueryFriends(pe *Pet) *PetQuery {
	query := &PetQuery{config: c.config}
	id := pe.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = pet.Table
	step.From.Column = pet.FieldID
	step.To.Table = pet.Table
	step.To.Column = pet.FieldID
	step.Edge.Rel = sql.M2M
	step.Edge.Inverse = false
	step.Edge.Table = pet.FriendsTable
	step.Edge.Columns = append(step.Edge.Columns, pet.FriendsPrimaryKey...)
	query.sql = sql.Neighbors(pe.driver.Dialect(), step)

	return query
}

// QueryOwner queries the owner edge of a Pet.
func (c *PetClient) QueryOwner(pe *Pet) *UserQuery {
	query := &UserQuery{config: c.config}
	id := pe.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = pet.Table
	step.From.Column = pet.FieldID
	step.To.Table = user.Table
	step.To.Column = user.FieldID
	step.Edge.Rel = sql.M2O
	step.Edge.Inverse = true
	step.Edge.Table = pet.OwnerTable
	step.Edge.Columns = append(step.Edge.Columns, pet.OwnerColumn)
	query.sql = sql.Neighbors(pe.driver.Dialect(), step)

	return query
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	return &UserCreate{config: c.config}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	return &UserUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	return c.UpdateOneID(u.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	return &UserUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	return &UserDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	return &UserDeleteOne{c.Delete().Where(user.ID(id))}
}

// Create returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{config: c.config}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	u, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return u
}

// QueryPets queries the pets edge of a User.
func (c *UserClient) QueryPets(u *User) *PetQuery {
	query := &PetQuery{config: c.config}
	id := u.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = user.Table
	step.From.Column = user.FieldID
	step.To.Table = pet.Table
	step.To.Column = pet.FieldID
	step.Edge.Rel = sql.O2M
	step.Edge.Inverse = false
	step.Edge.Table = user.PetsTable
	step.Edge.Columns = append(step.Edge.Columns, user.PetsColumn)
	query.sql = sql.Neighbors(u.driver.Dialect(), step)

	return query
}

// QueryFriends queries the friends edge of a User.
func (c *UserClient) QueryFriends(u *User) *UserQuery {
	query := &UserQuery{config: c.config}
	id := u.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = user.Table
	step.From.Column = user.FieldID
	step.To.Table = user.Table
	step.To.Column = user.FieldID
	step.Edge.Rel = sql.M2M
	step.Edge.Inverse = false
	step.Edge.Table = user.FriendsTable
	step.Edge.Columns = append(step.Edge.Columns, user.FriendsPrimaryKey...)
	query.sql = sql.Neighbors(u.driver.Dialect(), step)

	return query
}

// QueryGroups queries the groups edge of a User.
func (c *UserClient) QueryGroups(u *User) *GroupQuery {
	query := &GroupQuery{config: c.config}
	id := u.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = user.Table
	step.From.Column = user.FieldID
	step.To.Table = group.Table
	step.To.Column = group.FieldID
	step.Edge.Rel = sql.M2M
	step.Edge.Inverse = true
	step.Edge.Table = user.GroupsTable
	step.Edge.Columns = append(step.Edge.Columns, user.GroupsPrimaryKey...)
	query.sql = sql.Neighbors(u.driver.Dialect(), step)

	return query
}

// QueryManage queries the manage edge of a User.
func (c *UserClient) QueryManage(u *User) *GroupQuery {
	query := &GroupQuery{config: c.config}
	id := u.ID
	step := &sql.Step{}
	step.From.V = id
	step.From.Table = user.Table
	step.From.Column = user.FieldID
	step.To.Table = group.Table
	step.To.Column = group.FieldID
	step.Edge.Rel = sql.O2M
	step.Edge.Inverse = true
	step.Edge.Table = user.ManageTable
	step.Edge.Columns = append(step.Edge.Columns, user.ManageColumn)
	query.sql = sql.Neighbors(u.driver.Dialect(), step)

	return query
}
