# omitzero vs omitempty
When transforming Go structs to JSON, choosing which field to include is important to make sure that the output is lean, clean, and does make sense. There two commonly used struct tags to achieve this: `omitempty` and `omitzero`.

For a seasoned Go developer, `omitempty` is a more familiar term. However, with the introduction of Go 1.24, `omitzero` has been added to provide a more explicit way to omit fields with zero values, including structs.

- omitzero is newer addition in Go 1.24 
- omitzero is clear about the intent. Remove field with zero-value.
- omitempty will not omit structs even if all the fields are zero-value.

Consider this example. 
```go
  type Address struct {
    Street     string `json:"street"`
    City       string `json:"city"`
    State      string `json:"state"`
    PostalCode string `json:"postal_code"`
    Country    string `json:"country"`
  }

  type User struct {
    ID         uuid.UUID `json:"id"`
    FirstName  string    `json:"first_name"`
    LastName   string    `json:"last_name"`
    MiddleName string    `json:"middle_name,omitzero"`
    Address    Address   `json:"address,omitempty"` // <-- embedded Address struct with omitempty directive
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
  }

  // Create struct but left out Address field
  user := data.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
```

Here's the result:
```bash
❯ curl -i localhost:4000/v1/users/019b87ef-b562-7d67-b00b-8940309795de
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "address": {
    "street": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "created_at": "2026-01-04T16:34:47.13207+08:00",
  "updated_at": "2026-01-04T16:34:47.13207+08:00"
}
```

If `omitzero` was used instead, we'll get a cleaner result:
```go
  type User struct {
    Address    Address   `json:"address,omitzero"`
    ...
  }
```
```bash
{
  "id": "019b87ef-b562-7d67-b00b-8940309795de",
  "first_name": "John",
  "last_name": "Doe",
  "created_at": "2026-01-04T16:37:54.072378+08:00",
  "updated_at": "2026-01-04T16:37:54.072378+08:00"
}
```
- Other difference includes:
  - omitempty will not omit time.Time
  - omitempty will not omit arrays
  - omitempty will not omit empty slices and maps
- 
