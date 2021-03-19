package records

import (
	"testing"

	"github.com/pavlo67/data_exchange/components/structures"

	"github.com/pavlo67/data_exchange/components/tags"

	"github.com/pavlo67/common/common/auth"
	"github.com/stretchr/testify/require"
)

const authID1 = auth.ID("1")
const authID2 = auth.ID("2")
const authID3 = auth.ID("3")

const testNewTag = "testNewTag"

var embedded = []Content{
	{
		Title:    "56567",
		Summary:  "3333333",
		TypeKey:  "test...",
		Data:     "wqerwer",
		Embedded: []Content{{Data: "werwe"}},
	},
}

var item11 = Item{
	Content: Content{
		Title:    "345456",
		Summary:  "6578gj",
		TypeKey:  "test",
		Embedded: embedded,
		Data:     `{"AAA": "aaa", "BBB": 222}`,
	},
	ItemDescription: structures.ItemDescription{
		Tags: []tags.Item{"1", "333"},
	},
}

var item12 = Item{
	Content: Content{
		Title:   "345eeeee456rt",
		Summary: "6578eegj",
		TypeKey: "test1",
		Data:    `{"AAA": "awraa", "BBB": 22552}`,
	},
	ItemDescription: structures.ItemDescription{
		Tags: []tags.Item{"1", "333"},
	},
}

var item22 = Item{
	Content: Content{
		Title:    "34545ee6rt",
		Summary:  "6578weqreegj",
		TypeKey:  "test2",
		Data:     `wqerwer`,
		Embedded: append(embedded, embedded...),
	},
	ItemDescription: structures.ItemDescription{
		Tags: []tags.Item{"qw1", "333"},
	},
}

func readOkTest(t *testing.T, recordsOp Operator, item Item, identity auth.Identity) {
	itemReaded, err := recordsOp.Read(item.ID, &identity)
	require.NoError(t, err)
	require.NotNil(t, itemReaded)

	require.Equal(t, item.ID, itemReaded.ID)
	require.Equal(t, item.Content, itemReaded.Content)
	require.Equal(t, item.OwnerNSS, itemReaded.OwnerNSS)
	require.Equal(t, item.ViewerNSS, itemReaded.ViewerNSS)

	items, err := recordsOp.List(nil, &identity)
	require.NoError(t, err)

	found := false
	for _, itemListed := range items {
		if itemListed.ID == item.ID {
			found = true
			require.Equal(t, item.ID, itemListed.ID)
			require.Equal(t, item.Content, itemListed.Content)
			require.Equal(t, item.OwnerNSS, itemListed.OwnerNSS)
			require.Equal(t, item.ViewerNSS, itemListed.ViewerNSS)
		}
	}
	require.Truef(t, found, "%#v", items)

}

func readFailTest(t *testing.T, recordsOp Operator, itemID ID, identity auth.Identity) {
	itemReaded, err := recordsOp.Read(itemID, &identity)
	require.Error(t, err)
	require.Nil(t, itemReaded)

	items, err := recordsOp.List(nil, &identity)
	require.NoError(t, err)

	for _, itemListed := range items {
		if itemListed.ID == itemID {
			require.FailNow(t, "the item shouldn't be in list ", "%s -> %#v", itemID, itemListed)
		}
	}
}
