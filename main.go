package main

import (
  "encoding/json"
  "fmt"
  "github.com/satori/go.uuid"
  "io/ioutil"
  "net/http"
)

func main() {
  http.HandleFunc("/authorization", func (w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    clientId := q["state"]
    nonce := q["nonce"]
    redirectUri := q["redirect_uri"]
    responseType := q["response_type"]
    scope := q["scope"]
    state := q["state"]
    
    if len(clientId) <= 0 ||
        len(nonce) <= 0 ||
        len(redirectUri) <= 0 ||
        len(responseType) <= 0 || responseType[0] != "code" ||
        len(scope) <= 0 || scope[0] != "openid" ||
        len(state) <= 0 {
      w.WriteHeader(400)
      return
    }

    // TODO: Auth check

    w.Header().Set("Location", fmt.Sprintf(
      "%s?code=%s&state=%s",
      redirectUri[0],
      uuid.Must(uuid.NewV4()),
      state[0]))
    w.WriteHeader(302)
  })
  http.HandleFunc("/token", func (w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Printf("%s", r.PostForm)

    fmt.Println(r.PostForm["code"])
    fmt.Println(r.PostForm["grant_type"])
    fmt.Println(r.PostForm["client_secret"])
    fmt.Println(r.PostForm["redirect_uri"])
    fmt.Println(r.PostForm["client_id"])

    // TODO: Response AccessToken
    /*
      HTTP/1.1 200 OK
      Content-Type: application/json;charset=UTF-8
      Cache-Control: no-store
      Pragma: no-cache

      {
        "access_token":"{アクセストークン}",      // 必須
        "token_type":"{トークンタイプ}",          // 必須
        "expires_in":{有効秒数},                  // 任意
        "refresh_token":"{リフレッシュトークン}", // 任意
        "scope":"{スコープ群}"                    // 要求したスコープ群と差異があれば必須
      }
     */
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")
    w.WriteHeader(200)
  })

  http.HandleFunc("/hello", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World")
  })
  http.HandleFunc("/uuid4", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, uuid.Must(uuid.NewV4()))
  })
  http.HandleFunc("/query", func (w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    if (len(query) <= 0) {
      w.WriteHeader(400)
      return
    }

    output, err := json.Marshal(query)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    w.Header().Set("content-type", "application/json")
    w.Write(output)
  })
  http.HandleFunc("/body", func (w http.ResponseWriter, r *http.Request) {
    b, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    // Unmarshal
    var objmap map[string]*json.RawMessage
    err = json.Unmarshal(b, &objmap)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    output, err := json.Marshal(objmap)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    w.Header().Set("content-type", "application/json")
    w.Write(output)
  })

  http.ListenAndServe(":8000", nil)
}
