package concator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"yandex-catalog-handler/internal/entity"
	"yandex-catalog-handler/pkg/config"
)

const (
	Header = `<!DOCTYPE yml_catalog SYSTEM "shops.dtd">` + "\n"
)

type Concator struct {
	cfg   config.Config
	names map[string]float64
}

func New(cfg config.Config) *Concator {
	return &Concator{
		cfg:   cfg,
		names: make(map[string]float64),
	}
}

type Catalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Date    string   `xml:"date,attr"`
	Shop    struct {
		Name       string `xml:"name"`
		Company    string `xml:"company"`
		URL        string `xml:"url"`
		Currencies struct {
			Currency struct {
				ID   string `xml:"id,attr"`
				Rate string `xml:"rate,attr"`
			} `xml:"currency"`
		} `xml:"currencies"`
		Categories struct {
			Category []struct {
				ID       string `xml:"id,attr"`
				ParentId string `xml:"parentId,attr"`
			} `xml:"category"`
		} `xml:"categories"`
		DeliveryOptions struct {
			Option struct {
				Text        string `xml:",chardata"`
				Cost        string `xml:"cost,attr"`
				Days        string `xml:"days,attr"`
				OrderBefore string `xml:"order-before,attr"`
			} `xml:"option"`
		} `xml:"delivery-options"`
		Offers struct {
			Offer []struct {
				Text        string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				Available   string `xml:"available,attr"`
				URL         string `xml:"url"`
				Price       string `xml:"price"`
				CurrencyId  string `xml:"currencyId"`
				CategoryId  string `xml:"categoryId"`
				Pickup      string `xml:"pickup"`
				Delivery    string `xml:"delivery"`
				Name        string `xml:"name"`
				Vendor      string `xml:"vendor"`
				VendorCode  string `xml:"vendorCode"`
				Description string `xml:"description"`
				SalesNotes  string `xml:"sales_notes"`
				Barcode     string `xml:"barcode"`
				Param       struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"param"`
				Picture string `xml:"picture"`
			} `xml:"offer"`
		} `xml:"offers"`
	} `xml:"shop"`
}

func readFile(fileName string) (catalog Catalog, err error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	err = xml.Unmarshal(yamlFile, &catalog)
	if err != nil {
		return
	}

	return
}

func (c *Concator) Concate() (resultsByFile []entity.ResultByFile, err error) {
	files, err := ioutil.ReadDir(c.cfg.DataPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		err = fmt.Errorf("No one file in %s", c.cfg.DataPath)
		return
	}

	for _, f := range files {
		fileName := f.Name()

		log.Printf("Start read: %s\n", fileName)

		resultByFile := entity.ResultByFile{
			FileName: fmt.Sprint("%s%s", c.cfg.SourceUrl, fileName),
		}

		var catalog Catalog

		var filePath = fmt.Sprintf("%s/%s", c.cfg.DataPath, fileName)

		catalog, err = readFile(filePath)

		if err != nil {
			resultByFile.ErrorCause = err.Error()
			return
		}

		tmp := catalog.Shop.Offers.Offer[:0]

		for index, offer := range catalog.Shop.Offers.Offer {
			var price float64

			price, err = strconv.ParseFloat(offer.Price, 64)

			if err != nil {
				err = fmt.Errorf("Bad price value: %v can't convert from string to float")
				resultByFile.ErrorCause = err.Error()
				return
			}

			if oldPrice, ok := c.names[offer.Name]; ok {
				if oldPrice > price {
					price = oldPrice
				}

				resultByFile.WasRemove += 1
			} else {
				tmp = append(tmp, offer)
			}

			c.names[offer.Name] = price

			resultByFile.Was += 1

			if index%10000 == 0 {
				log.Printf("Processed: %d\n", index)
			}
		}

		catalog.Shop.Offers.Offer = tmp

		resultByFile.Now = len(catalog.Shop.Offers.Offer)

		log.Printf("End read: %s\n", fileName)

		resultsByFile = append(resultsByFile, resultByFile)

		err = os.Remove(filePath)

		if err != nil {
			return
		}

		file, _ := xml.MarshalIndent(catalog, "  ", "    ")

		file = []byte(xml.Header + Header + string(file))

		_ = ioutil.WriteFile(filePath, file, 0644)
	}

	return
}
