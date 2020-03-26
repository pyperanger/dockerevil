package main

import (
	"context"
  "fmt"
  "os"
  "sort"
  "flag"
  "os/user"
  "io/ioutil"
  "errors"
  //"golang.org/x/crypto/ssh"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
)


type Dockertar struct{
  Imagetag string
  Privatekey string
  Publickey string
}

var (
  exit = os.Exit;

  RED = "\x1b[31m;1m==>\x1b[0m"
  GREEN = "\x1b[32;1m==>\x1b[0m"

  host = flag.String("h", "", "Target Host: http://192.168.0.24:4243/")
  version = flag.String("v", "v1.25", "Remote Engine API version")
  cert = flag.String("c", homeKey()+"/.ssh/id_rsa", "Certificate Path")
)

func homeKey() string{
  user, err := user.Current()
  if err != nil {
    panic(err)
  }
  return user.HomeDir
}

func chkVersion(cli *client.Client) {
  version, err := cli.ServerVersion(context.Background());
  if err != nil{
    panic(err)
  }
  fmt.Printf("%s API Version: %s\n", GREEN, version.APIVersion);
  fmt.Printf("%s Docker Version: %s\n", GREEN, version.Version); 
  fmt.Printf("%s OS/Arch: %s/%s\n", GREEN, version.Os, version.Arch);
  fmt.Printf("%s Kernel Version: %s\n", GREEN, version.KernelVersion);
}

func imageChoose(cli *client.Client) (string, error){
  images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
  if err != nil{
    return "", err;
  }

  sort.SliceStable(images, func(i, j int) bool {
    return images[i].Size < images[j].Size
  })

  fmt.Printf("%s Smaller image: %s / Size: %d MB\n", GREEN, images[0].RepoTags, images[0].Size / 1024);

  return images[0].RepoTags[0], nil
}

func valideKey() (string, string, error){
  // check private and public keys
  // chance for os.Stat .. 
  priv, err := ioutil.ReadFile(*cert)
  if err != nil || len(priv) == 0 {
    return "","", erros.New("Unable to read file or is empty")
  }

  pub, err := ioutil.ReadFile(*cert+".pub")
  if err != nil || len(pub) == 0 {
    return "","", erros.New("Unable to read file or is empty")
  }

  return priv, pub, nil 
}

func dockerFile(cli *client.Client) (*Dockertar, error){ 
  docker := new(Dockertar)
  
  image, err  := imageChoose(cli)
  if err != nil {
    return docker, err
  }

  priv, pub, err := valideKey()
  if err != nil{
    panic(err)
  }

  docker.Imagetag = image
  docker.Privatekey = priv
  docker.Publickey = pub

  return docker, nil
}

func main() {
  flag.Parse()

  if *host == "" {
    flag.Usage();
    exit(1);
  } 

	cli, err := client.NewClient(*host,*version,nil,nil)
	if err != nil {
		panic(err)
	}

  fmt.Printf("%s Dockerevil - github.com/pyperanger/dockerevil\n", GREEN);
  chkVersion(cli);

  dockerevil , err := dockerFile(cli)
  if err != nil {
    panic(nil)
  }


  fmt.Printf("%T\n",dockerevil);

}
