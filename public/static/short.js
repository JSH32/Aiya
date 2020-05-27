function shorten() {
    let link = document.getElementById("linkarea").value

    axios({
        method: "post",
        url: "/api/create",
        data: {
            "link": link
        }
    }).then((data) => {
        let result = document.getElementById("linkarea")
        result.value = `${window.location.protocol}//${window.location.hostname}/${data.data}`

        // Copy to clipboard
        result.select()
        result.setSelectionRange(0, 99999)

        document.execCommand("copy")

        let count = document.getElementById("count").innerHTML
        let newNum = parseInt(count) + 1

        document.getElementById("count").innerHTML = newNum.toString()
    }).catch((data) => {
        document.getElementById("linkarea").value = ""
        document.getElementById("linkarea").setAttribute("placeholder", data.response.data)

        setTimeout(() => {
            document.getElementById("linkarea").setAttribute("placeholder", "Enter link")
        }, 2000);
    })
}