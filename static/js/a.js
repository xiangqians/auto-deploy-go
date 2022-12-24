(function () {
    let as = $('a[method]')
    for (let i = 0, l = as.length; i < l; i++) {
        let $a = $(as[i])
        // console.log($a)
        $a.click(function () {
            let formData = null
            let r = true
            let a = $a[0]
            let pre = a._pre_
            if (pre) {
                let rarr = pre()
                r = rarr && rarr[0] ? true : false
                let rl = 0
                if (r && (rl = rarr.length) > 1) {
                    formData = new FormData()
                    for (let ri = 1; ri < rl; ri++) {
                        let robj = rarr[ri]
                        for (let name in robj) {
                            formData.append(name, robj[name])
                        }
                    }
                }
            } else {
                let message = $a.attr("confirm")
                if (message) {
                    r = confirm(message)
                }
            }

            if (r) {
                let url = $a.attr("href")
                let method = $a.attr("method").trim().toUpperCase()
                // console.log(method, url, formData)
                $.ajax({
                    url: url,
                    type: method,
                    data: formData,
                    processData: false,
                    contentType: false,
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