package handler

//
//import (
//	"context"
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"strings"
//
//	"awesomeProjectCr/internal/config"
//	logger "github.com/sirupsen/logrus"
//)
//
//type key string
//
//const (
//	identityKey key = "identity"
//)
//
//func (h *Handler) Authentication(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		var jwtPayload user.JwtPayload
//
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			asErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
//			return
//		}
//
//		authToken := strings.Split(authHeader, " ")
//		if authToken[0] != "Bearer" && authToken[0] != "bearer" {
//			asErrorResponse(w, http.StatusUnauthorized, "Bearer is required")
//			return
//		}
//
//		token := strings.TrimPrefix(authHeader, "Bearer ")
//
//		parts := strings.Split(token, ".")
//		if len(parts) != 3 {
//			asErrorResponse(w, http.StatusUnauthorized, "Token invalid")
//			return
//		}
//
//		encodedPayload := parts[1]
//
//		decodedPayload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
//		if err != nil {
//			asErrorResponse(w, http.StatusUnauthorized, "Error decode: "+err.Error())
//			return
//		}
//
//		payloadString := string(decodedPayload)
//
//		err = json.Unmarshal([]byte(payloadString), &jwtPayload)
//		if err != nil {
//			asErrorResponse(w, http.StatusUnauthorized, "Error unmarshal: "+err.Error())
//			return
//		}
//
//		//tokenExp := int64(jwtPayload.Iat)
//		//
//		//currentTime := time.Now().Unix()
//		//
//		//if tokenExp < currentTime {
//		//	asErrorResponse(w, http.StatusUnauthorized, "token expired")
//		//	return
//		//}
//
//		//userData, err := h.userService.GetUserByID(r.Context(), jwtPayload.Vid)
//		//if err != nil {
//		//	asErrorResponse(w, http.StatusUnauthorized, "h.userService.GetUserById: "+err.Error())
//		//	return
//		//}
//
//		userData := user.User{
//			ID:          jwtPayload.Vid,
//			DisplayName: "",
//			Nickname:    "",
//			Email:       "",
//			PhoneNumber: "",
//			Avatar:      "",
//		}
//
//		ctx := context.WithValue(r.Context(), identityKey, &userData)
//
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func (h *Handler) CacheMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		path := r.URL.Path
//		key := fmt.Sprintf("%v:%v", config.CachePrefix, r.URL.RawQuery)
//		result := h.cache.Get(r.Context(), key)
//
//		data, _ := result.Bytes()
//		if data == nil {
//			ctx := context.WithValue(r.Context(), identityKey, nil)
//
//			next.ServeHTTP(w, r.WithContext(ctx))
//			return
//		}
//
//		switch path {
//		case "/ads/v1/cust-params":
//			var mapData []map[string]string
//
//			err := json.Unmarshal(data, &mapData)
//			if err != nil {
//				logger.Errorln(err.Error())
//			}
//			asNoFormatJsonResponse(w, http.StatusOK, "success", mapData)
//			return
//
//		case "/ads/v1/vmap":
//			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
//			w.Header().Set("Access-Control-Allow-Credentials", "true")
//
//			asXmlResponse(w, http.StatusOK, "success", data)
//			return
//
//		default:
//			asXmlResponse(w, http.StatusOK, "success", data)
//			return
//		}
//	})
//}
