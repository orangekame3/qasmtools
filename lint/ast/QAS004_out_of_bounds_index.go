package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// OutOfBoundsIndexRule implements QAS004 using AST-based analysis
type OutOfBoundsIndexRule struct {
	*ASTRuleBase
}

// NewOutOfBoundsIndexRule creates a new AST-based out of bounds index rule
func NewOutOfBoundsIndexRule() ASTRule {
	return &OutOfBoundsIndexRule{
		ASTRuleBase: NewASTRuleBase("QAS004"),
	}
}

// CheckAST performs AST-based out of bounds index analysis
func (r *OutOfBoundsIndexRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Find all declarations to get array sizes
	declarations := astutil.FindDeclarations(program)
	arrayInfo := r.buildArrayInfo(declarations)

	// Find all indexed identifiers and check bounds
	astutil.VisitAllNodes(program, func(node parser.Node) {
		if indexedId, ok := node.(*parser.IndexedIdentifier); ok {
			if violation := r.checkIndexBounds(indexedId, arrayInfo, ctx); violation != nil {
				violations = append(violations, violation)
			}
		}
	})

	return violations
}

// ArrayInfo holds information about array declarations
type ArrayInfo struct {
	Name string
	Size int
	Type string // "qubit" or "bit"
}

// buildArrayInfo creates a map of array names to their size information
func (r *OutOfBoundsIndexRule) buildArrayInfo(declarations *astutil.Declarations) map[string]*ArrayInfo {
	arrayInfo := make(map[string]*ArrayInfo)

	// Process quantum declarations
	for _, qubitDecl := range declarations.Quantum {
		if qubitDecl.Size != nil {
			if size, ok := astutil.GetArraySize(qubitDecl.Size); ok {
				arrayInfo[qubitDecl.Identifier] = &ArrayInfo{
					Name: qubitDecl.Identifier,
					Size: size,
					Type: "qubit",
				}
			}
		}
	}

	// Process classical declarations
	for _, classicalDecl := range declarations.Classical {
		if classicalDecl.Size != nil {
			// Array declaration with explicit size
			if size, ok := astutil.GetArraySize(classicalDecl.Size); ok {
				arrayInfo[classicalDecl.Identifier] = &ArrayInfo{
					Name: classicalDecl.Identifier,
					Size: size,
					Type: classicalDecl.Type,
				}
			}
		} else if strings.Contains(classicalDecl.Type, "[") {
			// Array size encoded in type string: "bit[n]"
			if size := r.extractSizeFromType(classicalDecl.Type); size > 0 {
				// Extract base type (bit from bit[3])
				baseType := strings.Split(classicalDecl.Type, "[")[0]
				arrayInfo[classicalDecl.Identifier] = &ArrayInfo{
					Name: classicalDecl.Identifier,
					Size: size,
					Type: baseType,
				}
			}
		}
	}

	return arrayInfo
}

// extractSizeFromType extracts array size from type strings like "bit[3]"
func (r *OutOfBoundsIndexRule) extractSizeFromType(typeStr string) int {
	// Look for pattern like "bit[3]"
	if start := strings.Index(typeStr, "["); start != -1 {
		if end := strings.Index(typeStr[start:], "]"); end != -1 {
			sizeStr := typeStr[start+1 : start+end]
			if size, err := strconv.Atoi(sizeStr); err == nil {
				return size
			}
		}
	}
	return 0
}

// checkIndexBounds checks if an indexed access is within bounds
func (r *OutOfBoundsIndexRule) checkIndexBounds(indexedId *parser.IndexedIdentifier, arrayInfo map[string]*ArrayInfo, ctx *CheckContext) *Violation {
	// Get array information
	info, exists := arrayInfo[indexedId.Name]
	if !exists {
		// Array not found, this might be a separate issue (undefined identifier)
		return nil
	}

	// Try to extract the index value
	indexValue, ok := r.extractIntegerIndex(indexedId.Index)
	if !ok {
		// Cannot determine index at compile time
		return nil
	}

	// Check bounds (arrays are 0-indexed)
	if indexValue < 0 || indexValue >= info.Size {
		message := fmt.Sprintf("Index out of bounds: accessing '%d' on '%s' of length %d.", 
			indexValue, info.Name, info.Size)
			
		return r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(indexedId).
			WithNodeName(fmt.Sprintf("%s[%d]", info.Name, indexValue)).
			AsError().
			Build()
	}

	return nil
}

// extractIntegerIndex attempts to extract an integer value from an index expression
func (r *OutOfBoundsIndexRule) extractIntegerIndex(expr parser.Expression) (int, bool) {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return int(e.Value), true
	case *parser.ParenthesizedExpression:
		return r.extractIntegerIndex(e.Expression)
	default:
		// Cannot determine value at compile time
		return 0, false
	}
}