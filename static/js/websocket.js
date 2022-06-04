const SendMessageType = "Send"
const ChangeNameType = "ChangeName"

window.onload = function() {
    var conn;
    var msg = document.getElementById("chat-text");
    var log = document.getElementById("log");

    var username = document.getElementById("username");
    username.value = "User #" + Math.floor(Math.random() * 1000);


    function appendLog(item) {
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
        var url = new URL("ws://" + document.location.host + "/ws");
        url.searchParams.append("username", username.value);

        conn = new WebSocket(url.toString());

        conn.onclose = function(evt) {
            var item = createWaringBlock();
            appendLog(item);
        };

        conn.onmessage = function(evt) {
            var message = JSON.parse(evt.data);
            var item = createMessageBlock(message);
            appendLog(item);
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

    if (username.value == user) {
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