# Bug reproducer: Strangeness triggered by map key length

It isn't yet clear what the specific cause is, but the observed trigger is the length of the key in a map.

The current belief is that this is a reliable reproducer, given the specific message structure, as production
observations indicate specific data failing consistently.

To demonstrate the issue, run the following (N.B. see .tool-versions for the node version):

```sh
npm install
# works
node build/src/index.js Cn8Ke2FhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYRAB
# does not work
node build/src/index.js CoABCnxhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhEAE=
```

The actual reproducer is [index.ts](src/index.ts).

See also [main.go](main.go), which generated the data used in this reproducer.

N.B. Both Go and Java Protobuf implementations have been tested and work as expected.

**Succeeded:**

| Byte Range | Field Number | Type     |
|------------|--------------|----------|
| 0-129      | 1            | protobuf |

Content (field 1):

| Byte Range | Field Number | Type   | Content                                                                                                                     |
|------------|--------------|--------|-----------------------------------------------------------------------------------------------------------------------------|
| 0-125      | 1            | string | aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
| 125-127    | 2            | varint | As uint: 1  <br>As sint: -1                                                                                                 |

**Failed:**

| Byte Range | Field Number | Type     |
|------------|--------------|----------|
| 0-131      | 1            | protobuf |

Content (field 1):

| Byte Range | Field Number | Type   | Content                                                                                                                      |
|------------|--------------|--------|------------------------------------------------------------------------------------------------------------------------------|
| 0-126      | 1            | string | aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
| 126-128    | 2            | varint | As uint: 1  <br>As sint: -1                                                                                                  |
