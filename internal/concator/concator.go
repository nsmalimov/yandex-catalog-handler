package concator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"yandex-catalog-handler/pkg/config"
)

const (
	Header = `<!DOCTYPE yml_catalog SYSTEM "shops.dtd">` + "\n"
)

type Result struct {
	All          int            `json:"all"`
	Unic         int            `json:"unic"`
	ResultByFile []ResultByFile `json:"result_by_file"`
}

type ResultByFile struct {
	FileName string `json:"filename"`
	Was      int    `json:"was"`
	Now      int    `json:"now"`
}

type Concator struct {
	cfg    config.Config
	result Result
}

func New(cfg config.Config) *Concator {
	return &Concator{
		cfg:    cfg,
		result: Result{},
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
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
				Rate string `xml:"rate,attr"`
			} `xml:"currency"`
		} `xml:"currencies"`
		Categories struct {
			Category []struct {
				Text     string `xml:",chardata"`
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
			Offer []Offer `xml:"offer"`
		} `xml:"offers"`
	} `xml:"shop"`
}

type Offer struct {
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

func GetPriceAndKey(offer Offer) (price float64, key string, err error) {
	price, err = strconv.ParseFloat(offer.Price, 64)
	if err != nil {
		err = fmt.Errorf("Bad price value: %v can't convert from string to float")
		return
	}

	key = offer.Vendor + offer.VendorCode

	return
}

func (c *Concator) PrepareData(files []os.FileInfo) (
	catalogs map[string]Catalog,
	prices map[string]float64,
	err error) {
	catalogs = make(map[string]Catalog)
	prices = make(map[string]float64)

	for _, f := range files {
		fileName := f.Name()

		if fileName == ".DS_Store" {
			continue
		}

		log.Printf("Start read: %s\n", fileName)

		var catalog Catalog

		var filePath = fmt.Sprintf("%s/%s", c.cfg.DataPath, fileName)

		catalog, err = readFile(filePath)
		if err != nil {
			return
		}

		catalogs[fileName] = catalog

		for index, offer := range catalog.Shop.Offers.Offer {
			var price float64
			var key string

			price, key, err = GetPriceAndKey(offer)

			if err != nil {
				return
			}

			if oldPrice, ok := prices[key]; ok {
				if price > oldPrice {
					prices[key] = price
				}
			} else {
				prices[key] = price
			}

			if index%10000 == 0 {
				log.Printf("Processed: %d\n", index)
			}

			c.result.All += 1
		}
	}

	c.result.Unic = len(prices)

	return
}

func (c *Concator) WriteToFile(catalog Catalog, fileName string) (err error) {
	filePath := fmt.Sprintf("%s/%s", c.cfg.DataPath, fileName)

	//err = os.Remove(filePath)
	//if err != nil {
	//	return
	//}

	var file []byte

	file, err = xml.MarshalIndent(catalog, "  ", "    ")
	if err != nil {
		return
	}

	headerBytes := []byte(xml.Header + Header)

	file = append(headerBytes, file...)

	err = ioutil.WriteFile(filePath+"_111", file, 0644)
	if err != nil {
		return
	}

	return
}

func (c *Concator) Concate() (err error) {
	files, err := ioutil.ReadDir(c.cfg.DataPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		err = fmt.Errorf("No one file in %s", c.cfg.DataPath)
		return
	}

	var catalogs map[string]Catalog
	var prices map[string]float64

	catalogs, prices, err = c.PrepareData(files)

	fmt.Println(c.result.All, c.result.Unic)

	alreadyWritten := make(map[string]interface{})

	if err != nil {
		return
	}

	countAll := 0

	for filename, catalog := range catalogs {
		resultByFile := ResultByFile{}

		tmp := catalog.Shop.Offers.Offer[:0]

		for _, offer := range catalog.Shop.Offers.Offer {
			var key string

			_, key, err = GetPriceAndKey(offer)

			countAll += 1

			if _, ok := alreadyWritten[key]; ok {
				continue
			} else {
				alreadyWritten[key] = nil

				offer.Price = fmt.Sprintf("%.2f", prices[key])

				tmp = append(tmp, offer)
			}
		}

		catalog.Shop.Offers.Offer = tmp

		err = c.WriteToFile(catalog, filename)

		if err != nil {
			return
		}

		resultByFile.FileName = filename
		resultByFile.Was = countAll
		resultByFile.Now = len(tmp)

		countAll = 0

		c.result.ResultByFile = append(c.result.ResultByFile, resultByFile)
	}

	for _, elem := range c.result.ResultByFile {
		fmt.Println(elem.FileName, elem.Was, elem.Now)
	}

	return
}
