// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/google/uuid"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleBlob() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the blob's edges.

	// create blob vertex with its edges.
	b := client.Blob.
		Create().
		SetUUID(uuid.UUID{}).
		SaveX(ctx)
	log.Println("blob created:", b)

	// query edges.

	// Output:
}
func ExampleGroup() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the group's edges.
	u0 := client.User.
		Create().
		SaveX(ctx)
	log.Println("user created:", u0)

	// create group vertex with its edges.
	gr := client.Group.
		Create().
		AddUsers(u0).
		SaveX(ctx)
	log.Println("group created:", gr)

	// query edges.
	u0, err = gr.QueryUsers().First(ctx)
	if err != nil {
		log.Fatalf("failed querying users: %v", err)
	}
	log.Println("users found:", u0)

	// Output:
}
func ExampleUser() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the user's edges.

	// create user vertex with its edges.
	u := client.User.
		Create().
		SaveX(ctx)
	log.Println("user created:", u)

	// query edges.

	// Output:
}
