package hermes

import (
	"testing"

	"github.com/sapcc/hermes/pkg/configdb"
	"github.com/stretchr/testify/require"
)

func Test_GetAudit(t *testing.T) {
	entry, err := GetAudit("", configdb.Mock{})
	require.Nil(t, err)
	require.NotNil(t, entry)

}
