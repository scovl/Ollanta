package treesitter

import (
	"embed"

	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MetaFS holds the JSON metadata files for all tree-sitter rules in this package.
//
//go:embed *.json
var MetaFS embed.FS

func init() {
	ollantarules.MustRegister(MetaFS, "*.json",
		BroadExceptPY,
		ComparisonToNonePY,
		EqEqEqJS,
		MutableDefaultArgumentPY,
		NoConsoleLogJS,
		NoLargeFunctionsJS,
		NoLargeFunctionsPY,
		TooManyParametersJS,
		TooManyParametersPY,
		// Wave 1 — Python
		UselessEqEqPY,
		UselessComparisonPY,
		DictModifyIteratingPY,
		ListModifyIteratingPY,
		ReturnInInitPY,
		PassBodyPY,
		HardcodedTmpPathPY,
		UnspecifiedOpenEncodingPY,
		// Wave 1 — JavaScript / TypeScript
		UselessEqEqJS,
		UselessAssignJS,
		LeftoverDebuggingJS,
		AssignedUndefinedJS,
		DetectEvalJS,
		MomentDeprecatedTS,
		// Wave 2 — Python
		InsecureHashPY,
		DangerousSubprocessPY,
		DangerousOsExecPY,
		SyncSleepInAsyncPY,
		MissingHashWithEqPY,
		UncheckedReturnsPY,
		OpenNeverClosedPY,
		UseDefusedXmlPY,
		// Wave 2 — JavaScript / TypeScript
		DetectChildProcessJS,
		DetectInsecureWebsocketJS,
		DetectPseudoRandomBytesJS,
		UselessTernaryTS,
		// Wave 3 — Python
		AvoidPyyamlLoadPY,
		PicklePY,
		MarshalPY,
		UnverifiedSSLContextPY,
		RegexDosPY,
		// Wave 3 — JavaScript
		DetectRedosJS,
		PathJoinResolveTraversalJS,
		SpawnGitCloneJS,
		IncompleteSanitizationJS,
	)
}
