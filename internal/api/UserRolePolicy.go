package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

const policy = `
package example

default allow = false

allow  {
    input.user.role=="admin"
}
`
const inputData = `
{
    "user":{
        "role": "admin"
    }
    
}
`

func ExecuteUserRolePolicy() bool {

	//policyBytes, err := os.ReadFile("./policy/userRole.rego")

	// if err != nil {

	// 	fmt.Println("Error reading policy:", err)

	// 	return false

	// }

	// policy := string(policyBytes)

	//fmt.Println("policy = ", policy)
	// Read the input data from the JSON file

	// inputDataBytes, err := os.ReadFile("./policy/userRoleInput.json")

	// if err != nil {

	// 	fmt.Println("Error reading input data:", err)

	// 	return false

	// }

	// inputData := string(inputDataBytes)

	//fmt.Println("inputData = ", inputData)
	// Create a Rego query and policy evaluation

	query := "data.example.allow"

	// a1 := `{

	// 		"user": "nayeem"

	// }`
	regoPolicy := rego.New(

		rego.Query(query),

		rego.Module("userRole.rego", policy),

		//rego.Input([]byte(inputData)),
	)
	//fmt.Println("regopolicy = ", regoPolicy)
	// Evaluate the policy
	ctx := context.Background()
	preparequery, err := regoPolicy.PrepareForEval(ctx)
	if err != nil {
		fmt.Println("Policy evaluation error:", err)

		return false

	}

	var prepareInput interface{}

	err = json.NewDecoder(bytes.NewBufferString(inputData)).Decode(&prepareInput)
	if err != nil {

		fmt.Println("Policy evaluation error:", err)

		return false

	}
	results, nil := preparequery.Eval(ctx, rego.EvalInput(prepareInput))
	//results, err := regoPolicy.Eval(ctx)

	if err != nil {

		fmt.Println("Policy evaluation error:", err)

		return false

	}
	//fmt.Println("results = ", results[0].Expressions[0].Value.(bool))
	// Check the result

	if len(results) == 1 && results[0].Expressions[0].Value.(bool) {

		fmt.Println("Access granted!")
		return true

	} else {

		fmt.Println("Access denied.")
		return false
	}

}
