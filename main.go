package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/wealdtech/go-ens/v3"
)

type List struct {
	Domains []string
}

type Domain struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Output struct {
	Available    []*Domain //available domains
	NotAvailable []*Domain // already registered domains
	NoResolver   []*Domain // no resolver found
	NoAddress    []*Domain // no address found
}

var (
	unregisteredDomainErrStr = "unregistered name"
	noResolverErrStr         = "no resolver"
	noAddressErrStr          = "no address"
)

func init() {
	viper.SetConfigFile("./config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to read config file")
	}
}

func main() {
	if viper.GetString("client-endpoint") == "" {
		panic("client-endpoint is required, please register in infura.io, create a project & copy the http endpoint to the config file")
	}

	lf, err := ioutil.ReadFile(viper.GetString("list-file"))
	if err != nil {
		panic(err)
	}

	list := &List{}
	err = json.Unmarshal(lf, list)
	if err != nil {
		panic(err)
	}

	client, err := ethclient.Dial(viper.GetString("client-endpoint"))
	if err != nil {
		panic(err)
	}

	o := &Output{}
	var wg sync.WaitGroup
	for i, d := range list.Domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			log.Printf("resolving %s.eth", d)
			address, err := ens.Resolve(client, fmt.Sprintf("%s.eth", d))
			if err != nil {
				switch err.Error() {
				case unregisteredDomainErrStr:
					o.Available = append(o.Available, &Domain{
						Name: fmt.Sprintf("%s.eth", d),
					})
				case noResolverErrStr:
					o.NoResolver = append(o.Available, &Domain{
						Name: fmt.Sprintf("%s.eth", d),
					})
				case noAddressErrStr:
					o.NoAddress = append(o.NoAddress, &Domain{
						Name: fmt.Sprintf("%s.eth", d),
					})
				default:
					log.Printf("error while resolving address %s.eth: %v", d, err)
				}
				return
			}

			o.NotAvailable = append(o.NotAvailable, &Domain{
				Name:    fmt.Sprintf("%s.eth", d),
				Address: address.String(),
			})
		}(d)

		if (i+1)%100 == 0 { // sleep to prevent API abuse
			time.Sleep(5 * time.Second)
		}
	}

	wg.Wait()
	log.Printf("all domain checked, storing to output file")
	err = storeOutput(o, viper.GetString("output-file"))
	if err != nil {
		log.Printf("error while storing the output file: %s", err)
		dumpData(o)
		return
	}
	log.Printf("%s saved", viper.GetString("output-file"))
}

func storeOutput(output *Output, file string) error {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func dumpData(o *Output) {
	fmt.Printf("=======\nDUMPING DATA:\n%#v\n", o)
}