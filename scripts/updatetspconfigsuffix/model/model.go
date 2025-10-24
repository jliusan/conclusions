package model

type YamlConfig struct {
	Parameters Parameters `yaml:"parameters"`
	Emit       []string   `yaml:"emit"`
	Linter     Linter     `yaml:"linter"`
	Options    Options    `yaml:"options"`
}

type Parameters struct {
	ServiceDir           DefaultParam `yaml:"service-dir"`
	ServiceDirectoryName DefaultParam `yaml:"service-directory-name"`
}

type DefaultParam struct {
	Default string `yaml:"default"`
}

type Linter struct {
	Extends []string `yaml:"extends"`
}

type Options struct {
	AutoRest AutoRestOptions `yaml:"@azure-tools/typespec-autorest"`
	CSharp   CSharpOptions   `yaml:"@azure-tools/typespec-csharp"`
	Go       GoOptions       `yaml:"@azure-tools/typespec-go"`
	Java     JavaOptions     `yaml:"@azure-tools/typespec-java"`
	Python   PythonOptions   `yaml:"@azure-tools/typespec-python"`
	TS       TSOptions       `yaml:"@azure-tools/typespec-ts"`
}

type AutoRestOptions struct {
	AzureResourceProviderFolder string `yaml:"azure-resource-provider-folder"`
	EmitCommonTypesSchema       string `yaml:"emit-common-types-schema"`
	ArmResourceFlattening       bool   `yaml:"arm-resource-flattening"`
	EmitterOutputDir            string `yaml:"emitter-output-dir"`
	OutputFile                  string `yaml:"output-file"`
	ArmTypesDir                 string `yaml:"arm-types-dir"`
	OmitUnreachableTypes        bool   `yaml:"omit-unreachable-types"`
	UseReadOnlyStatusSchema     bool   `yaml:"use-read-only-status-schema"`
}

type CSharpOptions struct {
	EmitterOutputDir  string `yaml:"emitter-output-dir"`
	Flavor            string `yaml:"flavor"`
	ClearOutputFolder bool   `yaml:"clear-output-folder"`
	ModelNamespace    bool   `yaml:"model-namespace"`
	Namespace         string `yaml:"namespace"`
}

type GoOptions struct {
	ServiceDir         string `yaml:"service-dir"`
	EmitterOutputDir   string `yaml:"emitter-output-dir"`
	Module             string `yaml:"module"`
	FixConstStuttering bool   `yaml:"fix-const-stuttering"`
	Flavor             string `yaml:"flavor"`
	GenerateSamples    bool   `yaml:"generate-samples"`
	GenerateFakes      bool   `yaml:"generate-fakes"`
	HeadAsBoolean      bool   `yaml:"head-as-boolean"`
	InjectSpans        bool   `yaml:"inject-spans"`
}

type JavaOptions struct {
	EmitterOutputDir      string `yaml:"emitter-output-dir"`
	Namespace             string `yaml:"namespace"`
	ServiceName           string `yaml:"service-name"`
	Flavor                string `yaml:"flavor"`
	UseObjectForUnknown   bool   `yaml:"use-object-for-unknown"`
	ClientSideValidations bool   `yaml:"client-side-validations"`
}

type PythonOptions struct {
	ServiceDir       string `yaml:"service-dir"`
	EmitterOutputDir string `yaml:"emitter-output-dir"`
	Namespace        string `yaml:"namespace"`
	Flavor           string `yaml:"flavor"`
	GenerateTest     bool   `yaml:"generate-test"`
	GenerateSample   bool   `yaml:"generate-sample"`
}

type TSOptions struct {
	TypespecTitleMap       map[string]string `yaml:"typespec-title-map"`
	ExperimentalExtensible bool              `yaml:"experimental-extensible-enums"`
	EmitterOutputDir       string            `yaml:"emitter-output-dir"`
	Flavor                 string            `yaml:"flavor"`
	PackageDetails         PackageDetails    `yaml:"package-details"`
}

type PackageDetails struct {
	Name string `yaml:"name"`
}
