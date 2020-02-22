package concator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"yandex-catalog-handler/pkg/config"
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

type Results struct {
	Results []Result `json:"results"`
}

type Result struct {
	FileName   string `json:"filename"`
	Was        int    `json:"was"`
	Now        int    `json:"now"`
	WasRemove  int    `json:"was_remove"`
	ErrorCause string `json:"error_cause"`
}

type Catalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Text    string   `xml:",chardata"`
	Date    string   `xml:"date,attr"`
	Shop    struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name"`
		Company    string `xml:"company"`
		URL        string `xml:"url"`
		Currencies struct {
			Text     string `xml:",chardata"`
			Currency struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
				Rate string `xml:"rate,attr"`
			} `xml:"currency"`
		} `xml:"currencies"`
		Categories struct {
			Text     string `xml:",chardata"`
			Category []struct {
				Text     string `xml:",chardata"`
				ID       string `xml:"id,attr"`
				ParentId string `xml:"parentId,attr"`
			} `xml:"category"`
		} `xml:"categories"`
		DeliveryOptions struct {
			Text   string `xml:",chardata"`
			Option struct {
				Text        string `xml:",chardata"`
				Cost        string `xml:"cost,attr"`
				Days        string `xml:"days,attr"`
				OrderBefore string `xml:"order-before,attr"`
			} `xml:"option"`
		} `xml:"delivery-options"`
		Offers struct {
			Text  string `xml:",chardata"`
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

//func RemoveElem(slice []Catalog, pos int)  {
//	return append(slice[:s], slice[s+1:]...)
//}

func (c *Concator) Concate() (results Results, err error) {
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

		result := Result{
			FileName: fmt.Sprint("%s%s", c.cfg.SourceUrl, fileName),
		}

		var catalog Catalog
		catalog, err = readFile(fmt.Sprintf("%s/%s", c.cfg.DataPath, fileName))

		if err != nil {
			result.ErrorCause = err.Error()
			return
		}

		tmp := catalog.Shop.Offers.Offer[:0]

		for index, offer := range catalog.Shop.Offers.Offer {
			var price float64

			price, err = strconv.ParseFloat(offer.Price, 64)

			if err != nil {
				err = fmt.Errorf("Bad price value: %v can't convert from string to float")
				result.ErrorCause = err.Error()
				return
			}

			if oldPrice, ok := c.names[offer.Name]; ok {
				if oldPrice > price {
					price = oldPrice
				}

				result.WasRemove += 1
			} else {
				tmp = append(tmp, offer)
			}

			c.names[offer.Name] = price

			result.Was += 1

			if index%10000 == 0 {
				log.Printf("Processed: %d\n", index)
			}
		}

		catalog.Shop.Offers.Offer = tmp

		result.Now = len(catalog.Shop.Offers.Offer)

		log.Printf("End read: %s\n", fileName)

		results.Results = append(results.Results, result)
	}

	return
}
