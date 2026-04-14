package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

const defaultLimit = 100

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "collect":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: go run ./skills/tanaoroshi collect <owner/repo> [owner/repo2 ...]")
			os.Exit(1)
		}
		collect(os.Args[2:])
	case "summary":
		summary(os.Args[2:])
	case "refs":
		refs(os.Args[2:])
	case "resolve":
		resolve(os.Args[2:])
	case "comments":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: go run ./skills/tanaoroshi comments <owner/repo:N> [owner/repo:N ...]")
			os.Exit(1)
		}
		comments(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: go run ./skills/tanaoroshi <command> [args...]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  collect <owner/repo> [...]                    Fetch open issues and PRs (JSON)")
	fmt.Fprintln(os.Stderr, "  summary                                      Compact summary without body (stdin)")
	fmt.Fprintln(os.Stderr, "  refs                                         Extract cross-references from body (stdin)")
	fmt.Fprintln(os.Stderr, "  resolve --issues <owner/repo:N> --prs <ref>  Check status of referenced items")
	fmt.Fprintln(os.Stderr, "  comments <owner/repo:N> [...]               Fetch issue/review comments")
}

// collect fetches open issues and PRs for each repo.
func collect(repos []string) {
	result := map[string]any{}
	limit := fmt.Sprintf("%d", defaultLimit)

	for _, repo := range repos {
		issues := ghList("issue", repo, limit,
			"number,title,author,createdAt,updatedAt,labels,assignees,body,url")
		prs := ghList("pr", repo, limit,
			"number,title,author,createdAt,updatedAt,labels,isDraft,reviewDecision,headRefName,body,url")

		if len(issues) >= defaultLimit {
			fmt.Fprintf(os.Stderr, "[warn] %s: issues hit limit (%d).\n", repo, defaultLimit)
		}
		if len(prs) >= defaultLimit {
			fmt.Fprintf(os.Stderr, "[warn] %s: PRs hit limit (%d).\n", repo, defaultLimit)
		}

		result[repo] = map[string]any{
			"issues": issues,
			"prs":    prs,
		}
	}

	writeJSON(result)
}

// summary reads collect JSON and outputs a compact summary without body.
func summary(args []string) {
	data := readInput(args)
	ignore := loadIgnoreList()

	ignored := 0
	result := map[string]any{}
	for repo, val := range data {
		rv, ok := val.(map[string]any)
		if !ok {
			continue
		}
		issues := toSlice(rv["issues"])
		prs := toSlice(rv["prs"])
		filteredIssues := filterIgnored(issues, ignore, repo)
		filteredPRs := filterIgnored(prs, ignore, repo)
		ignored += len(issues) - len(filteredIssues) + len(prs) - len(filteredPRs)
		result[repo] = map[string]any{
			"issues": stripBodies(filteredIssues),
			"prs":    stripBodies(filteredPRs),
		}
	}
	if ignored > 0 {
		fmt.Fprintf(os.Stderr, "[info] ignored %d item(s) from %s\n", ignored, ignoreFile)
	}

	writeJSON(result)
}

func stripBodies(items []any) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		cp := make(map[string]any, len(m))
		for k, v := range m {
			if k != "body" {
				cp[k] = v
			}
		}
		out = append(out, cp)
	}
	return out
}

var (
	reGHURL    = regexp.MustCompile(`https://github\.com/([^/\s]+/[^/\s]+)/(issues|pull)/(\d+)`)
	reCrossRef = regexp.MustCompile(`([a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+)#(\d+)`)
	reLocalRef = regexp.MustCompile(`(?:^|[\s(])#(\d+)`)
	reCloseFix = regexp.MustCompile(`(?i)(?:closes|fixes|close|fix)\s+#(\d+)`)
)

