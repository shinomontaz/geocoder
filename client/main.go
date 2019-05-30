package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Url string
}

/*
{
    "result": {
        "lat": 37.564114,
        "long": 55.569479
    },
    "error": null,
    "id": 1
}
*/
/*
type YandexJson struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Precision string `json:"precision"`
						}
					}
					Point struct {
						Pos string `json:"pos"`
					}
				}
			} `json:"featureMember"`
		}
	} `json:"response"`
}

*/
type JsonRpcResponse struct {
	Result struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	}
	Error string `json:"error"`
	ID    int64  `json:"id"`
}

func (c *Client) Get(address string) (x, y float64, err error) {
	var jsonStr = fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"RpcServer.Geocode",
		"params":[{
			"Address": "%s"
		}],
		"id": 1
		}`, address)

	req, err := http.NewRequest("POST", c.Url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return x, y, err
	}

	var jsonResponse JsonRpcResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return x, y, err
	}

	return x, y, nil
}
