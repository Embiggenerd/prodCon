function onLoad() {
    const source = new EventSource('/stream')
    source.onmessage = function (event) {
        const dataJSON = JSON.parse(event.data);

        console.log('dataz0', dataJSON.data)
        const dataSpan = document.getElementById("data");
        dataSpan.innerHTML = "message: " + dataJSON.data;
    };
}