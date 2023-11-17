package usb

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/google/gousb"
	"github.com/spf13/cobra"
	asciitree "github.com/thediveo/go-asciitree"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// option bindings
var bIsOutputTree bool

type UsbDevice struct {
	Bus          int         `json:"bus"`
	Port         int         `json:"port"`
	Address      int         `json:"address"`
	Class        gousb.Class `json:"class"`
	VendorID     gousb.ID    `json:"idVendor"`
	Manufacturer string      `json:"manufacturer"`
	SerialNumber string      `json:"serialNumber"`
	Product      string      `json:"product"`
	ProductID    gousb.ID    `json:"idProduct"`
}

func gousbDevice2MyUsbDeviceStruct(device *gousb.Device) UsbDevice {
	strSerialNumber, errSerialNumber := device.SerialNumber()
	if errSerialNumber != nil {
		strSerialNumber = "N/A"
	}
	strMfr, errMfr := device.Manufacturer()
	if errMfr != nil {
		strMfr = "N/A"
	}
	strProduct, errProduct := device.Product()
	if errProduct != nil {
		strProduct = "N/A"
	}
	return UsbDevice{
		Bus:          device.Desc.Bus,
		Port:         device.Desc.Port,
		Address:      device.Desc.Address,
		Class:        device.Desc.Class,
		VendorID:     device.Desc.Vendor,
		Manufacturer: trimStringDescriptorValue(strMfr),
		SerialNumber: trimStringDescriptorValue(strSerialNumber),
		Product:      trimStringDescriptorValue(strProduct),
		ProductID:    device.Desc.Product,
	}
}

func (device UsbDevice) ToString() string {
	return fmt.Sprintf("Bus=%d Port=%d Dev#=%d\nCls=%d\nManufacturer=%s (%s)\nProduct=%s (%s)\nSerialNumber=%s\n", device.Bus, device.Port, device.Address, device.Class, device.Manufacturer, device.VendorID.String(), device.Product, device.ProductID.String(), device.SerialNumber)
}

type asciiTreeNode struct {
	ID       string           `json:"-"`
	Label    string           `asciitree:"label" json:"-"`
	Props    []string         `asciitree:"properties" json:"-"`
	Children []*asciiTreeNode `asciitree:"children" json:"children"`
	UsbDevice
}

func (node asciiTreeNode) Equals(other asciiTreeNode) bool {
	return (node.ID == other.ID)
}

func gousbDevice2AsciiTreeNode(device *gousb.Device) asciiTreeNode {
	structDev := gousbDevice2MyUsbDeviceStruct(device)
	return asciiTreeNode{
		ID:    usbDeviceID(device),
		Label: fmt.Sprintf("%s (%s)", structDev.Product, structDev.ProductID.String()),
		Props: []string{
			fmt.Sprintf("Bus=%d Port=%d Addr=%d", structDev.Bus, structDev.Port, structDev.Address),
			fmt.Sprintf("Cls=%d", structDev.Class),
			fmt.Sprintf("Manufacturer=%s (%s)", structDev.Manufacturer, structDev.VendorID.String()),
			fmt.Sprintf("SerialNumber=%s", structDev.SerialNumber),
		},
		UsbDevice: structDev,
	}
}

var usbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all connected iOS devices",
	Long:  "List all connected iOS devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		util.InitLogger()
		usbdevices := []UsbDevice{}
		d := findUsbDevice(isAppleDevice, func(usbdevice *gousb.Device) bool { return (udid == "" || isMySerial(udid)(usbdevice)) })
		rootHubs := map[string]*asciiTreeNode{}
		for _, device := range d {
			defer device.Close()
			d := gousbDevice2MyUsbDeviceStruct(device)
			if bIsOutputTree {
				var currentNode *asciiTreeNode
				for _, parent := range getUsbDeviceLineage(device, -1) {
					nodeID := usbDeviceID(parent)
					if isRootHub(parent) {
						if v, bFound := rootHubs[nodeID]; bFound {
							currentNode = v
						} else {
							newNode := gousbDevice2AsciiTreeNode(parent)
							currentNode = &newNode
							rootHubs[nodeID] = currentNode
						}
					} else {
						newNode := gousbDevice2AsciiTreeNode(parent)
						if p := slices.IndexFunc(currentNode.Children, func(child *asciiTreeNode) bool { return child.Equals(newNode) }); p >= 0 {
							currentNode = currentNode.Children[p]
						} else {
							currentNode.Children = append(currentNode.Children, &newNode)
							currentNode = &newNode
						}
					}
				}
				newNode := gousbDevice2AsciiTreeNode(device)
				currentNode.Children = append(currentNode.Children, &newNode)
			}
			if bIsOutputJson {
				usbdevices = append(usbdevices, d)
			} else if !bIsOutputTree {
				fmt.Println(d.ToString())
			}
		}
		if bIsOutputJson {
			var b []byte
			if bIsOutputTree {
				b, _ = json.Marshal(rootHubs)
			} else {
				b, _ = json.Marshal(usbdevices)
			}
			fmt.Println(string(b))
		} else {
			if bIsOutputTree {
				fmt.Println(asciitree.RenderFancy(maps.Values(rootHubs)))
			}
		}
		return nil
	},
}

func initUsbListCmd() {
	usbRootCMD.AddCommand(usbListCmd)
	usbListCmd.Flags().StringVarP(&udid, "udid", "u", "", "target specific device by UDID")
	usbListCmd.Flags().BoolVarP(&bIsOutputJson, "json", "j", false, "output in JSON format")
	usbListCmd.Flags().BoolVarP(&bIsOutputTree, "tree", "t", false, "output in tree format")
}

/*
References:
https://pkg.go.dev/github.com/thediveo/go-asciitree
*/
