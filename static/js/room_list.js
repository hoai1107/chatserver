window.onload = function (){
    // TODO: Get all name of room and add to template
    let xhr = new XMLHttpRequest();
    let roomList = document.getElementById('list');
    
    xhr.onreadystatechange = function() {
        if (this.readyState === 4 && this.status === 200) {
            let data = JSON.parse(xhr.responseText);

            for(let x in data.rooms){
                let item = createRoomLink(data.rooms[x])
                roomList.appendChild(item)
            }
        }
    };

    xhr.open('GET', '/all_rooms', true);
    xhr.send()
}

function createRoomLink(name){
    let item = document.createElement('a');

    item.innerHTML = name;
    item.setAttribute('class', 'btn btn-secondary w-25');
    item.setAttribute('href', `/chat/${name}`);

    return item;
}

