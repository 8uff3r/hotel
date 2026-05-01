package sana

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	BaseURL = "http://mosaferin.hotl.ir:11010/ws_amaken.asmx"
)

type Config struct {
	KelidVahed     string
	KelidPeimankar string
	CodeVahed      int
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	config     Config
}

func NewClient(config Config) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    BaseURL,
		config:     config,
	}
}

func (c *Client) buildURL(function string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	return fmt.Sprintf("%s/%s?%s", c.baseURL, function, values.Encode())
}

func (c *Client) get(function string, params map[string]string) ([]byte, error) {
	url := c.buildURL(function, params)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) baseParams() map[string]string {
	return map[string]string{
		"KelidVahed":     c.config.KelidVahed,
		"KelidPeimankar": c.config.KelidPeimankar,
		"CodeVahed":      fmt.Sprintf("%d", c.config.CodeVahed),
	}
}

func parseAnavinResponse(body []byte) ([]AnavinItem, error) {
	var result AnavinResponse
	err := xml.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GereftanAnavinElaatSafar() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinElaatSafar", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanAnavinNesbat() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinNesbat", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanAnavinMeliat() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinMeliat", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanAnavinShahrhayeIran() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinShahrhayeIran", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanAnavinShahrhayeKhareji() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinShahrhayeKhareji", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanAnavinShoghl() ([]AnavinItem, error) {
	body, err := c.get("GereftanAnavinShoghl", c.baseParams())
	if err != nil {
		return nil, err
	}
	return parseAnavinResponse(body)
}

func (c *Client) GereftanPayam(ShomarePayam int) ([]PayamItem, error) {
	params := c.baseParams()
	params["ShomarePayam"] = fmt.Sprintf("%d", ShomarePayam)
	body, err := c.get("GereftanPayam", params)
	if err != nil {
		return nil, err
	}
	var result PayamResponse
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) SabtChidemanVahed(Rac string) (SabtChidemanResponse, error) {
	params := c.baseParams()
	params["Rac"] = Rac
	body, err := c.get("SabtChidemanVahed", params)
	if err != nil {
		return SabtChidemanResponse{}, err
	}
	response := string(body)
	return SabtChidemanResponse{
		Message: response,
		IsOK:    strings.Contains(response, "OK"),
	}, nil
}

type SabtMosaferinInput struct {
	NameMosafer      string
	FamilMosafer     string
	NamePedar        string
	ShomareShenasaee string
	TarikhTavalod    string
	ID_Jensiat       int
	ID_Shoghl        int
	ID_ElaatSafar    int
	TedadHamrah      int
	ID_Nesbat        int
	MosafereKhareji  bool
	ID_Meliat        int
	ID_Mabda         int
	ID_Maghsad       int
	ID_MahaleTavalod int
	TarikhVorod      string
	TarikhKhoroj     string
	SaatVorod        string
	SaatKhoroj       string
	ShomareOtagh     string
	RecordMosafer    int
	ShomareFaragir   string
	ShomarePaziresh  string
	Code_Moaref      int
	Name_Moaref      string
	Tel_Moaref       string
	NameKarbareSabt  string
	ID_NoeDadeh      MosafereDataType
}

func (c *Client) Sabt_Mosaferin(req SabtMosaferinInput) (SabtMosaferinResponse, error) {
	params := c.baseParams()
	params["NameMosafer"] = req.NameMosafer
	params["FamilMosafer"] = req.FamilMosafer
	params["NamePedar"] = req.NamePedar
	params["ShomareShenasaee"] = req.ShomareShenasaee
	params["TarikhTavalod"] = req.TarikhTavalod
	params["ID_Jensiat"] = fmt.Sprintf("%d", req.ID_Jensiat)
	params["ID_Shoghl"] = fmt.Sprintf("%d", req.ID_Shoghl)
	params["ID_ElaatSafar"] = fmt.Sprintf("%d", req.ID_ElaatSafar)
	params["TedadHamrah"] = fmt.Sprintf("%d", req.TedadHamrah)
	params["ID_Nesbat"] = fmt.Sprintf("%d", req.ID_Nesbat)
	if req.MosafereKhareji {
		params["MosafereKhareji"] = "1"
	} else {
		params["MosafereKhareji"] = "0"
	}
	params["ID_Meliat"] = fmt.Sprintf("%d", req.ID_Meliat)
	params["ID_Mabda"] = fmt.Sprintf("%d", req.ID_Mabda)
	params["ID_Maghsad"] = fmt.Sprintf("%d", req.ID_Maghsad)
	params["ID_MahaleTavalod"] = fmt.Sprintf("%d", req.ID_MahaleTavalod)
	params["TarikhVorod"] = req.TarikhVorod
	params["TarikhKhoroj"] = req.TarikhKhoroj
	params["SaatVorod"] = req.SaatVorod
	params["SaatKhoroj"] = req.SaatKhoroj
	params["ShomareOtagh"] = req.ShomareOtagh
	params["RecordMosafer"] = fmt.Sprintf("%d", req.RecordMosafer)
	params["ShomareFaragir"] = req.ShomareFaragir
	params["ShomarePaziresh"] = req.ShomarePaziresh
	params["Code_Moaref"] = fmt.Sprintf("%d", req.Code_Moaref)
	params["Name_Moaref"] = req.Name_Moaref
	params["Tel_Moaref"] = req.Tel_Moaref
	params["NameKarbareSabt"] = req.NameKarbareSabt
	params["ID_NoeDadeh"] = fmt.Sprintf("%d", req.ID_NoeDadeh)

	body, err := c.get("Sabt_Mosaferin", params)
	if err != nil {
		return SabtMosaferinResponse{}, err
	}
	response := string(body)

	shomareOtagh := ""
	shomarePaziresh := ""
	recordMosafer := 0
	codeError := 0

	if strings.Contains(response, "OK") {
		parts := strings.Split(response, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "ShomareOtagh=") {
				shomareOtagh = strings.TrimPrefix(part, "ShomareOtagh=")
			}
			if strings.HasPrefix(part, "ShomarePaziresh=") {
				shomarePaziresh = strings.TrimPrefix(part, "ShomarePaziresh=")
			}
			if strings.HasPrefix(part, "RecordMosafer=") {
				fmt.Sscanf(strings.TrimPrefix(part, "RecordMosafer="), "%d", &recordMosafer)
			}
			if strings.HasPrefix(part, "CodeError=") {
				fmt.Sscanf(strings.TrimPrefix(part, "CodeError="), "%d", &codeError)
			}
		}
	} else {
		parts := strings.Split(response, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "CodeError=") {
				fmt.Sscanf(strings.TrimPrefix(part, "CodeError="), "%d", &codeError)
			}
		}
	}

	return SabtMosaferinResponse{
		ShomareOtagh:    shomareOtagh,
		ShomarePaziresh: shomarePaziresh,
		RecordMosafer:   recordMosafer,
		CodeError:       codeError,
		IsOK:            strings.Contains(response, "OK"),
	}, nil
}

func (c *Client) GereftanAnavinErrorCode() ([]ErrorCodeItem, error) {
	params := map[string]string{
		"KelidPeimankar": c.config.KelidPeimankar,
	}
	body, err := c.get("GereftanAnavinErrorCode", params)
	if err != nil {
		return nil, err
	}
	var result ErrorCodeResponse
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
