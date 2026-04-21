# Style Standards Reference

Look for a `vs-standards` repo in the workspace or a sibling directory. If found, read the standards files directly. If not, fall back to Confluence or the Google Go Style Guide.

| Language | Local path (if vs-standards exists) | Confluence fallback | External fallback |
|----------|-------------------------------------|--------------------|--------------------|
| Go | `vs-standards/standards/go/` | VetSoft Go style guide (page 5153587807) | Google Go Style Guide |
| API | `vs-standards/standards/api/` | API Standards (page 5153980984) | - |
| Software Design | `vs-standards/standards/software-design/` | - | - |
| PHP | - | Search Confluence for "PHP Standard" | - |
| TypeScript | - | TypeScript Coding Standards (page 5421826381) | - |

## Go-specific instructions

For Go PRs, read `STANDARD.md` and all `.md` files under the `go/` subdirectories (errors, testing, database, infrastructure). Check each changed file against applicable standards.

## Finding vs-standards

Try these locations in order:
1. `../vs-standards/` (sibling directory)
2. `~/code/idexx-emu/vs-standards/` (common clone location)
3. If neither exists, use Confluence via MCP or skip standards checks with a note.
