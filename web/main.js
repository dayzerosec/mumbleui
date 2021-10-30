var state = {}
const socket = io({transports: ['websocket']});

socket.on("connect", () => {
    socket.emit('action', {"action":"user-list"})
});

socket.on("broadcast", (msg)=>{
    if(msg.name === undefined) return
    state[msg.name] = msg
})
socket.on("action", (msg)=>{
    switch(msg.action) {
        case "user-list":
            msg.data.forEach((v)=>{state[v.name] = {
                "name": v,
                "is_speaking": false,
            }})
            console.log(state)
            break;
    }
})