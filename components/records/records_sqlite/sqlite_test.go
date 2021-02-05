package records_sqlite

import (
	"os"
	"testing"

	"github.com/pavlo67/tools/components/records"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/serializer"
)

const serviceName = "notebook"

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + serviceName + "." + env + ".yaml"
	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgSQLite := config.Access{}
	err = cfg.Value("sqlite", &cfgSQLite)
	require.NoError(t, err)

	l.Infof("%#v", cfgSQLite)

	dataOp, cleanerOp, err := New(cfgSQLite, "storage", "")
	require.NoError(t, err)

	l.Debugf("%#v", dataOp)

	records.OperatorTestScenario(t, l)
}
