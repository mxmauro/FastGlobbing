# FastGlobbing

This project started as a partial port of the great https://github.com/Robert-van-Engelen/FastGlobbing library to Go.

The current matcher behavior intentionally differs in some areas from the original library, so this repository should be
treated as a Go adaptation with its own final matching rules.

#### IMPORTANT NOTES:

* Only the `Gitignore-Style Globs` is implemented. Go language already supports other pattern matching styles.
* Patterns use gitignore-style `/` path separators on every platform.
* On Microsoft Windows, input paths are matched case-insensitively and both `/` and `\` are treated as path separators 
  in both patterns and input paths.
* To match wildcard characters literally in a portable way, prefer bracket expressions such as `[*]`, `[?]`, and `[[]`.
* A trailing `/` in the pattern marks it as directory-only and matches only the directory boundary. Use `dir/**` to
  match the subtree under that directory.
* This matcher is string-based. To mimic `.gitignore` behavior for patterns ending with `/`, the caller is responsible
  for knowing whether the candidate path is a real directory.

## Pattern matching rules:

| Pattern  | Meaning                                                             |
|----------|---------------------------------------------------------------------|
| `*`      | matches anything except a `/`                                       |
| `?`      | matches any one character except a `/`                              |
| `[a-z]`  | matches one character in the selected range of characters           |
| `[^a-z]` | matches one character not in the selected range of characters       |
| `[!a-z]` | matches one character not in the selected range of characters       |
| `/`      | when used at the begin of a pattern, matches if pathname has no `/` |
| `**/`    | matches zero or more directories                                    |
| `/**`    | when at the end of a pattern, matches everything after the `/`      |
| `dir/`   | matches only the directory `dir`                                    |
| `dir/**` | matches anything under the directory `dir`                          |
| `[?]`    | matches a literal `?`                                               |
| `[*]`    | matches a literal `*`                                               |
| `[[]`    | matches a literal `[`                                               |
| `[]]`    | matches a literal `]`                                               |

#### Examples:

| Pattern    | Matches                                                          |
|------------|------------------------------------------------------------------|
| `*`        | `a`, `b`, `x/a`, `x/y/b`                                         |
| `a`        | `a`, `x/a`, `x/y/a` but not `b`, `x/b`, `a/a/b`                  |
| `/*`       | `a`, `b` but not `x/a`, `x/b`, `x/y/a`                           |
| `/a`       | `a` but not `x/a`, `x/y/a`                                       |
| `a?b`      | `axb`, `ayb` but not `a`, `b`, `ab`, `a/b`                       |
| `a[xy]b`   | `axb`, `ayb` but not `a`, `b`, `azb`                             |
| `a[a-z]b`  | `aab`, `abb`, `acb`, `azb` but not `a`, `b`, `a3b`, `aAb`, `aZb` |
| `a[^xy]b`  | `aab`, `abb`, `acb`, `azb` but not `a`, `b`, `axb`, `ayb`        |
| `a[^a-z]b` | `a3b`, `aAb`, `aZb` but not `a`, `b`, `aab`, `abb`, `acb`, `azb` |
| `a/*/b`    | `a/x/b`, `a/y/b` but not `a/b`, `a/x/y/b`                        |
| `**/a`     | `a`, `x/a`, `x/y/a` but not `b`, `x/b`                           |
| `a/**/b`   | `a/b`, `a/x/b`, `a/x/y/b` but not `x/a/b`, `a/b/x`               |
| `a/**`     | `a/x`, `a/y`, `a/x/y` but not `a`, `b/x`                         |
| `build/`   | `build` but not `build/x`, `build/x/y`, `src/build/x`            |
| `build/**` | `build/x`, `build/x/y` but not `build`                           |
| `a[?]b`    | `a?b` but not `a`, `b`, `ab`, `axb`, `a/b`                       |
| `a[*]b`    | `a*b` but not `ab`, `axb`, `a/b`                                 |
| `a[[]b`    | `a[b` but not `ab`, `axb`, `a/b`                                 |
| `a[]]b`    | `a]b` but not `ab`, `axb`, `a/b`                                 |

## LICENSE

BSD-3 and CPOL. See [LICENSE](./LICENSE) for details.
