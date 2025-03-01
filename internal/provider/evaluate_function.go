package provider

import (
	"context"
	"fmt"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"path/filepath"
	jsonnetUtil "terraform-provider-jsonnet/internal/jsonnet"
)

var (
	_ function.Function = &EvaluateFunction{}
)

func NewEvaluateFunction() function.Function {
	return EvaluateFunction{}
}

type EvaluateFunction struct {
}

type EvaluateOptions struct {
	JPaths []string `tfsdk:"jpaths"`
}

func (o EvaluateOptions) merge(other EvaluateOptions) EvaluateOptions {
	return EvaluateOptions{
		JPaths: append(o.JPaths, other.JPaths...),
	}
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
				MarkdownDescription: "The Jsonnet code to be evaluated",
				Name:                "code",
			},
		},
		VariadicParameter: function.ObjectParameter{
			AttributeTypes: map[string]attr.Type{
				"jpaths": types.ListType{
					ElemType: types.StringType,
				},
			},
			MarkdownDescription: "Additional options to be passed to the Jsonnet VM",
			Name:                "options",
		},
		Return: function.StringReturn{},
	}
}

func (j EvaluateFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var code string
	var options []EvaluateOptions
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &code, &options))

	if resp.Error != nil {
		return
	}

	mergedOptions := EvaluateOptions{}
	jsonnetPath := filepath.SplitList(os.Getenv("JSONNET_PATH"))
	for i := len(jsonnetPath) - 1; i >= 0; i-- {
		mergedOptions.JPaths = append(mergedOptions.JPaths, jsonnetPath[i])
	}

	for _, o := range options {
		mergedOptions = mergedOptions.merge(o)
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: mergedOptions.JPaths,
	})
	vm.NativeFunction(jsonnetUtil.UuidV5())
	preamble := `
      local stdTf = std {
        tf: {
          uuidv5: std.native('uuidv5')
        }
      };
      local std = stdTf;
	`
	snippet := fmt.Sprintf("%s%s", preamble, code)
	jsonStr, err := vm.EvaluateAnonymousSnippet("main.jsonnet", snippet)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, jsonStr))
}
