(function () {
    let as = $('a[_method="DELETE"]')
    console.log('as', as)
    for (let i = 0, l = as.length; i < l; i++) {
        let $a = $(as[i])
        $a.click(function () {
            let message = $a.attr("message")
            let r = confirm(message)
            if (r) {
                let url = $a.attr("href")
                let method = "DELETE"
                $.ajax({
                    url: url,
                    type: method,
                    success: function (resp) {
                        location.reload()
                    },
                    error: function (resp) {
                        console.error(resp)
                    }
                })
            }
            return false
        })
    }
})()