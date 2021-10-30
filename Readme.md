# Mumble UI
An attempt at creating an overlay that could be used inside of OBS to show who is speaking for the DAY[0] Podcast. It is meant to be reasonable extendable if you want to edit the code, but its written for my needs.

There are two components, the main webserver which serves `./web` as the webroot, and `./socket.io` which serves a Socket.IO (v2.20) socket that broadcasts events from the Mumble server into the browser.

## Message

**Broadcasts** - these are broadcast to all connects whenever the associated event happens
- `method=broadcast` - Whenever a user starts or stops speaking. 
  - The data is an `object`:
      ```
      {
          name: "the speakers mumble username",
          is_speaking: true|false
      }
      ```
- `method=join` - Whenever a new user is seen (on first connect all users will be joined)
  - The data is a `string` containing their username
- `method=leave` - Whenever a user leaves
  - The data is a `string` containing their username

**Actions** - these are messages that can be sent to the server
 - `method=action` 
 - data object:
     ```
     {
       "action":"user-list"
       "id":"",
       "data": null,
     }
     ```
   - All requests should fill in the `action` and `id` fields and can skip the `data` field.
   - Responses will use teh same `action` and `id` values as the request, with the response in `data`


## Arguments:
 - `-addr` - host:port for the target mumble server
 - `-user` - Username to connect as
 - `-pw` - Server password if any, currently no support for certificates 
 - `-verify=false` - Disable certificate validation against the server

The address, username, and password can also be provided through the following environment variables :
 - `MUMBLE_ADDR`
 - `MUMBLE_USER`
 - `MUMBLE_PW`

## Limitations:
I'm sure there are plenty more but a few of note: 
 - The bot joins the root of the mumble server, it cannot use a specific channel
 - No certificate based authentication (can't login to an existing account)
 - The display is configured to hide unknown users (see web/users.js) by default
