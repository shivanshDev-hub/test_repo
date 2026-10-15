# HTTP Status Codes Quick Reference

**A comprehensive catalog of standard HTTP status codes, structured by category**

---

**1xx Informational Responses**

~ Informational responses indicate that the request was received and understood. It is issued on a provisional basis while request processing continues.

```
- `100` -> **Continue**: The server has received the request headers and the client should proceed to send the request body.
- `101` -> **Switching Protocols**: The requester has asked the server to switch protocols and the server has agreed to do so.
- `102` -> **Processing** *(WebDAV)*: The server has received and is processing the request, but no response is available yet.
- `103` -> **Early Hints**: Used to return some response headers before the final HTTP message is delivered.
```

**2xx Success Responses**

~ This class of status codes indicates the action requested by the client was received, understood, accepted, and processed successfully.

```
- `200` -> **OK**: Standard response for successful HTTP requests.
- `201` -> **Created**: The request has been fulfilled, resulting in the creation of a new resource.
- `202` -> **Accepted**: The request has been accepted for processing, but processing has not yet completed.
- `203` -> **Non-Authoritative Information**: The server is a transforming proxy that received a `200 OK` from its origin but is returning a modified version.
- `204` -> **No Content**: The server successfully processed the request and is not returning any content.
- `205` -> **Reset Content**: The server successfully processed the request, but requires the requester to reset the document view.
- `206` -> **Partial Content**: The server is delivering only part of the resource due to a range header sent by the client.
- `207` -> **Multi-Status** *(WebDAV)*: Conveys information about multiple resources in situations where multiple status codes might be appropriate.
- `208` -> **Already Reported** *(WebDAV)*: The members of a DAV binding have already been enumerated in a preceding part of the (multistatus) response.
- `226` -> **IM Used**: The server has fulfilled a GET request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance.
```

**3xx Redirection Messages**

~ This class of status code indicates the client must take additional action to complete the request.
```
- `300` -> **Multiple Choices**: Indicates multiple options for the resource from which the client may choose.
- `301` -> **Moved Permanently**: This request and all future requests should be directed to the given URI.
- `302` -> **Found**: Tells the client to look at (browse to) another URL temporarily.
- `303` -> **See Other**: The response to the request can be found under another URI using the GET method.
- `304` -> **Not Modified**: Indicates that the resource has not been modified since the version specified by the request headers `If-Modified-Since` or `If-None-Match`.
- `305` -> **Use Proxy**: The requested resource is available only through a proxy, the address for which is provided in the response.
- `306` -> **Switch Proxy**: No longer used; originally meant "Subsequent requests should use the specified proxy."
- `307` -> **Temporary Redirect**: The request should be repeated with another URI, but future requests should still use the original URI.
- `308` -> **Permanent Redirect**: The request and all future requests should be repeated using another URI.
```

**4xx Client Error Responses**

