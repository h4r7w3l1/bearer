package string

import (
	"regexp"

	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

var stringRegex = regexp.MustCompile(`\A\w?['"]{1,3}(.*?)['"]{1,3}\z`)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "string":
		return handleTemplateString(node, detectorContext)
	case "concatenated_string":
		return common.ConcatenateChildStrings(node, detectorContext)
	case "binary_operator":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "augmented_assignment":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}

func handleTemplateString(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	text := ""
	isLiteral := true

	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		var childValue string
		var childIsLiteral bool
		namedChildren := child.NamedChildren()

		if len(namedChildren) == 0 {
			childValue = ""
			childIsLiteral = true
		} else {
			var err error
			childValue, childIsLiteral, err = common.GetStringValue(namedChildren[0], detectorContext)
			if err != nil {
				return err
			}
		}

		if childValue == "" && !childIsLiteral {
			childValue = common.NonLiteralValue
		}

		text += childValue

		if !childIsLiteral {
			isLiteral = false
		}

		return nil
	})

	text = stringRegex.ReplaceAllString(text, `$1`)

	return []interface{}{common.String{
		Value:     text,
		IsLiteral: isLiteral,
	}}, err
}
