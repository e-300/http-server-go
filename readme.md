This is an http server built in go emulating twitter
Im calling it chirpy instead...

Project isnt anywhere near where I want it and the code quality, naming conventions and certain algos used in the code need to be optimized, 

Just documenting the patterns and knowledge im learning for my personal record as this is my first project with go.



## HTTP Server Architecture
    1. TCP Listener
    2. Connection Handling
    3. Parsing the Request
    4. Router
    5. Handlers
    6. Response
    7. Middleware

## Patterns 

### Middleware    
    // take in a the next handler and return a new one
    func middlewareMetricsInc(next http.Handler) http.Handler{
        // http.HandlerFunc is a type that has a ServeHTTP method
        //  So we cast an anonymous function to http.HandlerFunc 
        //  https.HandlerFunc is a type 
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){ 
            // some logic runs before request hits main handler
            next.ServeHTTP(w,r) // acutal work being done by server
            // logic after handler finishes
        })
    }

### Thread saftey
    sync/atomic
    Atomic.Int32 for concurrent safe operations 

### Type Conversion vs. Interfaces
    http.Handler -> interface contract 
    http.HandlerFunc -> type conversion allows a plain function to satisfy contract by giving serveHTTP method      
                        automatically

### Router & fileserver
    Router maps endpoints to handlers 
    Router also maps the /app/ endpoint to the fileserver handler 
        fileserver is special becuase its another mapping within the routers map
        except fileservers mapping is the rest of the request url mapped to the server files or hd
        Basically a nested lookup, which is why we need to strip the prefix so  
        example. 
        | Request URL       | Mux Action           | FileServer Action  | Final Local Path      |
        |-------------------|----------------------|--------------------|-----------------------|
        | `/api/users`      | Send to `UserHandler`| N/A                | N/A                   |
        | `/app/index.html` | Send to `FileServer` | Match `index.html` | `./static/index.html` |
        | `/app/js/main.js` | Send to `FileServer` | Match `js/main.js` | `./static/js/main.js` |