// refs reads collect JSON and extracts unique cross-references.
func refs(args []string) {
	data := readInput(args)
	ignore := loadIgnoreList()

	seen := map[string]bool{}
	type refEntry struct {
		Source string `json:"source"`
		Ref    string `json:"ref"`
	}
	var entries []refEntry

	for repo, val := range data {
		rv, ok := val.(map[string]any)
		if !ok {
			continue
		}
		for _, items := range [][]any{toSlice(rv["issues"]), toSlice(rv["prs"])} {
			for _, item := range items {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				body, _ := m["body"].(string)
				if body == "" {
					continue
				}
				number := jsonNumber(m["number"])
				if isIgnored(ignore, repo, number) {
					continue
				}
				source := fmt.Sprintf("%s#%s", repo, number)

				extracted := extractRefs(body, repo)
				for _, ref := range extracted {
					if !seen[ref] {
						seen[ref] = true
						entries = append(entries, refEntry{Source: source, Ref: ref})
					}
				}
			}
		}
	}

	if entries == nil {
		entries = []refEntry{}
	}
	writeJSON(entries)
}

func extractRefs(body, defaultRepo string) []string {
	seen := map[string]bool{}
	var refs []string

	add := func(ref string) {
		if !seen[ref] {
			seen[ref] = true
			refs = append(refs, ref)
		}
	}

	for _, m := range reGHURL.FindAllStringSubmatch(body, -1) {
		add(fmt.Sprintf("%s#%s", m[1], m[3]))
	}
	for _, m := range reCrossRef.FindAllStringSubmatch(body, -1) {
		if !strings.HasPrefix(m[0], "https://") {
			add(fmt.Sprintf("%s#%s", m[1], m[2]))
		}
	}
	for _, m := range reCloseFix.FindAllStringSubmatch(body, -1) {
		add(fmt.Sprintf("%s#%s", defaultRepo, m[1]))
	}
	for _, m := range reLocalRef.FindAllStringSubmatch(body, -1) {
		ref := fmt.Sprintf("%s#%s", defaultRepo, m[1])
		if !seen[ref] {
			add(ref)
		}
	}

	return refs
}

func jsonNumber(v any) string {
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", n)
	case json.Number:
		return n.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// resolve checks the status of referenced issues/PRs.
func resolve(args []string) {
	var issueRefs, prRefs []string
	mode := ""

	for _, arg := range args {
		switch arg {
		case "--issues":
			mode = "issues"
		case "--prs":
			mode = "prs"
		default:
			switch mode {
			case "issues":
				issueRefs = append(issueRefs, arg)
			case "prs":
				prRefs = append(prRefs, arg)
			default:
				fmt.Fprintln(os.Stderr, "usage: go run ./skills/tanaoroshi resolve --issues <owner/repo:N> ... --prs <owner/repo:N> ...")
				os.Exit(1)
			}
		}
	}

	if len(issueRefs) == 0 && len(prRefs) == 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./skills/tanaoroshi resolve --issues <owner/repo:N> ... --prs <owner/repo:N> ...")
		os.Exit(1)
	}

	result := map[string]any{}

	for _, ref := range issueRefs {
		repo, number := parseRef(ref)
		data := ghView("issue", repo, number, "state,title,closedAt,url")
		if data != nil {
			data["type"] = "issue"
			result[ref] = data
		} else {
			result[ref] = map[string]any{"type": "issue", "state": "UNKNOWN", "title": "", "closedAt": nil, "url": nil}
		}
	}

	for _, ref := range prRefs {
		repo, number := parseRef(ref)
		data := ghView("pr", repo, number, "state,title,closedAt,mergedAt,url")
		if data != nil {
			data["type"] = "pr"
			result[ref] = data
		} else {
			result[ref] = map[string]any{"type": "pr", "state": "UNKNOWN", "title": "", "closedAt": nil, "mergedAt": nil, "url": nil}
		}
	}

	writeJSON(result)
}

// helpers

// readInput reads JSON from a file (first arg) or stdin if no args given.
func readInput(args []string) map[string]any {
	var b []byte
	var err error

	if len(args) > 0 {
		b, err = os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] failed to read file %s: %v\n", args[0], err)
			os.Exit(1)
		}
	} else {
		b, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] failed to read stdin: %v\n", err)
			os.Exit(1)
		}
	}

	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to parse JSON: %v\n", err)
		os.Exit(1)
	}
	return data
}

func toSlice(v any) []any {
	s, ok := v.([]any)
	if !ok {
		return nil
	}
	return s
}

func parseRef(ref string) (string, string) {
	idx := strings.LastIndex(ref, ":")
	if idx < 0 {
		fmt.Fprintf(os.Stderr, "[error] invalid ref format (expected owner/repo:N): %s\n", ref)
		os.Exit(1)
	}
	return ref[:idx], ref[idx+1:]
}

