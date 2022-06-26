const SendMessageType = "Send"
const ChangeNameType = "ChangeName"
const NotifyType = "Notify"

window.onload = function() {
    var conn;
    var msg = document.getElementById("chat-text");
    var log = document.getElementById("log");
    var usernameHolder = document.getElementById("username-holder");
    var changeForm = document.getElementById("change-name-form")
    var urlParams = new URLSearchParams(window.location.search)

    var username = document.getElementById("username");
    username.value = usernameHolder.innerHTML = "User #" + Math.floor(Math.random() * 1000);

    changeForm.onsubmit = function(e) {
        console.log(username.value);
        usernameHolder.innerHTML = username.value;

        conn.send(JSON.stringify({
            type: ChangeNameType,
            username: username.value,
            content: username.value
        }));

        return false;
    }

    function appendLog(message) {
        var item;

        switch (message.type) {
            case NotifyType:
                item = createNotification(message)
                break;
            case ChangeNameType:
                item = createNotification(message)
                break;
            case SendMessageType:
                item = createMessageBlock(message)
                break;
        }

        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function() {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        conn.send(JSON.stringify({
            type: SendMessageType,
            username: username.value,
            content: msg.value
        }));

        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        let url = new URL("ws://" + document.location.host + "/ws");
        let pathArray = window.location.pathname.split('/');
        const roomName = pathArray[2];
        console.log(roomName);

        url.searchParams.append("username", username.value);
        url.searchParams.append("room", roomName);

        conn = new WebSocket(url.toString());

        // Load old messages
        conn.onopen = function(evt){
            let xhr = new XMLHttpRequest();

            xhr.onreadystatechange = function() {
                if (this.readyState === 4 && this.status === 200) {
                    let messages = JSON.parse(xhr.responseText)
                    messages.pop()
                    messages.forEach((item) => {
                        appendLog(item)
                    })
                }
            };

            xhr.open('GET', `/chat/${roomName}/history`, true);
            xhr.send()
        }

        conn.onclose = function(evt) {
            var item = createWaringBlock();
            appendLog(item);
        };

        conn.onmessage = function(evt) {
            var message = JSON.parse(evt.data);
            appendLog(message);
        };

    } else {
        var item = document.createElement("div");
        item.classList.add(["msg", "bg-warning"]);
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};

function createMessageBlock(message) {
    var user = message.username;
    var content = message.content;

    var item = document.createElement("div");
    item.setAttribute("class", "vstack gap-1 message-block flex-grow-0")

    if (username.value === user) {
        item.classList.add("user");
    }

    var usernameBlock = document.createElement("div");
    usernameBlock.setAttribute("class", "badge rounded-pill text-bg-dark block-username")
    usernameBlock.innerHTML = user;

    var messageBlock = document.createElement("div");
    messageBlock.classList.add("msg");
    messageBlock.innerHTML = content;

    username.value !== user ? item.appendChild(usernameBlock) : null;
    item.appendChild(messageBlock);

    return item;
}

function createWaringBlock() {
    var item = document.createElement("div");
    item.classList.add("vstack", "gap-1", "message-block", "flex-grow-0");

    var messageBlock = document.createElement("div");
    messageBlock.setAttribute("class", "msg bg-danger text-white")
    messageBlock.innerHTML = "Connection closed.";

    item.appendChild(messageBlock);

    return item;
}

function createNotification(message) {
    var item = document.createElement("div")
    item.setAttribute("class", "text-center text-muted fst-italic")

    item.innerHTML = message.content;
    return item;
}