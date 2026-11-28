# services/admin

Content management domain (backoffice).

Intended for admin operations on story content: create, update, publish, and archive stories.

Planned layout (to be filled in later):

```
services/admin/
├── api/            JSON handlers for admin operations
├── web/            HTML handlers/templates for admin UI
├── usecase.go      business rules for content management
├── repository.go   data access for admin actions
└── model.go        shared domain types
```

Currently a placeholder — no implementation yet.
