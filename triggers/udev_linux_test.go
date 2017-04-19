package triggers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getValueForProperty(t *testing.T) {
	udevLine1 := `SUBSYSTEMS=="usb", ATTRS{idVendor}=="0400", ATTRS{idProduct}=="6000", MODE="0666", OWNER="james" # my device`

	assert.Equal(t, "0400", getValueForProperty(udevLine1, idVendorKey))
	assert.Equal(t, "6000", getValueForProperty(udevLine1, idProductKey))
	assert.Equal(t, "", getValueForProperty(udevLine1, runKey))

	udevLine2 := `SUBSYSTEMS=="usb", ATTRS{idVendor}=="0400", ATTRS{idProduct}=="6000", MODE="0666", RUN+="/opt/myscript" OWNER="james" # my device`

	assert.Equal(t, "0400", getValueForProperty(udevLine2, idVendorKey))
	assert.Equal(t, "6000", getValueForProperty(udevLine2, idProductKey))
	assert.Equal(t, "/opt/myscript", getValueForProperty(udevLine2, runKey))
}

const sampleFile = `
# sample udev file
# checking commands and comments to check the correct entries are retrieved

SUBSYSTEMS=="usb", ATTRS{idVendor}=="0400", ATTRS{idProduct}=="6000", MODE="0666", OWNER="james" # my device

#comment mixed in file
SUBSYSTEMS=="usb", ATTRS{idVendor}=="0400", ATTRS{idProduct}=="6000", MODE="0666", RUN+="/opt/myscript" OWNER="james" # my device

`

func Test_rulesFromFile(t *testing.T) {
	rules := rulesFromFile(bytes.NewBuffer([]byte(sampleFile)), "/etc/udev/rules.d/51-taskrunner-example")

	assert.Len(t, rules, 1)

	assert.Equal(t, "0400", rules[0].IdVendor)
	assert.Equal(t, "6000", rules[0].IdProduct)
	assert.Equal(t, "/opt/myscript", rules[0].Run)
}
