package provider

import (
	"context"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = &EvaluateFunction{}
)

func NewEvaluateFunction() function.Function {
	return EvaluateFunction{}
}

type EvaluateFunction struct {
}

func (j EvaluateFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "evaluate"
}

func (j EvaluateFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "evaluate",
		MarkdownDescription: "Evaluates the provided string as Jsonnet",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "jsonnet",
				MarkdownDescription: "The Jsonnet value to be evaluated",
			},
		},
		Return: function.StringReturn{},
	}
}

func (j EvaluateFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data))

	if resp.Error != nil {
		return
	}

	vm := jsonnet.MakeVM()
	jsonStr, err := vm.EvaluateAnonymousSnippet("main.jsonnet", data)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, jsonStr))
}
