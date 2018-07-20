package drink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	client = http.DefaultClient
)

const (
	DrinkAPIURL = "https://drink-17.herokuapp.com"
)

type DrinkShop struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type DrinkShopReponse struct {
	Options []*DrinkShop `json:"options"`
}

type MenuParams struct {
	Name        string `json:"name"`
	Duration    int    `json:"duration"`
	DrinkShopID int    `json:"drink_shop_id"`
	Email       string `json:"email"`
}

func GetDrinkShops() (*DrinkShopReponse, error) {
	drinkShopURL := fmt.Sprintf("%s/%s", DrinkAPIURL, "drink_shops.json")
	req, err := http.NewRequest("GET", drinkShopURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	out := &DrinkShopReponse{}
	json.Unmarshal(data, out)
	fmt.Println(out)
	return out, nil
}

func CreateMenu(params *MenuParams) error {
	if !strings.HasSuffix(params.Email, "17.media") {
		params.Email = fmt.Sprintf("%s@%s", params.Email, "17.media")
	}

	createMenuURL := fmt.Sprintf("%s/%s", DrinkAPIURL, "api/menus")
	data, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", createMenuURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("DRINKER_ACCESS_TOKEN")))

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
