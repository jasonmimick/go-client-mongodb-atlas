package main

import (
    "context"
    "github.com/Sectorbob/mlab-ns2/gae/ns/digest"
    "flag"
    //"github.com/davecgh/go-spew/spew"
    //"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
    "go.mongodb.org/atlas/mongodbatlas"

    "io/ioutil"
    "fmt"
    "log"
    "os"
    //"encoding/json"
    //"gopkg.in/yaml.v2"
    //osb "sigs.k8s.io/go-open-service-broker-client/v2"
)

const (
	publicKeyApiEnv  = "ATLAS_PUBLIC_KEY"
	privateKeyApiEnv = "ATLAS_PRIVATE_KEY"
	projectIDEnv  = "ATLAS_PROJECT_ID"
    appIDEnv = "ATLAS_APP_ID"
)

var (
    envPublicApiKey  = os.Getenv(publicKeyApiEnv)
	envPrivateApiKey = os.Getenv(privateKeyApiEnv)
	envProjectID  = os.Getenv(projectIDEnv)
	envAppID = os.Getenv(appIDEnv)
)


func main() {

    // Simple cli to manage
    // Realm app values
    // Usage
    // $realmval --groupid <GROUP_ID> --appid <APP_ID> <Key> [Value|-f <PathToValueFile>]
    //
    // If no appid, then pass --create-app <APPNAME|We generate a name>
    // Values needs to be valid JSON string or file

    var groupID string
    var appID string
    var publicApiKey string
    var privateApiKey string
    var key string
    var value string
    var verbose bool
    var deleteFlag bool

    flag.BoolVar(&verbose,"verbose",false,"Enable verbose output")
    flag.BoolVar(&deleteFlag,"delete",false,"Set to delelete the given --key")
    flag.StringVar(&groupID, "groupid", envProjectID, "MongoDB Atlas Project Id, env ATLAS_PROJECT_ID")
    flag.StringVar(&appID, "appid", envAppID, "MongoDB Realm App Id, env ATLAS_APP_ID")
    flag.StringVar(&key, "key", "", "Key for new value, or used as keyid without value or for --delete")
    flag.StringVar(&value, "value", "", "JSON string for your new value")
    flag.StringVar(&publicApiKey, "publicApiKey", envPublicApiKey, "MongoDB Atlas Public Api Key, or ATLAS_PUBLIC_KEY") 
    flag.StringVar(&privateApiKey, "privateApiKey",envPrivateApiKey, "MongoDB Atlas Private Api Key, or ATLAS_PUBLIC_KEY") 

    flag.Parse()

    if !verbose {
        log.SetOutput(ioutil.Discard)
    }


    t := digest.NewTransport(publicApiKey, privateApiKey)
    tc, err := t.Client()
    if err != nil {
        log.Fatalf(err.Error())
    }
    atlasclient := mongodbatlas.NewClient(tc)
    
    atlasclient.SetCurrentRealmAtlasApiKey ( &mongodbatlas.RealmAtlasApiKey{
        Username: publicApiKey,
        Password: privateApiKey,
    })

    log.Printf("atlasclient.GetCurrentRealmAtlasApiKey(): %+v", atlasclient.GetCurrentRealmAtlasApiKey())


    if (len(appID) > 0) && (len(key) > 0) {
        if len(value)==0 {
            log.Fatalf(fmt.Sprintf("Found key %s but no value",key))
        }
        realmValue,err := atlasclient.RealmValueFromString(key, value)
        if err != nil {
            log.Fatalf(err.Error())
        }

        if !deleteFlag {
            log.Printf("attempt delete realmValue %+v",realmValue)
            value, err := atlasclient.RealmValues.Delete(context.Background(),groupID,appID,realmValue.ID)
            if err != nil {
                log.Fatalf(err.Error())
            }
            log.Printf("delete done: %+v",value)
        } else {
            log.Printf("attempt create realmValue %+v",realmValue)
            value, _, err := atlasclient.RealmValues.Create(context.Background(),groupID,appID,realmValue)
            if err != nil {
                log.Fatalf(err.Error())
            }
            log.Printf("create done %+v",value)
        }
    }  

    if len(appID) > 0 {   // list values
        log.Printf("listing values")
        values, _, err := atlasclient.RealmValues.List(context.Background(),groupID,appID,nil)
        if err != nil {
            log.Fatalf(err.Error())
        }

        fmt.Printf("%+v",values)

    } else {    // list apps
        apps, _, err := atlasclient.RealmApps.List(context.Background(),groupID,nil)
        if err != nil {
            log.Fatalf(err.Error())
        }

        fmt.Printf("%+v",apps)
    }

}



