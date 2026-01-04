# Encoding Go Structs to JSON

- my notes from Let's Go Further book by Alex Edwards

```go
type User struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
```
### Customize JSON by struct tags
- change how key names that appear in the JSON object
- choose which fields to include in the response

```go
type User struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
```

- IMPORTANT: struct fields must be exported (starts with capital letter). Any fields which are not exported will be ignored when encoding a struct to JSON

```go 
type User struct {
	id         uuid.UUID `json:"id"`
  ...
}
```

```bash
❯ curl localhost:4000/v1/users/019b87ef-b562-7d67-b00b-8940309795de
{
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "M",
  "created_at": "2026-01-04T15:38:53.968434+08:00",
  "updated_at": "2026-01-04T15:38:53.968434+08:00"
}
```

### Hiding fields in the JSON object
- `-` (hyphen) directive can be used when you NEVER want a struct field to appear in the JSON output
```go
type User struct {
	MiddleName string    `json:"_"`
  ...
}
```

```bash
❯ curl localhost:4000/v1/users/019b87ef-b562-7d67-b00b-8940309795de
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "created_at": "2026-01-04T15:42:39.592092+08:00",
  "updated_at": "2026-01-04T15:42:39.592092+08:00"
}
```
- `omitzero` directive hides a field in the JSON output if and only if the value is zero value for the field type.

Here's a sample on no-value field:
```go
  user := data.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
    // MiddleName is omitted to demonstrate default behavior
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
```

```bash
❯ curl -i localhost:4000/v1/users/019b87ef-b562-7d67-b00b-8940309795de
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "",
  "created_at": "2026-01-04T16:18:49.2911+08:00",
  "updated_at": "2026-01-04T16:18:49.2911+08:00"
}
```

Adding `omitzero` hides zero-value field

```go
  type User struct {
    MiddleName string    `json:"middle_name,omitzero"`
    ...
  }

  user := data.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
```

Notice that instead of empty string, the field was removed altogether
```bash
❯ curl -i localhost:4000/v1/users/019b87ef-b562-7d67-b00b-8940309795de
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "created_at": "2026-01-04T16:24:26.914947+08:00",
  "updated_at": "2026-01-04T16:24:26.914947+08:00"
}
```
- why not use `omitempty` instead? Read my notes on comparing `omitzero` vs `omitempty`

# The string directive
- use this directive to force the output to be string
- works on `int*`,`uint*`,`float*` or `bool` types

Default `bool` behavior:
```go
  type User struct {
	  IsActive bool `json:"is_active"`
    ...
  }

  user := data.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
```
```bash
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": false,
  "created_at": "2026-01-04T16:46:47.141081+08:00",
  "updated_at": "2026-01-04T16:46:47.141082+08:00"
}
```

But using `string` directive:
```go
  type User struct {
	  IsActive bool `json:"is_active,string"`
    ...
  }

  user := data.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
```
```bash
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": "false",
  "created_at": "2026-01-04T16:47:11.037533+08:00",
  "updated_at": "2026-01-04T16:47:11.037533+08:00"
}
```