~ This class of status code is intended for situations in which the error seems to have been caused by the client.
```
- `400` -> **Bad Request**: The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing).
- `401` -> **Unauthorized**: Similar to `403 Forbidden`, but specifically for use when authentication is required and has failed or has not yet been provided.
- `402` -> **Payment Required**: Reserved for future use; originally intended for digital cash or micropayment schemes.
- `403` -> **Forbidden**: The request was valid, but the server is refusing action. The user might not have the necessary permissions for a resource.
- `404` -> **Not Found**: The requested resource could not be found but may be available in the future.
- `405` -> **Method Not Allowed**: A request method is not supported for the requested resource (e.g., a GET request on a form that requires POST).
- `406` -> **Not Acceptable**: The requested resource is capable of generating only content not acceptable according to the Accept headers sent in the request.
- `407` -> **Proxy Authentication Required**: The client must first authenticate itself with the proxy.
- `408` -> **Request Timeout**: The server timed out waiting for the request.
- `409` -> **Conflict**: Indicates that the request could not be processed because of conflict in the current state of the resource.
- `410` -> **Gone**: Indicates that the resource requested is no longer available and will not be available again.
- `411` -> **Length Required**: The request did not specify the length of its content, which is required by the requested resource.
- `412` -> **Precondition Failed**: The server does not meet one of the preconditions that the requester put on the request header fields.
- `413` -> **Payload Too Large**: The request is larger than the server is willing or able to process.
- `414` -> **URI Too Long**: The URI provided was too long for the server to process.
- `415` -> **Unsupported Media Type**: The request entity has a media type which the server or resource does not support.
- `416` -> **Range Not Satisfiable**: The client has asked for a portion of the file (byte range), but the server cannot supply that portion.
- `417` -> **Expectation Failed**: The server cannot meet the requirements of the Expect request-header field.
- `418` -> **I'm a teapot**: Defined in RFC 2324 as an April Fools' joke; returned by teapots requested to brew coffee.
- `421` -> **Misdirected Request**: The request was directed at a server that is not able to produce a response.
- `422` -> **Unprocessable Entity** *(WebDAV)*: The request was well-formed but was unable to be followed due to semantic errors.
- `423` -> **Locked** *(WebDAV)*: The resource that is being accessed is locked.
- `424` -> **Failed Dependency** *(WebDAV)*: The request failed because it depended on another request and that request failed.
- `425` -> **Too Early**: Indicates that the server is unwilling to risk processing a request that might be replayed.
- `426` -> **Upgrade Required**: The client should switch to a different protocol such as TLS/1.3.
- `428` -> **Precondition Required**: The origin server requires the request to be conditional.
- `429` -> **Too Many Requests**: The user has sent too many requests in a given amount of time ("rate limiting").
- `431` -> **Request Header Fields Too Large**: The server is unwilling to process the request because its header fields are too large.
- `451` -> **Unavailable For Legal Reasons**: A server operator has received a legal demand to deny access to a resource or to a set of resources.
```
**5xx Server Error Responses**

~ Response status codes beginning with the digit "5" indicate cases in which the server is aware that it has encountered an error or is otherwise incapable of performing the request.
```
- `500` -> **Internal Server Error**: A generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
- `501` -> **Not Implemented**: The server either does not recognize the request method, or it lacks the ability to fulfill the request.
- `502` -> **Bad Gateway**: The server was acting as a gateway or proxy and received an invalid response from the upstream server.
- `503` -> **Service Unavailable**: The server cannot handle the request (because it is overloaded or down for maintenance).
- `504` -> **Gateway Timeout**: The server was acting as a gateway or proxy and did not receive a timely response from the upstream server.
- `505` -> **HTTP Version Not Supported**: The server does not support the HTTP protocol version used in the request.
- `506` -> **Variant Also Negotiates**: Transparent content negotiation for the request results in a circular reference.
- `507` -> **Insufficient Storage** *(WebDAV)*: The server is unable to store the representation needed to complete the request.
- `508` -> **Loop Detected** *(WebDAV)*: The server detected an infinite loop while processing the request.
- `510` -> **Not Extended**: Further extensions to the request are required for the server to fulfill it.
- `511` -> **Network Authentication Required**: The client needs to authenticate to gain network access.
```

**Unofficial / Vendor-Specific Status Codes**

~ Commonly used custom codes encountered in web services and CDN providers:
```
- `419` -> **Page Expired** *(Laravel Framework)*: Used when a CSRF token is missing or expired.
- `420` -> **Method Failure** *(Spring Framework)* / **Enhance Your Calm** *(Twitter)*: Rate limiting warning.
- `499` -> **Client Closed Request** *(Nginx)*: Used when the client closes the connection while Nginx is processing the request.
- `520` -> **Web Server Returned an Unknown Error** *(Cloudflare)*: The origin server returned an empty, unknown, or unexplained response.
- `521` -> **Web Server Is Down** *(Cloudflare)*: The origin server refused the connection from Cloudflare.
- `522` -> **Connection Timed Out** *(Cloudflare)*: Cloudflare could not negotiate a TCP handshake with the origin server.
- `523` -> **Origin Is Unreachable** *(Cloudflare)*: Cloudflare could not reach the origin server (e.g., DNS issues).
- `524` -> **A Timeout Occurred** *(Cloudflare)*: Cloudflare established a TCP connection, but the origin did not reply with an HTTP response in time.
```
