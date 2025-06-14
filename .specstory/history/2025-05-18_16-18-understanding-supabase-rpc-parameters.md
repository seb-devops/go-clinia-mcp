<!-- Generated by SpecStory -->

# Understanding Supabase RPC Parameters (2025-05-18 12:18:07)

_**User**_

the second parameters is count. What does it means

---

_**Assistant**_

In the line:

```go
resStr := c.supabase.Rpc("match_site_pages", "get", map[string]interface{}{
	"query_embedding": embedding,
	"filter":          "source = 'clinia-doc'",
	"match_count":     matchCount,
})
```

You are referring to the second parameter of the `Rpc` method:

```go
c.supabase.Rpc("match_site_pages", "get", ...)
```

Here, `"match_site_pages"` is the name of the Supabase RPC (Remote Procedure Call) function you are calling, and `"get"` is the second parameter.

### What does the `"get"` parameter mean?

- In the context of the [supabase-go](https://github.com/supabase-community/supabase-go) client, the `Rpc` method is typically defined as:
  ```go
  func (c *Client) Rpc(fnName string, method string, params map[string]interface{}) (string, error)
  ```
  - `fnName`: The name of the Postgres function (RPC) to call.
  - `method`: The HTTP method to use (usually `"POST"` or `"GET"`).
  - `params`: The parameters to pass to the function.

- In your code, `"get"` is being used as the HTTP method.  
  - This is a bit unusual, as most Supabase RPC calls use `"POST"` to send parameters in the body.  
  - If you use `"get"`, the parameters may be sent as query parameters, which is only supported if your Supabase function is defined to accept them that way.

### Summary

- The second parameter (`"get"`) specifies the HTTP method used for the RPC call.
- It tells the Supabase client to use a GET request to call the `match_site_pages` function.
- If your function expects parameters in the body, you should use `"post"` instead.

**If you are unsure which method to use, check your Supabase function definition and the supabase-go client documentation.**  
Most of the time, `"post"` is the correct choice for RPCs that take parameters.

Let me know if you want to see the relevant part of the supabase-go client or need help changing this!

---

