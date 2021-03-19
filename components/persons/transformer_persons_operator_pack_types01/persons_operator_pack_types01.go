package transformer_persons_operator_pack_types01

import (
	"fmt"

	"github.com/pavlo67/data_exchange/components/structures"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures/types01"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/tools/components/persons"
)

var _ transformer.Operator = &transformerOperatorPackPersonsTypes01{}

type transformerOperatorPackPersonsTypes01 struct {
	personsOp persons.Operator
	identity  *auth.Identity
}

const onNew = "on transformerOperatorPackPersonsTypes01.New(): "

func New(personsOp persons.Operator, identity *auth.Identity) (transformer.Operator, error) {
	if personsOp == nil {
		return nil, errors.New(onNew + ": no persons.Operator")
	}
	//if personsOpCleaner == nil {
	//	return nil, errors.New(onNew + ": no persons.Cleaner")
	//}

	return &transformerOperatorPackPersonsTypes01{
		personsOp: personsOp,
		identity:  identity,
	}, nil
}

func (transformOp *transformerOperatorPackPersonsTypes01) Name() string {
	return string(InterfaceKey)
}

const onIn = "on transformerOperatorPackPersonsTypes01.In(): "

func (transformOp *transformerOperatorPackPersonsTypes01) In(pack structures.Pack, params common.Map) error {
	if pack == nil {
		return errors.New(onIn + "nil pack to import")
	}

	data := pack.Data()

	var persons01 []types01.Person

	switch v := data.(type) {
	case []types01.Person:
		persons01 = v
	case *[]types01.Person:
		if v != nil {
			persons01 = *v
		}
	default:
		return fmt.Errorf("wrong pack.Data() to import: %#v", data)
	}

	// TODO!!! clear old pack records

	for i, person01 := range persons01 {
		personItem := persons.Item{
			Identity: auth.Identity{
				Nickname: person01.Nickname,
				Roles:    person01.Roles,
			},
			ItemDescription: person01.ItemDescription,
		}
		personItem.SetCreds(person01.Creds)

		if _, err := transformOp.personsOp.Save(personItem, transformOp.identity); err != nil {
			return fmt.Errorf(onIn+": can't save item (%d / %#v), got %s", i, personItem, err)
		}
	}

	return nil
}

const onOut = "on transformerOperatorPackPersonsTypes01.Out(): "

func (transformOp *transformerOperatorPackPersonsTypes01) Out(selector *selectors.Term, params common.Map) (structures.Pack, error) {

	// TODO!!! create .Description
	persons01Pack := structures.PackAny{}

	var persons01 []types01.Person

	personsItems, err := transformOp.personsOp.List(selector, transformOp.identity)
	if err != nil {
		return nil, fmt.Errorf(onOut+": can't list items (%#v), got %s", selector, err)
	}

	for _, personsItem := range personsItems {
		// TODO!!! set URN (and save using transformOp) if absent
		persons01 = append(persons01, types01.Person{
			Nickname:        personsItem.Nickname,
			Roles:           personsItem.Roles,
			Creds:           personsItem.Creds(),
			ItemDescription: personsItem.ItemDescription,
		})
	}

	persons01Pack.PackData = persons01

	return persons01Pack, nil
}

const onStat = "on transformerOperatorPackPersonsTypes01.Stat(): "

func (transformOp *transformerOperatorPackPersonsTypes01) Stat(selector *selectors.Term, params common.Map) (interface{}, error) {
	return transformOp.personsOp.Stat(selector, transformOp.identity)
}

func (transformOp *transformerOperatorPackPersonsTypes01) Copy(selector *selectors.Term, params common.Map) (data interface{}, err error) {
	return nil, common.ErrNotImplemented
}
