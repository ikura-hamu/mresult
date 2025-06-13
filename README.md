# mresult

Go library to simplify writing test mock return values.

[![Go Reference](https://pkg.go.dev/badge/github.com/ikura-hamu/mresult.svg)](https://pkg.go.dev/github.com/ikura-hamu/mresult)

## Overview

Package mresult provides utilities for simplifying mock return values in Go tests. This package helps manage mock results in a clean and type-safe way, reducing boilerplate code when writing table-driven tests with mocks.

## Features

- **Type-safe mock results**: Generic types ensure compile-time safety
- **Simple API**: Easy-to-use generator functions for creating mock results
- **Multiple return values**: Support for functions returning 0, 1, or 2 values plus error
- **Execution tracking**: Built-in checks to ensure mock results are properly configured

## Types

The package provides three main result types:

- **[`MResult0`](https://pkg.go.dev/github.com/ikura-hamu/mresult#MResult0)**: For functions returning only an error
- **[`MResult[T]`](https://pkg.go.dev/github.com/ikura-hamu/mresult#MResult)**: For functions returning a value and an error
- **[`MResult2[T1, T2]`](https://pkg.go.dev/github.com/ikura-hamu/mresult#MResult2)**: For functions returning two values and an error

## Installation

```bash
go get github.com/ikura-hamu/mresult
```

## Quick Start

```go
import "github.com/ikura-hamu/mresult"

// Create generators for mock results
getUserR, getUserRErr := mresult.Generator[*User](t)

// Use in test cases
testCases := map[string]struct {
    getUserResult mresult.MResult[*User]
}{
    "success": {
        getUserResult: getUserR(t, &User{ID: 1}),
    },
    "error": {
        getUserResult: getUserRErr(t, errors.New("user not found")),
    },
}
```

## Examples

<details>
<summary>Without mresult</summary>

```go
package user_test

import (
	"errors"
	"testing"

	"yourproject/.../user"

	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Create a mock for the UserRepository interface
	mockRepo := NewMockUserRepository(ctrl)

	testCases := map[string]struct {
		args              user.CreateUserArgs
		executeGetUser    bool // Too many fields for a single test case to manage mock result.
		GetUserUser       *User
		GetUserError      error
		executeCreateUser bool // If you want to execute multiple methods, you need to add more fields!
		CreateUserUser    *User
		CreateUserError   error
		expectedError     error
	}{
		"success": {
			args: user.CreateUserArgs{
				Username: "testuser",
				Age:      30,
			},
			executeGetUser:    true,
			GetUserUser:       &User{ID: 1, Username: "testuser"},
			GetUserError:      nil,
			executeCreateUser: true,
			CreateUserUser:    &User{ID: 1, Username: "testuser"},
			CreateUserError:   nil,
		},
		"user already exists": {
			args: user.CreateUserArgs{
				Username: "existinguser",
				Age:      25,
			},
			executeGetUser:    true,
			GetUserUser:       &User{ID: 2, Username: "existinguser"},
			GetUserError:      nil,
			executeCreateUser: false,
		},
		"GetUser error": {
			args: user.CreateUserArgs{
				Username: "erroruser",
				Age:      20,
			},
			executeGetUser:    true,
			GetUserUser:       nil,
			GetUserError:      user.ErrDatabase,
			executeCreateUser: false,
		},
		"invalid age": {
			args: user.CreateUserArgs{
				Username: "invaliduser",
				Age:      -1,
			},
			executeGetUser: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// Set up expectations
			if testCase.executeGetUser {
				mockRepo.EXPECT().GetUser(gomock.Any()).Return(testCase.GetUserUser, testCase.GetUserError)
			}
			if testCase.executeCreateUser {
				mockRepo.EXPECT().Create(gomock.Any()).Return(testCase.CreateUserUser, testCase.CreateUserError)
			}

			// Call the function under test
			err := CreateUser(mockRepo, "testuser")

			// Assert that there were no errors
			if testCase.expectedError != nil {
				if !errors.Is(err, testCase.expectedError) {
					t.Errorf("expected error %v, got %v", testCase.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
```

</details>

<details>
<summary>With mresult</summary>

```go
package user_test

import (
	"errors"
	"testing"

	"yourproject/.../user"

	"github.com/ikura-hamu/mresult"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Create a mock for the UserRepository interface
	mockRepo := NewMockUserRepository(ctrl)

	// Use mresult to prepare mock results for GetUser and CreateUser methods
	getUserR, getUserRErr := mresult.Generator[*User](t)
	createUserR, createUserRErr := mresult.Generator[*User](t)

	testCases := map[string]struct {
		args             user.CreateUserArgs
		GetUserResult    mresult.MResult[*User] // Use mresult.MResult to handle mock results
		CreateUserResult mresult.MResult[*User] // Simplifies the test case structure
		expectedError    error
	}{
		"success": {
			args: user.CreateUserArgs{
				Username: "testuser",
				Age:      30,
			},
			GetUserResult:    getUserR(t, &User{ID: 1, Username: "testuser"}), // Mock GetUser to return a user
			CreateUserResult: createUserR(t, &User{ID: 1, Username: "testuser"}), // Mock CreateUser to return the user
		},
		"user already exists": {
			args: user.CreateUserArgs{
				Username: "existinguser",
				Age:      25,
			},
			GetUserResult: getUserR(t, &User{ID: 2, Username: "existinguser"}),
			// CreateUserResult not set - method won't be called
		},
		"GetUser error": {
			args: user.CreateUserArgs{
				Username: "erroruser",
				Age:      20,
			},
			GetUserResult: getUserRErr(t, user.ErrDatabase),
			// CreateUserResult not set - method won't be called
		},
		"invalid age": {
			args: user.CreateUserArgs{
				Username: "invaliduser",
				Age:      -1,
			},
			// Neither result set - no methods will be called
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// Set up expectations using mresult
			if testCase.GetUserResult.IsExecuted(t) { // Check if the result is configured
				mockRepo.EXPECT().GetUser(gomock.Any()).
					Return(testCase.GetUserResult.Val(t), testCase.GetUserResult.Err(t)) // Get value and error
			}
			if testCase.CreateUserResult.IsExecuted(t) {
				mockRepo.EXPECT().Create(gomock.Any()).
					Return(testCase.CreateUserResult.Val(t), testCase.CreateUserResult.Err(t))
			}

			// Call the function under test
			err := CreateUser(mockRepo, "testuser")

			// Assert that there were no errors
			if testCase.expectedError != nil {
				if !errors.Is(err, testCase.expectedError) {
					t.Errorf("expected error %v, got %v", testCase.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
```

</details>

With `mresult`, the test cases are cleaner and easier to manage. Each mock result is encapsulated in a type-safe manner, allowing you to focus on the logic of your tests without worrying about the underlying mock setup.

## Contributing

Contributions are welcome! Please feel free to submit Pull Requests and create Issues.

You can use either **Japanese** or **English** for:

- Issue descriptions and discussions
- Pull Request descriptions and comments

We appreciate contributions in any form, whether it's bug reports, feature requests, documentation improvements, or code contributions.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
