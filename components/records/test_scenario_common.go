package records

import (
	"testing"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
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
		Tags:     []string{"1", "332343"},
	},
}

var item11 = Item{
	Content: Content{
		Title:    "345456",
		Summary:  "6578gj",
		TypeKey:  "test",
		Embedded: embedded,
		Data:     `{"AAA": "aaa", "BBB": 222}`,
		Tags:     []string{"1", "333"},
	},
}

var item12 = Item{
	Content: Content{
		Title:   "345eeeee456rt",
		Summary: "6578eegj",
		TypeKey: "test1",
		Data:    `{"AAA": "awraa", "BBB": 22552}`,
		Tags:    []string{"1", "333"},
	},
}

var item22 = Item{
	Content: Content{
		Title:    "34545ee6rt",
		Summary:  "6578weqreegj",
		TypeKey:  "test2",
		Data:     `wqerwer`,
		Embedded: append(embedded, embedded...),
		Tags:     []string{"qw1", "333"},
	},
}

func readOkTest(t *testing.T, recordsOp Operator, item Item, options crud.Options) {
	itemReaded, err := recordsOp.Read(item.ID, &options)
	require.NoError(t, err)
	require.NotNil(t, itemReaded)

	require.Equal(t, item.ID, itemReaded.ID)
	require.Equal(t, item.Content, itemReaded.Content)
	require.Equal(t, item.OwnerID, itemReaded.OwnerID)
	require.Equal(t, item.ViewerID, itemReaded.ViewerID)

	items, err := recordsOp.List(&options)
	require.NoError(t, err)

	found := false
	for _, itemListed := range items {
		if itemListed.ID == item.ID {
			found = true
			require.Equal(t, item.ID, itemListed.ID)
			require.Equal(t, item.Content, itemListed.Content)
			require.Equal(t, item.OwnerID, itemListed.OwnerID)
			require.Equal(t, item.ViewerID, itemListed.ViewerID)
		}
	}
	require.Truef(t, found, "%#v", items)

}

func readFailTest(t *testing.T, recordsOp Operator, itemID ID, options crud.Options) {
	itemReaded, err := recordsOp.Read(itemID, &options)
	require.Error(t, err)
	require.Nil(t, itemReaded)

	items, err := recordsOp.List(&options)
	require.NoError(t, err)

	for _, itemListed := range items {
		if itemListed.ID == itemID {
			require.FailNow(t, "the item shouldn't be in list ", "%s -> %#v", itemID, itemListed)
		}
	}
}
