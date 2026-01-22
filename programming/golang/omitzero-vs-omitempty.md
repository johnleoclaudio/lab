# omitzero vs omitempty

**TL;DR** Go 1.24 introduced `omitzero` - a better alternative to `omitempty` for cleaner JSON output. It handles structs, time.Time, and slices more intuitively. Use it for new projects!

When transforming Go structs to JSON, choosing which field to include is important to make sure that the output is lean, tidy, and make sense. There are two commonly used `structs` tags to achieve this: `omitempty` and `omitzero`.

For a seasoned Go developer, `omitempty` is a more familiar term when dealing with omitting fields. However, with the introduction of Go 1.24, `omitzero` has been added to the language as an alternative to control your JSON output. 

**Note:** You'll need Go 1.24 or later to use `omitzero`.

It's important to understand the differences between these two tags to make an informed decision on which one to use in your codebase. 

### omitempty

Let's talk about `omitempty` first. Key features of `omitempty` are:

- `omitempty` will omit the field if it has a zero-value.
- Zero-values include: `0` for numbers, `""` for strings, `false` for booleans, `nil` for pointers, slices, maps, interfaces, and channels
- `omitempty` will not omit `structs` even if all the fields are zero-value.

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
    MiddleName string    `json:"middle_name,omitempty"`
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

### omitzero

Now, let's look at `omitzero`. First off, the name, it clearly communicates its intended purpose: to omit fields that have zero-values.

Comparing the example earlier, we will only change the struct tag for the `Address` field from `omitempty` to `omitzero`.

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

Not only the payload is leaner, but it also makes more sense. If the `Address` field is not provided, it should not appear in the JSON output at all.


Here's what `omitzero` handles better than `omitempty`:
- omitempty will not omit time.Time
- omitempty will not omit arrays
- omitempty will not omit empty slices and maps

### Key Takeaways

**For new code, prefer `omitzero`** - it provides more intuitive behavior and cleaner JSON output:

- Omits all zero-value fields consistently, including structs, time.Time, slices, and empty collections
- The name clearly communicates its purpose
- Results in leaner, more meaningful JSON payloads

**When to use `omitempty`**:
- Working with legacy codebases where changing behavior might break clients
- You specifically need to preserve empty structs or time.Time zero values in JSON output
- Maintaining backward compatibility with existing APIs

`omitzero` represents Go's evolution toward more predictable JSON marshaling behavior. For greenfield projects, it's the better default choice.

### Further Reading
- [Go 1.24 Release Notes](https://go.dev/doc/go1.24#json)
- [Go JSON Package Documentation](https://pkg.go.dev/encoding/json)
