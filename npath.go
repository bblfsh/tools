package tools

import (
	"fmt"

	"github.com/bblfsh/sdk/uast"
)

//CodeReference
//https://pmd.github.io/pmd-5.7.0/pmd-java/xref/net/sourceforge/pmd/lang/java/rule/codesize/NPathComplexityRule.html

type Npath struct{}

type NpathData struct {
	Name       string
	Complexity int
}

func (np Npath) Exec(n *uast.Node) error {
	result := NpathComplexity(n)
	fmt.Println(result)
	return nil
}

func (nd *NpathData) String() string {
	return fmt.Sprintf("FuncName:%s, Complexity:%d\n", nd.Name, nd.Complexity)
}

//NpathComplexity return a NpathData for each function in the tree
func NpathComplexity(n *uast.Node) []*NpathData {
	var result []*NpathData
	var funcs []*uast.Node
	var names []string

	if containsRole(n, uast.FunctionDeclarationBody) {
		funcs = append(funcs, n)
		names = append(names, "NoName")
	} else {
		funcDecs := deepChildrenOfRole(n, uast.FunctionDeclaration)
		for _, funcDec := range funcDecs {
			names = append(names, childrenOfRole(funcDec, uast.FunctionDeclarationName)[0].Token)
			funcs = append(funcs, childrenOfRole(funcDec, uast.FunctionDeclarationBody)[0])
		}
	}
	for i, function := range funcs {
		npath := visitFunctionBody(function)
		result = append(result, &NpathData{Name: names[i], Complexity: npath})
	}

	return result
}

func visitorSelector(n *uast.Node) int {
	// I need to add a error when the node dont have any rol
	// when I got 2 or more roles that are inside the switch this doesn't work
	for _, role := range n.Roles {
		switch role {
		case uast.If:
			return visitIf(n)
		case uast.While:
			return visitWhile(n)
		case uast.Switch:
			return visitSwitch(n)
		case uast.DoWhile:
			return visitDoWhile(n)
		case uast.For:
			return visitFor(n)
		case uast.ForEach:
			return visitForEach(n)
		case uast.Return:
			return visitReturn(n)
		case uast.Try:
			return visitTry(n)
		default:
		}
	}
	return visitNotCompNode(n)
}

func complexityMultOf(n *uast.Node) int {
	npath := 1
	for _, child := range n.Children {
		npath *= visitorSelector(child)
	}
	return npath
}

func visitFunctionBody(n *uast.Node) int {
	return complexityMultOf(n)
}

func visitNotCompNode(n *uast.Node) int {
	return complexityMultOf(n)
}

func visitIf(n *uast.Node) int {
	// (npath of if + npath of else (or 1) + bool_comp of if) * npath of next
	npath := 0
	ifBody := childrenOfRole(n, uast.IfBody)
	ifCondition := childrenOfRole(n, uast.IfCondition)
	ifElse := childrenOfRole(n, uast.IfElse)

	if len(ifElse) == 0 {
		npath++
	} else {
		// This if is a short circuit to avoid the two roles in the switch problem
		if containsRole(ifElse[0], uast.If) {
			npath += visitIf(ifElse[0])
		} else {
			npath += complexityMultOf(ifElse[0])
		}
	}
	npath *= complexityMultOf(ifBody[0])
	npath += expressionComp(ifCondition[0])

	return npath
}

func visitWhile(n *uast.Node) int {
	// (npath of while + bool_comp of while + npath of else (or 1)) * npath of next
	npath := 0
	whileCondition := childrenOfRole(n, uast.WhileCondition)
	whileBody := childrenOfRole(n, uast.WhileBody)
	whileElse := childrenOfRole(n, uast.IfElse)
	// Some languages like python can have an else in a while loop
	if len(whileElse) == 0 {
		npath++
	} else {
		npath += complexityMultOf(whileElse[0])
	}
	npath *= complexityMultOf(whileBody[0])
	npath += expressionComp(whileCondition[0])

	return npath
}

