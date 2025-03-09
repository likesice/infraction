package store

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// no validation tests. These should happen earlier in the call stack
func TestInsertGroup(t *testing.T) {
	storeReg := NewStore(db)

	// arrange
	now = func() time.Time { return time.Unix(int64(1), int64(1)) }
	query := `SELECT user_id FROM groups_users WHERE group_id = ?;`
	var testUser User

	err := db.QueryRow("INSERT INTO users (name) VALUES (?) RETURNING id, name;", "test_user_name").Scan(
		&testUser.Id,
		&testUser.Name,
	)
	if err != nil {
		t.Fail()
	}

	tests := map[string]struct {
		group          Group
		expectedValues []any
	}{
		"group with member": {group: Group{Members: []int64{testUser.Id}}, expectedValues: []any{nil, 1, 1000, 1000, 1}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			err := storeReg.GroupStore.Insert(&tc.group)
			// assert
			if tc.expectedValues[0] != nil {
				if !errors.Is(err, tc.expectedValues[0].(error)) {
					t.Fatalf("%v test case failed because of: %v", name, err)
				}
			}
			if tc.group.Id != int64(tc.expectedValues[1].(int)) ||
				tc.group.CreatedAt != int64(tc.expectedValues[2].(int)) ||
				tc.group.UpdatedAt != int64(tc.expectedValues[3].(int)) ||
				tc.group.Version != int64(tc.expectedValues[4].(int)) {
				t.Fatalf("%v test case failed because of wrong group value", name)
			}

			var userId int64
			err = db.QueryRow(query, tc.group.Id).Scan(&userId)

			if err != nil || userId != tc.group.Members[0] {
				t.Fatalf("%v test case failed because of groups_users mapping", name)
			}

		})
	}
}

// no validation tests. These should happen earlier in the call stack
func TestSelectAll(t *testing.T) {
	storeReg := NewStore(db)

	// arrange
	now = func() time.Time { return time.Unix(int64(1), int64(1)) }
	testUsers := []*User{{}, {}}

	for i, user := range testUsers {
		err := db.QueryRow("INSERT INTO users (name) VALUES (?) RETURNING id, name;",
			fmt.Sprintf("test_user_name_%d", i)).Scan(
			&user.Id,
			&user.Name,
		)
		if err != nil {
			t.Fail()
		}
	}

	testGroups := []*Group{{Members: []int64{testUsers[0].Id}},
		{Members: []int64{testUsers[1].Id}},
		{Members: []int64{testUsers[1].Id}}}

	for _, group := range testGroups {
		err := storeReg.GroupStore.Insert(group)
		if err != nil {
			t.Fail()
		}
	}

	tests := map[string]struct {
		user           *User
		expectedValues []any
	}{
		"single group":    {user: testUsers[0], expectedValues: []any{nil, 1, []*Group{testGroups[0]}}},
		"multiple groups": {user: testUsers[1], expectedValues: []any{nil, 2, []*Group{testGroups[1], testGroups[2]}}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			groups, err := storeReg.GroupStore.SelectAll(tc.user)
			// assert
			if tc.expectedValues[0] != nil {
				if !errors.Is(err, tc.expectedValues[0].(error)) {
					t.Fatalf("%v test case failed because of: %v", name, err)
				}
			}
			if len(groups) != tc.expectedValues[1] {
				t.Fatalf("%v test case failed because mismatche group length", name)
			}
		})
	}
}
