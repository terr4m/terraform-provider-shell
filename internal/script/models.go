package script

const (
	LifecycleEnv            string = "TF_SCRIPT_LIFECYCLE"
	InputsEnv               string = "TF_SCRIPT_INPUTS"
	StateOutputEnv          string = "TF_SCRIPT_STATE_OUTPUT"
	ScriptOutputFilePathEnv string = "TF_SCRIPT_OUTPUT"
	ScriptErrorFilePathEnv  string = "TF_SCRIPT_ERROR"
)

// Lifecycle represents a Terraform lifecycle stage.
type Lifecycle string

const (
	LifecyclePlan   Lifecycle = "plan"
	LifecycleCreate Lifecycle = "create"
	LifecycleRead   Lifecycle = "read"
	LifecycleUpdate Lifecycle = "update"
	LifecycleDelete Lifecycle = "delete"
)