func visitDoWhile(n *uast.Node) int {
	// (npath of do + bool_comp of do + 1) * npath of next
	npath := 1
	doWhileCondition := childrenOfRole(n, uast.DoWhileCondition)
	doWhileBody := childrenOfRole(n, uast.DoWhileBody)

	npath *= complexityMultOf(doWhileBody[0])
	npath += expressionComp(doWhileCondition[0])

	return npath
}

func visitFor(n *uast.Node) int {
	// (npath of for + bool_comp of for + 1) * npath of next
	npath := 1
	forBody := childrenOfRole(n, uast.ForBody)

	npath *= complexityMultOf(forBody[0])

	npath++
	return npath
}

func visitReturn(n *uast.Node) int {
	if aux := expressionComp(n); aux != 1 {
		return aux - 1
	}
	return 1
}

func visitSwitch(n *uast.Node) int {
	caseDefault := childrenOfRole(n, uast.SwitchDefault)
	switchCases := childrenOfRole(n, uast.SwitchCase)
	npath := 0

	if len(caseDefault) != 0 {
		npath += complexityMultOf(caseDefault[0])
	} else {
		npath++
	}
	for _, switchCase := range switchCases {
		npath += complexityMultOf(switchCase)
	}
	return npath
}

func visitForEach(n *uast.Node) int {
	forBody := childrenOfRole(n, uast.ForBody)
	npath := 1

	npath *= complexityMultOf(forBody[0])
	npath++
	return npath
}

func visitTry(n *uast.Node) int {
	/*
		In pmd they decided the complexity of a try is the summatory of the complexity
		of the tryBody,cathBody and finallyBody.I don't think this is the most acurate way
		of doing this.
	*/

	tryBody := childrenOfRole(n, uast.TryBody)
	tryCatch := childrenOfRole(n, uast.TryCatch)
	tryFinaly := childrenOfRole(n, uast.TryFinally)

	catchComp := 0
	if len(tryCatch) != 0 {
		for _, catch := range tryCatch {
			catchComp += complexityMultOf(catch)
		}
	}
	finallyComp := 0
	if len(tryFinaly) != 0 {
		finallyComp = complexityMultOf(tryFinaly[0])
	}
	npath := complexityMultOf(tryBody[0]) + catchComp + finallyComp

	return npath
}

func visitConditionalExpr(n *uast.Node) {
	// TODO ternary operators are not defined on the UAST yet
}

func expressionComp(n *uast.Node) int {
	orCount := deepCountChildrenOfRole(n, uast.OpBooleanAnd)
	andCount := deepCountChildrenOfRole(n, uast.OpBooleanOr)

	return orCount + andCount + 1
}

func containsRole(n *uast.Node, role uast.Role) bool {
	for _, r := range n.Roles {
		if role == r {
			return true
		}
	}
	return false
}

func childrenOfRole(n *uast.Node, role uast.Role) []*uast.Node {
	var children []*uast.Node
	for _, child := range n.Children {
		if containsRole(child, role) {
			children = append(children, child)
		}
	}
	return children
}

func deepChildrenOfRole(n *uast.Node, role uast.Role) []*uast.Node {
	var childList []*uast.Node
	for _, child := range n.Children {
		if containsRole(child, role) {
			childList = append(childList, child)
		}
		childList = append(childList, deepChildrenOfRole(child, role)...)
	}
	return childList
}

func countChildrenOfRole(n *uast.Node, role uast.Role) int {
	count := 0
	for _, child := range n.Children {
		if containsRole(child, role) {
			count++
		}
	}
	return count
}

func deepCountChildrenOfRole(n *uast.Node, role uast.Role) int {
	count := 0
	for _, child := range n.Children {
		if containsRole(child, role) {
			count++
		}
		count += deepCountChildrenOfRole(child, role)
	}
	return count
}
