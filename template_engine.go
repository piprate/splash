package splash

import (
	"bytes"
	"embed"
	"fmt"

	"path"
	"strings"
	"text/template"

	"github.com/onflow/cadence/format"
	"github.com/onflow/flow-go-sdk"
	"github.com/rs/zerolog/log"
)

var (
	funcMap = template.FuncMap{
		// increment function
		"inc": func(i int) int {
			return i + 1
		},
		// decrement function
		"dec": func(i int) int {
			return i - 1
		},
		// turn a string into Cadence safe form
		"safe": func(v string) string {
			return format.String(v)
		},
		"ufix64": UFix64ToString,
	}
)

type (
	TemplateEngine struct {
		client                   *Connector
		template                 *template.Template
		preloadedTemplates       map[string]string
		wellKnownAddresses       map[string]string
		wellKnownAddressesBinary map[string]flow.Address
	}
)

const (
	ParamsKey = "Parameters"
)

func NewTemplateEngine(client *Connector, templateFS embed.FS, paths []string, requiredWellKnownContracts []string) (*TemplateEngine, error) {
	goTemplate, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/transactions/*.cdc", "templates/scripts/*.cdc", "templates/scripts/**/*.cdc")
	if err != nil {
		return nil, err
	}

	eng := &TemplateEngine{
		client:                   client,
		template:                 goTemplate,
		preloadedTemplates:       make(map[string]string),
		wellKnownAddresses:       make(map[string]string),
		wellKnownAddressesBinary: make(map[string]flow.Address),
	}

	if err := eng.loadContractAddresses(requiredWellKnownContracts); err != nil {
		return nil, err
	}

	return eng, nil
}

func (e *TemplateEngine) loadContractAddresses(requiredWellKnownContracts []string) error {
	contracts := e.client.State.Contracts()
	network := e.client.Services.Network()
	networkName := network.Name
	deployedContracts, err := e.client.State.DeploymentContractsByNetwork(network)
	if err != nil {
		return err
	}
	for _, contract := range *contracts {
		for _, alias := range contract.Aliases {
			if alias.Network == networkName {
				e.wellKnownAddressesBinary[contract.Name] = alias.Address
			}
		}
	}
	for _, contract := range deployedContracts {
		e.wellKnownAddressesBinary[strings.Split(path.Base(contract.Location()), ".")[0]] = contract.AccountAddress
	}

	for _, requiredContractName := range requiredWellKnownContracts {
		if _, found := e.wellKnownAddressesBinary[requiredContractName]; !found {
			return fmt.Errorf("address not found for contract %s", requiredContractName)
		}
	}
	log.Debug().Str("addresses", fmt.Sprintf("%v", e.wellKnownAddresses)).Msg("Loaded contract addresses")

	for name, addr := range e.wellKnownAddressesBinary {
		e.wellKnownAddresses[name] = addr.HexWithPrefix()
	}

	return nil
}

func (e *TemplateEngine) WellKnownAddresses() map[string]string {
	return e.wellKnownAddresses
}

func (e *TemplateEngine) ContractAddress(contractName string) flow.Address {
	return e.wellKnownAddressesBinary[contractName]
}

func (e *TemplateEngine) GetStandardScript(scriptID string) string {
	s, found := e.preloadedTemplates[scriptID]
	if !found {
		buf := &bytes.Buffer{}
		if err := e.template.ExecuteTemplate(buf, scriptID, e.wellKnownAddresses); err != nil {
			panic(err)
		}

		s = buf.String()
		e.preloadedTemplates[scriptID] = s
	}

	return s
}

func (e *TemplateEngine) GetCustomScript(scriptID string, params interface{}) string {
	data := map[string]interface{}{
		ParamsKey: params,
	}
	for k, v := range e.wellKnownAddresses {
		data[k] = v
	}
	buf := &bytes.Buffer{}
	if err := e.template.ExecuteTemplate(buf, scriptID, data); err != nil {
		panic(err)
	}

	return buf.String()
}

func (e *TemplateEngine) NewTransaction(scriptID string) FlowTransactionBuilder {
	return e.client.Transaction(e.GetStandardScript(scriptID))
}

func (e *TemplateEngine) NewInlineTransaction(script string) FlowTransactionBuilder {
	return e.client.Transaction(script)
}

func (e *TemplateEngine) NewScript(scriptID string) FlowScriptBuilder {
	return e.client.Script(e.GetStandardScript(scriptID))
}

func (e *TemplateEngine) NewInlineScript(script string) FlowScriptBuilder {
	return e.client.Script(script)
}
