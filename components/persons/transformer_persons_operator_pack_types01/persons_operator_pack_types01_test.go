package transformer_persons_operator_pack_types01

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/structures/types01"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/tools/components/persons"
	"github.com/pavlo67/tools/components/persons/persons_sqlite"
)

func TestTransformPersonsOperatorPack(t *testing.T) {

	return
	_, cfgService, l := apps.PrepareTests(t, "../../../apps/_environments/", "test", "")

	components := []starter.Starter{
		{db_sqlite.Starter(), nil},
		{persons_sqlite.Starter(), nil},
		{Starter(), common.Map{"persons_cleaner_key": persons.InterfaceCleanerKey}},
	}

	label := "PERSONS_OPERATOR_PACK/TEST BUILD"
	joinerOp, err := starter.Run(components, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	transformOp, _ := joinerOp.Interface(InterfaceKey).(transformer.Operator)
	require.NotNil(t, transformOp)

	dataInitial := structures.PackAny{
		PackDescription: structures.PackDescription{
			Fields: structures.Fields{},
			ItemDescription: structures.ItemDescription{
				URN:   ns.URN("test:test" + strconv.FormatInt(time.Now().UnixNano(), 10)),
				Label: "title",
				// ErrorsMap: nil,
				// History:   nil,
				CreatedAt: time.Now(),
				// UpdatedAt: nil,
			},
		},
		PackData: []types01.Person{
			{
				Nickname: "wqerwqer",
				Roles:    nil,
				Creds:    auth.Creds{auth.CredsEmail: "aaa@bbb.ccc"},

				ItemDescription: structures.ItemDescription{
					Info: common.Map{"xxx": "yyy", "zzz": 777},
					URN:  "urn1",
					// History:   nil,
					CreatedAt: time.Now(),
					// UpdatedAt: nil,
				},
			},
			{
				Nickname: "wqerwqer2",
				Roles:    rbac.Roles{rbac.RoleUser},
				Creds:    auth.Creds{auth.CredsEmail: "aaa2@bbb.ccc"},

				ItemDescription: structures.ItemDescription{
					Info: common.Map{"xxx2": "yyy", "zzz2": 222},
					URN:  "urn2",
					// History:   nil,
					CreatedAt: time.Now(),
					// UpdatedAt: nil,
				},
			},
		},
	}

	var params common.Map

	err = personsCleanerOp.Clean(nil)
	require.NoError(t, err)

	copyFinal, statFinal, dataFinal := transformer.TestOperator(t, transformOp, params, dataInitial, true)

	//copyFinal, _ := transformOp.Copy(nil, params)
	t.Logf("COPY (INTERNAL) FINAL: %#v", copyFinal)

	//statFinal, _ := transformOp.Stat(nil, params)
	if statFinalStringer, ok := statFinal.(fmt.Stringer); ok {
		t.Logf("STAT (INTERNAL) FINAL: %s", statFinalStringer.String())
	} else {
		t.Logf("STAT (INTERNAL) FINAL: %#v", statFinal)
	}

	//dataFinal, _ := transformOp.Out(nil, params)
	t.Logf("DATA (OUT) FINAL: %#v", dataFinal)

}
