package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const issuer = "anon-solicitor"

// GetJWT will return a valid JWT containing the provided information
func GetJWT(secret string, payload map[string]interface{}) string {
	header := "{\"alg\": \"HS256\", \"typ\": \"JWT\"}"

	payload["iss"] = issuer
	payload["exp"] = time.Now().UTC().Add(30 * time.Minute).Unix()
	pl, _ := json.Marshal(payload)

	h := hmac.New(sha256.New, []byte(secret))
	s1 := base64.URLEncoding.EncodeToString([]byte(header))
	s2 := base64.URLEncoding.EncodeToString([]byte(string(pl)))
	h.Write([]byte(s1 + "." + s2))
	sha := hex.EncodeToString(h.Sum(nil))
	s3 := base64.URLEncoding.EncodeToString([]byte(sha))

	return fmt.Sprintf("%v.%v.%v", s1, s2, s3)
}

// CheckToken will take a token an attempt to validate it using the given secret
func CheckToken(token string, secret string) (map[string]interface{}, error) {
	sections := strings.Split(token, ".")
	if len(sections) != 3 {
		// not correct format
		return map[string]interface{}{}, errors.New(fmt.Sprint("Invalid Format: ", token))
	}

	checkString := fmt.Sprintf("%v.%v", sections[0], sections[1])
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(checkString))
	sha := hex.EncodeToString(h.Sum(nil))
	s3 := base64.URLEncoding.EncodeToString([]byte(sha))

	if s3 != sections[2] {
		// signature doesn't match - do not authorize
		return map[string]interface{}{}, errors.New("invalid signature")
	}

	pMap := map[string]interface{}{}

	payload, err := base64.URLEncoding.DecodeString(sections[1])
	if err != nil {
		return map[string]interface{}{}, errors.New("couldn't decode payload")
	}

	err = json.Unmarshal([]byte(payload), &pMap)
	if err != nil {
		return map[string]interface{}{}, errors.New("couldn't parse payload after decoding")
	}

	for k, v := range pMap {
		switch k {
		case "iss":
			val, _ := v.(string)
			if val != issuer {
				return map[string]interface{}{}, errors.New("invalid issuer")
			}
		case "exp":
			val, ok := v.(float64)
			if !ok {
				return map[string]interface{}{}, errors.New(fmt.Sprint("couldnt parse expiration value: ", val))
			}
			exp := time.Unix(int64(val), 0)
			if time.Now().UTC().Unix() >= exp.Unix() {
				return map[string]interface{}{}, errors.New(fmt.Sprint("expired session, ", val))
			}
		default:
			log.Println("skipping validation for key : ", k, " -- value: ", v)
		}
	}

	return pMap, nil
}
