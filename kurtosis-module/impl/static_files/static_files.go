package static_files

import (
	"github.com/kurtosis-tech/stacktrace"
	"path"
	"text/template"
)

const (
	// The path on the module container where static files are housed
	staticFilesDirpath = "/static-files"

	// Geth + CL genesis generation
	genesisGenerationConfigDirpath = staticFilesDirpath + "/genesis-generation-config"

	elGenesisGenerationConfigDirpath                        = genesisGenerationConfigDirpath + "/el"
	ChainspecAndGethGenesisGenerationConfigTemplateFilepath = elGenesisGenerationConfigDirpath + "/geth-genesis-config.yaml.tmpl"
	NethermindGenesisGenerationJsonTemplateFilepath         = elGenesisGenerationConfigDirpath + "/nethermind-genesis.json.tmpl"

	clGenesisGenerationConfigDirpath             = genesisGenerationConfigDirpath + "/cl"
	CLGenesisGenerationConfigTemplateFilepath    = clGenesisGenerationConfigDirpath + "/config.yaml.tmpl"
	CLGenesisGenerationMnemonicsTemplateFilepath = clGenesisGenerationConfigDirpath + "/mnemonics.yaml.tmpl"

	// Forkmon config
	ForkmonConfigTemplateFilepath = staticFilesDirpath + "/forkmon-config/config.toml.tmpl"
)

func ParseTemplate(filepath string) (*template.Template, error) {
	tmpl, err := template.New(
		// For some reason, the template name has to match the basename of the file:
		//  https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template
		path.Base(filepath),
	).ParseFiles(
		filepath,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred parsing template file '%v'", filepath)
	}
	return tmpl, nil
}
