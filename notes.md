# Notes

## Response `Cache-Control` directives condition implementation order

### Ignore caching

- `no-store`
  - `must-understand` (unless)
- `private`

### Cache, don't re-use response, always revalidate

- `no-cache`

### Cache, re-use response, only revalidate on expiration

- *may serve stale* on revalidation failure
  - `max-age`
  - `s-maxage` (Auth override)
- *must not serve stale* on revalidation failure, rather return error
  - `must-revalidate` (Auth override)

### Heusteric caching

- `public` (Auth override)
  - `proxy-revalidate`

### For when you think about transforming the response

- `no-transform`
