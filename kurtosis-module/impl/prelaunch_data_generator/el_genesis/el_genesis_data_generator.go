package el_genesis

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-core-api-lib/api/golang/lib/services"
	"github.com/kurtosis-tech/stacktrace"
	"os"
	"text/template"
	"time"
)

const (
	// The prefix that the directory for containing information about this EL genesis generation run will have
	//  inside the shared directory
	elGenesisGenerationInstanceSharedDirpathPrefix = "el-genesis-"

	configDirname                      = "config"

	outputDirname = "output"
)


func GenerateELGenesisData(
	serviceCtx *services.ServiceContext,
	chainspecAndGethGenesisGenerationConfigTemplate *template.Template,
	nethermindGenesisConfigJsonTemplate *template.Template,
	genesisUnixTimestamp uint64,
	networkId string,
	depositContractAddress string,
	totalTerminalDifficulty uint64,
) (
	*ELGenesisData,
	error,
) {
	sharedDir := serviceCtx.GetSharedDirectory()
	generationInstanceSharedDir := sharedDir.GetChildPath(fmt.Sprintf(
		"%v%v",
		elGenesisGenerationInstanceSharedDirpathPrefix,
		time.Now().Unix(),
	))
	configSharedDir := generationInstanceSharedDir.GetChildPath(configDirname)
	outputSharedDir := generationInstanceSharedDir.GetChildPath(outputDirname)

	allSharedDirsToCreate := []*services.SharedPath{
		generationInstanceSharedDir,
		configSharedDir,
		outputSharedDir,
	}
	for _, sharedDirToCreate := range allSharedDirsToCreate {
		toCreateDirpathOnModuleContainer := sharedDirToCreate.GetAbsPathOnThisContainer()
		if err := os.Mkdir(toCreateDirpathOnModuleContainer, os.ModePerm); err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred creating directory '%v'", toCreateDirpathOnModuleContainer)
		}
	}

	chainspecFilepathOnModuleContainer, gethGenesisFilepathOnModuleContaienr, err := generateChainspecAndGethGenesis(
		chainspecAndGethGenesisGenerationConfigTemplate,
		configSharedDir,
		networkId,
		genesisUnixTimestamp,
		depositContractAddress,
		totalTerminalDifficulty,
		serviceCtx,
		outputSharedDir,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred generating the chainspec & Geth config files")
	}

	nethermindGenesisFilepathOnModuleContainer, err := generateNethermindGenesis(
		nethermindGenesisConfigJsonTemplate,
		networkId,
		genesisUnixTimestamp,
		totalTerminalDifficulty,
		serviceCtx,
		outputSharedDir,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred generating the Nethermind genesis file")
	}
	
	result := newELGenesisData(
		outputSharedDir.GetAbsPathOnThisContainer(),
		chainspecFilepathOnModuleContainer,
		gethGenesisFilepathOnModuleContaienr,
		nethermindGenesisFilepathOnModuleContainer,
	)

	return result, nil
}
