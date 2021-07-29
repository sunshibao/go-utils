# ZKits Requester Library #

## About ##

This package is a library of ZKits project. 
This library provides an efficient and easy-to-use HTTP client.

## Usage ##

 1. Import package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-requester
    ```

 2. Create a client to send some requests.

    ```go
    package main
    
    import (
       "fmt"
    
       "github.com/edoger/zkits-requester"
    )
    
    func main() {
       client := requester.Default()
       // client.Get("https://test.com", nil)
       res, err := client.New("https://test.com").Get()
       if err != nil {
           panic(err)
       }
       fmt.Println(res.String())
       
       data := map[string]interface{}{"key": "value"}
       // client.PostJSON("https://test.com", data)
       res, err = client.New("https://test.com").
           WithJSONBody(data).
           Post()
       if err != nil {
           panic(err)
       }
       
       var obj interface{}
       // Bind response json to object.
       if err = res.JSON(&obj); err != nil {
           panic(err)
       }
    }
    ```

 3. Upload file.

    ```go
    package main
    
    import (
       "fmt"
    
       "github.com/edoger/zkits-requester"
    )
    
    func main() {
       client := requester.Default()
       // client.UploadFile("https://test.com", "upload", "path/to/file")
       res, err := client.New("https://test.com").
           WithFormDataFile("upload", "path/to/file").
           Upload()
       if err != nil {
           panic(err)
       }
       fmt.Println(res.String())
    }
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
