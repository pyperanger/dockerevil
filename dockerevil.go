package main

import (
	"context"
  "fmt"
  "os"
  "flag"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
)

type server struct {
  types.ImageSumary{};
  //  version ImageSI 
}

var (
  exit = os.Exit;

  RED = "\x1b[31m==>\x1b[0m"
  GREEN = "\x1b[32m==>\x1b[0m"

  host = flag.String("h", "", "Target Host: http://192.168.0.24:4243/")
  version = flag.String("v", "v1.24", "Remote Engine API version")
  cert = flag.String("c", "~/.ssh/id_rsa", "Certificate Path")
)

func lowerImage(cli *Client) (img []types.ImageSummary){
  


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

  images, err := cli.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	for _, image := range images {
    fmt.Printf("%T %s %s %d\n",image, image.RepoTags, image.RepoDigests, image.Size)
	}
  
}
