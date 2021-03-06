package mongo

import (
	"errors"
	"fmt"
	"github.com/keshav-aggarwal/ticks24/config"
	"gopkg.in/mgo.v2"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type AppDB struct {
	Session *mgo.Session
	Config  *config.AppDatabase
}

var appdbI *AppDB

func GetAppDB() (*AppDB, error) {
	if appdbI == nil {
		return nil, errors.New("ERROR : No existing Connection to database")
	}
	return appdbI, nil
}

func (this *AppDB) Connect() (*mgo.Session, error) {
	if this.Session == nil {
		s, er := this.Setup()
		if er != nil {
			return nil, errors.New("ERROR : Failed to Connect to Database : (\n\t" + "\n)")
		}
		this.Session = s.Session
	}
	return this.Session.Clone(), nil
}

func (this AppDB) Setup() (AppDB, error) {
	if this.Session == nil {
		var err error
		url := "mongodb://" + this.Config.Ip + ":" + strconv.Itoa(int(this.Config.Port)) + "/" + this.Config.DatabaseName
		this.Session, err = mgo.Dial(url)
		if err != nil {
			er := this.startMongod()
			if er != nil {
				//tracelog.Errorf(err, "auth", "Connect", fmt.Sprint("Could not connect to AuthDatabase...   :(   Please test if '", url, "' is running."))
				return this, errors.New("ERROR : connection to database failed using " + url + " (\n\t" + err.Error() + ")")

			}
			return this.Setup()
		}
		appdbI = &this
	}
	return this, nil
}

func (this AppDB) startMongod() error {

	if this.Config.Ip == "127.0.0.1" || this.Config.Ip == "localhost" {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("mongod", "--dbpath", "C:\\data\\db")
		} else {
			return errors.New("ERROR : Could not start the mongod service")
		}
		fmt.Println("Starting Mongod server ...")
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
			return err
		}
		time.Sleep(time.Second * 10)
		fmt.Println("... Mongod server is Running now")
		return nil
	}
	return errors.New("It is a remote server you will have to start it yourself.")
}
