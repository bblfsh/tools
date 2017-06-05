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

	if n.ContainsRole(uast.FunctionDeclarationBody) {
		funcs = append(funcs, n)
		names = append(names, "NoName")
	} else {
		funcDecs := n.DeepChildrenOfRole(uast.FunctionDeclaration)
		for _, funcDec := range funcDecs {
			names = append(names, funcDec.ChildrenOfRole(uast.FunctionDeclarationName)[0].Token)
			funcs = append(funcs, funcDec.ChildrenOfRole(uast.FunctionDeclarationBody)[0])
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
	ifBody := n.ChildrenOfRole(uast.IfBody)
	ifCondition := n.ChildrenOfRole(uast.IfCondition)
	ifElse := n.ChildrenOfRole(uast.IfElse)

	if len(ifElse) == 0 {
		npath++
	} else {
		// This if is a short circuit to avoid the two roles in the switch problem
		if ifElse[0].ContainsRole(uast.If) {
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
	whileCondition := n.ChildrenOfRole(uast.WhileCondition)
	whileBody := n.ChildrenOfRole(uast.WhileBody)
	whileElse := n.ChildrenOfRole(uast.IfElse)
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
	doWhileCondition := n.ChildrenOfRole(uast.DoWhileCondition)
	doWhileBody := n.ChildrenOfRole(uast.DoWhileBody)

	npath *= complexityMultOf(doWhileBody[0])
	npath += expressionComp(doWhileCondition[0])

	return npath
}

func visitFor(n *uast.Node) int {
	// (npath of for + bool_comp of for + 1) * npath of next
	npath := 1
	forBody := n.ChildrenOfRole(uast.ForBody)

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

	caseDefault := n.ChildrenOfRole(uast.SwitchDefault)
	switchCases := n.ChildrenOfRole(uast.SwitchCase)
	switchCondition := n.ChildrenOfRole(uast.SwitchCaseCondition)
	npath := 0
	/*
		In pmd the expressionComp function returns always our value -1
		but in other places of the code the fuction works exactly as our function
		I suposed this happens because java AST differs with the UAST
	*/
	npath += expressionComp(switchCondition[0]) - 1
	if len(caseDefault) != 0 {
		npath += complexityMultOf(caseDefault[0])
	}
	for _, switchCase := range switchCases {
		npath += complexityMultOf(switchCase)
	}
	return npath
}

func visitForEach(n *uast.Node) int {
	forBody := n.ChildrenOfRole(uast.ForBody)
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

	tryBody := n.ChildrenOfRole(uast.TryBody)
	tryCatch := n.ChildrenOfRole(uast.TryCatch)
	tryFinaly := n.ChildrenOfRole(uast.TryFinally)

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
	orCount := n.DeepCountChildrenOfRole(uast.OpBooleanAnd)
	andCount := n.DeepCountChildrenOfRole(uast.OpBooleanOr)

	return orCount + andCount + 1
}
