# SMSBroadcast-go



The SMSBroadcast Go Package makes it simpler to integrate SMS functionality into your Go applications via the SMSBroadcast RESTful API.

At present, only **sending** of SMS is supported.

## Prerequisites
* Go >= 1.16 (probably lower, but I haven't gone back to test yet)

## Getting started

Using an existing Go app, or creating a new one, make sure you are using go modules.

``` bash
$ mkdir ~/broadcastsms
$ cd ~/broadcastsms
$ go mod init broadcastsms
```

Then add the package address as an import

```
package main

import "github.com/grainmarket/smsbroadcast-go"

func main()  {
  client, err := smsbroadcast.NewClient("", "", &smsbroadcast.ClientOptions{})
  if err != nil {
    panic(err)
  }
}
```

Make sure the module dependency is met

``` bash
$ go mod tidy
```

### Authentication

We recommend that you store your credentials in the SMSBROADCAST_USERNAME and the SMSBROADCAST_PASSWORD environment variables, so as to avoid the possibility of accidentally committing them to source control. If you do this, you can initialise the client with no arguments (as shown in the previous example) and it will automatically fetch them from the environment variables.

Alternatively, you can define the auth credentials when initializing the Client:

```
package main

import "github.com/grainmarket/smsbroadcast-go"

func main()  {
  client, err := smsbroadcast.NewClient("USERNAME", "PASSWORD", &smsbroadcast.ClientOptions{})
  if err != nil {
    panic(err)
  }
}
```
## Examples
### Sending A Message

```
  result, err := client.Send(smsbroadcast.Message{
    To: "the_source_number",
    From: "the_destination_number",
    Message: "Hello, world!",
		Ref: "123ABC",
  })
```

result is a MsgResponse with the following properties:

| Property      | Type   | Description                                                                                                                            |
|---------------|--------|----------------------------------------------------------------------------------------------------------------------------------------|
| **Status**    | int    | The http status code of the request                                                                                                    |
| **Summary**   | string | SMSBroadcast's status of the request (OK/BAD/ERROR)                                                                                    |
| **Recipient** | string | The receiving mobile number. This will be shown in international format regardless of the format it was submitted in.                  |
| **Reference** | string | SMS Reference Number or Error Message. Will display our reference number for the SMS message, or the reason for a failed SMS message.  |
