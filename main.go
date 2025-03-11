package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type UserAgentInfo struct {
	Product  string `json:"product"`
	Version  string `json:"version"`
	RawValue string `json:"raw_value"`
}

type IPInfo struct {
	IP         string        `json:"ip"`
	IPDecimal  int64         `json:"ip_decimal"`
	Country    string        `json:"country"`
	CountryISO string        `json:"country_iso"`
	CountryEU  bool          `json:"country_eu"`
	RegionName string        `json:"region_name"`
	RegionCode string        `json:"region_code"`
	City       string        `json:"city"`
	Latitude   float64       `json:"latitude"`
	Longitude  float64       `json:"longitude"`
	TimeZone   string        `json:"time_zone"`
	ASN        string        `json:"asn"`
	ASNOrg     string        `json:"asn_org"`
	UserAgent  UserAgentInfo `json:"user_agent"`
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		for _, i := range ips {
			cleanIP := strings.TrimSpace(i)
			parsedIP := net.ParseIP(cleanIP)
			if parsedIP != nil {
				return cleanIP
			}
		}
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		parsedIP := net.ParseIP(ip)
		if parsedIP != nil {
			return ip
		}
	}

	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		return ip
	}

	return "Unknown"
}

func getGeoInfo(ip string) (IPInfo, error) {
	resp, err := http.Get("https://ipwhois.app/json/" + ip)
	if err != nil {
		return IPInfo{}, err
	}
	defer resp.Body.Close()

	var data struct {
		IP         string  `json:"ip"`
		IPDecimal  int64   `json:"ip_decimal"`
		Country    string  `json:"country"`
		CountryISO string  `json:"country_code"`
		CountryEU  bool    `json:"is_in_european_union"`
		RegionName string  `json:"region"`
		RegionCode string  `json:"region_code"`
		City       string  `json:"city"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		TimeZone   string  `json:"timezone"`
		ASN        string  `json:"asn"`
		ASNOrg     string  `json:"asn_org"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return IPInfo{}, err
	}

	return IPInfo{
		IP:         data.IP,
		IPDecimal:  data.IPDecimal,
		Country:    data.Country,
		CountryISO: data.CountryISO,
		CountryEU:  data.CountryEU,
		RegionName: data.RegionName,
		RegionCode: data.RegionCode,
		City:       data.City,
		Latitude:   data.Latitude,
		Longitude:  data.Longitude,
		TimeZone:   data.TimeZone,
		ASN:        data.ASN,
		ASNOrg:     data.ASNOrg,
	}, nil
}

func parseUserAgent(userAgent string) UserAgentInfo {
	parts := strings.Split(userAgent, "/")
	product, version := "Unknown", "Unknown"

	if len(parts) > 1 {
		product = parts[0]
		version = parts[1]
	}

	return UserAgentInfo{
		Product:  product,
		Version:  version,
		RawValue: userAgent,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	userAgent := parseUserAgent(r.UserAgent())

	geoInfo, err := getGeoInfo(ip)
	if err != nil {
		http.Error(w, "Không lấy được thông tin địa lý", http.StatusInternalServerError)
		fmt.Println("Lỗi lấy thông tin IP:", err)
		return
	}

	geoInfo.UserAgent = userAgent

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geoInfo)
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
