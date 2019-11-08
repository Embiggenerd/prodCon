function onLoad() {
    const source = new EventSource('/stream')
    source.onmessage = function (event) {
        console.log('dataz0', event.data)
        const dataSpan = document.getElementById("data");
        dataSpan.innerHTML = "message: " + event.data;
    };
}