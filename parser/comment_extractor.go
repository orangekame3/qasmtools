package parser

import (
	"sort"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// CommentExtractor extracts comments from token stream
type CommentExtractor struct {
	comments []Comment
	content  string
}

// NewCommentExtractor creates a new comment extractor
func NewCommentExtractor(content string) *CommentExtractor {
	return &CommentExtractor{
		comments: make([]Comment, 0),
		content:  content,
	}
}

// ExtractComments extracts comments from the token stream
func (ce *CommentExtractor) ExtractComments(stream antlr.TokenStream) []Comment {
	// Convert to CommonTokenStream to access hidden tokens
	commonStream, ok := stream.(*antlr.CommonTokenStream)
	if !ok {
		return ce.comments
	}

	// Get all tokens including hidden ones
	tokens := commonStream.GetAllTokens()
	for _, token := range tokens {
		if token.GetChannel() == antlr.TokenHiddenChannel {
			text := token.GetText()
			if strings.HasPrefix(text, "//") || strings.HasPrefix(text, "/*") {
				ce.extractComment(token)
			}
		}
	}

	return ce.comments
}

// isCommentToken checks if a token is a comment
func (ce *CommentExtractor) isCommentToken(token antlr.Token) bool {
	// Check if token is on a hidden channel (typically comments)
	return token.GetChannel() != antlr.TokenDefaultChannel
}

// extractComment extracts comment information from a token
func (ce *CommentExtractor) extractComment(token antlr.Token) {
	text := token.GetText()
	if text == "" {
		return
	}

	// Determine comment type
	commentType := "line"
	if strings.HasPrefix(text, "/*") {
		commentType = "block"
	}

	// Clean up comment text
	cleanText := ce.cleanCommentText(text, commentType)

	comment := Comment{
		BaseNode: BaseNode{
			Position: Position{
				Line:   token.GetLine(),
				Column: token.GetColumn() + 1,
				Offset: token.GetStart(),
			},
			EndPos: Position{
				Line:   token.GetLine(),
				Column: token.GetColumn() + len(text),
				Offset: token.GetStop() + 1,
			},
		},
		Text:     cleanText,
		Type:     commentType,
		Attached: false, // Will be determined later
	}

	ce.comments = append(ce.comments, comment)
}

// cleanCommentText removes comment markers and cleans up text
func (ce *CommentExtractor) cleanCommentText(text, commentType string) string {
	switch commentType {
	case "line":
		// Remove // prefix
		if strings.HasPrefix(text, "//") {
			text = strings.TrimPrefix(text, "//")
		}
		// Remove leading/trailing whitespace
		text = strings.TrimSpace(text)

	case "block":
		// Remove /* */ markers
		if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimSuffix(text, "*/")
		}
		// Remove leading/trailing whitespace
		text = strings.TrimSpace(text)
	}

	return text
}

// AssociateCommentsWithStatements associates comments with nearby statements
func (ce *CommentExtractor) AssociateCommentsWithStatements(program *Program) {
	if len(ce.comments) == 0 {
		return
	}

	// Initialize LineComments map if nil
	if program.LineComments == nil {
		program.LineComments = make(map[int][]Comment)
	}

	// Sort comments by line number
	sort.Slice(ce.comments, func(i, j int) bool {
		return ce.comments[i].Position.Line < ce.comments[j].Position.Line
	})

	// Associate each comment with statements
	for i := range ce.comments {
		comment := &ce.comments[i]

		// Check if comment is on same line as a statement (trailing comment)
		var attached bool
		for _, stmt := range program.Statements {
			stmtLine := stmt.Pos().Line
			if comment.Position.Line == stmtLine {
				comment.Attached = true
				attached = true
				break
			}
		}

		// If not attached to any statement, it's a leading comment
		if !attached {
			// Find the next statement after this comment
			var nextStmtLine int
			for _, stmt := range program.Statements {
				if stmt.Pos().Line > comment.Position.Line {
					nextStmtLine = stmt.Pos().Line
					break
				}
			}

			// If there's a next statement and the comment is within 3 lines of it,
			// associate it with that statement
			if nextStmtLine > 0 && nextStmtLine-comment.Position.Line <= 3 {
				comment.Attached = false
			}
		}

		// Add comment to line-based mapping
		line := comment.Position.Line
		program.LineComments[line] = append(program.LineComments[line], *comment)
	}

	// Set the comments on the program
	program.Comments = ce.comments
}

// GetCommentsForLine returns comments associated with a specific line
func GetCommentsForLine(program *Program, line int) []Comment {
	if program.LineComments == nil {
		return nil
	}
	return program.LineComments[line]
}

// GetLeadingComments returns comments that appear before a statement
func GetLeadingComments(program *Program, stmt Statement) []Comment {
	stmtLine := stmt.Pos().Line
	leadingComments := make([]Comment, 0)

	// Look for comments in the lines immediately before the statement
	for line := stmtLine - 3; line < stmtLine; line++ {
		if line <= 0 {
			continue
		}
		if comments, exists := program.LineComments[line]; exists {
			for _, comment := range comments {
				if !comment.Attached {
					leadingComments = append(leadingComments, comment)
				}
			}
		}
	}

	return leadingComments
}

// GetTrailingComments returns comments that appear on the same line as a statement
func GetTrailingComments(program *Program, stmt Statement) []Comment {
	stmtLine := stmt.Pos().Line
	trailingComments := make([]Comment, 0)

	if comments, exists := program.LineComments[stmtLine]; exists {
		for _, comment := range comments {
			if comment.Attached {
				trailingComments = append(trailingComments, comment)
			}
		}
	}

	return trailingComments
}
