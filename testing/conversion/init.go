package conversion

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	"github.com/deckhouse/deckhouse/testing/library"
	"github.com/deckhouse/deckhouse/testing/library/values_store"
)

type Converter struct {
	moduleName string
	modulePath string
	values     *values_store.ValuesStore

	FinalValues  *values_store.ValuesStore
	FinalVersion string

	Error error
}

func (c *Converter) ValuesGet(path string) library.KubeResult {
	return c.values.Get(path)
}

func (c *Converter) ValuesSet(path string, value interface{}) {
	c.values.SetByPath(path, value)
}

func (c *Converter) ValuesSetFromYaml(path, value string) {
	c.values.SetByPathFromYAML(path, []byte(value))
}

func (c *Converter) Convert(version string) {
	conv := conversion.Registry().Get(c.moduleName, version)
	Expect(conv).ShouldNot(BeNil(), "Conversion for module %s and version %s should be registered", c.moduleName, version)

	convVer, convValues, convError := conv.Convert(version, c.values.Values)

	convValuesJSON, err := json.Marshal(convValues)
	if err != nil {
		c.Error = err
		return
	}

	c.FinalValues = values_store.NewStoreFromRawJSON(convValuesJSON)
	c.FinalVersion = convVer
	c.Error = convError
}

func (c *Converter) ConvertToLatest(fromVersion string) {
	hasModule := conversion.Registry().HasModule(c.moduleName)
	Expect(hasModule).Should(BeTrue(), "Module %s should have registered conversions", c.moduleName)

	hasVersion := conversion.Registry().HasVersion(c.moduleName, fromVersion)
	Expect(hasVersion).Should(BeTrue(), "Module %s should have registered conversion for version %s", c.moduleName, hasVersion)

	convVer, convValues, convError := conversion.ConvertToLatest(c.moduleName, fromVersion, c.values.Values)

	convValuesJSON, err := json.Marshal(convValues)
	if err != nil {
		c.Error = err
		return
	}

	c.FinalValues = values_store.NewStoreFromRawJSON(convValuesJSON)
	c.FinalVersion = convVer
	c.Error = convError
}

func SetupConverter(values string) *Converter {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	modulePath := filepath.Dir(wd)

	moduleName, err := library.GetModuleNameByPath(modulePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	moduleName = strcase.ToLowerCamel(moduleName)
	//moduleValuesKey := addonutils.ModuleNameToValuesKey(moduleName)

	initialValues, err := library.InitValues(modulePath, []byte(values))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initialValuesJSON, err := json.Marshal(initialValues)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	converter := &Converter{
		moduleName: moduleName,
		modulePath: modulePath,
	}

	BeforeEach(func() {
		converter.values = values_store.NewStoreFromRawJSON(initialValuesJSON)
	})

	return converter
}
