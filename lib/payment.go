package lib

import(
	"api/config"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func SpinGenerateSignature(trNumber, name string) string {
	str := config.MerchantID +`||`+
	config.Partner + `||` + 
	`474e50c41d661e651cf0c094d0551b886e3503d25a78a847854f5dc8e7d034a9` + `||` + 
	trNumber + `||` + name
	encryptedByte := sha256.Sum256([]byte(str))
	signature := hex.EncodeToString(encryptedByte[:])
	log.Info("signature :", signature)
	return signature
}

func Spin(trNumber string, name string, params map[string]string) (int, string, error) {
	spin := make(map[string]map[string]string)
	url := make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/orders"
	spin["CREATE_ORDER"] = url
	url = make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/pay/otp"
	spin["CREATE_OTP"] = url
	url = make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/pay"
	spin["PAY_ORDER"] = url
	signature := SpinGenerateSignature(trNumber, name)
	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	spinUrl := spin[name]
	req, err := http.NewRequest(spinUrl["method"], spinUrl["url"], payload)
	if err != nil {
		log.Error("Error1", err.Error())
		return http.StatusBadGateway, "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("auth-merchant", config.MerchantID)
	req.Header.Add("auth-partner", config.Partner)
	req.Header.Add("auth-signature", signature)

	log.Info(formatRequest(req))

	res, err := http.DefaultClient.Do(req)
	log.Info(res.StatusCode)
	// if res.StatusCode != 200 {
	// 	log.Error("Error : ", res.StatusCode)
	// 	return res.StatusCode, "", err
	// }
	if err != nil {
		log.Error("Error2", err.Error())
		return http.StatusBadGateway, "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Error3", err.Error())
		return http.StatusBadGateway, "", err
	}
	log.Info(string(body))
	// var sec map[string]interface{}
	// if err = json.Unmarshal(body, &sec); err != nil {
	// 	log.Error("Error4", err.Error())
	// 	return err.Error()
	// }

	return http.StatusOK, string(body), nil
}

func GenerateReference(prefix string, id string) (string){
	x := 6
	y := len(id)
	z := x-y
	r := prefix + strings.Repeat("0", z) + id + strconv.FormatInt(time.Now().Unix(), 10)
	return r
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
	  name = strings.ToLower(name)
	  for _, h := range headers {
		request = append(request, fmt.Sprintf("%v: %v", name, h))
	  }
	}
	
	// If this is a POST, add post data
	if r.Method == "POST" {
	   r.ParseForm()
	   request = append(request, "\n")
	   request = append(request, r.Form.Encode())
	} 
	 // Return the request as a string
	 return strings.Join(request, "\n")
}