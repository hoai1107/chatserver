var form = document.getElementById("add-room-form")
form.onsubmit = function(e) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/create', true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    xhr.onreadystatechange = function() { // Call a function when the state changes.
        if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
            return false;
        }
    }

    var roomName = document.getElementById("room-name")
    var json = {
        name: roomName.value
    }
    xhr.send(JSON.stringify(json));
}