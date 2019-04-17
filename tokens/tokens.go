package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const issuer = "anon-solicitor"

// GetJWT will return a valid JWT containing the provided information
func GetJWT(secret string, tok string) string {
	header := "{\"alg\": \"HS256\", \"typ\": \"JWT\"}"
	payload := fmt.Sprintf("{\"iss\":\"%v\",\"exp\":\"%v\",\"tok\":\"%v\"}",
		issuer, time.Now().UTC().Add(30*time.Minute).Unix(), tok)

	h := hmac.New(sha256.New, []byte(secret))
	s1 := base64.URLEncoding.EncodeToString([]byte(header))
	s2 := base64.URLEncoding.EncodeToString([]byte(payload))
	h.Write([]byte(s1 + "." + s2))
	sha := hex.EncodeToString(h.Sum(nil))
	s3 := base64.URLEncoding.EncodeToString([]byte(sha))

	return fmt.Sprintf("%v.%v.%v", s1, s2, s3)
}

// CheckToken will take a token an attempt to validate it using the given secret
func CheckToken(token string, secret string) (string, error) {
	sections := strings.Split(token, ".")
	if len(sections) != 3 {
		// not correct format
		return "", errors.New("Invalid Format")
	}

	checkString := fmt.Sprintf("%v.%v", sections[0], sections[1])
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(checkString))
	sha := hex.EncodeToString(h.Sum(nil))
	s3 := base64.URLEncoding.EncodeToString([]byte(sha))

	if s3 != sections[2] {
		// signature doesn't match - do not authorize
		return "", errors.New("invalid signature")
	}

	pMap := map[string]interface{}{}

	payload, err := base64.URLEncoding.DecodeString(sections[1])
	if err != nil {
		return "", errors.New("couldn't decode payload")
	}

	err = json.Unmarshal([]byte(payload), &pMap)
	if err != nil {
		return "", errors.New("couldn't parse payload after decoding")
	}

	for k, v := range pMap {
		switch k {
		case "iss":
			val, _ := v.(string)
			if val != issuer {
				return "", errors.New("invalid issuer")
			}
		case "exp":
			val, _ := v.(string)
			vall, _ := strconv.Atoi(val)
			if time.Now().UTC().Unix() >= int64(vall) {
				return "", errors.New(fmt.Sprint("expired session, ", vall))
			}
		case "tok":
			fmt.Println("tok")
			ret, ok := v.(string)
			if !ok {
				return "", errors.New("coldn't parse tok from token")
			}
			return ret, nil
		case "userid":
			fmt.Println(v)
		default:
			fmt.Println("unknown key : ", k, " -- value: ", v)
		}
	}

	return "", nil
}
