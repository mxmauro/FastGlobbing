# FastGobbling

This is a partial port of the great https://github.com/Robert-van-Engelen/FastGlobbing library to Go.

#### IMPORTANT NOTES:

* Only the `Gitignore-Style Globs` is implemented. Go language already supports other pattern matching styles.
* The `DOTGLOB` behavior is not implemented. `?` does not match a dot.
* On Microsoft Windows, forward slashes are converted to backslashes internally.
* Pattern matching is case sensitive except on Microsoft Windows, which defaults to case-sensitive.

## Pattern matching rules:

```
*        matches anything except a /
?        matches any one character except a /
[a-z]    matches one character in the selected range of characters
[^a-z]   matches one character not in the selected range of characters
[!a-z]   matches one character not in the selected range of characters
/        when used at the begin of a pattern, matches if pathname has no /
**/      matches zero or more directories
/**      when at the end of a pattern, matches everything after the /
\?       matches a ? (or any character specified after the backslash)
```

#### Examples:

```
*          a, b, x/a, x/y/b
a          a, x/a, x/y/a but not b, x/b, a/a/b
/*         a, b but not x/a, x/b, x/y/a
/a         a but not x/a, x/y/a
a?b        axb, ayb but not a, b, ab, a/b
a[xy]b     axb, ayb but not a, b, azb
a[a-z]b    aab, abb, acb, azb but not a, b, a3b, aAb, aZb
a[^xy]b    aab, abb, acb, azb but not a, b, axb, ayb
a[^a-z]b   a3b, aAb, aZb but not a, b, aab, abb, acb, azb
a/*/b      a/x/b, a/y/b but not a/b, a/x/y/b
**/a       a, x/a, x/y/a but not b, x/b
a/**/b     a/b, a/x/b, a/x/y/b but not x/a/b, a/b/x
a/**       a/x, a/y, a/x/y but not a, b/x
a\?b       a?b but not a, b, ab, axb, a/b
```

## LICENSE

This library uses the same CPOL license as the original library. The CPOL license allows users to freely use,
modify, and distribute the software in open source and commercial applications.
