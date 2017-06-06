package tools

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/bblfsh/sdk/protocol"
	"github.com/bblfsh/sdk/uast"
	"github.com/stretchr/testify/require"
)

func TestExpresionComplex(t *testing.T) {
	require := require.New(t)

	n := &uast.Node{InternalType: "ifCondition", Roles: []uast.Role{uast.IfCondition}, Children: []*uast.Node{
		{InternalType: "bool_and", Roles: []uast.Role{uast.OpBooleanAnd}},
		{InternalType: "bool_xor", Roles: []uast.Role{uast.OpBooleanXor}},
	}}
	n2 := &uast.Node{InternalType: "ifCondition", Roles: []uast.Role{uast.IfCondition}, Children: []*uast.Node{
		{InternalType: "bool_and", Roles: []uast.Role{uast.OpBooleanAnd}, Children: []*uast.Node{
			{InternalType: "bool_or", Roles: []uast.Role{uast.OpBooleanOr}, Children: []*uast.Node{
				{InternalType: "bool_xor", Roles: []uast.Role{uast.OpBooleanXor}},
			}},
		}},
	}}

	result := expressionComp(n)
	expect := 2
	require.Equal(expect, result)

	result = expressionComp(n2)
	expect = 3
	require.Equal(expect, result)
}

func TestNpathComplexity(t *testing.T) {
	require := require.New(t)
	var result []int
	var expect []int

	andBool := &uast.Node{InternalType: "bool_and", Roles: []uast.Role{uast.OpBooleanAnd}}
	orBool := &uast.Node{InternalType: "bool_or", Roles: []uast.Role{uast.OpBooleanOr}}
	statement := &uast.Node{InternalType: "Statement", Roles: []uast.Role{uast.Statement}}

	n := &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		statement,
	}}

	npathData := NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 1)
	/*
			if(3conditions){
				Statement
				Statement
			}else if(3conditions){
				Statement
				Statement
		  }else{
				Statement
				Statement
		  } Npath = 7
	*/
	ifCondition := &uast.Node{InternalType: "Condition", Roles: []uast.Role{uast.IfCondition}, Children: []*uast.Node{
		andBool,
		orBool,
	}}
	ifBody := &uast.Node{InternalType: "Body", Roles: []uast.Role{uast.IfBody}, Children: []*uast.Node{
		statement,
		statement,
	}}
	elseIf := &uast.Node{InternalType: "elseIf", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{
		&uast.Node{InternalType: "If", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
			ifCondition,
			ifBody,
		}},
	}}
	ifElse := &uast.Node{InternalType: "else", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{
		ifBody,
	}}
	nIf := &uast.Node{InternalType: "if", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
		elseIf,
		ifElse,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nIf,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 7)

	// This case looks like the previous one, but we have the ElseIF and the If Roles in the same uast.Node
	/*
		if(3conditions){
			Statement
			Statement
		}else if(3conditions){
			Statement
			Statement
		}else{
			Statement
			Statement
		} Npath = 7
	*/

	elseIf2Roles := &uast.Node{InternalType: "elseIf", Roles: []uast.Role{uast.IfElse, uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
	}}

	nIf2Roles := &uast.Node{InternalType: "if", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
		elseIf2Roles,
		ifElse,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nIf2Roles,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 7)

	/*
	  if(condition){
	    Statement
	    Statement
	  }Npath = 2
	*/
	nSimpleIf := &uast.Node{InternalType: "If", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
		{InternalType: "ifCondition", Roles: []uast.Role{uast.IfCondition}, Children: []*uast.Node{}},
		ifBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nSimpleIf,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	/*
		The same if structure of the example above
		but repeated three times, in sequencial way
		Npath = 343
	*/

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nIf,
		nIf,
		nIf,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 343)

	/*
		if(3conditions){
			if(3conditions){
				if(3conditions){
					Statement
					Statemenet
				}else{
					Statement
					Statement
				}
			}
		} Npath = 10
	*/
	nestedIfBody := &uast.Node{InternalType: "bodyº", Roles: []uast.Role{uast.IfBody}, Children: []*uast.Node{
		{InternalType: "if2", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
			ifCondition,
			{InternalType: "body2", Roles: []uast.Role{uast.IfBody}, Children: []*uast.Node{
				{InternalType: "if3", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
					ifCondition,
					ifBody,
					ifElse,
				}},
			}},
		}},
	}}
	nNestedIf := &uast.Node{InternalType: "if1", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
		ifCondition,
		nestedIfBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nNestedIf,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 10)

	/*
		while(2condtions){
			Statement
			Statement
			Statement
		}else{
			Statement
			Statement
		} Npath = 3
	*/
	whileCondition := &uast.Node{InternalType: "WhileCondition", Roles: []uast.Role{uast.WhileCondition}, Children: []*uast.Node{
		andBool,
	}}
	whileBody := &uast.Node{InternalType: "WhileBody", Roles: []uast.Role{uast.WhileBody}, Children: []*uast.Node{
		statement,
		statement,
		statement,
	}}
	whileElse := &uast.Node{InternalType: "WhileElse", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nWhile := &uast.Node{InternalType: "While", Roles: []uast.Role{uast.While}, Children: []*uast.Node{
		whileCondition,
		whileBody,
		whileElse,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nWhile,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 3)

	/*
		while(2conditions){
			while(2conditions){
				while(2conditions){
					Statement
					Statement
				}
			}
		} Npath = 7
	*/
	nestedWhileBody := &uast.Node{InternalType: "WhileBody1", Roles: []uast.Role{uast.WhileBody}, Children: []*uast.Node{
		{InternalType: "While2", Roles: []uast.Role{uast.While}, Children: []*uast.Node{
			whileCondition,
			{InternalType: "WhileBody2", Roles: []uast.Role{uast.WhileBody}, Children: []*uast.Node{
				{InternalType: "While3", Roles: []uast.Role{uast.While}, Children: []*uast.Node{
					whileCondition,
					whileBody,
				}},
			}},
		}},
	}}
	nNestedWhile := &uast.Node{InternalType: "While1", Roles: []uast.Role{uast.While}, Children: []*uast.Node{
		whileCondition,
		nestedWhileBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nNestedWhile,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 7)

	/*
			 for(init;2condition;update){
			 	Statement
				Statement
		 	 } Npath = 2
	*/
	forCondition := &uast.Node{InternalType: "forCondition", Roles: []uast.Role{uast.ForExpression}, Children: []*uast.Node{
		orBool,
	}}
	forInit := &uast.Node{InternalType: "forInit", Roles: []uast.Role{uast.ForInit}}
	forUpdate := &uast.Node{InternalType: "forUpdate", Roles: []uast.Role{uast.ForUpdate}}
	forBody := &uast.Node{InternalType: "forBody", Roles: []uast.Role{uast.ForBody}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nFor := &uast.Node{InternalType: "for", Roles: []uast.Role{uast.For}, Children: []*uast.Node{
		forInit,
		forCondition,
		forUpdate,
		forBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nFor,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	/*
		for(init;2conditions;update){
			for(init;2conditions;update){
				for(init;2condtions;update){
					Statement
					Statement
				}
			}
		} Npath = 4
	*/
	nestedForBody := &uast.Node{InternalType: "forBody1", Roles: []uast.Role{uast.ForBody}, Children: []*uast.Node{
		{InternalType: "for2", Roles: []uast.Role{uast.For}, Children: []*uast.Node{
			forInit,
			forCondition,
			forUpdate,
			{InternalType: "forBody2", Roles: []uast.Role{uast.ForBody}, Children: []*uast.Node{
				{InternalType: "for3", Roles: []uast.Role{uast.For}, Children: []*uast.Node{
					forInit,
					forCondition,
					forUpdate,
					forBody,
				}},
			}},
		}},
	}}
	nNestedFor := &uast.Node{InternalType: "for1", Roles: []uast.Role{uast.For}, Children: []*uast.Node{
		forInit,
		forCondition,
		forUpdate,
		nestedForBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nNestedFor,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 4)

	/*
		do{
			Statement
			Statement
		}while(3conditions)
		Npath = 4
	*/
	doWhileCondition := &uast.Node{InternalType: "doWhileCondition", Roles: []uast.Role{uast.DoWhileCondition}, Children: []*uast.Node{
		orBool,
		orBool,
	}}
	doWhileBody := &uast.Node{InternalType: "doWhileBody", Roles: []uast.Role{uast.DoWhileBody}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nDoWhile := &uast.Node{InternalType: "doWhile", Roles: []uast.Role{uast.DoWhile}, Children: []*uast.Node{
		doWhileBody,
		doWhileCondition,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nDoWhile,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 4)

	/*
		do{
			do{
				do{
					Statement
					Statement
				}while(3conditions)
			}while{3conditions}
		}while(3condtions)
		Npath = 10
	*/
	nestedDoWhileBody := &uast.Node{InternalType: "doWhileBody1", Roles: []uast.Role{uast.DoWhileBody}, Children: []*uast.Node{
		{InternalType: "doWhile2", Roles: []uast.Role{uast.DoWhile}, Children: []*uast.Node{
			{InternalType: "doWhileBody2", Roles: []uast.Role{uast.DoWhileBody}, Children: []*uast.Node{
				{InternalType: "doWhile3", Roles: []uast.Role{uast.DoWhile}, Children: []*uast.Node{
					doWhileBody,
					doWhileCondition,
				}},
			}},
			doWhileCondition,
		}},
	}}
	nNestedDoWhile := &uast.Node{InternalType: "doWhile1", Roles: []uast.Role{uast.DoWhile}, Children: []*uast.Node{
		nestedDoWhileBody,
		doWhileCondition,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nNestedDoWhile,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 10)

	/*
		switch(){
		case:
			Statement
			Statement
		case:
			Statement
			Statement
		default:
			Statement
			Statement
		} Npath = 3
	*/
	switchCondition := &uast.Node{InternalType: "switchCondition", Roles: []uast.Role{uast.SwitchCaseCondition}, Children: []*uast.Node{
		orBool,
		andBool,
	}}
	switchCaseBody := &uast.Node{InternalType: "switchCaseBody", Roles: []uast.Role{uast.SwitchCaseBody}, Children: []*uast.Node{
		statement,
		statement,
	}}
	switchCase := &uast.Node{InternalType: "switchCase", Roles: []uast.Role{uast.SwitchCase}, Children: []*uast.Node{
		switchCondition,
	}}
	defaultCase := &uast.Node{InternalType: "defaultCase", Roles: []uast.Role{uast.SwitchDefault}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nSwitch := &uast.Node{InternalType: "switch", Roles: []uast.Role{uast.Switch}, Children: []*uast.Node{
		switchCase,
		switchCaseBody,
		switchCase,
		switchCaseBody,
		defaultCase,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nSwitch,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 3)

	/*
		switch(){
		case:
			Statement
			Statement
		case:
			Statement
			Statement
		default:
			switch(){
			case:
				Statement
				Statement
			case:
				Statement
				Statement
			default:
				Statement
				Statement
		} Npath = 9
	*/
	nestedDefaultCase := &uast.Node{InternalType: "defaultCase", Roles: []uast.Role{uast.SwitchDefault}, Children: []*uast.Node{
		nSwitch,
	}}
	nNestedSwitch := &uast.Node{InternalType: "switch", Roles: []uast.Role{uast.Switch}, Children: []*uast.Node{
		switchCase,
		switchCaseBody,
		switchCase,
		switchCaseBody,
		nestedDefaultCase,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nNestedSwitch,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 5)

	/*
		return
	*/
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		{InternalType: "Return", Roles: []uast.Role{uast.Return}},
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 1)

	/*
		statement
		statement
		return 3condition
	*/
	nReturn := &uast.Node{InternalType: "Return", Roles: []uast.Role{uast.Return}, Children: []*uast.Node{
		orBool,
		andBool,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		statement,
		statement,
		nReturn,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	nForEach := &uast.Node{InternalType: "ForEach", Roles: []uast.Role{uast.ForEach}, Children: []*uast.Node{
		forInit,
		forCondition,
		forBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nForEach,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	tryBody := &uast.Node{InternalType: "TryBody", Roles: []uast.Role{uast.TryBody}, Children: []*uast.Node{
		statement,
		statement,
	}}

	tryCatch := &uast.Node{InternalType: "TryCatch", Roles: []uast.Role{uast.TryCatch}, Children: []*uast.Node{
		statement,
		statement,
	}}

	nTry := &uast.Node{InternalType: "Try", Roles: []uast.Role{uast.Try}, Children: []*uast.Node{
		tryBody,
		tryCatch,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nTry,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	tryFinally := &uast.Node{InternalType: "TryFinally", Roles: []uast.Role{uast.TryFinally}, Children: []*uast.Node{
		nSimpleIf,
	}}
	nTryFinally := &uast.Node{InternalType: "Try", Roles: []uast.Role{uast.Try}, Children: []*uast.Node{
		tryBody,
		tryCatch,
		tryCatch,
		tryCatch,
		tryFinally,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{
		nTryFinally,
	}}

	npathData = NpathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 6)

	require.Equal(expect, result)
}

func TestNpathMultiFunc(t *testing.T) {
	require := require.New(t)
	var result []int
	expect := []int{7, 7, 7}
	andBool := &uast.Node{InternalType: "bool_and", Roles: []uast.Role{uast.OpBooleanAnd}}
	orBool := &uast.Node{InternalType: "bool_or", Roles: []uast.Role{uast.OpBooleanOr}}
	statement := &uast.Node{InternalType: "Statement", Roles: []uast.Role{uast.Statement}}

	ifCondition := &uast.Node{InternalType: "Condition", Roles: []uast.Role{uast.IfCondition}, Children: []*uast.Node{
		andBool,
		orBool,
	}}
	ifBody := &uast.Node{InternalType: "Body", Roles: []uast.Role{uast.IfBody}, Children: []*uast.Node{
		statement,
		statement,
	}}
	elseIf := &uast.Node{InternalType: "elseIf", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{
		&uast.Node{InternalType: "If", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
			ifCondition,
			ifBody,
		}},
	}}
	ifElse := &uast.Node{InternalType: "else", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{
		ifBody,
	}}
	nIf := &uast.Node{InternalType: "if", Roles: []uast.Role{uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
		elseIf,
		ifElse,
	}}
	funcBody := &uast.Node{InternalType: "funcBody", Roles: []uast.Role{uast.FunctionDeclarationBody}, Children: []*uast.Node{nIf}}

	func1 := &uast.Node{InternalType: "func1", Roles: []uast.Role{uast.FunctionDeclaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName1", Roles: []uast.Role{uast.FunctionDeclarationName}, Children: []*uast.Node{}, Token: "Name1"},
		funcBody,
	}}
	func2 := &uast.Node{InternalType: "func2", Roles: []uast.Role{uast.FunctionDeclaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName2", Roles: []uast.Role{uast.FunctionDeclarationName}, Children: []*uast.Node{}, Token: "Name2"},
		funcBody,
	}}
	func3 := &uast.Node{InternalType: "func3", Roles: []uast.Role{uast.FunctionDeclaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName3", Roles: []uast.Role{uast.FunctionDeclarationName}, Children: []*uast.Node{}, Token: "Name3"},
		funcBody,
	}}

	n := &uast.Node{InternalType: "module", Children: []*uast.Node{
		func1,
		func2,
		func3,
	}}
	npathData := NpathComplexity(n)
	for _, v := range npathData {
		result = append(result, v.Complexity)
	}
	require.Equal(expect, result)
}
func TestZeroFunction(t *testing.T) {
	require := require.New(t)
	// Empty tree
	n := &uast.Node{InternalType: "module"}
	comp := NpathComplexity(n)
	require.Equal(0, len(comp))
}

func TestRealUAST(t *testing.T) {
	fileNames := []string{
		"fixtures/npath/ifelse.json",
		"fixtures/npath/do_while.json",
		"fixtures/npath/while.json",
		"fixtures/npath/for.json",
		"fixtures/npath/someFuncs.json",
		"fixtures/npath/switch.json",
	}

	require := require.New(t)
	var result []int
	for _, name := range fileNames {
		file, err := os.Open(name)
		require.NoError(err)

		reader := bufio.NewReader(file)
		dec := json.NewDecoder(reader)
		req := &protocol.ParseUASTResponse{}
		err = dec.Decode(req)
		require.NoError(err)
		n := req.UAST
		npathData := NpathComplexity(n)
		for _, v := range npathData {
			result = append(result, v.Complexity)
		}
	}

	expect := []int{2, 2, 2, 2, 2, 6, 2, 6, 3, 5, 4}

	require.Equal(expect, result)

}
