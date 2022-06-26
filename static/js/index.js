var form = document.getElementById("add-room-form")

form.onsubmit = function(e) {
    e.preventDefault()
    var roomName = document.getElementById("room-name")
    var json = {
        name: roomName.value
    }

    let xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (this.readyState === 4 && this.status === 200) {
            window.location = 'http://' + window.location.host + `/chat/${roomName.value}`
        }
    }

    xhr.open('POST', '/create', true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.send(JSON.stringify(json));
}