func ghList(kind, repo, limit, fields string) []any {
	out, err := exec.Command("gh", kind, "list",
		"--repo", repo,
		"--state", "open",
		"--limit", limit,
		"--json", fields,
	).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warn] gh %s list --repo %s failed: %v\n", kind, repo, err)
		return []any{}
	}
	var items []any
	if err := json.Unmarshal(out, &items); err != nil {
		fmt.Fprintf(os.Stderr, "[warn] failed to parse gh %s list output for %s: %v\n", kind, repo, err)
		return []any{}
	}
	return items
}

func ghView(kind, repo, number, fields string) map[string]any {
	out, err := exec.Command("gh", kind, "view", number,
		"--repo", repo,
		"--json", fields,
	).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warn] Failed to fetch %s: %s:%s\n", kind, repo, number)
		return nil
	}
	var data map[string]any
	if err := json.Unmarshal(out, &data); err != nil {
		fmt.Fprintf(os.Stderr, "[warn] Failed to parse %s view output for %s:%s\n", kind, repo, number)
		return nil
	}
	return data
}

// comments fetches comments for specified issues/PRs.
// For PRs, both issue comments and review comments are fetched and merged by time.
// Args: owner/repo:N [owner/repo:N ...]
func comments(args []string) {
	type comment struct {
		Author    string `json:"author"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		Type      string `json:"type"` // "comment" or "review"
	}
	type entry struct {
		Repo     string    `json:"repo"`
		Number   string    `json:"number"`
		Comments []comment `json:"comments"`
	}

	var entries []entry
	for _, arg := range args {
		repo, number := parseRef(arg)
		if repo == "" || number == "" {
			fmt.Fprintf(os.Stderr, "[warn] invalid ref: %s (expected owner/repo:N)\n", arg)
			continue
		}

		var all []comment

		// Fetch issue comments (paginated, up to 100)
		issueEndpoint := fmt.Sprintf("repos/%s/issues/%s/comments?per_page=100", repo, number)
		if out, err := exec.Command("gh", "api", issueEndpoint).Output(); err == nil {
			var raw []map[string]any
			if json.Unmarshal(out, &raw) == nil {
				for _, c := range raw {
					all = append(all, extractComment(c, "comment"))
				}
			}
		}

		// Fetch PR review comments (fails silently for issues)
		reviewEndpoint := fmt.Sprintf("repos/%s/pulls/%s/comments?per_page=100", repo, number)
		if out, err := exec.Command("gh", "api", reviewEndpoint).Output(); err == nil {
			var raw []map[string]any
			if json.Unmarshal(out, &raw) == nil {
				for _, c := range raw {
					all = append(all, extractComment(c, "review"))
				}
			}
		}

		// Sort chronologically
		sort.Slice(all, func(i, j int) bool {
			return all[i].CreatedAt < all[j].CreatedAt
		})

		entries = append(entries, entry{Repo: repo, Number: number, Comments: all})
	}

	if entries == nil {
		entries = []entry{}
	}
	writeJSON(entries)
}

func extractComment(c map[string]any, typ string) struct {
	Author    string `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
	Type      string `json:"type"`
} {
	author := ""
	if u, ok := c["user"].(map[string]any); ok {
		if login, ok := u["login"].(string); ok {
			author = login
		}
	}
	body, _ := c["body"].(string)
	createdAt, _ := c["created_at"].(string)
	return struct {
		Author    string `json:"author"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		Type      string `json:"type"`
	}{
		Author:    author,
		Body:      body,
		CreatedAt: createdAt,
		Type:      typ,
	}
}

const ignoreFile = "skills/tanaoroshi/ignore"

func loadIgnoreList() map[string]bool {
	b, err := os.ReadFile(ignoreFile)
	if err != nil {
		return map[string]bool{}
	}
	m := map[string]bool{}
	for line := range strings.SplitSeq(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m[line] = true
	}
	return m
}

func isIgnored(ignore map[string]bool, repo, number string) bool {
	return ignore[fmt.Sprintf("%s#%s", repo, number)]
}

func filterIgnored(items []any, ignore map[string]bool, repo string) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if isIgnored(ignore, repo, jsonNumber(m["number"])) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}
