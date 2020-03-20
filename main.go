package main


import (
  "encoding/json"
  twitter "github.com/dghubble/go-twitter/twitter"
  "github.com/fatih/structs"
  "github.com/dghubble/oauth1"
  "fmt"
  "io/ioutil"
  "log"
  "encoding/csv"
  "os"
  "strconv"
  "time"

)

func main() {

  client, err := getClient()
  if err != nil {
      log.Println("Error getting Twitter Client")
      log.Println(err)
  }
  params := &twitter.FriendListParams{
		ScreenName:          "TARGET_SCREEN_NAME",
		Count:               200,
		Cursor:              -1,
	}

  var result []twitter.User
  for params.Cursor != 0 {
		friends, _, err := client.Friends.List(params)
    if(err != nil) {
      log.Println(err)
      // wait for rate limit unlocked
      time.Sleep(30 * time.Second)
      continue;
    }
    log.Println(friends.Users)
    result = append(result, friends.Users...)
    params.Cursor = friends.NextCursor
    break;
    time.Sleep(1 * time.Minute)
	}
  log.Println("gathered friends count:" + strconv.Itoa(len(result)))

  csvFile, err := os.Create("result_friends.csv")

  if err != nil {
  	log.Fatalf("failed creating file: %s", err)
  }

  csvwriter := csv.NewWriter(csvFile)
  _ = csvwriter.Write([]string{"Name", "ScreenName", "FollowersCount", "Description"})
  for _, user := range result {
    userMap := structs.Map(user)
  	_ = csvwriter.Write([]string{userMap["Name"].(string), userMap["ScreenName"].(string), strconv.Itoa(userMap["FollowersCount"].(int)), userMap["Description"].(string)})
  }
  csvwriter.Flush()
  csvFile.Close()

}

func getClient() (*twitter.Client, error) {
  raw, error := ioutil.ReadFile("twitterAccount.json")
  if error != nil {
      fmt.Println(error.Error())
      return nil, error
  }

  var twitterAccount TwitterAccount
  json.Unmarshal(raw, &twitterAccount)
  config := oauth1.NewConfig(twitterAccount.ConsumerKey, twitterAccount.ConsumerSecret)
  token := oauth1.NewToken(twitterAccount.AccessToken, twitterAccount.AccessTokenSecret)

  httpClient := config.Client(oauth1.NoContext, token)
  client := twitter.NewClient(httpClient)

  verifyParams := &twitter.AccountVerifyParams{
      SkipStatus:   twitter.Bool(true),
      IncludeEmail: twitter.Bool(true),
  }
  _, _, err := client.Accounts.VerifyCredentials(verifyParams)
  if err != nil {
      return nil, err
  }
  return client, nil
}

type TwitterAccount struct {
    AccessToken       string `json:"accessToken"`
    AccessTokenSecret string `json:"accessTokenSecret"`
    ConsumerKey       string `json:"consumerKey"`
    ConsumerSecret    string `json:"consumerSecret"`
}
