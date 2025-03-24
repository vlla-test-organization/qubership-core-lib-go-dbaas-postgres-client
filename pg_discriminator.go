package pgdbaas

import "strings"

type pgDiscriminator struct {
	RoReplica bool
	Role      string
}

func (d *pgDiscriminator) GetValue() string {
	keyString := ""
	if d.Role != "" {
		keyString += d.Role + ":"
	}
	if d.RoReplica {
		keyString += "roReplica=true:"
	} else {
		keyString += "roReplica=false:"
	}
	if keyString != "" {
		keyString = strings.TrimSuffix(keyString, ":")
	}
	return keyString
}
