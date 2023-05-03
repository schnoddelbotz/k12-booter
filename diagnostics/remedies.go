package diagnostics

type remedy struct {
	DiagIssueID    string
	FixDescription string
	ManualFix      string // #EDU - e.g. for SSH perms, doc shell chmod ... etc.
	Fix            *func() *error
}

// idea here: diag may report issue
// diag should provide command to let k12-booter fix it, e.g.
// remedy ssh-permissions
// ... as "CLI" command should fix the issue.
