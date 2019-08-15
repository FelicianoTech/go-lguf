package lguf

import (
	"encoding/binary"
	"fmt"

	"github.com/google/gousb"
)

const (
	MaxBrightness uint16 = 54000
	MinBrightness uint16 = 540

	VendorID  gousb.ID = 0x043e
	ProductID gousb.ID = 0x9a40
)

type Connection struct {
	ctx        *gousb.Context
	dev        *gousb.Device
	cfg        *gousb.Config
	minEnforce bool
}

func NewConnection() (*Connection, error) {

	ctx := gousb.NewContext()
	dev, err := ctx.OpenDeviceWithVIDPID(VendorID, ProductID)
	if err != nil {
		return nil, fmt.Errorf("Could not open device: %v", err)
	}

	dev.SetAutoDetach(true)

	cfgID, err := dev.ActiveConfigNum()
	if err != nil {
		return nil, fmt.Errorf("Error getting config ID: %v", err)
	}

	cfg, err := dev.Config(cfgID)
	if err != nil {
		return nil, fmt.Errorf("Error getting config: %v", err)
	}

	return &Connection{ctx, dev, cfg, true}, nil
}

func (c *Connection) Brightness() (uint16, error) {

	data := make([]byte, 8)

	cResult, err := c.dev.Control(gousb.ControlIn|gousb.ControlClass|gousb.ControlInterface, 0x01, 768, 1, data)
	if err != nil {
		return 0, fmt.Errorf("Error receiving brightness. The result was: %v and the error was: %v", cResult, err)
	}

	brightness := binary.LittleEndian.Uint16(data[0:])
	if err := c.checkBrightness(brightness); err != nil {
		return 0, err
	}

	return brightness, nil
}

func (c *Connection) checkBrightness(brightness uint16) error {
	if brightness > MaxBrightness {
		return fmt.Errorf("Brightness (%v) is over the max value.", brightness)
	} else if brightness < MinBrightness && c.minEnforce {
		return fmt.Errorf("Brightness (%v) is lower than the minimum value and enforcement is on.", brightness)
	}

	return nil
}

func (c *Connection) LowerBrightness(amount uint16) error {

	brightness, err := c.Brightness()
	if err != nil {
		return err
	}

	if brightness-amount < 0 {
		return fmt.Errorf("Error: Lowering the brightness by %v would cause a negative value.", amount)
	}

	brightness -= amount

	return c.SetBrightness(brightness)
}

func (c *Connection) RaiseBrightness(amount uint16) error {

	// uint16 max value is 65535
	// We want to prevent wrapping around the brightness
	if amount > (65535 - MaxBrightness) {
		return fmt.Errorf("Error: Raising the brightness by %v would cause an overflowed value.", amount)
	}

	brightness, err := c.Brightness()
	if err != nil {
		return err
	}

	brightness += amount

	return c.SetBrightness(brightness)
}

func (c *Connection) SetBrightness(brightness uint16) error {

	if err := c.checkBrightness(brightness); err != nil {
		return err
	}

	data := make([]byte, 6)
	binary.LittleEndian.PutUint16(data, brightness)

	cResult, err := c.dev.Control(gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface, 0x09, 768, 1, data)
	if err != nil {
		return fmt.Errorf("Error setting brightness. The result was: %v and the error was: %v", cResult, err)
	}

	return nil
}

func (c *Connection) Close() {

	err := c.cfg.Close()
	if err != nil {
		fmt.Printf("Error: Couldn't close config successfully. %v", err)
	}

	err = c.dev.Close()
	if err != nil {
		fmt.Printf("Error: Couldn't close device successfully. %v", err)
	}

	err = c.ctx.Close()
	if err != nil {
		fmt.Printf("Error: Couldn't close context successfully. %v", err)
	}
}
